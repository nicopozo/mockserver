package service

import (
	newrelic "github.com/newrelic/go-agent"
	"github.com/nicopozo/mockserver/internal/model"
	log2 "github.com/nicopozo/mockserver/internal/utils/log"
)

type IMockService interface {
	SearchResponseForMethodAndPath(method, path string, txn newrelic.Transaction,
		logger log2.ILogger) (*model.Response, error)
}

type MockService struct {
	RuleService IRuleService
}

func (service MockService) SearchResponseForMethodAndPath(method, path string, txn newrelic.Transaction,
	logger log2.ILogger) (*model.Response, error) {
	logger.Debug(service, nil, "Entering MockService Execute()")

	result, err := service.RuleService.SearchByMethodAndPath(method, path, txn, logger)
	if err != nil {
		logger.Error(service, nil, err, "error searching responses")
		return nil, err
	}

	return &result.Responses[0], nil
}
