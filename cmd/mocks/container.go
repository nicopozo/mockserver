package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/nicopozo/mockserver/internal/controller"
	"github.com/nicopozo/mockserver/internal/repository"
	"github.com/nicopozo/mockserver/internal/service"
	"go.uber.org/dig"
)

var errInvalidDataSource = errors.New("invalid datasource type")

// MockContainer is a specialized container that uses dig.In to automatically
// resolve multiple dependencies based on tag/type mapping.
type MockContainer struct {
	dig.In

	Controllers struct {
		dig.In
		MockController *controller.MockController
		RuleController *controller.RuleController
		LogController  *controller.LogController
	}
}

// BuildContainer initialize the dependency injection container.
func BuildContainer() MockContainer {
	container := dig.New()

	providers := []any{
		// Repositories
		newRuleRepository,
		newLogRepository,

		// Services
		service.NewLogService,
		service.NewRuleService,
		service.NewMockService,

		// Controllers
		controller.NewMockController,
		controller.NewRuleController,
		controller.NewLogController,
	}

	for _, provider := range providers {
		if err := container.Provide(provider); err != nil {
			panic(fmt.Sprintf("failed to register dependency: %v", err))
		}
	}

	var api MockContainer

	if err := container.Invoke(func(c MockContainer) {
		api = c
	}); err != nil {
		panic(fmt.Sprintf("failed to resolve container: %v", err))
	}

	return api
}

// newRuleRepository contains the logic to select the appropriate repository implementation
// based on the MOCKS_DATASOURCE environment variable.
func newRuleRepository() (repository.RuleRepository, error) {
	dataSource := os.Getenv("MOCKS_DATASOURCE")

	switch dataSource {
	case "file", "":
		filePath := os.Getenv("MOCKS_FILE")

		repo, err := repository.NewRuleFileRepository(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to create rule file repository: %w", err)
		}

		return repo, nil

	case "mysql", "postgres":
		db, err := repository.GetDB()
		if err != nil {
			return nil, fmt.Errorf("error connecting to %s DB: %w", dataSource, err)
		}

		return repository.NewRuleSQLRepository(db), nil
	}

	return nil, errInvalidDataSource
}

// newLogRepository selects the appropriate log storage implementation.
func newLogRepository() (repository.LogRepository, error) {
	dataSource := os.Getenv("MOCKS_DATASOURCE")

	switch dataSource {
	case "mysql", "postgres":
		db, err := repository.GetDB()
		if err != nil {
			return nil, fmt.Errorf("error connecting to %s DB for logs: %w", dataSource, err)
		}

		return repository.NewLogSQLRepository(db), nil

	default:
		// Default to in-memory for "file" mode or if not specified
		return repository.NewLogMemoryRepository(), nil
	}
}
