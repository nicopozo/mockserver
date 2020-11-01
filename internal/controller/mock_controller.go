package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	mockscontext "github.com/nicopozo/mockserver/internal/context"
	ruleserrors "github.com/nicopozo/mockserver/internal/errors"
	"github.com/nicopozo/mockserver/internal/model"
	"github.com/nicopozo/mockserver/internal/service"
)

type MockController struct {
	MockService service.IMockService
}

func (controller *MockController) Execute(context *gin.Context) {
	reqContext := mockscontext.New(context)
	logger := mockscontext.Logger(reqContext)

	logger.Debug(controller, nil, "Entering MockController Execute()")

	path := context.Param("rule")

	response, err := controller.MockService.SearchResponseForRequest(reqContext, context.Request, path)
	if err != nil {
		switch err.(type) { //nolint:errorlint
		case ruleserrors.RuleNotFoundError:
			logger.Debug(controller, nil, "No rule found for path: %v and method: %s",
				path, context.Request.Method)

			errorResult := model.NewError(model.ResourceNotFoundError,
				"No rule found for path: %v and method: %s. %v", path, context.Request.Method, err.Error())
			context.JSON(http.StatusNotFound, errorResult)
		case ruleserrors.InvalidRulesError:
			logger.Debug(controller, nil, "No rule found for path: %v and method: %s",
				path, context.Request.Method)

			errorResult := model.NewError(model.ValidationError, err.Error())
			context.JSON(http.StatusNotFound, errorResult)
		default:
			logger.Error(controller, nil, err,
				"Failed to execute rule for method %v: and path %v", context.Request.Method, path)

			errorResult := model.NewError(model.InternalError, "Error occurred when getting rule. %s", err.Error())
			context.JSON(http.StatusInternalServerError, errorResult)
		}

		return
	}

	context.Header("content-type", response.ContentType)

	time.Sleep(time.Duration(response.Delay) * time.Millisecond)
	context.String(response.HTTPStatus, response.Body)
}
