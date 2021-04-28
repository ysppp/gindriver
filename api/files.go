package api

import (
	"fmt"
	"gindriver/lib"
	"gindriver/models"
	"gindriver/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//全部文件页面
func GetAllFiles(c *gin.Context) {
	username, ret := c.Get("SessionUser")
	if !ret {
		c.JSON(http.StatusUnauthorized, utils.ErrorWrapper(fmt.Errorf("not authorized")))
		return
	}
	//获取用户信息
	user := models.GetUserInfoByName(username)
	folderId, _ := strconv.ParseUint(c.DefaultQuery("fId", "0"), 10, 64)

	//获取当前目录所有文件
	files := models.GetUserFile(folderId, user.FileStoreId)
	//获取当前目录所有文件夹
	fileFolders := models.GetFileFolder(folderId, user.FileStoreId)

	//获取父级的文件夹信息
	parentFolder := models.GetParentFolder(folderId)

	//获取当前目录所有父级
	currentAllParent := models.GetCurrentAllParent(parentFolder, make([]models.FileFolder, 0))

	//获取当前目录信息
	currentFolder := models.GetCurrentFolder(folderId)

	//获取用户文件使用明细数量
	fileDetailUse := models.GetFileDetailUse(user.FileStoreId)

	c.JSON(http.StatusOK, gin.H{
		"currAll":          "active",
		"user":             user,
		"fId":              currentFolder.FolderId,
		"fName":            currentFolder.FolderName,
		"files":            files,
		"fileFolders":      fileFolders,
		"parentFolder":     parentFolder,
		"currentAllParent": currentAllParent,
		"fileDetailUse":    fileDetailUse,
	})
}

//处理新建文件夹
func AddFolder(c *gin.Context) {
	username, ret := c.Get("SessionUser")
	if !ret {
		c.JSON(http.StatusUnauthorized, utils.ErrorWrapper(fmt.Errorf("not authorized")))
		return
	}
	//获取用户信息
	user := models.GetUserInfoByName(username)

	folderName := c.PostForm("fileFolderName")
	parentId, _ := strconv.ParseUint(c.DefaultPostForm("parentFolderId", "0"), 10, 64)

	//新建文件夹数据
	models.CreateFolder(folderName, parentId, user.FileStoreId)

	//获取父文件夹信息
	//parent := models.GetParentFolder(parentId)

	//c.Redirect(http.StatusMovedPermanently, "/cloud/files?fId=" + parentId + "&fName=" + parent.FolderName)
}

func DownloadFile(c *gin.Context) {
	fId := c.Query("fId")

	file := models.GetFileInfo(fId)
	if file.FileHash == "" {
		return
	}

	//从oss获取文件
	fileData := lib.DownloadOss(file.FileHash, file.PostFix)
	//下载次数+1
	//models.DownloadNumAdd(fId)

	c.Header("Content-disposition", "attachment;filename=\""+file.FileName+file.PostFix+"\"")
	c.Data(http.StatusOK, "application/octect-stream", fileData)
}

//删除文件
func DeleteFile(c *gin.Context) {
	username, ret := c.Get("SessionUser")
	if !ret {
		c.JSON(http.StatusUnauthorized, utils.ErrorWrapper(fmt.Errorf("not authorized")))
		return
	}
	//获取用户信息
	user := models.GetUserInfoByName(username)

	fId := c.DefaultQuery("fId", "")
	folderId := c.Query("folder")
	if fId == "" {
		return
	}

	//删除数据库文件数据
	models.DeleteUserFile(fId, folderId, user.FileStoreId)

	c.Redirect(http.StatusMovedPermanently, "/cloud/files?fid="+folderId)
}

//删除文件夹
func DeleteFileFolder(c *gin.Context) {
	fId, err := strconv.ParseUint(c.DefaultQuery("fId", ""), 10, 64)
	if err != nil {
		fmt.Println("error: %s", err)
		return
	}
	//获取要删除的文件夹信息 取到父级目录重定向
	folderInfo := models.GetCurrentFolder(fId)

	//删除文件夹并删除文件夹中的文件信息
	models.DeleteFileFolder(fId)

	c.Redirect(http.StatusMovedPermanently, "/cloud/files?fId="+strconv.FormatUint(folderInfo.ParentFolderId, 10))
}

//修改文件夹名
func UpdateFileFolder(c *gin.Context) {
	fileFolderName := c.PostForm("fileFolderName")
	fileFolderId, _ := strconv.ParseUint(c.PostForm("fileFolderId"), 10, 64)

	fileFolder := models.GetCurrentFolder(fileFolderId)

	models.UpdateFolderName(fileFolderId, fileFolderName)

	c.Redirect(http.StatusMovedPermanently, "/cloud/files?fId="+strconv.FormatUint(fileFolder.ParentFolderId, 10))
}
