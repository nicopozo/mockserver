package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	mockscontext "github.com/nicopozo/mockserver/internal/context"
	ruleserrors "github.com/nicopozo/mockserver/internal/errors"
	"github.com/nicopozo/mockserver/internal/model"
	"github.com/nicopozo/mockserver/internal/service"
	jsonutils "github.com/nicopozo/mockserver/internal/utils/json"
)

type RuleController struct {
	RuleService service.IRuleService
}

// @Tags Rules
// @Summary Create a Rule
// @Description Create a Rule for serving a mock response
// @ID create-rule
// @Accept  json
// @Produce  json
// @Param rule body model.Rule true "The rule to be created"
// @Success 201 {object} model.Rule "Rule successfully created"
// @Failure 400 {object} model.Error "Validation of the rule failed"
// @Failure 500 {object} model.Error "Internal server error"
// @Router /rules [post]
// Create a Rule.
func (controller *RuleController) Create(context *gin.Context) {
	reqContext := mockscontext.New(context)
	logger := mockscontext.Logger(reqContext)
	logger.Debug(controller, nil, "Entering RuleController Save()")

	rule, err := model.UnmarshalRule(context.Request.Body)
	if err != nil {
		errorResult := model.NewError(model.ValidationError, "Invalid JSON. %s", err.Error())
		logger.Error(controller, nil, err, "Error unmarshalling Rule JSON")
		context.JSON(http.StatusBadRequest, errorResult)

		return
	}

	rule, err = controller.RuleService.Save(reqContext, rule)

	if err != nil {
		switch err.(type) { //nolint:errorlint
		case ruleserrors.InvalidRulesError:
			logger.Error(controller, nil, err, "Invalid rule: %v", jsonutils.Marshal(rule))

			errorResult := model.NewError(model.ValidationError, "%s", err.Error())
			context.JSON(http.StatusBadRequest, errorResult)
		default:
			errorResult := model.NewError(model.InternalError, "Error occurred when saving rule. %s", err.Error())

			logger.Error(controller, nil, err, "Error occurred when saving rule")
			context.JSON(http.StatusInternalServerError, errorResult)

			return
		}

		return
	}

	context.JSON(http.StatusCreated, rule)
}

// @Tags Rules
// @Summary Update a Rule
// @Description Update an existing Rule for serving a mock response
// @ID update-rule
// @Accept  json
// @Produce  json
// @Param key path string true "Key of rule to update"
// @Param rule body model.Rule true "The rule to be updated"
// @Success 200 {object} model.Rule "Rule successfully updated"
// @Failure 400 {object} model.Error "Validation of the rule failed"
// @Failure 500 {object} model.Error "Internal server error"
// @Router /rules/{key} [put]
// Update a Rule.
func (controller *RuleController) Update(context *gin.Context) {
	reqContext := mockscontext.New(context)
	logger := mockscontext.Logger(reqContext)
	logger.Debug(controller, nil, "Entering RuleController Save()")

	key := context.Param("key")

	rule, err := model.UnmarshalRule(context.Request.Body)
	if err != nil {
		errorResult := model.NewError(model.ValidationError, "Invalid JSON. %s", err.Error())
		logger.Error(controller, nil, err, "Error unmarshalling Rule JSON")
		context.JSON(http.StatusBadRequest, errorResult)

		return
	}

	rule, err = controller.RuleService.Update(reqContext, key, rule)

	if err != nil {
		switch err.(type) { //nolint:errorlint
		case ruleserrors.RuleNotFoundError:
			logger.Debug(controller, nil, "No rule found with key: %v", key)

			errorResult := model.NewError(model.ResourceNotFoundError, "%s", err.Error())
			context.JSON(http.StatusNotFound, errorResult)
		case ruleserrors.InvalidRulesError:
			logger.Error(controller, nil, err, "Invalid rule: %v", jsonutils.Marshal(rule))

			errorResult := model.NewError(model.ValidationError, "%s", err.Error())
			context.JSON(http.StatusBadRequest, errorResult)
		default:
			errorResult := model.NewError(model.InternalError, "Error occurred when saving rule. %s", err.Error())

			logger.Error(controller, nil, err, "Error occurred when saving rule")
			context.JSON(http.StatusInternalServerError, errorResult)
		}

		return
	}

	context.JSON(http.StatusOK, rule)
}

