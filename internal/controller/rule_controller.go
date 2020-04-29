package controller

import (
	"net/http"

	ruleserrors "github.com/nicopozo/mockserver/internal/errors"
	"github.com/nicopozo/mockserver/internal/model"
	"github.com/nicopozo/mockserver/internal/service"

	"github.com/gin-gonic/gin"
	newrelic "github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrgin/v1"
)

type RuleController struct {
	RuleService service.IRuleService
}

func (controller *RuleController) Create(c *gin.Context) {
	logger := GetLogger(c)
	logger.Debug(controller, nil, "Entering RuleController Save()")

	rule, err := model.UnmarshalRule(c.Request.Body)
	if err != nil {
		errorResult := model.NewError(model.ValidationError, "Invalid JSON. %s", err.Error())
		logger.Error(controller, nil, err, "Error unmarshalling Rule JSON")
		c.JSON(http.StatusBadRequest, errorResult)

		return
	}

	txn := nrgin.Transaction(c)
	segment := newrelic.StartSegment(txn, "ControllerSend")

	defer endSegment(controller, segment, logger)

	err = controller.RuleService.Save(rule, txn, logger)

	if err != nil {
		errorResult := model.NewError(model.InternalError, "Error occurred when saving rule. %s", err.Error())

		logger.Error(controller, nil, err, "Error occurred when saving rule")
		c.JSON(http.StatusInternalServerError, errorResult)

		return
	}

	c.JSON(http.StatusCreated, rule)
}

func (controller *RuleController) Get(context *gin.Context) {
	logger := GetLogger(context)
	logger.Debug(controller, nil, "Entering RuleController Get()")

	key := context.Param("key")

	txn := nrgin.Transaction(context)
	segment := newrelic.StartSegment(txn, "ControllerGet")

	defer endSegment(controller, segment, logger)

	task, err := controller.RuleService.Get(key, txn, logger)

	if err != nil {
		switch err.(type) {
		case ruleserrors.RuleNotFoundError:
			logger.Debug(controller, nil, "No rule found with key: %v", key)

			errorResult := model.NewError(model.ResourceNotFoundError, "%s", err.Error())
			context.JSON(http.StatusNotFound, errorResult)
		default:
			logger.Error(controller, nil, err,
				"Failed to get task with key: %v", key)

			errorResult := model.NewError(model.InternalError, "Error occurred when getting rule. %s", err.Error())
			context.JSON(http.StatusInternalServerError, errorResult)
		}

		return
	}

	context.JSON(http.StatusOK, task)
}
