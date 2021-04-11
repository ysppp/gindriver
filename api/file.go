package api

import (
	"fmt"
	"gindriver/config"
	"gindriver/lib"
	"gindriver/models"
	"gindriver/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
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

//func UploadHandler(c *gin.Context) {
//	_ , ret := c.Get("SessionUser")
//	if !ret {
//		c.JSON(http.StatusUnauthorized, utils.ErrorWrapper(fmt.Errorf("not authorized")))
//		return
//	}
//	//if user != "admin" {
//	//	c.JSON(http.StatusForbidden, utils.ErrorWrapper(fmt.Errorf("forbidden")))
//	//	return
//	//}
//	file, err := c.FormFile("files")
//	if err != nil {
//		c.JSON(http.StatusBadRequest, utils.ErrorWrapper(err))
//		return
//	}
//	err = c.SaveUploadedFile(file, fmt.Sprintf("./public/%s", file.Filename))
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, utils.ErrorWrapper(err))
//		return
//	}
//	c.JSON(http.StatusCreated, utils.SuccessWrapper(fmt.Sprintf("file saved at ./public/%s", file.Filename)))
//}

func Upload(c *gin.Context) {
	username, ret := c.Get("SessionUser")
	if !ret {
		c.JSON(http.StatusUnauthorized, utils.ErrorWrapper(fmt.Errorf("not authorized")))
		return
	}
	//获取用户信息
	user := models.GetUserInfoByName(username)

	fId, _ := strconv.ParseUint(c.DefaultQuery("fId", "0"), 10, 64)
	//获取当前目录信息
	currentFolder := models.GetCurrentFolder(fId)
	//获取当前目录所有的文件夹信息
	fileFolders := models.GetFileFolder(fId, user.FileStoreId)
	//获取父级的文件夹信息
	parentFolder := models.GetParentFolder(fId)
	//获取当前目录所有父级
	currentAllParent := models.GetCurrentAllParent(parentFolder, make([]models.FileFolder, 0))
	//获取用户文件使用明细数量
	fileDetailUse := models.GetFileDetailUse(user.FileStoreId)

	c.JSON(http.StatusOK, gin.H{
		"user":             user,
		"currUpload":       "active",
		"fId":              currentFolder.FolderId,
		"fName":            currentFolder.FolderName,
		"fileFolders":      fileFolders,
		"parentFolder":     parentFolder,
		"currentAllParent": currentAllParent,
		"fileDetailUse":    fileDetailUse,
	})
}

//处理上传文件
func UploadHandler(c *gin.Context) {
	username, ret := c.Get("SessionUser")
	if !ret {
		c.JSON(http.StatusUnauthorized, utils.ErrorWrapper(fmt.Errorf("not authorized")))
		return
	}
	//获取用户信息
	user := models.GetUserInfoByName(username)

	Fid, _ := strconv.ParseUint(c.GetHeader("id"), 10, 64)
	//conf := lib.LoadServerConfig()
	//接收上传文件
	file, err := c.FormFile("files")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorWrapper(err))
		fmt.Printf("Error: user:%s, Fid: %s, err: %s", user.Name, Fid, err)
		return
	}
	//判断当前文件夹是否有同名文件
	if ok := models.CurrFileExists(Fid, file.Filename); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": 501,
		})
		return
	}

	//判断用户的容量是否足够
	if ok := models.CapacityIsEnough(file.Size, user.FileStoreId); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": 503,
		})
		return
	}

	if err != nil {
		fmt.Println("文件上传错误", err.Error())
		return
	}
	//defer file.Close()

	//文件保存本地的路径
	location := config.Config.UploadLocation + file.Filename

	//newFile, err := os.Create(location)
	//在本地创建一个新的文件
	err = c.SaveUploadedFile(file, location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorWrapper(err))
		fmt.Println("文件创建失败", err.Error())
		return
	}

	//将上传文件拷贝至新创建的文件中
	//fileSize, err := io.Copy(newFile, file)
	//if err != nil {
	//	fmt.Println("文件拷贝错误", err.Error())
	//	return
	//}

	//将光标移至开头
	//_, _ = newFile.Seek(0, 0)
	f, _ := os.Open(location)
	fileHash := utils.GetSHA256HashCode(f)

	defer f.Close()

	//通过hash判断文件是否已上传过oss
	if ok := models.FileOssExists(fileHash); ok {
		//上传至阿里云oss
		go lib.UploadOss(f.Name(), fileHash)
		//新建文件信息
		models.CreateFile(f.Name(), fileHash, file.Size, Fid, user.FileStoreId)
		//上传成功减去相应剩余容量
		models.SubtractSize(file.Size/1024, user.FileStoreId)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
	})
}
