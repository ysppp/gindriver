package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func ErrorWrapper(err error) gin.H {
	return gin.H{
		"status":  "error",
		"message": fmt.Sprint(err),
	}
}

func SuccessWrapper(str string) gin.H {
	return gin.H{
		"status":  "success",
		"message": str,
	}
}
