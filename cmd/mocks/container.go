package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jmoiron/sqlx"
	"github.com/nicopozo/mockserver/internal/configs"
	"github.com/nicopozo/mockserver/internal/controller"
	"github.com/nicopozo/mockserver/internal/repository"
	"github.com/nicopozo/mockserver/internal/service"
	"go.uber.org/dig"
)

var (
	errInvalidDataSource    = errors.New("invalid datasource type")
	errDBNotInitialized     = errors.New("database connection not initialized")
	errDynamoNotInitialized = errors.New("dynamodb client not initialized")
)

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
		// Config
		configs.New,

		// Persistence
		provideSQLDB,
		provideDynamoClient,

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

// RepositoryDeps contains dependencies for repositories.
type RepositoryDeps struct {
	dig.In
	Config *configs.Config
	DB     *sqlx.DB         `optional:"true"`
	Dynamo *dynamodb.Client `optional:"true"`
}

// provideSQLDB provides a SQL DB connection if the datasource is SQL-based.
func provideSQLDB(cfg *configs.Config) (*sqlx.DB, error) {
	if !cfg.IsSQL() {
		return nil, nil //nolint:nilnil
	}

	db, err := repository.NewSQLDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create SQL DB: %w", err)
	}

	return db, nil
}

// provideDynamoClient provides a DynamoDB client if the datasource is DynamoDB-based.
func provideDynamoClient(cfg *configs.Config) (*dynamodb.Client, error) {
	if !cfg.IsDynamo() {
		return nil, nil //nolint:nilnil
	}

	client, err := repository.NewDynamoClient(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create DynamoDB client: %w", err)
	}

	return client, nil
}

// newRuleRepository contains the logic to select the appropriate repository implementation
// based on the configuration.
func newRuleRepository(deps RepositoryDeps) (repository.RuleRepository, error) {
	switch {
	case deps.Config.DataSource == "file" || deps.Config.DataSource == "":
		repo, err := repository.NewRuleFileRepository(deps.Config)
		if err != nil {
			return nil, fmt.Errorf("failed to create rule file repository: %w", err)
		}

		return repo, nil

	case deps.Config.IsSQL():
		if deps.DB == nil {
			return nil, errDBNotInitialized
		}

		return repository.NewRuleSQLRepository(deps.DB), nil

	case deps.Config.IsDynamo():
		if deps.Dynamo == nil {
			return nil, errDynamoNotInitialized
		}

		return repository.NewDynamoRuleRepository(deps.Dynamo, deps.Config), nil
	}

	return nil, errInvalidDataSource
}

// newLogRepository selects the appropriate log storage implementation.
func newLogRepository(deps RepositoryDeps) (repository.LogRepository, error) {
	switch {
	case deps.Config.IsSQL():
		if deps.DB == nil {
			return nil, errDBNotInitialized
		}

		return repository.NewLogSQLRepository(deps.DB), nil

	case deps.Config.IsDynamo():
		if deps.Dynamo == nil {
			return nil, errDynamoNotInitialized
		}

		return repository.NewDynamoLogRepository(deps.Dynamo, deps.Config), nil

	default:
		// Default to in-memory for "file" mode or if not specified
		return repository.NewLogMemoryRepository(), nil
	}
}
