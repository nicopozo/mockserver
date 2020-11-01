package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	guuid "github.com/google/uuid"
	mockscontext "github.com/nicopozo/mockserver/internal/context"
	mockserrors "github.com/nicopozo/mockserver/internal/errors"
	"github.com/nicopozo/mockserver/internal/model"
	jsonutils "github.com/nicopozo/mockserver/internal/utils/json"
	stringutils "github.com/nicopozo/mockserver/internal/utils/string"
)

//nolint:lll
//go:generate mockgen -destination=../utils/test/mocks/rule_repository_mock.go -package=mocks -source=./rule_repository.go

type IRuleRepository interface {
	Save(ctx context.Context, rule *model.Rule, isUpdate bool) (*model.Rule, error)
	Get(ctx context.Context, key string) (*model.Rule, error)
	Search(ctx context.Context, params map[string]interface{}, paging model.Paging) (*model.RuleList, error)
	SearchByMethodAndPath(ctx context.Context, method string, path string) (*model.Rule, error)
	Delete(ctx context.Context, key string) error
}

type RuleElasticRepository struct {
	client *elasticsearch.Client
}

func (repository *RuleElasticRepository) Save(ctx context.Context, rule *model.Rule,
	isUpdate bool) (*model.Rule, error) {
	logger := mockscontext.Logger(ctx)

	var err error

	if !isUpdate {
		rule.Key = fmt.Sprintf("%v", guuid.New())
	}

	_, err = repository.createPatterns(ctx, rule)
	if err != nil {
		return nil, err
	}

	req := esapi.IndexRequest{
		Index:      "rules",
		DocumentID: rule.Key,
		Body:       strings.NewReader(jsonutils.Marshal(rule)),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), repository.getElasticClient())
	if err != nil {
		logger.Error(repository, nil, err, "Error getting response: %s", err)
	}

	if res != nil {
		defer closeBody(res.Body, repository, logger)
	}

	if res.IsError() {
		logger.Error(repository, nil, nil, "[%s] Error indexing document Key: %v", res.Status(), stringutils.Hash(rule.Key))

		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(res.Body)
		newStr := buf.String()

		return nil, errors.New("error saving rule - " + newStr)
	}

	return rule, nil
}

func (repository *RuleElasticRepository) Get(ctx context.Context, key string) (*model.Rule, error) {
	logger := mockscontext.Logger(ctx)

	var err error

	getRuleReq := esapi.GetRequest{
		DocumentID: strings.ToLower(key),
		Index:      "rules",
	}

	getRuleResp, err := getRuleReq.Do(context.Background(), repository.getElasticClient())
	if getRuleResp != nil {
		defer closeBody(getRuleResp.Body, repository, logger)
	}

	if err != nil {
		logger.Error(repository, nil, err, "error getting rule from Elastic Search")

		return nil, fmt.Errorf("error reading Elastic Search Response, %w", err)
	}

	if getRuleResp.IsError() {
		if getRuleResp.StatusCode == http.StatusNotFound {
			msg := fmt.Sprintf("no rule found for key: %v", key)
			logger.Debug(repository, nil, msg)

			return nil, mockserrors.RuleNotFoundError{Message: msg}
		}

		logger.Error(repository, nil, errors.New("http status != 200: Actual: "+getRuleResp.String()),
			"error getting expression from Elastic Search")

		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(getRuleResp.Body)
		newStr := buf.String()

		return nil, errors.New("error getting expressions from Elastic Search - " + newStr)
	}

	var rule *model.Rule

	rule, err = model.UnmarshalESRule(getRuleResp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading Elastic Search Response, %w", err)
	}

	return rule, nil
}

