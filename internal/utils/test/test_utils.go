package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/nicopozo/mockserver/internal/model"
	"github.com/nicopozo/mockserver/internal/utils/test/mocks"
)

type Body struct {
	io.Reader
}

func (Body) Close() error { return nil }

func GetJSONFromFile(filename string) (string, error) {
	contentBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("error file, %w", err)
	}

	contentString := string(contentBytes)

	return contentString, nil
}

func GetGinContext() (*gin.Context, *mocks.MockGinResponseWriter) {
	request := http.Request{}

	theURL := url.URL{RawQuery: "callback:thecallback.com"}

	request.URL = &theURL
	ginContext := gin.Context{Request: &request}
	responseWriter := mocks.MockGinResponseWriter{}
	ginContext.Writer = &responseWriter

	return &ginContext, &responseWriter
}

func GetGinContextWithBody(requestBodyFile string) (*gin.Context, *mocks.MockGinResponseWriter, error) {
	bodyStr, err := GetJSONFromFile(requestBodyFile)
	if err != nil {
		return nil, nil, err
	}

	body := Body{bytes.NewBufferString(bodyStr)}
	u := url.URL{RawQuery: "callback:thecallback.com"}

	ginContext, response := getGinContext()
	ginContext.Request.Body = body
	ginContext.Request.URL = &u

	return ginContext, response, nil
}

func GetErrorFromResponse(response []byte) (*model.Error, error) {
	errorResponse := &model.Error{}

	err := json.Unmarshal(response, &errorResponse)
	if err != nil {
		return errorResponse, fmt.Errorf("error reading erro from response, %w", err)
	}

	return errorResponse, nil
}

func getGinContext() (*gin.Context, *mocks.MockGinResponseWriter) {
	request := http.Request{}

	theURL := url.URL{RawQuery: "callback:thecallback.com"}

	request.URL = &theURL

	ginContext := gin.Context{Request: &request}
	responseWriter := mocks.MockGinResponseWriter{}
	ginContext.Writer = &responseWriter

	return &ginContext, &responseWriter
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
