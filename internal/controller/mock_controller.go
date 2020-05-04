package controller

import (
	"net/http"
	"strings"
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

	logger.Debug(controller, nil, "Entering RuleController Get()")

	method := strings.ToUpper(context.Request.Method)
	path := context.Param("rule")

	response, err := controller.MockService.SearchResponseForMethodAndPath(reqContext, method, path)

	if err != nil {
		switch err.(type) {
		case ruleserrors.RuleNotFoundError:
			logger.Debug(controller, nil, "No rule found for path: %v and method: %s", path, method)

			errorResult := model.NewError(model.ResourceNotFoundError,
				"No rule found for path: %v and method: %s. %v", path, method, err.Error())
			context.JSON(http.StatusNotFound, errorResult)
		default:
			logger.Error(controller, nil, err,
				"Failed to execute rule for method %v: and path %v", method, path)

			errorResult := model.NewError(model.InternalError, "Error occurred when getting rule. %s", err.Error())
			context.JSON(http.StatusInternalServerError, errorResult)
		}

		return
	}

	context.Header("content-type", response.ContentType)

	time.Sleep(time.Duration(response.Delay) * time.Millisecond)
	context.String(response.HTTPStatus, response.Body)
}
