package service

import (
	"context"

	fwdcontext "github.com/nicopozo/mockserver/internal/context"

	"github.com/nicopozo/mockserver/internal/model"
)

type IMockService interface {
	SearchResponseForMethodAndPath(ctx context.Context, method, path string) (*model.Response, error)
}

type MockService struct {
	RuleService IRuleService
}

func (service MockService) SearchResponseForMethodAndPath(ctx context.Context, method, path string) (*model.Response, error) {
	logger := fwdcontext.Logger(ctx)

	logger.Debug(service, nil, "Entering MockService Execute()")

	result, err := service.RuleService.SearchByMethodAndPath(ctx, method, path)
	if err != nil {
		logger.Error(service, nil, err, "error searching responses")
		return nil, err
	}

	return &result.Responses[0], nil
}
