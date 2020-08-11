package api

import (
	"gindriver/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserInfo(c *gin.Context) {
	c.JSON(http.StatusOK, utils.SuccessWrapper("OK!"))
}
