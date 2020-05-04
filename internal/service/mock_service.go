package service

import (
	"context"

	mockscontext "github.com/nicopozo/mockserver/internal/context"
	"github.com/nicopozo/mockserver/internal/model"
)

//nolint:lll
//go:generate mockgen -destination=../utils/test/mocks/mock_service_mock.go -package=mocks -source=./mock_service.go

type IMockService interface {
	SearchResponseForMethodAndPath(ctx context.Context, method, path string) (*model.Response, error)
}

type MockService struct {
	RuleService IRuleService
}

func (service MockService) SearchResponseForMethodAndPath(ctx context.Context,
	method, path string) (*model.Response, error) {
	logger := mockscontext.Logger(ctx)

	logger.Debug(service, nil, "Entering MockService Execute()")

	result, err := service.RuleService.SearchByMethodAndPath(ctx, method, path)
	if err != nil {
		logger.Error(service, nil, err, "error searching responses")
		return nil, err
	}

	return &result.Responses[0], nil
}
