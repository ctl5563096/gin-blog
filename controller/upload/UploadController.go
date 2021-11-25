package upload

import (
	"gin-blog/models/oss"
	"gin-blog/pkg/app"
	"gin-blog/pkg/e"
	uploadOss "gin-blog/pkg/oss"
	"github.com/gin-gonic/gin"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const POS = "."

const DIR = "/"

// Avatar 上传用户头像
func Avatar(c *gin.Context)  {
	f, err := c.FormFile("avatar")
	if err != nil {
		app.FailWithMessage(e.GetMsg(e.MISS_PARAMS),e.MISS_PARAMS,c)
		return
	}else {
		ext :=strings.ToLower(path.Ext(f.Filename))

		// 限定上传文件类型
		if ext!=".png"&&ext!=".jpg"&&ext!=".gif"&&ext!=".jpeg"{
			app.FailWithMessage(e.GetMsg(e.IMAGE_TYPE_ERROR),e.IMAGE_TYPE_ERROR,c)
			return
		}

		fileName,_,result := GetAvatarName(ext,"resource/public/pic/avatar")
		if result != true {
			app.Fail(c)
			return
		}
		err = c.SaveUploadedFile(f,fileName)
		if err !=  nil {
			app.FailWithMessage(e.GetMsg(e.UPLOAD_FAIL),e.UPLOAD_FAIL,c)
			return
		}

		data 		 := make(map[string]interface{})
		data["path"]  = fileName

		app.OkWithData(data,c)
		return
	}
}

// GetAvatarName 获取头像名称
func GetAvatarName(ext string ,pathOrigin string) (path string,randomName string,result bool) {
	random := strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Int())
	datetime := time.Now().Format("200612")

	// 上传路径加上时间日期
	dstPath := pathOrigin + DIR + datetime

	// 文件的相对路径
	dstFile := dstPath + DIR + datetime + random +  strings.Replace(ext, DIR, POS, 1)

	// 检测下这个文件夹是否存在 如果没有就新建这个文件夹
	_, err := os.Stat(dstPath)
	res := os.IsNotExist(err)
	if res == true {
		err := os.MkdirAll(dstPath, os.ModePerm)
		if err != nil {
			return "","",false
		}
	}
	return dstFile, random,true
}

// Photo 上传图片
func Photo(c *gin.Context)  {
	name,_   := c.GetQuery("type")
	f, err := c.FormFile(name)
	if err != nil {
		app.FailWithMessage(e.GetMsg(e.MISS_PARAMS),e.MISS_PARAMS,c)
		return
	}else {
		ext :=strings.ToLower(path.Ext(f.Filename))

		// 限定上传文件类型
		if ext!=".png"&&ext!=".jpg"&&ext!=".gif"&&ext!=".jpeg"{
			app.FailWithMessage(e.GetMsg(e.IMAGE_TYPE_ERROR),e.IMAGE_TYPE_ERROR,c)
			return
		}
		originPath := "resource/public/pic/" + name
		fileName,randomName,result := GetAvatarName(ext,originPath)
		if result != true {
			app.Fail(c)
			return
		}
		err = c.SaveUploadedFile(f,fileName)
		ossUrl := uploadOss.UploadALiYunOss(c,name,randomName,ext,fileName)
		// 这里把oss和本地的url存起来
		oss.SaveRelation(ossUrl.(string),originPath)
		if err !=  nil {
			app.FailWithMessage(e.GetMsg(e.UPLOAD_FAIL),e.UPLOAD_FAIL,c)
			return
		}

		data 		 := make(map[string]interface{})
		data["path"]  = ossUrl

		app.OkWithData(data,c)
		return
	}
}