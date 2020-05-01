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
	mockscontext "github.com/nicopozo/mockserver/internal/context"
	ruleserrors "github.com/nicopozo/mockserver/internal/errors"
	"github.com/nicopozo/mockserver/internal/model"
	jsonutils "github.com/nicopozo/mockserver/internal/utils/json"
	stringutils "github.com/nicopozo/mockserver/internal/utils/string"
)

//nolint:lll
//go:generate mockgen -destination=../utils/test/mocks/rule_repository_mock.go -package=mocks -source=./rule_repository.go

type IRuleRepository interface {
	Save(ctx context.Context, rule *model.Rule) error
	Get(ctx context.Context, key string) (*model.Rule, error)
	Search(ctx context.Context, params map[string]interface{}, paging model.Paging) (*model.RuleList, error)
	SearchByMethodAndPath(ctx context.Context, method string, path string) (*model.Rule, error)
}

type RuleElasticRepository struct {
	client *elasticsearch.Client
}

func (repository *RuleElasticRepository) Save(ctx context.Context, rule *model.Rule) error {
	logger := mockscontext.Logger(ctx)

	var err error

	rule.Key = getKey(rule)

	_, err = repository.createPatterns(ctx, rule)
	if err != nil {
		return err
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
		return errors.New("error saving rule - " + newStr)
	}

	return nil
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
		return nil, err
	}

	if getRuleResp.IsError() {
		if getRuleResp.StatusCode == http.StatusNotFound {
			msg := fmt.Sprintf("no rule found for key: %v", key)
			logger.Debug(repository, nil, msg)

			return nil, ruleserrors.RuleNotFoundError{Message: msg}
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
		return nil, err
	}

	return rule, err
}

func (repository *RuleElasticRepository) Search(ctx context.Context, params map[string]interface{},
	paging model.Paging) (*model.RuleList, error) {
	logger := mockscontext.Logger(ctx)

	// Build the request body.
	var buf bytes.Buffer

	terms := make([]map[string]interface{}, len(params))

	for key, value := range params {
		term := map[string]interface{}{
			"term": map[string]interface{}{
				key: value,
			},
		}
		terms = append(terms, term)
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": terms,
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
		return nil, err
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
		return nil, err
	}

	result := &model.RuleList{}
	if esResponse.Hits != nil {
		paging.Total = int64(esResponse.Hits.Total.Value)
		result.Paging = paging
		for _, esRule := range esResponse.Hits.Hits {
			result.Results = append(result.Results, esRule.Source)
		}
	}

	return result, nil
}

func (repository *RuleElasticRepository) Delete(ctx context.Context, application, name string) error {
	//logger := mockscontext.Logger(ctx)

	/*var err error

	gorelic.WrapDatastoreSegment("KVS", "DELETE", txn, func() {
		err = repository.getRulesKvsClient().Delete("getKey(application, name)")
	})

	if err != nil {
		logger.Error(repository, nil, err, "error deleting rule from KVS")
		return err
	}
	*/
	return nil
}

func (repository *RuleElasticRepository) SearchByMethodAndPath(ctx context.Context, method string,
	path string) (*model.Rule, error) {
	logger := mockscontext.Logger(ctx)

	var err error

	getExprReq := esapi.GetRequest{
		DocumentID: strings.ToLower(method),
		Index:      "expressions",
	}

	getExprResp, err := getExprReq.Do(context.Background(), repository.getElasticClient())
	if getExprResp != nil {
		defer closeBody(getExprResp.Body, repository, logger)
	}

	if err != nil {
		logger.Error(repository, nil, err, "error getting expressions from Elastic Search")
		return nil, err
	}

	if getExprResp.IsError() {
		if getExprResp.StatusCode == http.StatusNotFound {
			msg := fmt.Sprintf("no expression found for method: %v", method)
			logger.Debug(repository, nil, msg)
			return nil, ruleserrors.RuleNotFoundError{Message: msg}
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
		return nil, err
	}

	var ruleKey string

	for _, pattern := range patternList.Patterns {
		var regex = regexp.MustCompile(pattern.PathExpression)
		if regex.MatchString(path) {
			ruleKey = pattern.RuleKeys[0]
			break
		}
	}

	if ruleKey == "" {
		return nil, ruleserrors.RuleNotFoundError{
			Message: "no rule found for path: " + path,
		}
	}

	getRuleReq := esapi.GetRequest{
		DocumentID: strings.ToLower(ruleKey),
		Index:      "rules",
	}

	getRuleResp, err := getRuleReq.Do(context.Background(), repository.getElasticClient())
	if getRuleResp != nil {
		defer closeBody(getRuleResp.Body, repository, logger)
	}

	if err != nil {
		logger.Error(repository, nil, err, "error getting rule from Elastic Search")
		return nil, err
	}

	if getRuleResp.IsError() {
		if getRuleResp.StatusCode == http.StatusNotFound {
			msg := fmt.Sprintf("no rule found for key: %v", ruleKey)
			logger.Debug(repository, nil, msg)

			return nil, ruleserrors.RuleNotFoundError{Message: msg}
		}
		logger.Error(repository, nil, errors.New("http status != 200: Actual: "+getRuleResp.String()),
			"error getting expressions from Elastic Search")
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(getRuleResp.Body)
		newStr := buf.String()

		return nil, errors.New("error getting expressions from Elastic Search - " + newStr)
	}

	var rule *model.Rule

	rule, err = model.UnmarshalESRule(getRuleResp.Body)
	if err != nil {
		return nil, err
	}

	return rule, err
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
		return nil, err
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
			return nil, err
		}
		list = *l

	} else {
		list.Patterns = make([]*model.Pattern, 0)
	}

	expression := createExpression(rule.Path)

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

	saveReq := esapi.IndexRequest{
		Index:      "expressions",
		DocumentID: rule.Method,
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
		return nil, errors.New("error saving Pattern - " + newStr)
	}

	return pattern, nil
}

func createExpression(path string) string {
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

func getKey(rule *model.Rule) string {
	return fmt.Sprintf("%v_%v_%v", rule.Application, rule.Method, stringutils.Hash(rule.Name))
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
