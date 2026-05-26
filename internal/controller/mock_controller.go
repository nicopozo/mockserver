package controller

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	mockscontext "github.com/nicopozo/mockserver/internal/context"
	ruleserrors "github.com/nicopozo/mockserver/internal/errors"
	"github.com/nicopozo/mockserver/internal/model"
	"github.com/nicopozo/mockserver/internal/service"
	httputils "github.com/nicopozo/mockserver/internal/utils/http"
	"github.com/nicopozo/mockserver/internal/utils/log"
	"github.com/oklog/ulid/v2"
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

func (controller *MockController) Execute(writer http.ResponseWriter, request *http.Request) {
	reqContext := mockscontext.New(request)
	logger := mockscontext.Logger(reqContext)

	logger.Debug(controller, nil, "Entering MockController Execute()")

	path := request.PathValue("rule")
	if path == "" {
		// Fallback for cases where it might not be matched as a named param
		path = request.URL.Path
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	reqBody := controller.extractExecutionBody(logger, request.Body)

	// Build base log entry from the incoming request.
	logEntry := controller.buildLogEntry(request, path, reqBody)

	// Generate a log ID upfront so the webhook callback can reference it.
	logID := ulid.Make().String()
	logEntry.ID = logID

	// Callback to record webhook result into the same log entry.
	onWebhookResult := func(result model.WebhookResult) {
		controller.LogService.Update(logID, func(entry *model.LogEntry) {
			entry.WebhookResults = append(entry.WebhookResults, result)
		})
	}

	response, assertionResult, err := controller.MockService.SearchResponseForRequest(
		reqContext, request, path, reqBody, onWebhookResult)

	// Attach assertion errors to the log entry regardless of outcome.
	logEntry.AssertionErrors = assertionResult.AssertionErrors

	if err != nil {
		controller.handleExecutionError(writer, request, logger, path, logEntry, err)

		return
	}

	writer.Header().Set("Content-Type", response.ContentType)

	time.Sleep(time.Duration(response.Delay) * time.Millisecond)
	writer.WriteHeader(response.HTTPStatus)

	_, _ = writer.Write([]byte(response.Body))

	controller.recordLog(logEntry, response.HTTPStatus, response.Body)
}

func (controller *MockController) handleExecutionError(
	writer http.ResponseWriter,
	request *http.Request,
	logger log.ILogger,
	path string,
	logEntry model.LogEntry,
	err error,
) {
	if errors.As(err, &ruleserrors.RuleNotFoundError{}) {
		logger.Debug(controller, nil, "No rule found for path: %v and method: %s",
			path, request.Method)

		errorResult := model.NewError(model.ResourceNotFoundError,
			"No rule found for path: %v and method: %s. %v", path, request.Method, err.Error())

		httputils.WriteJSON(writer, http.StatusNotFound, errorResult)
		controller.recordLog(logEntry, http.StatusNotFound, errorResult.Message)

		return
	}

	if errors.As(err, &ruleserrors.InvalidRulesError{}) {
		logger.Debug(controller, nil, "No valid rule found for path: %v and method: %s",
			path, request.Method)

		errorResult := model.NewError(model.ValidationError, "%s", err.Error())

		httputils.WriteJSON(writer, http.StatusNotFound, errorResult)
		controller.recordLog(logEntry, http.StatusNotFound, errorResult.Message)

		return
	}

	if errors.As(err, &ruleserrors.AssertionError{}) {
		logger.Debug(controller, nil, "One or more assertions failed.",
			path, request.Method)

		errorResult := model.NewError(model.ValidationError, "%s", err.Error())

		httputils.WriteJSON(writer, http.StatusBadRequest, errorResult)
		controller.recordLog(logEntry, http.StatusBadRequest, errorResult.Message)

		return
	}

	logger.Error(controller, nil, err,
		"Failed to execute rule for method %v: and path %v", request.Method, path)

	errorResult := model.NewError(model.InternalError, "Error occurred when getting rule. %s", err.Error())

	httputils.WriteJSON(writer, http.StatusInternalServerError, errorResult)
	controller.recordLog(logEntry, http.StatusInternalServerError, errorResult.Message)
}

// buildLogEntry constructs a LogEntry from the incoming request.
func (controller *MockController) buildLogEntry(request *http.Request, path, reqBody string) model.LogEntry {
	headers := make(map[string]string, len(request.Header))

	for key, values := range request.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	queryParams := make(map[string]string)

	for key, values := range request.URL.Query() {
		if len(values) > 0 {
			queryParams[key] = values[0]
		}
	}

	fullURL := path
	if request.URL.RawQuery != "" {
		fullURL = path + "?" + request.URL.RawQuery
	}

	return model.LogEntry{
		Method:         request.Method,
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
