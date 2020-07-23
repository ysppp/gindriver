package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadHandler(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf("Query: %s", c.Query("test")))
}

func NoRouterHandler(c *gin.Context) {
	c.String(http.StatusNotFound, "Not Found")
}

