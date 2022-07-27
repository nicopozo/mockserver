package main

import (
	"fmt"
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
	router.Static("/mock-service/admin", "/Users/npozo/Proyectos/mockserver/web/dist")

	ruleController := newRuleController()
	router.POST("/mock-service/rules", ruleController.Create)
	router.GET("/mock-service/rules/:key", ruleController.Get)
	router.GET("/mock-service/rules", ruleController.Search)
	router.DELETE("/mock-service/rules/:key", ruleController.Delete)
	router.PUT("/mock-service/rules/:key", ruleController.Update)
	router.PUT("/mock-service/rules/:key/status", ruleController.UpdateStatus)

	mockController := newMockController()
	router.Any("/mock-service/mock/*rule", mockController.Execute)

	router.GET("/ping", ping)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func newRuleController() controller.RuleController {
	db, err := repository.GetDB()
	if err != nil {
		panic(fmt.Sprintf("Error connecting to mysql DB: %s", err.Error()))
	}

	ruleRepository := repository.NewRuleMySQLRepository(db)
	ruleService := &service.RuleService{
		RuleRepository: ruleRepository,
	}
	ruleController := controller.RuleController{
		RuleService: ruleService,
	}

	return ruleController
}

func newMockController() controller.MockController {
	db, err := repository.GetDB()
	if err != nil {
		panic(fmt.Sprintf("Error connecting to mysql DB: %s", err.Error()))
	}

	ruleRepository := repository.NewRuleMySQLRepository(db)
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
