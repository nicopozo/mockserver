package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/nicopozo/mockserver/api"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func mapRoutes(router *gin.Engine, api MockContainer) {
	router.Static("/mock-service/admin", "web/dist")

	ruleController := api.Controllers.RuleController
	router.POST("/mock-service/rules", ruleController.Create)
	router.GET("/mock-service/rules/:key", ruleController.Get)
	router.GET("/mock-service/rules", ruleController.Search)
	router.DELETE("/mock-service/rules/:key", ruleController.Delete)
	router.PUT("/mock-service/rules/:key", ruleController.Update)
	router.PUT("/mock-service/rules/:key/status", ruleController.UpdateStatus)
	router.GET("/mock-service/rules/export", ruleController.Export)
	router.POST("/mock-service/rules/import", ruleController.Import)

	mockController := api.Controllers.MockController
	router.Any("/mock-service/mock/*rule", mockController.Execute)

	logController := api.Controllers.LogController
	router.GET("/mock-service/logs", logController.GetLogs)
	router.DELETE("/mock-service/logs", logController.ClearLogs)

	router.GET("/ping", ping)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
