package oss

import (
	"fmt"
	"gin-blog/pkg/app"
	"gin-blog/pkg/e"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"os"
	"sync"
	"time"
)

type Server struct {
}

var once sync.Once
var ossServerInstance *Server

func NewOssServer() *Server {
	once.Do(func() {
		ossServerInstance = &Server{}
	})
	return ossServerInstance
}

func (that *Server) Init(accessKeyId, accessKeySecret, bucketName, endPoint string) (client *oss.Client, bucket *oss.Bucket, err error) {
	client, err = oss.New(endPoint, accessKeyId, accessKeySecret)
	if err != nil {
		return
	}
	if bucketName != "" {
		bucket, err = client.Bucket(bucketName)
		if err != nil {
			return
		}
	}
	return
}

// UploadALiYunOss 上传至阿里云oss
func UploadALiYunOss(c *gin.Context, file string, fileName string, ext string, originPath string) interface{} {
	accessKeyId := os.Getenv("ACCESSKEYID")
	accessKeySecret := os.Getenv("ACCESSKEYIDSECRET")
	bucketName := os.Getenv("BUCKETNAME")
	endPoint := os.Getenv("ENDPOINT")
	_, bucket, err := NewOssServer().Init(accessKeyId, accessKeySecret, bucketName, endPoint)
	if err != nil {
		fmt.Println(err)
	}
	newOssFileName := time.Now().Format("200612")
	objectKey := file + "/" + newOssFileName + "/" + fileName + ext
	err = bucket.PutObjectFromFile(objectKey, originPath)
	if err != nil {
		fmt.Println(err.Error())
		app.FailWithMessage(e.GetMsg(e.UPLOAD_FAIL)+err.Error(), e.UPLOAD_FAIL, c)
		return nil
	}
	var ossUrl string
	ossUrl = "https://ctl-blog-1.oss-cn-beijing.aliyuncs.com/" + objectKey
	return ossUrl
}