func (repository *RuleElasticRepository) Search(ctx context.Context, params map[string]interface{},
	paging model.Paging) (*model.RuleList, error) {
	logger := mockscontext.Logger(ctx)

	// Build the request body.
	var buf bytes.Buffer

	terms := make([]map[string]interface{}, 0)

	for key, value := range params {
		if key == "method" {
			value = strings.ToLower(value.(string))
		}

		term := map[string]interface{}{
			"wildcard": map[string]interface{}{
				key: map[string]interface{}{
					"value": fmt.Sprintf("*%v*", value),
				},
			},
		}
		terms = append(terms, term)
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": terms,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Perform the search request.
	client := repository.getElasticClient()
	searchRuleResp, err := client.Search(
		client.Search.WithContext(context.Background()),
		client.Search.WithIndex("rules"),
		client.Search.WithSize(int(paging.Limit)),
		client.Search.WithFrom(int(paging.Offset)),
		client.Search.WithBody(&buf),
		client.Search.WithTrackTotalHits(true),
		client.Search.WithPretty(),
	)

	if searchRuleResp != nil {
		defer closeBody(searchRuleResp.Body, repository, logger)
	}

	if err != nil {
		logger.Error(repository, nil, err, "error getting rule from Elastic Search")

		return nil, fmt.Errorf("error reading Elastic Search Response, %w", err)
	}

	if searchRuleResp.IsError() {
		logger.Error(repository, nil, errors.New("http status != 200: Actual: "+searchRuleResp.String()),
			"error getting rules from Elastic Search")

		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(searchRuleResp.Body)
		newStr := buf.String()

		return nil, errors.New("error searching rules from Elastic Search - " + newStr)
	}

	var esResponse *model.ESSearchResult

	esResponse, err = model.UnmarshalSearchESRule(searchRuleResp.Body)
	if err != nil {
		logger.Error(repository, nil, err, "error unmarshalling results from Elastic Search")

		return nil, fmt.Errorf("error reading Elastic Search Response, %w", err)
	}

	result := &model.RuleList{}
	result.Results = []*model.Rule{}

	if esResponse.Hits != nil {
		paging.Total = int64(esResponse.Hits.Total.Value)
		result.Paging = paging

		for _, esRule := range esResponse.Hits.Hits {
			result.Results = append(result.Results, esRule.Source)
		}
	}

	return result, nil
}

func (repository *RuleElasticRepository) Delete(ctx context.Context, key string) error {
	logger := mockscontext.Logger(ctx)

	rule, err := repository.Get(ctx, key)
	if err != nil {
		return err
	}

	patternList, err := repository.getExpressionsByMethod(ctx, rule.Method)
	if err != nil {
		return err
	}

	for index, pattern := range patternList.Patterns {
		var regex = regexp.MustCompile(pattern.PathExpression)
		if regex.MatchString(rule.Path) {
			for i, ruleKey := range pattern.RuleKeys {
				if ruleKey == key {
					pattern.RuleKeys = append(pattern.RuleKeys[:i], pattern.RuleKeys[i+1:]...)
					patternList.Patterns[index] = pattern

					break
				}
			}

			break
		}
	}

	if err = repository.saveExpression(ctx, rule.Method, patternList); err != nil {
		return err
	}

	req := esapi.DeleteRequest{
		Index:      "rules",
		DocumentID: key,
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), repository.getElasticClient())
	if err != nil {
		logger.Error(repository, nil, err, "Error getting response: %s", err)
	}

	if res != nil {
		defer closeBody(res.Body, repository, logger)
	}

	if res.StatusCode == http.StatusNotFound || res.StatusCode == http.StatusOK {
		return nil
	}

	if res.IsError() {
		logger.Error(repository, nil, nil, "[%s] Error document document Key: %v", res.Status(), stringutils.Hash(key))

		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(res.Body)
		newStr := buf.String()

		return errors.New("error deleting rule - " + newStr)
	}

	return nil
}

func (repository *RuleElasticRepository) SearchByMethodAndPath(ctx context.Context, method string,
	path string) (*model.Rule, error) {
	var err error

	patternList, err := repository.getExpressionsByMethod(ctx, method)
	if err != nil {
		return nil, err
	}

	for _, pattern := range patternList.Patterns {
		var regex = regexp.MustCompile(pattern.PathExpression)

		if regex.MatchString(path) {
			for _, ruleKey := range pattern.RuleKeys {
				rule, _ := repository.Get(ctx, ruleKey)
				if rule != nil && rule.Status == model.RuleStatusEnabled {
					return rule, nil
				}
			}
		}
	}

	return nil, mockserrors.RuleNotFoundError{
		Message: fmt.Sprintf("no rule found for path: %s and method %s", path, method),
	}
}

func (repository *RuleElasticRepository) getExpressionsByMethod(ctx context.Context,
	method string) (*model.PatternList, error) {
	logger := mockscontext.Logger(ctx)

	var err error

	getExprReq := esapi.GetRequest{
		DocumentID: strings.ToUpper(method),
		Index:      "expressions",
	}

	getExprResp, err := getExprReq.Do(context.Background(), repository.getElasticClient())
	if getExprResp != nil {
		defer closeBody(getExprResp.Body, repository, logger)
	}

	if err != nil {
		logger.Error(repository, nil, err, "error getting expressions from Elastic Search")

		return nil, fmt.Errorf("error reading Elastic Search Response, %w", err)
	}

	if getExprResp.IsError() {
		if getExprResp.StatusCode == http.StatusNotFound {
			msg := fmt.Sprintf("no expression found for method: %v", method)
			logger.Debug(repository, nil, msg)

			return nil, mockserrors.RuleNotFoundError{Message: msg}
		}

		logger.Error(repository, nil, nil, "error getting expressions from Elastic Search")

		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(getExprResp.Body)
		newStr := buf.String()

		return nil, errors.New("error getting expressions from Elastic Search - " + newStr)
	}

	var patternList *model.PatternList

	patternList, err = model.UnmarshalESPatternList(getExprResp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading Elastic Search Response, %w", err)
	}

	return patternList, nil
}

