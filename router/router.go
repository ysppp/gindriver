package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() *gin.Engine {
	app := gin.Default()

	app.Static("/", "public")
	app.POST("/upload", UploadHandler)
	app.NoRoute(NoRouterHandler)

	return app
}

func UploadHandler(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf("Query: %s", c.Query("test")))
}

func NoRouterHandler(c *gin.Context) {
	c.String(http.StatusNotFound, "Not Found")
}
