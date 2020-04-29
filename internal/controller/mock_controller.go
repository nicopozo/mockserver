package controller

import (
	"net/http"
	"strings"

	ruleserrors "github.com/nicopozo/mockserver/internal/errors"
	"github.com/nicopozo/mockserver/internal/model"
	"github.com/nicopozo/mockserver/internal/service"

	newrelic "github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrgin/v1"

	"github.com/gin-gonic/gin"
)

type MockController struct {
	MockService service.IMockService
}

func (controller *MockController) Execute(context *gin.Context) {
	logger := GetLogger(context)
	logger.Debug(controller, nil, "Entering RuleController Get()")

	txn := nrgin.Transaction(context)
	segment := newrelic.StartSegment(txn, "ControllerGet")

	defer endSegment(controller, segment, logger)

	method := strings.ToLower(context.Request.Method)
	path := context.Param("rule")

	response, err := controller.MockService.SearchResponseForMethodAndPath(method, path, txn, logger)

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
	context.String(response.HTTPStatus, response.Body)
}