func (repository *RuleElasticRepository) createPatterns(ctx context.Context, rule *model.Rule) (*model.Pattern, error) {
	logger := mockscontext.Logger(ctx)

	var err error

	getReq := esapi.GetRequest{
		DocumentID: rule.Method,
		Index:      "expressions",
	}

	getResp, err := getReq.Do(context.Background(), repository.getElasticClient())
	if getResp != nil {
		defer closeBody(getResp.Body, repository, logger)
	}

	if err != nil {
		logger.Error(repository, nil, err, "error getting expressions from Elastic Search")

		return nil, fmt.Errorf("error reading Elastic Search Response, %w", err)
	}

	if getResp.IsError() {
		if getResp.StatusCode != http.StatusNotFound {
			logger.Error(repository, nil, nil, "error getting expressions from Elastic Search")

			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(getResp.Body)
			newStr := buf.String()

			return nil, errors.New("error getting expressions from Elastic Search - " + newStr)
		}
	}

	var list model.PatternList

	if getResp.StatusCode != http.StatusNotFound {
		l, err := model.UnmarshalESPatternList(getResp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading Elastic Search Response, %w", err)
		}

		list = *l
	} else {
		list.Patterns = make([]*model.Pattern, 0)
	}

	expression := CreateExpression(rule.Path)

	var pattern *model.Pattern

	foundPatterIndex := -1

	for i, p := range list.Patterns {
		if p.PathExpression == expression {
			foundPatterIndex = i

			break
		}
	}

	if foundPatterIndex >= 0 {
		foundKey := false

		for _, key := range list.Patterns[foundPatterIndex].RuleKeys {
			if key == rule.Key {
				foundKey = true

				break
			}
		}

		if foundKey {
			return list.Patterns[foundPatterIndex], nil
		}

		list.Patterns[foundPatterIndex].RuleKeys = append(list.Patterns[foundPatterIndex].RuleKeys, rule.Key)

		pattern = list.Patterns[foundPatterIndex]
	} else {
		pattern = &model.Pattern{
			PathExpression: expression,
			RuleKeys:       []string{rule.Key},
		}
		list.Patterns = append(list.Patterns, pattern)
	}

	err = repository.saveExpression(ctx, strings.ToUpper(rule.Method), &list)
	if err != nil {
		return nil, err
	}

	return pattern, nil
}

func (repository *RuleElasticRepository) saveExpression(ctx context.Context, method string,
	list *model.PatternList) error {
	logger := mockscontext.Logger(ctx)

	saveReq := esapi.IndexRequest{
		Index:      "expressions",
		DocumentID: method,
		Body:       strings.NewReader(jsonutils.Marshal(list)),
		Refresh:    "true",
	}

	saveResp, err := saveReq.Do(context.Background(), repository.getElasticClient())
	if err != nil {
		logger.Error(repository, nil, err, "error saving Pattern in Elastic Search")
	}

	if saveResp != nil {
		defer closeBody(saveResp.Body, repository, logger)
	}

	if saveResp.IsError() {
		logger.Error(repository, nil, nil, "error saving Pattern in Elastic Search")

		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(saveResp.Body)
		newStr := buf.String()

		return errors.New("error saving Pattern - " + newStr)
	}

	return nil
}

func (repository *RuleElasticRepository) getElasticClient() *elasticsearch.Client {
	if repository.client == nil {
		cfg := elasticsearch.Config{
			Addresses: []string{
				"http://localhost:9200",
			},
		}

		es, err := elasticsearch.NewClient(cfg)
		if err != nil {
			fmt.Println("Que hacer aqui!!")
		}

		repository.client = es
	}

	return repository.client
}

func CreateExpression(path string) string {
	var paramRegex = regexp.MustCompile("{.+?}/")
	params := paramRegex.FindAllString(path, -1)

	for _, param := range params {
		path = strings.ReplaceAll(path, param, "[^/]+?/")
	}

	paramRegex = regexp.MustCompile("{.+?}")
	params = paramRegex.FindAllString(path, -1)

	for _, param := range params {
		path = strings.ReplaceAll(path, param, "[^/]+")
	}

	return fmt.Sprintf("^%s$", path)
}
