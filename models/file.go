package models

import (
	"fmt"
	"gindriver/lib"
	"gindriver/utils"
	"path"
	"strconv"
	"strings"
	"time"
)

type File struct {
	FileId      uint64 `gorm:"column:fileid;type:numeric ;primary_key"`
	FileName    string `gorm:"column:filename;type:varchar(255)"`
	FileHash    string `gorm:"column:filehash;type:varchar(255)"`
	FileStoreId uint64 `gorm:"column:filestoreid;type:numeric"` //	文件所属仓库ID
	FilePath    string `gorm:"column:filepath;type:varchar(255)"`

	Size    int64  `gorm:"column:size;type:integer"`
	SizeStr string `gorm:"column:sizestr;type:varchar(50)"` //	文件大小单位
	Type    int    `gorm:"column:type;type:integer"`
	PostFix string `gorm:"column:postfix;type:varchar(255)"` //	文件后缀

	ParentFolderId uint64 `gorm:"column:parentfolderid;type:numeric"` //	父文件夹ID
	UploadTime     string `gorm:"column:uploadtime;type:timestamp"`
}

func CreateFile(fileName, filePath, fileHash string, fileSize int64, fId, fileStoreId uint64) {
	var sizeStr string

	fileId := uint64(randomUint64())
	// fileName may be like "/a/b/c.txt"
	// fileSuffix will be ".txt", filePrefix will be "/a/b/c"
	//获取文件后缀
	fileSuffix := path.Ext(filePath)
	//获取文件名
	//filePrefix := filePath[0 : len(filePath)-len(fileSuffix)]

	if fileSize < 1048576 {
		sizeStr = strconv.FormatInt(fileSize/1024, 10) + "KB"
	} else {
		sizeStr = strconv.FormatInt(fileSize/102400, 10) + "MB"
	}

	myFile := File{
		FileId:      fileId,
		FileName:    fileName,
		FileHash:    fileHash,
		FileStoreId: fileStoreId,
		FilePath:    filePath,
		//DownloadNum:    0,
		UploadTime:     time.Now().Format("2006-01-02 15:04:05"),
		ParentFolderId: fId,
		Size:           fileSize / 1024,
		SizeStr:        sizeStr,
		Type:           utils.GetFileTypeInt(fileSuffix),
		PostFix:        strings.ToLower(fileSuffix),
	}

	utils.Database.Create(&myFile)
}

//获取用户的文件
func GetUserFile(parentId, storeId uint64) (files []File) {
	utils.Database.Find(&files, "FileStoreId = ? and ParentFolderId = ?", storeId, parentId)
	return
}

//文件上传成功减去相应容量
func SubtractSize(size int64, fileStoreId uint64) {
	var fileStore FileStore
	utils.Database.First(&fileStore, fileStoreId)

	fileStore.CurrentSize = fileStore.CurrentSize + size/1024
	fileStore.MaxSize = fileStore.MaxSize - size/1024
	utils.Database.Save(&fileStore)
}

//获取用户文件数量
func GetUserFileCount(fileStoreId uint64) (fileCount int64) {
	var file []File
	utils.Database.Find(&file, "FileStoreId = ?", fileStoreId).Count(&fileCount)
	return
}

//获取用户文件使用明细情况
func GetFileDetailUse(fileStoreId uint64) map[string]int64 {
	var files []File
	var (
		docCount   int64
		imgCount   int64
		videoCount int64
		musicCount int64
		otherCount int64
	)

	fileDetailUseMap := make(map[string]int64, 0)

	//文档类型
	docCount = utils.Database.Find(&files, "FileStoreId = ? AND Type = ?", fileStoreId, 1).RowsAffected
	fileDetailUseMap["docCount"] = docCount
	////图片类型
	imgCount = utils.Database.Find(&files, "FileStoreId = ? and Type = ?", fileStoreId, 2).RowsAffected
	fileDetailUseMap["imgCount"] = imgCount
	//视频类型
	videoCount = utils.Database.Find(&files, "FileStoreId = ? and Type = ?", fileStoreId, 3).RowsAffected
	fileDetailUseMap["videoCount"] = videoCount
	//音乐类型
	musicCount = utils.Database.Find(&files, "FileStoreId = ? and Type = ?", fileStoreId, 4).RowsAffected
	fileDetailUseMap["musicCount"] = musicCount
	//其他类型
	otherCount = utils.Database.Find(&files, "FileStoreId = ? and Type = ?", fileStoreId, 5).RowsAffected
	fileDetailUseMap["otherCount"] = otherCount

	return fileDetailUseMap
}

//根据文件类型获取文件
func GetTypeFile(fileType, fileStoreId uint64) (files []File) {
	utils.Database.Find(&files, "FileStoreId=? and Type=?", fileStoreId, fileType)
	return
}

//判断当前文件夹是否有同名文件
func CurrFileExists(folderId uint64, filename string) bool {
	var file File
	//获取文件后缀
	fileSuffix := strings.ToLower(path.Ext(filename))
	//获取文件名
	filePrefix := filename[0 : len(filename)-len(fileSuffix)]

	utils.Database.Find(&file, "ParentFolderId=? and FileName=? and PostFix=?", folderId, filePrefix, fileSuffix)

	if file.Size > 0 {
		return false
	}
	return true
}

//通过hash判断文件是否已上传过oss
func FileOssExists(fileHash string, folderId uint64) bool {
	var file File
	utils.Database.Find(&file, "FileHash = ? and parentFolderId = ?", fileHash, folderId)
	if file.FileId != 0 {
		return false
	}
	return true
}

//通过fileId获取文件信息
func GetFileInfo(fId uint64) (file File) {
	utils.Database.First(&file, fId)
	return
}

//删除数据库文件数据
func DeleteUserFile(fId, folderId, storeId uint64) {
	fmt.Printf("FileId = %d and FileStoreId = %d and ParentFolderId = %d", fId, storeId, folderId)
	file := GetFileInfo(fId)
	lib.DeleteOss(file.FileHash, file.PostFix)
	utils.Database.Where("FileId = ? and FileStoreId = ? and ParentFolderId = ?", fId, storeId, folderId).Delete(File{})
}

func UpdateUserFile(fileId uint64, fileName string) {
	var File File
	utils.Database.Model(&File).Where("FileId=?", fileId).Update("FileName", fileName)
}

func MoveUserFile(fileId, folderId uint64) {
	var File File
	utils.Database.Model(&File).Where("FileId=?", fileId).Update("ParentFolderId", folderId)
}
