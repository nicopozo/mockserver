package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nicopozo/mockserver/internal/model"
)

// @Title Mock Server
// @Description Mock Server is intended to serve mocks for different APIs during development process.
// We need to create rules containing the response we expect for a given endpoint
// @Version 2.0.0
// @Host localhost:8080
// @BasePath /mock-server
// @Schemes http
// @contact.name Nicolas Pozo
// @contact.email nicopozo@gmail.com
// main.
func main() {
	router := gin.New()
	router.NoRoute(noRouteHandler)

	router.Use(cors.Default())

	mapRoutes(router)

	if err := router.Run(":8080"); err != nil {
		panic(err.Error())
	}
}

func noRouteHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, model.NewError(model.ResourceNotFoundError, "no handler found for path"))
}
