package service

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
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

type IMockService interface {
	SearchResponseForRequest(ctx context.Context, request *http.Request, path string) (*model.Response, error)
}

type MockService struct {
	RuleService IRuleService
}

func (service *MockService) SearchResponseForRequest(ctx context.Context,
	request *http.Request, path string) (*model.Response, error) {
	logger := mockscontext.Logger(ctx)

	logger.Debug(service, nil, "Entering MockService Execute()")

	method := strings.ToUpper(request.Method)

	rule, err := service.RuleService.SearchByMethodAndPath(ctx, method, path)
	if err != nil {
		logger.Error(service, nil, err, "error searching responses")

		return nil, fmt.Errorf("error searching rule, %w", err)
	}

	reqBodyBuilder := new(strings.Builder)

	_, err = io.Copy(reqBodyBuilder, request.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body, %w", err)
	}

	reqBody := reqBodyBuilder.String()

	response, err := service.getResponseFromRule(rule, request, reqBody, path)
	if err != nil {
		return nil, err
	}

	body, err := service.applyVariables(request, reqBody, response, rule, path)
	if err != nil {
		return nil, err
	}

	response.Body = body

	return response, nil
}

func (service *MockService) applyVariables(request *http.Request, reqBody string, response *model.Response,
	rule *model.Rule, path string) (string, error) {
	var err error

	respBody := response.Body

	for _, v := range rule.Variables {
		switch v.Type {
		case model.VariableTypeHeader:
			respBody = service.applyHeaderVariables(request, respBody, v.Name, v.Key)
		case model.VariableTypeBody:
			respBody, err = service.applyBodyVariables(reqBody, respBody, v.Name, v.Key)
			if err != nil {
				return respBody, err
			}
		case model.VariableTypeHash:
			respBody = service.applyHashVariables(respBody, v.Name)
		case model.VariableTypeRandom:
			respBody = service.applyRandomVariables(respBody, v.Name)
		case model.VariableTypeQuery:
			respBody, err = service.applyQueryVariables(request, respBody, v.Name, v.Key)
			if err != nil {
				return respBody, err
			}
		case model.VariableTypePath:
			respBody, err = service.applyPathVariables(respBody, rule.Path, path)
			if err != nil {
				return respBody, err
			}
		}
	}

	return respBody, err
}

func (service *MockService) applyPathVariables(responseBody, rulePath, reqPath string) (string, error) {
	params, err := service.getPathParams(rulePath, reqPath)
	if err != nil {
		return "", err
	}

	for key, value := range params {
		responseBody = strings.ReplaceAll(responseBody, fmt.Sprintf("{%s}", key), value)
	}

	return responseBody, nil
}

func (service *MockService) applyBodyVariables(reqBody, respBody string, name string,
	key string) (string, error) {
	value, err := service.getBodyVariableValue(key, reqBody)
	if err != nil {
		return "", err
	}

	respBody = strings.ReplaceAll(respBody, fmt.Sprintf("{%s}", name), value)

	return respBody, nil
}

func (service *MockService) applyRandomVariables(respBody string, name string) string {
	return strings.ReplaceAll(respBody, fmt.Sprintf("{%s}", name), service.getRandomVariableValue())
}

func (service *MockService) applyHashVariables(respBody string, name string) string {
	return strings.ReplaceAll(respBody, fmt.Sprintf("{%s}", name), service.getHashVariableValue())
}

func (service *MockService) applyQueryVariables(request *http.Request, body string, name string,
	key string) (string, error) {
	queryValue, err := service.getQueryVariableValue(key, request)
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(body, fmt.Sprintf("{%s}", name), queryValue), nil
}

func (service *MockService) applyHeaderVariables(request *http.Request, body string, name string, key string) string {
	header := service.getHeaderVariableValue(key, request)
	if header != "" {
		body = strings.ReplaceAll(body, fmt.Sprintf("{%s}", name), header)
	}

	return body
}

