package api

import (
	"fmt"
	"gindriver/lib"
	"gindriver/models"
	"gindriver/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type structOfShare struct {
	FileId uint64 `json:"fileId"`
	Code   string `json:"code"`
	Hash   string `json:"hash"`
	Url    string `json:"url"`
}

//创建分享文件
func ShareFile(c *gin.Context) {
	username, ret := c.Get("SessionUser")
	if !ret {
		c.JSON(http.StatusUnauthorized, utils.ErrorWrapper(fmt.Errorf("not authorized")))
		return
	}
	//获取用户信息
	user := models.GetUserInfoByName(username)

	json := structOfShare{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}
	fileId := json.FileId
	url := json.Url
	//获取内容
	code := utils.GetRandomString(4)

	hash := models.CreateShare(code, user.Name, uint64(fileId))

	c.JSON(http.StatusOK, gin.H{
		"url":  url + "?f=" + hash,
		"code": code,
	})
}

//分享文件页面
func SharePass(c *gin.Context) {
	json := structOfShare{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}
	hash := json.Hash
	//获取分享信息
	shareInfo := models.GetShareInfo(hash)
	//获取文件信息
	file := models.GetFileInfo(shareInfo.FileId)

	c.JSON(http.StatusOK, gin.H{
		"id":       shareInfo.FileId,
		"username": shareInfo.UserName,
		"fileType": file.Type,
		"filename": file.FileName,
		"hash":     shareInfo.ShareHash,
	})
}

//下载分享文件
func DownloadShareFile(c *gin.Context) {
	json := structOfShare{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}
	fileId := json.FileId
	code := json.Code
	hash := json.Hash

	fileInfo := models.GetFileInfo(fileId)

	//校验提取码
	if ok := models.VerifyShareCode(fileId, strings.ToLower(code), hash); !ok {
		c.JSON(http.StatusBadRequest, utils.SuccessWrapper("Error code!"))
		return
	}

	//从oss获取文件
	fileData := lib.DownloadOss(fileInfo.FileHash, fileInfo.PostFix)
	//下载次数+1
	//model.DownloadNumAdd(fileId)

	c.Header("Content-disposition", "attachment;filename=\""+fileInfo.FileName+fileInfo.PostFix+"\"")
	c.Data(http.StatusOK, "application/octect-stream", fileData)
}
