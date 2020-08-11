package router

import (
	"fmt"
	"gindriver/api"
	"gindriver/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() *gin.Engine {
	app := gin.Default()

	// Static file handler
	app.StaticFile("/", "public/index.html")
	app.Static("/pubkey/", "public/pubkeys")

	// API router
	apiRouter := app.Group("/api/auth/")
	apiRouter.POST("/register/begin", api.BeginRegistration)
	apiRouter.PATCH("/register/:name/finish", api.FinishRegistration)
	apiRouter.GET("/login/:name/begin", api.BeginLogin)
	apiRouter.PATCH("/login/:name/finish", api.FinishLogin)

	// Auth required router
	apiAuthRequiredRouter := app.Group("/api/user/")
	apiAuthRequiredRouter.Use(middleware.LoginRequired())
	{
		apiAuthRequiredRouter.GET("/info", api.UserInfo)
	}

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
