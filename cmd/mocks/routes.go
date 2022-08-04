package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/nicopozo/mockserver/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func mapRoutes(router *gin.Engine) {
	router.Static("/mock-service/admin", "web/dist")

	applicationContainer := container{}

	ruleController := applicationContainer.RuleController()
	router.POST("/mock-service/rules", ruleController.Create)
	router.GET("/mock-service/rules/:key", ruleController.Get)
	router.GET("/mock-service/rules", ruleController.Search)
	router.DELETE("/mock-service/rules/:key", ruleController.Delete)
	router.PUT("/mock-service/rules/:key", ruleController.Update)
	router.PUT("/mock-service/rules/:key/status", ruleController.UpdateStatus)

	mockController := applicationContainer.MockController()
	router.Any("/mock-service/mock/*rule", mockController.Execute)

	router.GET("/ping", ping)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
