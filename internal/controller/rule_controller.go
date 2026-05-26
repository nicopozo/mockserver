package controller

import (
	"errors"
	"net/http"

	mockscontext "github.com/nicopozo/mockserver/internal/context"
	ruleserrors "github.com/nicopozo/mockserver/internal/errors"
	"github.com/nicopozo/mockserver/internal/model"
	"github.com/nicopozo/mockserver/internal/service"
	httputils "github.com/nicopozo/mockserver/internal/utils/http"
	jsonutils "github.com/nicopozo/mockserver/internal/utils/json"
)

const exportLimit = 10000

type RuleController struct {
	RuleService service.RuleService
}

func NewRuleController(ruleService service.RuleService) *RuleController {
	return &RuleController{
		RuleService: ruleService,
	}
}

// Create a Rule.
func (controller *RuleController) Create(writer http.ResponseWriter, request *http.Request) {
	reqContext := mockscontext.New(request)
	logger := mockscontext.Logger(reqContext)
	logger.Debug(controller, nil, "Entering RuleController Save()")

	rule, err := model.UnmarshalRule(request.Body)
	if err != nil {
		logger.Error(controller, nil, err, "Error unmarshalling Rule JSON")
		httputils.WriteError(writer, model.ValidationError, "Invalid JSON. %s", err.Error())

		return
	}

	serviceRule, err := controller.RuleService.Save(reqContext, *rule)
	if err != nil {
		switch err.(type) { //nolint:errorlint
		case ruleserrors.InvalidRulesError:
			logger.Error(controller, nil, err, "Invalid rule: %v", jsonutils.Marshal(rule))
			httputils.WriteError(writer, model.ValidationError, "%s", err.Error())
		default:
			logger.Error(controller, nil, err, "Error occurred when saving rule")
			httputils.WriteError(writer, model.InternalError, "Error occurred when saving rule. %s", err.Error())
		}

		return
	}

	httputils.WriteJSON(writer, http.StatusCreated, serviceRule)
}

// Update a Rule.
func (controller *RuleController) Update(writer http.ResponseWriter, request *http.Request) {
	reqContext := mockscontext.New(request)
	logger := mockscontext.Logger(reqContext)
	logger.Debug(controller, nil, "Entering RuleController Save()")

	key := request.PathValue("key")

	rule, err := model.UnmarshalRule(request.Body)
	if err != nil {
		logger.Error(controller, nil, err, "Error unmarshalling Rule JSON")
		httputils.WriteError(writer, model.ValidationError, "Invalid JSON. %s", err.Error())

		return
	}

	serviceRule, err := controller.RuleService.Update(reqContext, key, *rule)
	if err != nil {
		switch err.(type) { //nolint:errorlint
		case ruleserrors.RuleNotFoundError:
			logger.Debug(controller, nil, "No rule found with key: %v", key)
			httputils.WriteError(writer, model.ResourceNotFoundError, "%s", err.Error())
		case ruleserrors.InvalidRulesError:
			logger.Error(controller, nil, err, "Invalid rule: %v", jsonutils.Marshal(rule))
			httputils.WriteError(writer, model.ValidationError, "%s", err.Error())
		default:
			logger.Error(controller, nil, err, "Error occurred when saving rule")
			httputils.WriteError(writer, model.InternalError, "Error occurred when saving rule. %s", err.Error())
		}

		return
	}

	httputils.WriteJSON(writer, http.StatusOK, serviceRule)
}

// UpdateStatus updates a Rule Status.
func (controller *RuleController) UpdateStatus(writer http.ResponseWriter, request *http.Request) {
	reqContext := mockscontext.New(request)
	logger := mockscontext.Logger(reqContext)
	logger.Debug(controller, nil, "Entering RuleController Save()")

	key := request.PathValue("key")

	ruleStatus, err := model.UnmarshalRuleStatus(request.Body)
	if err != nil {
		logger.Error(controller, nil, err, "Error unmarshalling Rule JSON")
		httputils.WriteError(writer, model.ValidationError, "Invalid JSON. %s", err.Error())

		return
	}

	serviceRule, err := controller.RuleService.UpdateStatus(reqContext, key, *ruleStatus)
	if err != nil {
		switch err.(type) { //nolint:errorlint
		case ruleserrors.RuleNotFoundError:
			logger.Debug(controller, nil, "No rule found with key: %v", key)
			httputils.WriteError(writer, model.ResourceNotFoundError, "%s", err.Error())
		case ruleserrors.InvalidRulesError:
			logger.Error(controller, nil, err, "Invalid status: %v", ruleStatus.Status)
			httputils.WriteError(writer, model.ValidationError, "%s", err.Error())
		default:
			logger.Error(controller, nil, err, "Error occurred when saving rule")
			httputils.WriteError(writer, model.InternalError, "Error occurred when saving rule. %s", err.Error())
		}

		return
	}

	httputils.WriteJSON(writer, http.StatusOK, serviceRule)
}

