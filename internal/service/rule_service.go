package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	mockscontext "github.com/nicopozo/mockserver/internal/context"
	mockserrors "github.com/nicopozo/mockserver/internal/errors"
	"github.com/nicopozo/mockserver/internal/model"
	"github.com/nicopozo/mockserver/internal/repository"
)

//go:generate mockgen -destination=../utils/test/mocks/rule_service_mock.go -package=mocks -source=./rule_service.go

type IRuleService interface {
	Save(ctx context.Context, rule *model.Rule) (*model.Rule, error)
	Update(ctx context.Context, key string, rule *model.Rule) (*model.Rule, error)
	UpdateStatus(ctx context.Context, key string, rule *model.RuleStatus) (*model.Rule, error)
	Get(ctx context.Context, key string) (*model.Rule, error)
	Search(ctx context.Context, params map[string]interface{}, paging model.Paging) (*model.RuleList, error)
	SearchByMethodAndPath(ctx context.Context, method, path string) (*model.Rule, error)
	Delete(ctx context.Context, key string) error
}

type RuleService struct {
	RuleRepository repository.IRuleRepository
}

func (ruleService *RuleService) Save(ctx context.Context, rule *model.Rule) (*model.Rule, error) {
	logger := mockscontext.Logger(ctx)

	logger.Debug(ruleService, nil, "Entering RuleService Save()")

	if err := validateRule(rule); err != nil {
		logger.Error(ruleService, nil, err, "Rule Validation failed")
		return nil, err
	}

	return ruleService.RuleRepository.Save(ctx, formatRule(rule), false)
}

func (ruleService *RuleService) Update(ctx context.Context, key string, rule *model.Rule) (*model.Rule, error) {
	logger := mockscontext.Logger(ctx)

	logger.Debug(ruleService, nil, "Entering RuleService Update()")

	_, err := ruleService.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	err = ruleService.Delete(ctx, key)
	if err != nil {
		return nil, err
	}

	rule.Key = key

	return ruleService.RuleRepository.Save(ctx, formatRule(rule), true)
}

func (ruleService *RuleService) UpdateStatus(ctx context.Context, key string,
	ruleStatus *model.RuleStatus) (*model.Rule, error) {
	logger := mockscontext.Logger(ctx)

	logger.Debug(ruleService, nil, "Entering RuleService Update()")

	if err := ruleStatus.Validate(); err != nil {
		return nil, err
	}

	rule, err := ruleService.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	rule.Status = ruleStatus.Status

	return ruleService.RuleRepository.Save(ctx, formatRule(rule), true)
}

func (ruleService *RuleService) Get(ctx context.Context, key string) (*model.Rule, error) {
	logger := mockscontext.Logger(ctx)

	logger.Debug(ruleService, nil, "Entering TaskService Get()")

	result, err := ruleService.RuleRepository.Get(ctx, key)
	if err != nil {
		logger.Error(ruleService, nil, err, "error getting task")
		return nil, err
	}

	return result, nil
}

func (ruleService *RuleService) Search(ctx context.Context, params map[string]interface{},
	paging model.Paging) (*model.RuleList, error) {
	logger := mockscontext.Logger(ctx)

	logger.Debug(ruleService, nil, "Entering RuleService Search()")

	return ruleService.RuleRepository.Search(ctx, params, paging)
}

func (ruleService *RuleService) SearchByMethodAndPath(ctx context.Context, method, path string) (*model.Rule, error) {
	logger := mockscontext.Logger(ctx)

	logger.Debug(ruleService, nil, "Entering RuleService Search()")

	result, err := ruleService.RuleRepository.SearchByMethodAndPath(ctx, method, path)
	if err != nil {
		logger.Error(ruleService, nil, err, "error searching rules")
		return nil, err
	}

	return result, nil
}

func (ruleService *RuleService) Delete(ctx context.Context, key string) error {
	logger := mockscontext.Logger(ctx)

	logger.Debug(ruleService, nil, "Entering TaskService Get()")

	return ruleService.RuleRepository.Delete(ctx, key)
}

func validateRule(rule *model.Rule) error {
	if rule == nil {
		return mockserrors.InvalidRulesError{
			Message: "rule cannot be nil",
		}
	}

	if rule.Name == "" {
		return mockserrors.InvalidRulesError{
			Message: "name cannot be empty",
		}
	}

	if rule.Path == "" {
		return mockserrors.InvalidRulesError{
			Message: "path cannot be empty",
		}
	}

	if rule.Status != "" && rule.Status != model.RuleStatusEnabled && rule.Status != model.RuleStatusDisabled {
		return mockserrors.InvalidRulesError{
			Message: "invalid status - only 'enabled' or 'disabled' are valid values",
		}
	}

	if err := validateHTTPMethod(rule.Method); err != nil {
		return err
	}

	if rule.Strategy != "" && rule.Strategy != model.RuleStrategyNormal && rule.Strategy != model.RuleStrategyRandom &&
		rule.Strategy != model.RuleStrategySequential {
		return mockserrors.InvalidRulesError{
			Message: fmt.Sprintf("invalid rule strategy - only '%s', '%s' or '%s' are valid values",
				model.RuleStrategyNormal, model.RuleStrategyRandom, model.RuleStrategySequential),
		}
	}

	return validateResponses(rule.Responses)
}

func validateResponses(responses []model.Response) error {
	if len(responses) == 0 {
		return mockserrors.InvalidRulesError{
			Message: "at least one response required",
		}
	}

	for _, response := range responses {
		if response.HTTPStatus < http.StatusOK || response.HTTPStatus > 599 {
			return mockserrors.InvalidRulesError{
				Message: fmt.Sprintf("%v is not a valid HTTP Status", response.HTTPStatus),
			}
		}
	}

	return nil
}

func validateHTTPMethod(method string) error {
	if method == "" {
		return mockserrors.InvalidRulesError{
			Message: "method cannot be empty",
		}
	}

	switch strings.ToUpper(method) {
	case http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete,
		http.MethodConnect, http.MethodOptions, http.MethodTrace:
		return nil
	default:
		return mockserrors.InvalidRulesError{
			Message: fmt.Sprintf("%s is not a valid HTTP Method", method),
		}
	}
}

func formatRule(rule *model.Rule) *model.Rule {
	rule.Method = strings.ToUpper(rule.Method)

	if rule.Status == "" {
		rule.Status = model.RuleStatusEnabled
	}

	if rule.Strategy == "" {
		rule.Strategy = model.RuleStrategyNormal
	}

	return rule
}
