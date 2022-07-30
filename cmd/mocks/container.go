package main

import (
	"fmt"
	"os"

	"github.com/nicopozo/mockserver/internal/controller"
	"github.com/nicopozo/mockserver/internal/repository"
	"github.com/nicopozo/mockserver/internal/service"
)

type container struct {
	ruleController controller.RuleController
	mockController controller.MockController
	ruleRepository repository.IRuleRepository
}

func (c *container) RuleController() controller.RuleController {
	ruleRepository := c.RuleRepository()
	ruleService := &service.RuleService{
		RuleRepository: ruleRepository,
	}
	ruleController := controller.RuleController{
		RuleService: ruleService,
	}

	return ruleController
}

func (c *container) MockController() controller.MockController {
	ruleRepository := c.RuleRepository()
	ruleService := &service.RuleService{
		RuleRepository: ruleRepository,
	}
	mockService := &service.MockService{
		RuleService: ruleService,
	}
	mockController := controller.MockController{
		MockService: mockService,
	}

	return mockController
}

func (c *container) RuleRepository() repository.IRuleRepository {
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
