package controller

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	mockscontext "github.com/nicopozo/mockserver/internal/context"
	ruleserrors "github.com/nicopozo/mockserver/internal/errors"
	"github.com/nicopozo/mockserver/internal/model"
	"github.com/nicopozo/mockserver/internal/service"
	"github.com/nicopozo/mockserver/internal/utils/log"
)

type MockController struct {
	MockService service.MockService
	LogService  service.LogService
}

func NewMockController(mockService service.MockService, logService service.LogService) *MockController {
	return &MockController{
		MockService: mockService,
		LogService:  logService,
	}
}

func (controller *MockController) Execute(context *gin.Context) {
	reqContext := mockscontext.New(context)
	logger := mockscontext.Logger(reqContext)

	logger.Debug(controller, nil, "Entering MockController Execute()")

	path := context.Param("rule")
	reqBody := controller.extractExecutionBody(logger, context.Request.Body)

	// Build base log entry from the incoming request.
	logEntry := controller.buildLogEntry(context, path, reqBody)

	response, err := controller.MockService.SearchResponseForRequest(reqContext, context.Request, path, reqBody)
	if err != nil {
		controller.handleExecutionError(context, logger, path, logEntry, err)

		return
	}

	context.Header("content-type", response.ContentType)

	time.Sleep(time.Duration(response.Delay) * time.Millisecond)
	context.String(response.HTTPStatus, response.Body)

	controller.recordLog(logEntry, response.HTTPStatus, response.Body)
}

func (controller *MockController) handleExecutionError(
	context *gin.Context,
	logger log.ILogger,
	path string,
	logEntry model.LogEntry,
	err error,
) {
	if errors.As(err, &ruleserrors.RuleNotFoundError{}) {
		logger.Debug(controller, nil, "No rule found for path: %v and method: %s",
			path, context.Request.Method)

		errorResult := model.NewError(model.ResourceNotFoundError,
			"No rule found for path: %v and method: %s. %v", path, context.Request.Method, err.Error())
		context.JSON(http.StatusNotFound, errorResult)

		controller.recordLog(logEntry, http.StatusNotFound, errorResult.Message)

		return
	}

	if errors.As(err, &ruleserrors.InvalidRulesError{}) {
		logger.Debug(controller, nil, "No valid rule found for path: %v and method: %s",
			path, context.Request.Method)

		errorResult := model.NewError(model.ValidationError, "%s", err.Error())
		context.JSON(http.StatusNotFound, errorResult)

		controller.recordLog(logEntry, http.StatusNotFound, errorResult.Message)

		return
	}

	if errors.As(err, &ruleserrors.AssertionError{}) {
		logger.Debug(controller, nil, "One or more assertions failed.",
			path, context.Request.Method)

		errorResult := model.NewError(model.ValidationError, "%s", err.Error())
		context.JSON(http.StatusBadRequest, errorResult)

		controller.recordLog(logEntry, http.StatusBadRequest, errorResult.Message)

		return
	}

	logger.Error(controller, nil, err,
		"Failed to execute rule for method %v: and path %v", context.Request.Method, path)

	errorResult := model.NewError(model.InternalError, "Error occurred when getting rule. %s", err.Error())
	context.JSON(http.StatusInternalServerError, errorResult)

	controller.recordLog(logEntry, http.StatusInternalServerError, errorResult.Message)
}

// buildLogEntry constructs a LogEntry from the incoming Gin context.
func (controller *MockController) buildLogEntry(context *gin.Context, path, reqBody string) model.LogEntry {
	headers := make(map[string]string, len(context.Request.Header))
	for key, values := range context.Request.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	queryParams := make(map[string]string)

	for key, values := range context.Request.URL.Query() {
		if len(values) > 0 {
			queryParams[key] = values[0]
		}
	}

	fullURL := path
	if context.Request.URL.RawQuery != "" {
		fullURL = path + "?" + context.Request.URL.RawQuery
	}

	return model.LogEntry{
		Method:         context.Request.Method,
		URL:            fullURL,
		RequestBody:    reqBody,
		RequestHeaders: headers,
		QueryParams:    queryParams,
	}
}

// recordLog saves the completed log entry (with response info) to the LogService.
func (controller *MockController) recordLog(entry model.LogEntry, responseStatus int, responseBody string) {
	if controller.LogService == nil {
		return
	}

	entry.ResponseStatus = responseStatus
	entry.ResponseBody = responseBody
	controller.LogService.Add(entry)
}

func (controller *MockController) extractExecutionBody(logger log.ILogger, body io.Reader) string {
	var bodyContents []byte

	if body == nil {
		logger.Warn(controller, nil, "Received body was NIL")

		return "<nil>"
	}

	bodyContents, err := io.ReadAll(body)
	if err != nil {
		logger.Error(controller, nil, err, "Error extracting body from execution: "+err.Error())

		return "<error>"
	}

	return string(bodyContents)
}
