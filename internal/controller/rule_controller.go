package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	mockscontext "github.com/nicopozo/mockserver/internal/context"
	ruleserrors "github.com/nicopozo/mockserver/internal/errors"
	"github.com/nicopozo/mockserver/internal/model"
	"github.com/nicopozo/mockserver/internal/service"
)

type RuleController struct {
	RuleService service.IRuleService
}

func (controller *RuleController) Create(c *gin.Context) {
	reqContext := mockscontext.New(c)
	logger := mockscontext.Logger(reqContext)
	logger.Debug(controller, nil, "Entering RuleController Save()")

	rule, err := model.UnmarshalRule(c.Request.Body)
	if err != nil {
		errorResult := model.NewError(model.ValidationError, "Invalid JSON. %s", err.Error())
		logger.Error(controller, nil, err, "Error unmarshalling Rule JSON")
		c.JSON(http.StatusBadRequest, errorResult)

		return
	}

	err = controller.RuleService.Save(reqContext, rule)

	if err != nil {
		errorResult := model.NewError(model.InternalError, "Error occurred when saving rule. %s", err.Error())

		logger.Error(controller, nil, err, "Error occurred when saving rule")
		c.JSON(http.StatusInternalServerError, errorResult)

		return
	}

	c.JSON(http.StatusCreated, rule)
}

func (controller *RuleController) Get(context *gin.Context) {
	reqContext := mockscontext.New(context)
	logger := mockscontext.Logger(reqContext)

	logger.Debug(controller, nil, "Entering RuleController Get()")

	key := context.Param("key")

	task, err := controller.RuleService.Get(reqContext, key)

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

func (controller *RuleController) Search(context *gin.Context) {
	reqContext := mockscontext.New(context)
	logger := mockscontext.Logger(reqContext)

	logger.Debug(controller, nil, "Entering RuleController Search()")

	paging, err := getPagingFromRequest(context.Request)
	if err != nil {
		logger.Error(controller, nil, err, "Error searching rules. Error parsing pagination params")
		errorResult := model.NewError(model.ValidationError, "Error parsing pagination params: %s", err.Error())
		context.JSON(http.StatusBadRequest, errorResult)

		return
	}

	params := getParametersFromRequest(context.Request)

	ruleList, err := controller.RuleService.Search(reqContext, params, *paging)

	if err != nil {
		logger.Error(controller, nil, err, "Failed to search rules")

		errorResult := model.NewError(model.InternalError, "Error occurred when searching rules. %s", err.Error())
		context.JSON(http.StatusInternalServerError, errorResult)

		return
	}

	context.JSON(http.StatusOK, ruleList)
}
