package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/nicopozo/mockserver/docs"
	"github.com/nicopozo/mockserver/internal/controller"
	"github.com/nicopozo/mockserver/internal/repository"
	"github.com/nicopozo/mockserver/internal/service"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func mapRoutes(router *gin.Engine) {
	ruleController := newRuleController()
	router.POST("/mock-server/rules", ruleController.Create)
	router.GET("/mock-server/rules/:key", ruleController.Get)
	router.GET("/mock-server/rules", ruleController.Search)
	router.DELETE("/mock-server/rules/:key", ruleController.Delete)
	router.PUT("/mock-server/rules/:key", ruleController.Update)

	mockController := newMockController()
	router.Any("/mock-server/mock/*rule", mockController.Execute)

	router.GET("/ping", ping)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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

func ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
