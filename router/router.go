package router

import (
	"fmt"
	"gindriver/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() *gin.Engine {
	app := gin.Default()

	// Static file handler
	app.Static("/", "public")

	// API router
	app.POST("/upload", UploadHandler)
	app.POST("/api/register/begin", api.BeginRegistration)

	// 404 Handler
	app.NoRoute(NoRouterHandler)

	return app
}

func UploadHandler(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf("Query: %s", c.Query("test")))
}

func NoRouterHandler(c *gin.Context) {
	c.String(http.StatusNotFound, "Not Found")
}
