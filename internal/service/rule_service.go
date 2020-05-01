package service

import (
	"context"
	"strings"

	fwdcontext "github.com/nicopozo/mockserver/internal/context"
	"github.com/nicopozo/mockserver/internal/model"
	"github.com/nicopozo/mockserver/internal/repository"
)

type IRuleService interface {
	Save(ctx context.Context, rule *model.Rule) error
	Get(ctx context.Context, key string) (*model.Rule, error)
	Search(ctx context.Context, params map[string]interface{}, paging model.Paging) (*model.RuleList, error)
	SearchByMethodAndPath(ctx context.Context, method, path string) (*model.Rule, error)
}

type RuleService struct {
	RuleRepository repository.IRuleRepository
}

func (service *RuleService) Save(ctx context.Context, rule *model.Rule) error {
	logger := fwdcontext.Logger(ctx)

	logger.Debug(service, nil, "Entering RuleService Save()")

	if err := service.validateRule(rule); err != nil {
		return err
	}

	rule = service.formatRule(rule)

	return service.RuleRepository.Save(ctx, rule)
}

func (service *RuleService) Get(ctx context.Context, key string) (*model.Rule, error) {
	logger := fwdcontext.Logger(ctx)

	logger.Debug(service, nil, "Entering TaskService Get()")

	result, err := service.RuleRepository.Get(ctx, key)
	if err != nil {
		logger.Error(service, nil, err, "error getting task")
		return nil, err
	}

	return result, nil
}

func (service *RuleService) Search(ctx context.Context, params map[string]interface{},
	paging model.Paging) (*model.RuleList, error) {
	logger := fwdcontext.Logger(ctx)

	logger.Debug(service, nil, "Entering RuleService Search()")

	return service.RuleRepository.Search(ctx, params, paging)
}

func (service *RuleService) SearchByMethodAndPath(ctx context.Context, method, path string) (*model.Rule, error) {
	logger := fwdcontext.Logger(ctx)

	logger.Debug(service, nil, "Entering RuleService Search()")

	result, err := service.RuleRepository.SearchByMethodAndPath(ctx, method, path)
	if err != nil {
		logger.Error(service, nil, err, "error searching rules")
		return nil, err
	}

	return result, nil
}

func (service *RuleService) validateRule(rule *model.Rule) error {
	return nil
}

func (service *RuleService) formatRule(rule *model.Rule) *model.Rule {
	rule.Method = strings.ToLower(rule.Method)
	return rule
}
