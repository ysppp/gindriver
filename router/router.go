package router

import (
	"gindriver/api"
	"gindriver/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() *gin.Engine {
	app := gin.Default()

	// Static file handler
	//app.StaticFile("/", "public/index.html")
	// Dev frontend
	app.StaticFile("/", "frontend/dist/index.html")
	app.StaticFile("/umi.css", "frontend/dist/umi.css")
	app.StaticFile("/umi.js", "frontend/dist/umi.js")
	app.Static("/pubkeys/", "public/pubkeys")

	// API router
	apiRouter := app.Group("/api/auth/")
	apiRouter.POST("/register/begin", api.BeginRegistration)
	apiRouter.PATCH("/register/:name/finish", api.FinishRegistration)
	apiRouter.GET("/login/:name/begin", api.BeginLogin)
	apiRouter.PATCH("/login/:name/finish", api.FinishLogin)

	apiAuthRequiredRouter := app.Group("/api/user/")
	apiAuthRequiredRouter.Use(middleware.LoginRequired())
	{
		apiAuthRequiredRouter.GET("/:name", api.UserInfo)
	}

	apiFileRouter := app.Group("/api/file")
	apiFileRouter.Use(middleware.LoginRequired())
	{
		apiFileRouter.POST("/download", api.DownloadFile)
		apiFileRouter.POST("/folder/add", api.AddFolder)
		apiFileRouter.POST("/folder/update", api.UpdateFileFolder)
		apiFileRouter.GET("/getAll", api.GetAllFiles)
		apiFileRouter.POST("/upload", api.UploadHandler)
		apiFileRouter.POST("/update", api.UpdateFile)
		apiFileRouter.POST("/delete", api.DeleteFile)
		apiFileRouter.POST("/share/add", api.ShareFile)
		apiFileRouter.POST("/share/download", api.DownloadShareFile)
		apiFileRouter.POST("/share/show", api.SharePass)
	}
	app.GET("/api/files", api.GetFilesByType)
	// 404 Handler
	app.NoRoute(NoRouterHandler)

	return app
}

func NoRouterHandler(c *gin.Context) {
	c.String(http.StatusNotFound, "Not Found")
}