func (service *MockService) getResponseFromRule(rule *model.Rule, request *http.Request, body string,
	path string) (*model.Response, error) {
	strategy := rule.Strategy

	switch strategy {
	case model.RuleStrategyNormal:
		return &rule.Responses[0], nil
	case model.RuleStrategyScene:
		var scene *model.Variable

		for _, variable := range rule.Variables {
			if variable.Name == model.RuleStrategyScene {
				scene = variable

				break
			}
		}

		if scene == nil {
			return nil, mockserrors.InvalidRulesError{
				Message: "rule doesn't have any variable names 'scene'",
			}
		}

		sceneName, err := service.getVariableValue(*scene, request, body, rule, path)
		if err != nil {
			return nil, err
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

		for i, resp := range rule.Responses {
			if resp.Scene == sceneName {
				respIndex = i

				break
			}

			if strings.ToLower(resp.Scene) == "default" {
				respIndex = i
			}
		}

		if respIndex >= 0 {
			return &rule.Responses[respIndex], nil
		}

		return nil, mockserrors.InvalidRulesError{
			Message: fmt.Sprintf("rule doesn't have an scene called %s", sceneName),
		}
	case model.RuleStrategyRandom:
		responsesLen := len(rule.Responses)

		rand.Seed(time.Now().UnixNano())
		i := rand.Int63n(int64(responsesLen)) // nolint:gosec

		return &rule.Responses[i], nil
	}

	return nil, mockserrors.InvalidRulesError{
		Message: "rule doesn't have a valid strategy",
	}
}

func (service *MockService) getHeaderVariableValue(key string, request *http.Request) string {
	return request.Header.Get(key)
}

func (service *MockService) getBodyVariableValue(key, body string) (string, error) {
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

func (service *MockService) getHashVariableValue() string {
	n := rand.Int63n(max) //nolint:gosec
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", n))) //nolint:errcheck
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}

func (service *MockService) getRandomVariableValue() string {
	rand.Seed(time.Now().UnixNano())
	n := rand.Int63n(max) //nolint:gosec

	return strconv.FormatInt(n, 10)
}

func (service *MockService) getQueryVariableValue(key string, request *http.Request) (string, error) {
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

func (service *MockService) getPathVariableValue(key, rulePath, reqPath string) (string, error) {
	pathVariables, err := service.getPathParams(rulePath, reqPath)
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

func (service *MockService) getVariableValue(v model.Variable, request *http.Request, body string,
	rule *model.Rule, path string) (string, error) {
	switch v.Type {
	case model.VariableTypeHeader:
		return service.getHeaderVariableValue(v.Key, request), nil
	case model.VariableTypeBody:
		return service.getBodyVariableValue(v.Key, body)
	case model.VariableTypeHash:
		return service.getHashVariableValue(), nil
	case model.VariableTypeRandom:
		return service.getRandomVariableValue(), nil
	case model.VariableTypeQuery:
		return service.getQueryVariableValue(v.Key, request)
	case model.VariableTypePath:
		return service.getPathVariableValue(v.Key, rule.Path, path)
	}

	return "", mockserrors.InvalidRulesError{
		Message: fmt.Sprintf("%s is invalid variable type", v.Type),
	}
}

func (service *MockService) getPathParams(rulePath, reqPath string) (map[string]string, error) {
	u, err := url.Parse(reqPath)
	if err != nil {
		return nil, fmt.Errorf("error parsing url, %w", err)
	}

	values := strings.Split(u.Path, "/")
	pathParts := strings.Split(rulePath, "/")

	params := make(map[string]string)

	for i, part := range pathParts {
		if part != "" {
			first := string(part[0])
			last := string(part[len(part)-1])

			if first == "{" && last == "}" {
				key := part[1 : len(part)-1]
				params[key] = values[i]
			}
		}
	}

	return params, nil
}
