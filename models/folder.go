package models

import (
	"gindriver/utils"
	"time"
)

type FileFolder struct {
	FolderId   uint64 `gorm:"column:folderid;type:numeric ;primary_key"`
	FolderName string `gorm:"column:foldername;type:varchar(255)"`

	ParentFolderId uint64 `gorm:"column:parentfolderid;type:numeric"` //	父文件夹ID
	FileStoreId    uint64 `gorm:"column:filestoreid;type:numeric"`    //	文件所属仓库ID
	Time           string `gorm:"column:time;type:timestamp"`
}

//新建文件夹
func CreateFolder(folderName string, parentId, fileStoreId uint64) {
	fileFolder := FileFolder{
		FolderId:       uint64(randomUint64()),
		FolderName:     folderName,
		ParentFolderId: parentId,
		FileStoreId:    fileStoreId,
		Time:           time.Now().Format("2006-01-02 15:04:05"),
	}
	utils.Database.Create(&fileFolder)
}

//获取父类的id
func GetParentFolder(fId uint64) (fileFolder FileFolder) {
	utils.Database.Find(&fileFolder, "FolderId=?", fId)
	return
}

//获取目录所有文件夹
func GetFileFolder(parentId, fileStoreId uint64) (fileFolders []FileFolder) {
	utils.Database.Order("time desc").Find(&fileFolders, "ParentFolderId=? and FileStoreId=?", parentId, fileStoreId)
	return
}

//获取当前的目录信息
func GetCurrentFolder(fId uint64) (fileFolder FileFolder) {
	utils.Database.Find(&fileFolder, "FolderId=?", fId)
	return
}

//获取当前路径所有的父级
func GetCurrentAllParent(folder FileFolder, folders []FileFolder) []FileFolder {
	var parentFolder FileFolder
	if folder.ParentFolderId != 0 {
		utils.Database.Find(&parentFolder, "FolderId=?", folder.ParentFolderId)
		folders = append(folders, parentFolder)
		//递归查找当前所有父级
		return GetCurrentAllParent(parentFolder, folders)
	}

	//反转切片
	for i, j := 0, len(folders)-1; i < j; i, j = i+1, j-1 {
		folders[i], folders[j] = folders[j], folders[i]
	}

	return folders
}

//获取用户文件夹数量
func GetUserFileFolderCount(fileStoreId uint64) (fileFolderCount int64) {
	var fileFolder []FileFolder
	utils.Database.Find(&fileFolder, "FileStoreId=?", fileStoreId).Count(&fileFolderCount)
	return
}

//删除文件夹信息
func DeleteFileFolder(fId uint64) bool {
	var fileFolder FileFolder
	var fileFolder2 FileFolder
	//删除文件夹信息
	utils.Database.Where("FolderId=?", fId).Delete(FileFolder{})
	//删除文件夹中文件信息
	utils.Database.Where("ParentFolderId=?", fId).Delete(File{})
	//删除文件夹中文件夹信息
	utils.Database.Find(&fileFolder, "ParentFolderId=?", fId)
	utils.Database.Where("ParentFolderId=?", fId).Delete(FileFolder{})

	utils.Database.Find(&fileFolder2, "ParentFolderId=?", fileFolder.FolderId)
	if fileFolder2.FolderId != 0 { //递归删除文件下的文件夹
		return DeleteFileFolder(fileFolder.FolderId)
	}

	return true
}

//修改文件夹名
func UpdateFolderName(fId uint64, fName string) {
	var fileFolder FileFolder
	utils.Database.Model(&fileFolder).Where("FolderId=?", fId).Update("FileFolderName", fName)
}
