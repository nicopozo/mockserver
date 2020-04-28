package service

import (
	"strings"

	newrelic "github.com/newrelic/go-agent"
	"github.com/nicopozo/mockserver/internal/model"
	"github.com/nicopozo/mockserver/internal/repository"
	log "github.com/nicopozo/mockserver/internal/utils/log"
)

type IRuleService interface {
	Save(rule *model.Rule, txn newrelic.Transaction, logger log.ILogger) error
	Get(key string, txn newrelic.Transaction, logger log.ILogger) (*model.Rule, error)
	Search(params map[string]interface{}, paging model.Paging, txn newrelic.Transaction,
		logger log.ILogger) (*model.RuleList, error)
	SearchByMethodAndPath(method, path string, txn newrelic.Transaction, logger log.ILogger) (*model.Rule, error)
}

type RuleService struct {
	RuleRepository repository.IRuleRepository
}

func (service *RuleService) Save(rule *model.Rule, txn newrelic.Transaction, logger log.ILogger) error {
	logger.Debug(service, nil, "Entering RuleService Save()")

	if err := service.validateRule(rule); err != nil {
		return err
	}

	rule = service.formatRule(rule)

	return service.RuleRepository.Save(rule, txn, logger)
}

func (service *RuleService) Get(key string, txn newrelic.Transaction,
	logger log.ILogger) (*model.Rule, error) {
	logger.Debug(service, nil, "Entering TaskService Get()")

	result, err := service.RuleRepository.Get(key, txn, logger)
	if err != nil {
		logger.Error(service, nil, err, "error getting task")
		return nil, err
	}

	return result, nil
}

func (service *RuleService) Search(params map[string]interface{}, paging model.Paging, txn newrelic.Transaction,
	logger log.ILogger) (*model.RuleList, error) {
	panic("implement me")
}

func (service *RuleService) SearchByMethodAndPath(method, path string, txn newrelic.Transaction,
	logger log.ILogger) (*model.Rule, error) {
	logger.Debug(service, nil, "Entering RuleService Search()")

	result, err := service.RuleRepository.SearchByMethodAndPath(method, path, txn, logger)
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
