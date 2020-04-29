package repository

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	ruleserrors "github.com/nicopozo/mockserver/internal/errors"

	jsonutils "github.com/nicopozo/mockserver/internal/utils/json"

	"github.com/elastic/go-elasticsearch/v7/esapi"

	stringutils "github.com/nicopozo/mockserver/internal/utils/string"

	"github.com/elastic/go-elasticsearch/v7"
	newrelic "github.com/newrelic/go-agent"
	"github.com/nicopozo/mockserver/internal/model"
	"github.com/nicopozo/mockserver/internal/utils/log"
)

//nolint:lll
//go:generate mockgen -destination=../utils/test/mocks/rule_repository_mock.go -package=mocks -source=./rule_repository.go

type IRuleRepository interface {
	Save(rule *model.Rule, txn newrelic.Transaction, logger log.ILogger) error
	Get(key string, txn newrelic.Transaction, logger log.ILogger) (*model.Rule, error)
	Search(params map[string]interface{}, paging model.Paging, txn newrelic.Transaction,
		logger log.ILogger) (*model.RuleList, error)
	SearchByMethodAndPath(method string, path string, txn newrelic.Transaction,
		logger log.ILogger) (*model.Rule, error)
}

type RuleElasticRepository struct {
	client *elasticsearch.Client
}

func (repository *RuleElasticRepository) Save(rule *model.Rule, txn newrelic.Transaction, logger log.ILogger) error {
	var err error

	rule.Key = getKey(rule)

	_, err = repository.createPatterns(rule, txn, logger)
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

func (repository *RuleElasticRepository) Get(key string, txn newrelic.Transaction,
	logger log.ILogger) (*model.Rule, error) {
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

func (repository *RuleElasticRepository) Search(params map[string]interface{}, paging model.Paging, txn newrelic.Transaction,
	logger log.ILogger) (*model.RuleList, error) {
	return nil, nil
	/*	var err error

		var response *godsclient.SearchResponse

		queryBuilder := &querybuilders.AndQueryBuilder{}
		for key, value := range params {
			queryBuilder = queryBuilder.And(querybuilders.Eq(key, value))
		}

		gorelic.WrapDatastoreSegment("DS", "SEARCH", txn, func() {
			response, err = repository.getDSClient().
				SearchBuilder().
				AddSort(sortbuilders.Field("id", fieldtype.Number).Order(sortorder.Desc)).
				WithQuery(queryBuilder).
				WithOffset(paging.Offset).
				WithSize(paging.Limit).
				Execute()
		})

		if err != nil {
			logger.Error(repository, nil, err, "error executing Document Search query")
			return nil, err
		}

		results := make([]*model.Rule, 0, len(response.Documents))

		for _, dsItem := range response.Documents {
			repository := bytes.NewReader(dsItem)
			rule := &model.Rule{}

			if err := jsonutils.Unmarshal(repository, rule); err != nil {
				logger.Error(repository, nil, err, "error parsing reconcilable from Document Search")
				return nil, err
			}

			results = append(results, rule)
		}

		paging.Total = response.Total

		return &model.RuleList{Paging: paging, Results: results}, nil*/
}

func (repository *RuleElasticRepository) Delete(application, name string, txn newrelic.Transaction, logger log.ILogger) error {
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

func (repository *RuleElasticRepository) SearchByMethodAndPath(method string, path string, txn newrelic.Transaction,
	logger log.ILogger) (*model.Rule, error) {
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

func (repository *RuleElasticRepository) createPatterns(rule *model.Rule, txn newrelic.Transaction,
	logger log.ILogger) (*model.Pattern, error) {
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

func closeBody(body io.ReadCloser, repository *RuleElasticRepository, logger log.ILogger) {
	err := body.Close()
	if err != nil {
		logger.Error(repository, nil, err, "error closing response body")
	}
}
