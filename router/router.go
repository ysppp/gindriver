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
	app.StaticFile("/", "public/index.html")
	app.Static("/pubkey/", "public/pubkeys")

	// API router
	app.POST("/api/upload", UploadHandler)
	app.POST("/api/register/begin", api.BeginRegistration)
	app.PATCH("/api/register/:name/finish", api.FinishRegistration)
	app.GET("/api/login/:name/begin", api.BeginLogin)
	app.PATCH("/api/login/:name/finish", api.FinishLogin)

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
