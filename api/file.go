package api

import (
	"fmt"
	"gindriver/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

//type MyFile struct {
//	Id             int
//	FileName       string //文件名
//	FileHash       string //文件哈希值
//	FileStoreId    int    //文件仓库id
//	FilePath       string //文件存储路径
//	DownloadNum    int    //下载次数
//	UploadTime     string //上传时间
//	ParentFolderId int    //父文件夹id
//	Size           int64  //文件大小
//	SizeStr        string //文件大小单位
//	Type           int    //文件类型
//	Postfix        string //文件后缀
//}
//
////添加文件数据
//func CreateFile(filename, fileHash string, fileSize int64, fId string, fileStoreId int) {
//	var sizeStr string
//	//获取文件后缀
//	fileSuffix := path.Ext(filename)
//	//获取文件名
//	filePrefix := filename[0 : len(filename)-len(fileSuffix)]
//	fid, _ := strconv.Atoi(fId)
//
//	if fileSize < 1048576 {
//		sizeStr = strconv.FormatInt(fileSize/1024, 10) + "KB"
//	} else {
//		sizeStr = strconv.FormatInt(fileSize/102400, 10) + "MB"
//	}
//
//	myFile := MyFile{
//		FileName:       filePrefix,
//		FileHash:       fileHash,
//		FileStoreId:    fileStoreId,
//		FilePath:       "",
//		DownloadNum:    0,
//		UploadTime:     time.Now().Format("2006-01-02 15:04:05"),
//		ParentFolderId: fid,
//		Size:           fileSize / 1024,
//		SizeStr:        sizeStr,
//		Type:           util.GetFileTypeInt(fileSuffix),
//		Postfix:        strings.ToLower(fileSuffix),
//	}
//	sqlDB.DB.Create(&myFile)
//}

func UploadHandler(c *gin.Context) {
	user, ret := c.Get("SessionUser")
	if !ret {
		c.JSON(http.StatusUnauthorized, utils.ErrorWrapper(fmt.Errorf("not authorized")))
		return
	}
	if user != "admin" {
		c.JSON(http.StatusForbidden, utils.ErrorWrapper(fmt.Errorf("forbidden")))
		return
	}

	file, err := c.FormFile("files")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorWrapper(err))
		return
	}
	err = c.SaveUploadedFile(file, fmt.Sprintf("./public/%s", file.Filename))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorWrapper(err))
		return
	}
	c.JSON(http.StatusCreated, utils.SuccessWrapper(fmt.Sprintf("file saved at ./public/%s", file.Filename)))
}