// @Tags Rules
// @Summary Update a Rule Status
// @Description Update an existing Rule for serving a mock response
// @ID update-rule
// @Accept  json
// @Produce  json
// @Param key path string true "Key of rule to update"
// @Param rule body model.RuleStatus true "The rule to be updated"
// @Success 200 {object} model.Rule "Rule successfully updated"
// @Failure 400 {object} model.Error "Validation of the rule failed"
// @Failure 500 {object} model.Error "Internal server error"
// @Router /rules/{key}/status [put]
// Update a Rule.
func (controller *RuleController) UpdateStatus(context *gin.Context) {
	reqContext := mockscontext.New(context)
	logger := mockscontext.Logger(reqContext)
	logger.Debug(controller, nil, "Entering RuleController Save()")

	key := context.Param("key")

	ruleStatus, err := model.UnmarshalRuleStatus(context.Request.Body)
	if err != nil {
		errorResult := model.NewError(model.ValidationError, "Invalid JSON. %s", err.Error())
		logger.Error(controller, nil, err, "Error unmarshalling Rule JSON")
		context.JSON(http.StatusBadRequest, errorResult)

		return
	}

	rule, err := controller.RuleService.UpdateStatus(reqContext, key, ruleStatus)
	if err != nil {
		switch err.(type) { //nolint:errorlint
		case ruleserrors.RuleNotFoundError:
			logger.Debug(controller, nil, "No rule found with key: %v", key)

			errorResult := model.NewError(model.ResourceNotFoundError, "%s", err.Error())
			context.JSON(http.StatusNotFound, errorResult)
		case ruleserrors.InvalidRulesError:
			logger.Error(controller, nil, err, "Invalid status: %v", ruleStatus.Status)

			errorResult := model.NewError(model.ValidationError, "%s", err.Error())
			context.JSON(http.StatusBadRequest, errorResult)
		default:
			errorResult := model.NewError(model.InternalError, "Error occurred when saving rule. %s", err.Error())

			logger.Error(controller, nil, err, "Error occurred when saving rule")
			context.JSON(http.StatusInternalServerError, errorResult)
		}

		return
	}

	context.JSON(http.StatusOK, rule)
}

// @Tags Rules
// @Summary Get Rule by Key
// @Description Get a Rule, if not found return 404
// @ID get-rule
// @Produce json
// @Param key path string true "Key generated by service"
// @Success 200 {object} model.Rule "Result"
// @Failure 400 {object} model.Error
// @Failure 404 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /rules/{key} [get]
// Get a Rule.
func (controller *RuleController) Get(context *gin.Context) {
	reqContext := mockscontext.New(context)
	logger := mockscontext.Logger(reqContext)

	logger.Debug(controller, nil, "Entering RuleController Get()")

	key := context.Param("key")

	task, err := controller.RuleService.Get(reqContext, key)
	if err != nil {
		switch err.(type) { //nolint:errorlint
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

// @Tags Rules
// @Summary Search Rule
// @Description Search Rule by key, name, application, method or status
// @ID search-rule
// @Produce json
// @Param key query string false "Rule key generated by service"
// @Param name query string false "Name of the key"
// @Param application query string false "Application"
// @Param method query string false "Method"
// @Param status query string false "Enabled/Disabled"
// @Param limit query number false "Max expected number of results" default(30)
// @Param offset query string false "number of results to be skipped" default(0)
// @Success 200 {object} model.RuleList "Result"
// @Failure 400 {object} model.Error
// @Failure 404 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /rules [get]
// Search Rules.
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

// @Tags Rules
// @Summary Delete Rule by key
// @Description Delete Rule by Key
// @ID delete-rule
// @Produce json
// @Param key path string true "Key generated by service"
// @Success 204
// @Failure 500 {object} model.Error
// @Router /rules/{key} [delete]
// Delete Rule.
func (controller *RuleController) Delete(context *gin.Context) {
	reqContext := mockscontext.New(context)
	logger := mockscontext.Logger(reqContext)

	logger.Debug(controller, nil, "Entering RuleController Delete()")

	key := context.Param("key")

	err := controller.RuleService.Delete(reqContext, key)
	if err != nil {
		switch err.(type) { //nolint:errorlint
		case ruleserrors.RuleNotFoundError:
			context.Status(http.StatusNoContent)
		default:
			logger.Error(controller, nil, err,
				"Failed to delete rule with key: %v", key)

			errorResult := model.NewError(model.InternalError, "Error occurred when deleting rule. %s", err.Error())
			context.JSON(http.StatusInternalServerError, errorResult)
		}

		return
	}

	context.Status(http.StatusNoContent)
}
