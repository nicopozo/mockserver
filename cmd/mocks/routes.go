package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicopozo/mockserver/internal/controller"
	"github.com/nicopozo/mockserver/internal/repository"
	"github.com/nicopozo/mockserver/internal/service"
)

func mapRoutes(router *gin.Engine) {
	ruleController := newRuleController()
	router.POST("/reconciliations/mocks/rules", ruleController.Create)
	router.GET("/reconciliations/mocks/rules/:key", ruleController.Get)

	mockController := newMockController()
	router.Any("/reconciliations/mocks/mock/*rule", mockController.Execute)

	router.GET("/ping", ping())
}

func newRuleController() controller.RuleController {
	ruleRepository := &repository.RuleElasticRepository{}
	ruleService := &service.RuleService{
		RuleRepository: ruleRepository,
	}
	ruleController := controller.RuleController{
		RuleService: ruleService,
	}

	return ruleController
}

func newMockController() controller.MockController {
	ruleRepository := &repository.RuleElasticRepository{}
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

func ping() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	}
}
