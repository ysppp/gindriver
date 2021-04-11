package api

import (
	"fmt"
	"gindriver/models"
	"gindriver/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func OtherFiles(c *gin.Context) {
	username, ret := c.Get("SessionUser")
	if !ret {
		c.JSON(http.StatusUnauthorized, utils.ErrorWrapper(fmt.Errorf("not authorized")))
		return
	}
	//获取用户信息
	user := models.GetUserInfoByName(username)
	//获取用户文件使用明细数量
	fileDetailUse := models.GetFileDetailUse(user.FileStoreId)
	//获取音频类型文件
	otherFiles := models.GetTypeFile(5, user.FileStoreId)

	c.JSON(http.StatusOK, gin.H{
		"user":          user,
		"fileDetailUse": fileDetailUse,
		"otherFiles":    otherFiles,
		"otherCount":    len(otherFiles),
		"currOther":     "active",
		"currClass":     "active",
	})
}
