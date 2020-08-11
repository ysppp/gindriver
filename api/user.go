package api

import (
	"fmt"
	"gindriver/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserInfo(c *gin.Context) {
	user, ret := c.Get("SessionUser")
	if !ret {
		c.JSON(http.StatusUnauthorized, utils.ErrorWrapper(fmt.Errorf("not authorized")))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessWrapper(fmt.Sprint(user)))
}
