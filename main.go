package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Printf("Hello world! %d", http.StatusOK)

	app := gin.Default()

	app.Static("/", "public")

	app.POST("/upload", UploadHandler)

	app.NoRoute(NoRouterHandler)

	err := app.Run(ListenAddr)

	if err != nil {
		fmt.Printf("err: %s", err)
	}
}
