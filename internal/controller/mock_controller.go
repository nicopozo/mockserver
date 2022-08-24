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
}

func (controller *MockController) Execute(context *gin.Context) {
	reqContext := mockscontext.New(context)
	logger := mockscontext.Logger(reqContext)

	logger.Debug(controller, nil, "Entering MockController Execute()")

	path := context.Param("rule")
	reqBody := controller.extractExecutionBody(logger, context.Request.Body)

	response, err := controller.MockService.SearchResponseForRequest(reqContext, context.Request, path, reqBody)
	if err != nil {
		if errors.As(err, &ruleserrors.RuleNotFoundError{}) {
			logger.Debug(controller, nil, "No rule found for path: %v and method: %s",
				path, context.Request.Method)

			errorResult := model.NewError(model.ResourceNotFoundError,
				"No rule found for path: %v and method: %s. %v", path, context.Request.Method, err.Error())
			context.JSON(http.StatusNotFound, errorResult)

			return
		}

		if errors.As(err, &ruleserrors.InvalidRulesError{}) {
			logger.Debug(controller, nil, "No valid rule found for path: %v and method: %s",
				path, context.Request.Method)

			errorResult := model.NewError(model.ValidationError, err.Error())
			context.JSON(http.StatusNotFound, errorResult)

			return
		}

		if errors.As(err, &ruleserrors.AssertionError{}) {
			logger.Debug(controller, nil, "One or more assertions failed.",
				path, context.Request.Method)

			errorResult := model.NewError(model.ValidationError, err.Error())
			context.JSON(http.StatusBadRequest, errorResult)

			return
		}

		logger.Error(controller, nil, err,
			"Failed to execute rule for method %v: and path %v", context.Request.Method, path)

		errorResult := model.NewError(model.InternalError, "Error occurred when getting rule. %s", err.Error())
		context.JSON(http.StatusInternalServerError, errorResult)

		return
	}

	context.Header("content-type", response.ContentType)

	time.Sleep(time.Duration(response.Delay) * time.Millisecond)
	context.String(response.HTTPStatus, response.Body)
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