// Get a Rule.
func (controller *RuleController) Get(writer http.ResponseWriter, request *http.Request) {
	reqContext := mockscontext.New(request)
	logger := mockscontext.Logger(reqContext)

	logger.Debug(controller, nil, "Entering RuleController Get()")

	key := request.PathValue("key")

	task, err := controller.RuleService.Get(reqContext, key)
	if err != nil {
		if errors.As(err, &ruleserrors.RuleNotFoundError{}) {
			logger.Debug(controller, nil, "No rule found with key: %v", key)
			httputils.WriteError(writer, model.ResourceNotFoundError, "%s", err.Error())

			return
		}

		logger.Error(controller, nil, err, "Failed to get task with key: %v", key)
		httputils.WriteError(writer, model.InternalError, "Error occurred when getting rule. %s", err.Error())

		return
	}

	// Debug: log the full rule JSON returned
	taskJSON := jsonutils.Marshal(task)
	logger.Debug(controller, map[string]string{"rule_json": taskJSON}, "Rule returned from get")

	httputils.WriteJSON(writer, http.StatusOK, task)
}

// Search Rules.
func (controller *RuleController) Search(writer http.ResponseWriter, request *http.Request) {
	reqContext := mockscontext.New(request)
	logger := mockscontext.Logger(reqContext)

	logger.Debug(controller, nil, "Entering RuleController Search()")

	paging, err := getPagingFromRequest(request)
	if err != nil {
		logger.Error(controller, nil, err, "Error searching rules. Error parsing pagination params")
		httputils.WriteError(writer, model.ValidationError, "Error parsing pagination params: %s", err.Error())

		return
	}

	params := getParametersFromRequest(request)

	ruleList, err := controller.RuleService.Search(reqContext, params, *paging)
	if err != nil {
		if errors.As(err, &ruleserrors.InvalidRulesError{}) {
			httputils.WriteError(writer, model.ValidationError, "Invalid parameters. %s", err.Error())

			return
		}

		logger.Error(controller, nil, err, "Failed to search rules")
		httputils.WriteError(writer, model.InternalError, "Error occurred when searching rules. %s", err.Error())

		return
	}

	httputils.WriteJSON(writer, http.StatusOK, ruleList)
}

// Delete Rule.
func (controller *RuleController) Delete(writer http.ResponseWriter, request *http.Request) {
	reqContext := mockscontext.New(request)
	logger := mockscontext.Logger(reqContext)

	logger.Debug(controller, nil, "Entering RuleController Delete()")

	key := request.PathValue("key")

	err := controller.RuleService.Delete(reqContext, key)
	if err != nil {
		if errors.As(err, &ruleserrors.RuleNotFoundError{}) {
			writer.WriteHeader(http.StatusNoContent)

			return
		}

		logger.Error(controller, nil, err, "Failed to delete rule with key: %v", key)
		httputils.WriteError(writer, model.InternalError, "Error occurred when deleting rule. %s", err.Error())

		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

// Export all Rules.
func (controller *RuleController) Export(writer http.ResponseWriter, request *http.Request) {
	reqContext := mockscontext.New(request)
	logger := mockscontext.Logger(reqContext)

	logger.Debug(controller, nil, "Entering RuleController Export()")

	// Use a large limit to get everything
	paging := model.Paging{Limit: exportLimit, Offset: 0}

	ruleList, err := controller.RuleService.Search(reqContext, nil, paging)
	if err != nil {
		logger.Error(controller, nil, err, "Failed to export rules")
		httputils.WriteError(writer, model.InternalError, "Error occurred when exporting rules. %s", err.Error())

		return
	}

	httputils.WriteJSON(writer, http.StatusOK, ruleList.Results)
}

// Import Rules.
func (controller *RuleController) Import(writer http.ResponseWriter, request *http.Request) {
	reqContext := mockscontext.New(request)
	logger := mockscontext.Logger(reqContext)

	logger.Debug(controller, nil, "Entering RuleController Import()")

	rules, err := model.UnmarshalRules(request.Body)
	if err != nil {
		logger.Error(controller, nil, err, "Error unmarshalling Rules JSON array")
		httputils.WriteError(writer, model.ValidationError, "Invalid JSON array. %s", err.Error())

		return
	}

	stats := map[string]int{
		"created": 0,
		"updated": 0,
		"failed":  0,
	}

	for _, rule := range rules {
		var err error
		if rule.Key != "" {
			// Try to update first
			_, err = controller.RuleService.Update(reqContext, rule.Key, rule)
			if err == nil {
				stats["updated"]++

				continue
			}

			// If error is not "not found", record failure
			if !errors.As(err, &ruleserrors.RuleNotFoundError{}) {
				logger.Error(controller, nil, err, "Failed to update rule during import: %s", rule.Key)

				stats["failed"]++

				continue
			}
		}

		// Create if key doesn't exist or update failed because not found
		_, err = controller.RuleService.Save(reqContext, rule)
		if err != nil {
			logger.Error(controller, nil, err, "Failed to create rule during import: %s", rule.Name)

			stats["failed"]++
		} else {
			stats["created"]++
		}
	}

	httputils.WriteJSON(writer, http.StatusOK, stats)
}
