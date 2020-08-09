package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func ErrorWrapper(err error) gin.H {
	return gin.H{
		"error": fmt.Sprint(err),
	}
}
