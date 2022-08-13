package main

import (
	"fmt"
	"os"

	"github.com/nicopozo/mockserver/internal/controller"
	"github.com/nicopozo/mockserver/internal/repository"
	"github.com/nicopozo/mockserver/internal/service"
)

type container struct {
	mockController *controller.MockController
	ruleController *controller.RuleController

	mockService service.MockService
	ruleService service.RuleService

	ruleRepository repository.RuleRepository
}

func (c *container) MockController() *controller.MockController {
	if c.mockController == nil {
		mockService := c.MockService()
		mockController := &controller.MockController{
			MockService: mockService,
		}

		c.mockController = mockController
	}

	return c.mockController
}

func (c *container) RuleController() *controller.RuleController {
	if c.ruleController == nil {
		ruleService := c.RuleService()
		ruleController := &controller.RuleController{
			RuleService: ruleService,
		}

		c.ruleController = ruleController
	}

	return c.ruleController
}

func (c *container) MockService() service.MockService {
	if c.mockService == nil {
		ruleService := c.RuleService()

		mockService, err := service.NewMockService(ruleService)
		if err != nil {
			panic("error creating Mock Service")
		}

		c.mockService = mockService
	}

	return c.mockService
}

func (c *container) RuleService() service.RuleService {
	if c.ruleService == nil {
		ruleRepository := c.RuleRepository()

		ruleService, err := service.NewRuleService(ruleRepository)
		if err != nil {
			panic("error creating Rule Service")
		}

		c.ruleService = ruleService
	}

	return c.ruleService
}

func (c *container) RuleRepository() repository.RuleRepository {
	if c.ruleRepository == nil {
		dataSource := os.Getenv("MOCKS_DATASOURCE")

		switch dataSource {
		case "file", "":
			filePath := os.Getenv("MOCKS_FILE")

			repo, err := repository.NewRuleFileRepository(filePath)
			if err != nil {
				panic(fmt.Sprintf("Error creating File Repository: %s", err.Error()))
			}

			c.ruleRepository = repo

			return c.ruleRepository

		case "elastic":
			c.ruleRepository = repository.NewRuleElasticRepository()

			return c.ruleRepository
		case "mysql":
			db, err := repository.GetDB()
			if err != nil {
				panic(fmt.Sprintf("Error connecting to mysql DB: %s", err.Error()))
			}

			c.ruleRepository = repository.NewRuleMySQLRepository(db)

			return c.ruleRepository
		}

		panic("invalid datasource type: " + dataSource)
	}

	return c.ruleRepository
}
