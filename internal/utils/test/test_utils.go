package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/nicopozo/mockserver/internal/model"
	"github.com/nicopozo/mockserver/internal/utils/test/mocks"

	"github.com/gin-gonic/gin"
)

type Body struct {
	io.Reader
}

func (Body) Close() error { return nil }

func GetJSONFromFile(filename string) (string, error) {
	contentBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	contentString := string(contentBytes)

	return contentString, nil
}

func GetGinContext() (*gin.Context, *mocks.MockGinResponseWriter) {
	request := http.Request{}

	theURL := url.URL{RawQuery: "callback:thecallback.com"}

	request.URL = &theURL
	context := gin.Context{Request: &request}
	responseWriter := mocks.MockGinResponseWriter{}
	context.Writer = &responseWriter

	return &context, &responseWriter
}

func GetGinContextWithBody(requestBodyFile string) (*gin.Context, *mocks.MockGinResponseWriter, error) {
	bodyStr, err := GetJSONFromFile(requestBodyFile)
	if err != nil {
		return nil, nil, err
	}

	body := Body{bytes.NewBufferString(bodyStr)}
	url := url.URL{RawQuery: "callback:thecallback.com"}

	context, response := getGinContext()
	context.Request.Body = body
	context.Request.URL = &url

	return context, response, nil
}

func GetErrorFromResponse(response []byte) (*model.Error, error) {
	errorResponse := &model.Error{}
	err := json.Unmarshal(response, &errorResponse)

	return errorResponse, err
}

func getGinContext() (*gin.Context, *mocks.MockGinResponseWriter) {
	request := http.Request{}

	theURL := url.URL{RawQuery: "callback:thecallback.com"}

	request.URL = &theURL

	context := gin.Context{Request: &request}
	responseWriter := mocks.MockGinResponseWriter{}
	context.Writer = &responseWriter

	return &context, &responseWriter
}

func GetRuleFromResponse(response []byte) (*model.Rule, error) {
	rule := &model.Rule{}
	err := json.Unmarshal(response, &rule)

	return rule, err
}

func GetRuleListFromResponse(response []byte) (*model.RuleList, error) {
	list := &model.RuleList{}
	err := json.Unmarshal(response, &list)

	return list, err
}
