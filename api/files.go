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

type structOfJson struct {
	FileFolderName string `json:"fileFolderName"`
	ParentFolderId uint64 `json:"parentFolderId"`
	FileId         uint64 `json:"fId"`
	FolderId       uint64 `json:"folderId"`
	FileName       string `json:"fileName"`
	ToFolderId     uint64 `json:"toFolderId"`
}

func GetFilesByType(c *gin.Context) {
	username, ret := c.Get("SessionUser")
	if !ret {
		c.JSON(http.StatusUnauthorized, utils.ErrorWrapper(fmt.Errorf("not authorized")))
		return
	}
	//获取用户信息
	user := models.GetUserInfoByName(username)

	var Files []models.File

	fileType := c.Query("type")
	switch fileType {
	case "1":
		Files = models.GetTypeFile(1, user.FileStoreId)
	case "2":
		Files = models.GetTypeFile(2, user.FileStoreId)
	case "3":
		Files = models.GetTypeFile(3, user.FileStoreId)
	case "4":
		Files = models.GetTypeFile(4, user.FileStoreId)
	}

	fileDetailUse := models.GetFileDetailUse(user.FileStoreId)

	c.JSON(http.StatusOK, gin.H{
		"currAll":       "active",
		"user":          user,
		"files":         Files,
		"fileDetailUse": fileDetailUse,
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

	json := structOfJson{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}
	folderName := json.FileFolderName
	parentId := json.ParentFolderId

	if !models.IsFolderNameOK(parentId, user.FileStoreId, folderName) {
		c.JSON(http.StatusBadRequest, utils.ErrorWrapper(fmt.Errorf("folder name exists!")))
		return
	}
	//新建文件夹数据
	models.CreateFolder(folderName, parentId, user.FileStoreId)
}

func DownloadFile(c *gin.Context) {
	json := structOfJson{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}
	fId := json.FileId

	file := models.GetFileInfo(fId)
	if file.FileHash == "" {
		return
	}

	//从oss获取文件
	fileData := lib.DownloadOss(file.FileHash, file.PostFix, true)
	//下载次数+1
	//models.DownloadNumAdd(fId)

	c.Header("Content-disposition", "attachment;filename=\""+file.FileName+file.PostFix+"\"")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Data(http.StatusOK, "application/octect-stream", fileData)
}

func UpdateFile(c *gin.Context) {
	json := structOfJson{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()})
	}
	fileId := json.FileId
	fileName := json.FileName
	models.UpdateUserFile(fileId, fileName)
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
	json := structOfJson{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()})
	}
	fId := json.FileId
	folderId := json.ParentFolderId

	//删除数据库文件数据
	models.DeleteUserFile(fId, folderId, user.FileStoreId)
}

//删除文件夹
func DeleteFileFolder(c *gin.Context) {
	json := structOfJson{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()})
	}
	folderId := json.FolderId
	//获取要删除的文件夹信息 取到父级目录重定向
	//folderInfo := models.GetCurrentFolder(folderId)

	//删除文件夹并删除文件夹中的文件信息
	models.DeleteFileFolder(folderId)
}

//修改文件夹名
func UpdateFileFolder(c *gin.Context) {
	json := structOfJson{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()})
	}
	fileFolderName := json.FileFolderName
	fileFolderId := json.FolderId
	folder := models.GetCurrentFolder(fileFolderId)
	if !models.IsFolderNameOK(folder.ParentFolderId, folder.FileStoreId, fileFolderName) {
		c.JSON(http.StatusBadRequest, utils.ErrorWrapper(fmt.Errorf("folder name exists!")))
		return
	}
	models.UpdateFolderName(fileFolderId, fileFolderName)
}

func MoveFile(c *gin.Context) {
	json := structOfJson{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()})
	}

	fileId := json.FileId
	foldId := json.FolderId

	models.MoveUserFile(fileId, foldId)
}

func MoveFolder(c *gin.Context) {
	json := structOfJson{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()})
	}

	toFolderId := json.ToFolderId
	folderId := json.FolderId
	folder := models.GetCurrentFolder(folderId)
	if !models.IsFolderNameOK(toFolderId, folder.FileStoreId, folder.FolderName) {
		c.JSON(http.StatusBadRequest, utils.ErrorWrapper(fmt.Errorf("目标文件夹存在同名文件夹！")))
		return
	}
	models.MoveFolder(folderId, toFolderId)
}
