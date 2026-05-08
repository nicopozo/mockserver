package testutils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/nicopozo/mockserver/internal/model"
)

func GetJSONFromFile(filename string) (string, error) {
	contentBytes, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("error file, %w", err)
	}

	contentString := string(contentBytes)

	return contentString, nil
}

func GetHTTPContext() (*httptest.ResponseRecorder, *http.Request) {
	request := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/?callback=thecallback.com", nil)
	responseWriter := httptest.NewRecorder()

	return responseWriter, request
}

func GetHTTPContextWithBody(requestBodyFile string) (*httptest.ResponseRecorder, *http.Request, error) {
	bodyStr, err := GetJSONFromFile(requestBodyFile)
	if err != nil {
		return nil, nil, err
	}

	request := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"/?callback=thecallback.com",
		bytes.NewBufferString(bodyStr))
	responseWriter := httptest.NewRecorder()

	return responseWriter, request, nil
}

func GetErrorFromResponse(response []byte) (*model.Error, error) {
	errorResponse := &model.Error{}

	err := json.Unmarshal(response, &errorResponse)
	if err != nil {
		return errorResponse, fmt.Errorf("error reading erro from response, %w", err)
	}

	return errorResponse, nil
}

func GetRuleFromResponse(response []byte) (*model.Rule, error) {
	rule := &model.Rule{}

	err := json.Unmarshal(response, &rule)
	if err != nil {
		return nil, fmt.Errorf("error reading rule from response, %w", err)
	}

	return rule, nil
}

func GetRuleListFromResponse(response []byte) (*model.RuleList, error) {
	list := &model.RuleList{}

	err := json.Unmarshal(response, &list)
	if err != nil {
		return nil, fmt.Errorf("error reading rule list from response, %w", err)
	}

	return list, nil
}
