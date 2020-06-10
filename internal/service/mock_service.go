package service

import (
	"context"
	"net/url"
	"strings"

	mockscontext "github.com/nicopozo/mockserver/internal/context"
	"github.com/nicopozo/mockserver/internal/model"
)

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

	u, err := url.Parse(path)
	if err != nil {
		logger.Error(service, nil, err, "error parsing path")
		return nil, err
	}

	response := result.Responses[0]

	values := strings.Split(u.Path, "/")
	pathParts := strings.Split(result.Path, "/")

	for i, part := range pathParts {
		if part != "" {
			first := string(part[0])
			last := string(part[len(part)-1])

			if first == "{" && last == "}" {
				response.Body = strings.Replace(response.Body, part, values[i], -1)
			}
		}
	}

	return &response, nil
}
