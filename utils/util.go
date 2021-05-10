package utils

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	//"crypto/sha256"
)

func GetFileTypeInt(filePrefix string) int {
	filePrefix = strings.ToLower(filePrefix)
	if filePrefix == ".doc" || filePrefix == ".docx" || filePrefix == ".txt" || filePrefix == ".pdf" {
		return 1
	}
	if filePrefix == ".jpg" || filePrefix == ".png" || filePrefix == ".gif" || filePrefix == ".jpeg" {
		return 2
	}
	if filePrefix == ".mp4" || filePrefix == ".avi" || filePrefix == ".mov" || filePrefix == ".rmvb" || filePrefix == ".rm" || filePrefix == ".flv" {
		return 3
	}
	if filePrefix == ".mp3" || filePrefix == ".cda" || filePrefix == ".wav" || filePrefix == ".wma" || filePrefix == ".ogg" || filePrefix == ".flac" {
		return 4
	}

	return 5
}

func GetRandomString(n int) string {
	randBytes := make([]byte, n/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

//SHA256生成哈希值
func GetSHA256HashCode(file *os.File) string {
	//创建一个基于SHA256算法的hash.Hash接口的对象
	hash := sha256.New()
	_, _ = io.Copy(hash, file)
	//计算哈希值
	bytes := hash.Sum(nil)
	//将字符串编码为16进制格式,返回字符串
	hashCode := hex.EncodeToString(bytes)
	//返回哈希值
	return hashCode

}

// md5加密
func EncodeMd5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func FilterUsername(str string) bool {
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9]{4,16}$", str); !ok {
		return false
	}
	return true
}

func SaveFile(filename string, i []byte) error {
	return ioutil.WriteFile(filename, i, 0644)
}

func ReadFile(filename string) []byte {
	if c, err := ioutil.ReadFile(filename); err == nil {
		return c
	}
	return nil
}
