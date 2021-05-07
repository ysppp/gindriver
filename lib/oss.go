package lib

import (
	"encoding/base64"
	"fmt"
	"gindriver/config"
	"gindriver/utils"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io/ioutil"
	"os"
	"path"
)

//上传文件至阿里云
func UploadOss(filename, fileHash string) {
	//获取文件后缀
	fileSuffix := path.Ext(filename)
	//conf := LoadServerConfig()
	// 创建OSSClient实例。
	//client, err := oss.New(config.Config.Oss.EndPoint, config.Config.Oss.AccessKeyId, config.Config.Oss.AccessKeySecret)
	client, err := oss.New(config.Config.Oss.EndPoint, config.Config.Oss.AccessKeyId, config.Config.Oss.AccessKeySecret)
	if err != nil {
		fmt.Println("创建实例Error:", err)
		return
	}

	// 获取存储空间。
	bucket, err := client.Bucket("filefree")
	if err != nil {
		fmt.Println("获取存储空间Error:", err)
		return
	}

	// 上传本地文件。
	err = bucket.PutObjectFromFile("files/"+fileHash+fileSuffix, filename)
	if err != nil {
		fmt.Println("本地文件上传Error:", err)
		return
	}
}

//从oss下载文件
func DownloadOss(fileHash, fileType string, isGetOriginPic bool) []byte {
	//conf := LoadServerConfig()
	// 创建OSSClient实例。
	client, err := oss.New(config.Config.Oss.EndPoint, config.Config.Oss.AccessKeyId, config.Config.Oss.AccessKeySecret)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// 获取存储空间。
	bucket, err := client.Bucket(config.Config.Oss.BucketName)
	if err != nil {
		fmt.Println("Error:", err)
	}

	bucketName := "files/"
	b := utils.GetFileTypeInt(fileType)
	if b == 2 && !isGetOriginPic {
		bucketName = "heic/"
	}

	// 下载文件到流。
	body, err := bucket.GetObject(bucketName + fileHash + fileType)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// 数据读取完成后，获取的流必须关闭，否则会造成连接泄漏，导致请求无连接可用，程序无法正常工作。
	defer body.Close()

	data, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return data
}

//从oss删除文件
func DeleteOss(fileHash, fileType string) {
	//conf := LoadServerConfig()
	// 创建OSSClient实例。
	client, err := oss.New(config.Config.Oss.EndPoint, config.Config.Oss.AccessKeyId, config.Config.Oss.AccessKeySecret)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 获取存储空间。
	bucket, err := client.Bucket(config.Config.Oss.BucketName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = bucket.DeleteObject("files/" + fileHash + fileType)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = bucket.DeleteObject("heic/" + fileHash + fileType)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func ProcessHeic(fileHash, fileSuffix string) {
	// 创建OSSClient实例。
	client, err := oss.New(config.Config.Oss.EndPoint, config.Config.Oss.AccessKeyId, config.Config.Oss.AccessKeySecret)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 指定原图所在Bucket。
	bucketName := "filefree"
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	// 原图名称。若图片不在Bucket根目录，需携带文件访问路径，例如example/example.jpg。
	sourceImageName := "files/" + fileHash + fileSuffix
	// 指定处理后图片存放的Bucket，该Bucket需与源Bucket在相同地域。
	//targetBucketName := "heic"
	// 指定处理后的图片名称。
	targetImageName := "heic/" + fileHash + fileSuffix
	// 将图片缩放为固定宽高100 px后转存到指定存储空间。
	style := "image/format,webp"
	process := fmt.Sprintf("%s|sys/saveas,o_%v", style, base64.URLEncoding.EncodeToString([]byte(targetImageName)))
	result, err := bucket.ProcessObject(sourceImageName, process)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	} else {
		fmt.Println(result)
	}
}
