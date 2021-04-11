package api

import (
	"fmt"
	"gindriver/lib"
	"gindriver/models"
	"gindriver/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

//创建分享文件
func ShareFile(c *gin.Context) {
	username, ret := c.Get("SessionUser")
	if !ret {
		c.JSON(http.StatusUnauthorized, utils.ErrorWrapper(fmt.Errorf("not authorized")))
		return
	}
	//获取用户信息
	user := models.GetUserInfoByName(username)

	fId := c.Query("id")
	url := c.Query("url")
	//获取内容
	code := utils.GetRandomString(4)

	fileId, _ := strconv.Atoi(fId)
	hash := models.CreateShare(code, user.Name, uint64(fileId))

	c.JSON(http.StatusOK, gin.H{
		"url":  url + "?f=" + hash,
		"code": code,
	})
}

//分享文件页面
func SharePass(c *gin.Context) {
	f := c.Query("f")

	//获取分享信息
	shareInfo := models.GetShareInfo(f)
	//获取文件信息
	file := models.GetFileInfo(strconv.Itoa(int(shareInfo.FileId)))

	c.HTML(http.StatusOK, "share.html", gin.H{
		"id":       shareInfo.FileId,
		"username": shareInfo.UserName,
		"fileType": file.Type,
		"filename": file.FileName + file.PostFix,
		"hash":     shareInfo.ShareHash,
	})
}

//下载分享文件
func DownloadShareFile(c *gin.Context) {
	fileId := c.Query("id")
	code := c.Query("code")
	hash := c.Query("hash")

	fileInfo := models.GetFileInfo(fileId)

	//校验提取码
	if ok := models.VerifyShareCode(fileId, strings.ToLower(code)); !ok {
		c.Redirect(http.StatusMovedPermanently, "/file/share?f="+hash)
		return
	}

	//从oss获取文件
	fileData := lib.DownloadOss(fileInfo.FileHash, fileInfo.PostFix)
	//下载次数+1
	//model.DownloadNumAdd(fileId)

	c.Header("Content-disposition", "attachment;filename=\""+fileInfo.FileName+fileInfo.PostFix+"\"")
	c.Data(http.StatusOK, "application/octect-stream", fileData)
}