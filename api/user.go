package api

import (
	"fmt"
	"gindriver/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserInfo(c *gin.Context) {
	user, ret := c.Get("SessionUser")
	fmt.Println(user)
	if !ret {
		c.JSON(http.StatusUnauthorized, utils.ErrorWrapper(fmt.Errorf("not authorized")))
		return
	}
	username := c.Param("name")
	if user != username {
		c.JSON(http.StatusForbidden, utils.ErrorWrapper(fmt.Errorf("access denied")))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"username": user,
	})
}
