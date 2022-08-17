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

type RuleService interface {
	Save(ctx context.Context, rule model.Rule) (model.Rule, error)
	Update(ctx context.Context, key string, rule model.Rule) (model.Rule, error)
	UpdateStatus(ctx context.Context, key string, rule model.RuleStatus) (model.Rule, error)
	Get(ctx context.Context, key string) (model.Rule, error)
	Search(ctx context.Context, params map[string]interface{}, paging model.Paging) (model.RuleList, error)
	SearchByMethodAndPath(ctx context.Context, method, path string) (model.Rule, error)
	Delete(ctx context.Context, key string) error
}

type ruleService struct {
	RuleRepository repository.RuleRepository
}

func NewRuleService(ruleRepository repository.RuleRepository) (RuleService, error) {
	if ruleRepository == nil {
		return nil, fmt.Errorf("rule repository cannot be nil") //nolint:goerr113
	}

	return &ruleService{
		RuleRepository: ruleRepository,
	}, nil
}

func (ruleService *ruleService) Save(ctx context.Context, rule model.Rule) (model.Rule, error) {
	logger := mockscontext.Logger(ctx)

	logger.Debug(ruleService, nil, "Entering ruleService Save()")

	if err := validateRule(rule); err != nil {
		logger.Error(ruleService, nil, err, "Rule Validation failed")

		return model.Rule{}, err
	}

	rule = formatRule(rule)

	repoRule, err := ruleService.RuleRepository.Create(ctx, &rule)
	if err != nil {
		return model.Rule{}, fmt.Errorf("error creating rule - %w", err)
	}

	return *repoRule, nil
}

func (ruleService *ruleService) Update(ctx context.Context, key string, rule model.Rule) (model.Rule, error) {
	logger := mockscontext.Logger(ctx)

	logger.Debug(ruleService, nil, "Entering ruleService Update()")

	if err := validateRule(rule); err != nil {
		logger.Error(ruleService, nil, err, "Rule Validation failed")

		return model.Rule{}, fmt.Errorf("error validating rule - %w", err)
	}

	rule.Key = key
	rule = formatRule(rule)

	repoRule, err := ruleService.RuleRepository.Update(ctx, &rule)
	if err != nil {
		return model.Rule{}, fmt.Errorf("error updating rule - %w", err)
	}

	return *repoRule, nil
}

func (ruleService *ruleService) UpdateStatus(ctx context.Context, key string,
	ruleStatus model.RuleStatus) (model.Rule, error) {
	logger := mockscontext.Logger(ctx)

	logger.Debug(ruleService, nil, "Entering ruleService Update()")

	if err := ruleStatus.Validate(); err != nil {
		return model.Rule{}, fmt.Errorf("error validating rule, %w", err)
	}

	rule, err := ruleService.Get(ctx, key)
	if err != nil {
		return model.Rule{}, err
	}

	rule.Status = ruleStatus.Status

	rule = formatRule(rule)

	repoRule, err := ruleService.RuleRepository.Update(ctx, &rule)
	if err != nil {
		return model.Rule{}, fmt.Errorf("error updating rule status - %w", err)
	}

	return *repoRule, nil
}

func (ruleService *ruleService) Get(ctx context.Context, key string) (model.Rule, error) {
	logger := mockscontext.Logger(ctx)

	logger.Debug(ruleService, nil, "Entering TaskService Get()")

	repoRule, err := ruleService.RuleRepository.Get(ctx, key)
	if err != nil {
		logger.Error(ruleService, nil, err, "error getting task")

		return model.Rule{}, fmt.Errorf("error getting rule, %w", err)
	}

	return *repoRule, nil
}

func (ruleService *ruleService) Search(ctx context.Context, params map[string]interface{},
	paging model.Paging) (model.RuleList, error) {
	logger := mockscontext.Logger(ctx)

	logger.Debug(ruleService, nil, "Entering ruleService Search()")

	ReposRules, err := ruleService.RuleRepository.Search(ctx, params, paging)
	if err != nil {
		return model.RuleList{}, fmt.Errorf("error searching rule - %w", err)
	}

	return *ReposRules, nil
}

func (ruleService *ruleService) SearchByMethodAndPath(ctx context.Context, method, path string) (model.Rule, error) {
	logger := mockscontext.Logger(ctx)

	logger.Debug(ruleService, nil, "Entering ruleService Search()")

	result, err := ruleService.RuleRepository.SearchByMethodAndPath(ctx, method, path)
	if err != nil {
		logger.Error(ruleService, nil, err, "error searching rules")

		return model.Rule{}, fmt.Errorf("error searching rule, %w", err)
	}

	return *result, nil
}

func (ruleService *ruleService) Delete(ctx context.Context, key string) error {
	logger := mockscontext.Logger(ctx)

	logger.Debug(ruleService, nil, "Entering TaskService Get()")

	err := ruleService.RuleRepository.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("error deleting rule - %w", err)
	}

	return nil
}

//nolint:cyclop
func validateRule(rule model.Rule) error {
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

	if !strings.HasPrefix(rule.Path, "/") {
		return mockserrors.InvalidRulesError{
			Message: "path must start with '/'",
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
		rule.Strategy != model.RuleStrategySequential && rule.Strategy != model.RuleStrategyScene {
		return mockserrors.InvalidRulesError{
			Message: fmt.Sprintf("invalid rule strategy - only '%s', '%s', '%s' or '%s' are valid values",
				model.RuleStrategyNormal, model.RuleStrategyRandom, model.RuleStrategySequential,
				model.RuleStrategyScene),
		}
	}

	if err := validateVariables(rule.Variables); err != nil {
		return err
	}

	if err := validateAssertions(rule.Assertions); err != nil {
		return err
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

func validateVariables(variables []*model.Variable) error {
	if variables == nil {
		return nil
	}

	for _, variable := range variables {
		if err := variable.Validate(); err != nil {
			return fmt.Errorf("error validating variable, %w", err)
		}
	}

	return nil
}

func validateAssertions(assertions []*model.Assertion) error {
	if assertions == nil {
		return nil
	}

	for _, assertion := range assertions {
		if err := assertion.Validate(); err != nil {
			return fmt.Errorf("error validating assertions, %w", err)
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

func formatRule(rule model.Rule) model.Rule {
	rule.Method = strings.ToUpper(rule.Method)

	if rule.Status == "" {
		rule.Status = model.RuleStatusEnabled
	}

	if rule.Strategy == "" {
		rule.Strategy = model.RuleStrategyNormal
	}

	return rule
}
