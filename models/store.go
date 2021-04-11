package models

import (
	"gindriver/utils"
)

type FileStore struct {
	FileStoreId uint64 `gorm:"column:filestoreid;type:numeric ;primary_key"` //	文件所属仓库ID
	Id          uint64 `gorm:"column:id;type:numeric"`                       // user id
	CurrentSize int64  `gorm:"column:currentsize;type:integer"`              // 当前容量，单位（KB）
	MaxSize     int64  `gorm:"column:maxsize;type:integer"`                  // 最大容量，单位（KB）
}

func (s FileStore) Insert() (fileStore *FileStore, err error) {
	result := utils.Database.Create(&s)

	fileStore = &s

	if result.Error != nil {
		err = result.Error
		return
	}
	return
}

func NewFileStore(Id uint64, MaxSize int64) *FileStore {
	fileStore := &FileStore{}
	fileStore.FileStoreId = uint64(randomUint64())
	fileStore.Id = Id
	fileStore.MaxSize = MaxSize
	fileStore.CurrentSize = 0
	return fileStore
}

//根据用户id获取仓库信息
func GetUserFileStore(userId uint64) (fileStore FileStore) {
	utils.Database.Find(&fileStore, "id = ?", userId)
	return
}

//判断用户容量是否足够
func CapacityIsEnough(fileSize int64, fileStoreId uint64) bool {
	var fileStore FileStore
	utils.Database.First(&fileStore, fileStoreId)
	if fileStore.MaxSize-(fileSize/1024) < 0 {
		return false
	}

	return true
}
