package service

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	mockscontext "github.com/nicopozo/mockserver/internal/context"
	mockserrors "github.com/nicopozo/mockserver/internal/errors"
	"github.com/nicopozo/mockserver/internal/model"
	jsonutils "github.com/nicopozo/mockserver/internal/utils/json"
	"github.com/yalp/jsonpath"
)

const max = 9999999999

//go:generate mockgen -destination=../utils/test/mocks/mock_service_mock.go -package=mocks -source=./mock_service.go

type MockService interface {
	SearchResponseForRequest(ctx context.Context, request *http.Request, path, body string) (model.Response, error)
}

func NewMockService(ruleService RuleService) (MockService, error) {
	if ruleService == nil {
		return nil, fmt.Errorf("rule service cannot be nil") //nolint:goerr113
	}

	return &mockService{
		RuleService: ruleService,
	}, nil
}

type mockService struct {
	RuleService RuleService
}

func (svc *mockService) SearchResponseForRequest(ctx context.Context,
	request *http.Request, path, body string,
) (model.Response, error) {
	logger := mockscontext.Logger(ctx)

	logger.Debug(svc, nil, "Entering mockService Execute()")

	method := strings.ToUpper(request.Method)

	rule, err := svc.RuleService.SearchByMethodAndPath(ctx, method, path)
	if err != nil {
		logger.Error(svc, nil, err, "error searching responses")

		return model.Response{}, fmt.Errorf("error searching rule, %w", err)
	}

	assertionResult := svc.applyAssertionsFromRule(rule)

	assertionResult.Print(ctx)

	if assertionResult.Fail {
		return model.Response{}, assertionResult.GetError() //nolint:wrapcheck
	}

	response, err := svc.getResponseFromRule(rule, request, body, path)
	if err != nil {
		return model.Response{}, err
	}

	variables, err := svc.getVariableValues(*request, body, rule, path)
	if err != nil {
		return model.Response{}, err
	}

	body = svc.applyVariables(response.Body, variables)

	response.Body = body

	return response, nil
}

func (svc *mockService) applyVariables(respBody string, variables []*model.Variable) string {
	for _, variable := range variables {
		respBody = strings.ReplaceAll(respBody, fmt.Sprintf("{%s}", variable.Name), variable.Value)
	}

	return respBody
}

//nolint:cyclop,funlen
func (svc *mockService) getResponseFromRule(rule model.Rule, request *http.Request, body string,
	path string,
) (model.Response, error) {
	strategy := rule.Strategy

	switch strategy {
	case model.RuleStrategyNormal:
		return rule.Responses[0], nil
	case model.RuleStrategyScene:
		var scene *model.Variable

		for _, variable := range rule.Variables {
			if variable.Name == model.RuleStrategyScene {
				scene = variable

				break
			}
		}

		if scene == nil {
			return model.Response{}, mockserrors.InvalidRulesError{
				Message: "rule doesn't have any variable names 'scene'",
			}
		}

		sceneName, err := svc.getVariableValue(*scene, request, body, rule, path)
		if err != nil {
			return model.Response{}, err
		}

		if sceneName != "" {
			first := string(sceneName[0])
			last := string(sceneName[len(sceneName)-1])

			// if it is a BODY variable, it is returned as JSON. So, I delete the "" from the beginning and the end
			if first == "\"" && last == "\"" {
				sceneName = sceneName[1 : len(sceneName)-1]
			}
		}

		respIndex := -1

		for index, resp := range rule.Responses {
			if resp.Scene == sceneName {
				respIndex = index

				break
			}

			if strings.ToLower(resp.Scene) == "default" {
				respIndex = index
			}
		}

		if respIndex >= 0 {
			return rule.Responses[respIndex], nil
		}

		return model.Response{}, mockserrors.InvalidRulesError{
			Message: fmt.Sprintf("rule doesn't have an scene called %s", sceneName),
		}
	case model.RuleStrategyRandom:
		responsesLen := len(rule.Responses)

		rand.Seed(time.Now().UnixNano())
		i := rand.Int63n(int64(responsesLen)) // nolint:gosec

		return rule.Responses[i], nil
	}

	return model.Response{}, mockserrors.InvalidRulesError{
		Message: "rule doesn't have a valid strategy",
	}
}

