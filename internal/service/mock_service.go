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

	mockscontext "github.com/nicopozo/mockserver/internal/context"
	"github.com/nicopozo/mockserver/internal/model"
	jsonutils "github.com/nicopozo/mockserver/internal/utils/json"
	"github.com/nicopozo/mockserver/internal/utils/log"
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

	result, err := service.RuleService.SearchByMethodAndPath(ctx, method, path)
	if err != nil {
		logger.Error(service, nil, err, "error searching responses")
		return nil, err
	}

	response := result.Responses[0]
	reqBody := new(strings.Builder)

	_, err = io.Copy(reqBody, request.Body)
	if err != nil {
		return nil, err
	}

	body, err := service.applyVariables(request, reqBody.String(), response.Body, result.Variables)
	if err != nil {
		return nil, err
	}

	body, err = service.replacePathParam(body, result.Path, path, logger)
	if err != nil {
		return nil, err
	}

	response.Body = body

	return &response, nil
}

func (service *MockService) replacePathParam(responseBody, rulePath, reqPath string,
	logger log.ILogger) (string, error) {
	u, err := url.Parse(reqPath)
	if err != nil {
		logger.Error(service, nil, err, "error parsing path")
		return "", err
	}

	values := strings.Split(u.Path, "/")
	pathParts := strings.Split(rulePath, "/")

	for i, part := range pathParts {
		if part != "" {
			first := string(part[0])
			last := string(part[len(part)-1])

			if first == "{" && last == "}" {
				responseBody = strings.Replace(responseBody, part, values[i], -1)
			}
		}
	}

	return responseBody, nil
}

func (service *MockService) applyVariables(request *http.Request, reqBody, respBody string,
	variables []*model.Variable) (string, error) {
	var err error

	for _, v := range variables {
		switch v.Type {
		case model.VariableTypeHeader:
			respBody = service.applyHeaderVariables(request, respBody, v.Name, v.Key)
		case model.VariableTypeBody:
			respBody, err = service.applyBodyVariables(reqBody, respBody, v.Name, v.Key)
			if err != nil {
				break
			}
		case model.VariableTypeHash:
			respBody = service.applyHashVariables(respBody, v.Name)
		case model.VariableTypeRandom:
			respBody = service.applyRandomVariables(respBody, v.Name)
		case model.VariableTypeQuery:
			respBody, err = service.applyQueryVariables(request, respBody, v.Name, v.Key)
			if err != nil {
				break
			}
		}
	}

	return respBody, err
}

func (service *MockService) applyBodyVariables(reqBody, respBody string, name string,
	key string) (string, error) {
	apply, err := jsonpath.Prepare(key)
	if err != nil {
		return "", fmt.Errorf("invalid JSON path for key %s - %w", key, err)
	}

	var reqMap interface{}
	if err := json.Unmarshal([]byte(reqBody), &reqMap); err != nil {
		return "", err
	}

	value, err := apply(reqMap)
	if err != nil {
		return "", err
	}

	respBody = strings.Replace(respBody, fmt.Sprintf("{%s}", name), jsonutils.Marshal(value), -1)

	return respBody, nil
}

func (service *MockService) applyRandomVariables(respBody string, name string) string {
	n := rand.Int63n(max)
	respBody = strings.Replace(respBody, fmt.Sprintf("{%s}", name), strconv.FormatInt(n, 10), -1)

	return respBody
}

func (service *MockService) applyHashVariables(respBody string, name string) string {
	n := rand.Int63n(max)
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", n))) //nolint:errcheck
	bs := h.Sum(nil)
	respBody = strings.Replace(respBody, fmt.Sprintf("{%s}", name), fmt.Sprintf("%x", bs), -1)

	return respBody
}

func (service *MockService) applyQueryVariables(request *http.Request, body string, name string,
	key string) (string, error) {
	queries, err := url.ParseQuery(request.URL.RawQuery)
	if err != nil {
		return "", fmt.Errorf("error parsing queries %w", err)
	}

	for queryName, queryValue := range queries {
		if queryName == key {
			body = strings.Replace(body, fmt.Sprintf("{%s}", name), queryValue[0], -1)
			break
		}
	}

	return body, nil
}

func (service *MockService) applyHeaderVariables(request *http.Request, body string, name string, key string) string {
	header := request.Header[key]
	if len(header) > 0 {
		body = strings.Replace(body, fmt.Sprintf("{%s}", name), header[0], -1)
	}

	return body
}
