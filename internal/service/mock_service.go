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
	request *http.Request, path, body string) (model.Response, error) {
	logger := mockscontext.Logger(ctx)

	logger.Debug(svc, nil, "Entering mockService Execute()")

	method := strings.ToUpper(request.Method)

	rule, err := svc.RuleService.SearchByMethodAndPath(ctx, method, path)
	if err != nil {
		logger.Error(svc, nil, err, "error searching responses")

		return model.Response{}, fmt.Errorf("error searching rule, %w", err)
	}

	response, err := svc.getResponseFromRule(rule, request, body, path)
	if err != nil {
		return model.Response{}, err
	}

	body, err = svc.applyVariables(request, body, response, rule, path)
	if err != nil {
		return model.Response{}, err
	}

	response.Body = body

	return response, nil
}

//nolint:cyclop
func (svc *mockService) applyVariables(request *http.Request, reqBody string, response model.Response,
	rule *model.Rule, path string) (string, error) {
	var err error

	respBody := response.Body

	for _, variable := range rule.Variables {
		switch variable.Type {
		case model.VariableTypeHeader:
			respBody = svc.applyHeaderVariables(request, respBody, variable.Name, variable.Key)
		case model.VariableTypeBody:
			respBody, err = svc.applyBodyVariables(reqBody, respBody, variable.Name, variable.Key)
			if err != nil {
				return respBody, err
			}
		case model.VariableTypeHash:
			respBody = svc.applyHashVariables(respBody, variable.Name)
		case model.VariableTypeRandom:
			respBody = svc.applyRandomVariables(respBody, variable.Name)
		case model.VariableTypeQuery:
			respBody, err = svc.applyQueryVariables(request, respBody, variable.Name, variable.Key)
			if err != nil {
				return respBody, err
			}
		case model.VariableTypePath:
			respBody, err = svc.applyPathVariable(respBody, rule.Path, path, variable.Name, variable.Key)
			if err != nil {
				return respBody, err
			}
		}
	}

	return respBody, err
}

func (svc *mockService) applyPathVariable(responseBody, rulePath, reqPath, variableName,
	variableKey string) (string, error) {
	params, err := svc.getPathParams(rulePath, reqPath)
	if err != nil {
		return "", err
	}

	for paramKey, paramValue := range params {
		if paramKey == variableKey {
			return strings.ReplaceAll(responseBody, fmt.Sprintf("{%s}", variableName), paramValue), nil
		}
	}

	return responseBody, nil
}

func (svc *mockService) applyBodyVariables(reqBody, respBody string, name string,
	key string) (string, error) {
	value, err := svc.getBodyVariableValue(key, reqBody)
	if err != nil {
		return "", err
	}

	respBody = strings.ReplaceAll(respBody, fmt.Sprintf("{%s}", name), value)

	return respBody, nil
}

func (svc *mockService) applyRandomVariables(respBody string, name string) string {
	return strings.ReplaceAll(respBody, fmt.Sprintf("{%s}", name), svc.getRandomVariableValue())
}

func (svc *mockService) applyHashVariables(respBody string, name string) string {
	return strings.ReplaceAll(respBody, fmt.Sprintf("{%s}", name), svc.getHashVariableValue())
}

func (svc *mockService) applyQueryVariables(request *http.Request, body string, name string,
	key string) (string, error) {
	queryValue, err := svc.getQueryVariableValue(key, request)
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(body, fmt.Sprintf("{%s}", name), queryValue), nil
}

func (svc *mockService) applyHeaderVariables(request *http.Request, body string, name string, key string) string {
	header := svc.getHeaderVariableValue(key, request)
	if header != "" {
		body = strings.ReplaceAll(body, fmt.Sprintf("{%s}", name), header)
	}

	return body
}

//nolint:cyclop,funlen
func (svc *mockService) getResponseFromRule(rule *model.Rule, request *http.Request, body string,
	path string) (model.Response, error) {
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

	value, err := apply(reqMap)
	if err != nil {
		return "", fmt.Errorf("error getting value from JSON PATH. Request: %v: %w", reqMap, err)
	}

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
	rule *model.Rule, path string) (string, error) {
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
