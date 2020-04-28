package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicopozo/mockserver/internal/model"
)

func main() {
	router := gin.New()
	router.NoRoute(noRouteHandler)

	mapRoutes(router)

	if err := router.Run(":8080"); err != nil {
		panic(err.Error())
	}
}

func noRouteHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, model.NewError(model.ResourceNotFoundError, "no handler found for path"))
}