func (svc *mockService) getHeaderVariableValue(key string, request *http.Request) string {
	return request.Header.Get(key)
}

func (svc *mockService) getBodyVariableValue(key, body string) (string, error) {
	apply, err := jsonpath.Prepare(key)
	if err != nil {
		return "", fmt.Errorf("invalid JSON path for key %s - %w", key, err)
	}

	var reqMap interface{}
	if err := json.Unmarshal([]byte(body), &reqMap); err != nil {
		return "", fmt.Errorf("error unmarshalling request body. Body: %s: %w", body, err)
	}

	value, _ := apply(reqMap)

	return jsonutils.Marshal(value), nil
}

func (svc *mockService) getHashVariableValue() string {
	n := rand.Int63n(max) //nolint:gosec
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", n)))
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}

func (svc *mockService) getRandomVariableValue() string {
	rand.Seed(time.Now().UnixNano())
	n := rand.Int63n(max) //nolint:gosec

	return strconv.FormatInt(n, 10) //nolint:gomnd
}

func (svc *mockService) getQueryVariableValue(key string, request *http.Request) (string, error) {
	queries, err := url.ParseQuery(request.URL.RawQuery)
	if err != nil {
		return "", fmt.Errorf("error parsing queries %w", err)
	}

	for queryName, queryValue := range queries {
		if queryName == key {
			return queryValue[0], nil
		}
	}

	return "", mockserrors.InvalidRulesError{
		Message: fmt.Sprintf("no query param found with key %s", key),
	}
}

func (svc *mockService) getPathVariableValue(key, rulePath, reqPath string) (string, error) {
	pathVariables, err := svc.getPathParams(rulePath, reqPath)
	if err != nil {
		return "", err
	}

	if pathVariables[key] == "" {
		return "", mockserrors.InvalidRulesError{
			Message: fmt.Sprintf("no path param found with key %s", key),
		}
	}

	return pathVariables[key], nil
}

func (svc *mockService) getVariableValue(variable model.Variable, request *http.Request, body string,
	rule model.Rule, path string,
) (string, error) {
	switch variable.Type {
	case model.VariableTypeHeader:
		return svc.getHeaderVariableValue(variable.Key, request), nil
	case model.VariableTypeBody:
		return svc.getBodyVariableValue(variable.Key, body)
	case model.VariableTypeHash:
		return svc.getHashVariableValue(), nil
	case model.VariableTypeRandom:
		return svc.getRandomVariableValue(), nil
	case model.VariableTypeQuery:
		return svc.getQueryVariableValue(variable.Key, request)
	case model.VariableTypePath:
		return svc.getPathVariableValue(variable.Key, rule.Path, path)
	}

	return "", mockserrors.InvalidRulesError{
		Message: fmt.Sprintf("%s is invalid variable type", variable.Type),
	}
}

func (svc *mockService) getPathParams(rulePath, reqPath string) (map[string]string, error) {
	u, err := url.Parse(reqPath)
	if err != nil {
		return nil, fmt.Errorf("error parsing url, %w", err)
	}

	values := strings.Split(u.Path, "/")
	pathParts := strings.Split(rulePath, "/")

	params := make(map[string]string)

	for index, part := range pathParts {
		if part != "" {
			first := string(part[0])
			last := string(part[len(part)-1])

			if first == "{" && last == "}" {
				key := part[1 : len(part)-1]
				params[key] = values[index]
			}
		}
	}

	return params, nil
}

func (svc *mockService) getVariableValues(request http.Request, body string, rule model.Rule,
	path string,
) ([]*model.Variable, error) {
	variables := rule.Variables
	for idx := range variables {
		value, err := svc.getVariableValue(*variables[idx], &request, body, rule, path)
		if err != nil {
			return nil, err
		}

		variables[idx].Value = value
	}

	return variables, nil
}

func (svc *mockService) applyAssertionsFromRule(rule model.Rule) model.AssertionResult {
	result := model.AssertionResult{Fail: false}

	for _, variable := range rule.Variables {
		for _, assertion := range variable.Assertions {
			if msg, ok := assertion.Assert(variable); !ok {
				result.AddAssertionError(assertion.FailOnError, msg)
			}
		}
	}

	return result
}
