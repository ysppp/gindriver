package models

import (
	"gindriver/utils"
	"strings"
	"time"
)

type Share struct {
	ShareId   uint64 `gorm:"column:shareid;type:numeric ;primary_key"`
	Code      string `gorm:"column:code;type:varchar(10)"`
	FileId    uint64 `gorm:"column:fileid;type:numeric"`
	UserName  string `gorm:"column:username;type:varchar(50)"`
	ShareHash string `gorm:"column:sharehash;type:varchar(255)"`
}

//创建分享
func CreateShare(code, username string, fId uint64) string {
	share := Share{
		Code:      strings.ToLower(code),
		FileId:    fId,
		UserName:  username,
		ShareHash: utils.EncodeMd5(code + string(time.Now().Unix())),
	}
	utils.Database.Create(&share)

	return share.ShareHash
}

//查询分享
func GetShareInfo(hash string) (share Share) {
	utils.Database.Find(&share, "hash = ?", hash)
	return
}

//校验提取码
func VerifyShareCode(fId uint64, code, sharehash string) bool {
	var share Share
	utils.Database.Find(&share, "fileid = ? and code = ? and sharehash = ?", fId, code, sharehash)
	if share.ShareId == 0 {
		return false
	}
	return true
}
