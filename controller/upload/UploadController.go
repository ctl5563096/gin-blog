package upload

import (
	"fmt"
	"gin-blog/pkg/app"
	"gin-blog/pkg/e"
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

const PATH = "resource/public/pic/avatar"

// Avatar 上传用户头像
func Avatar(c *gin.Context)  {
	f, err := c.FormFile("avatar")
	if err != nil {
		app.Fail(c)
		return
	}else {
		ext :=strings.ToLower(path.Ext(f.Filename))

		// 限定上传文件类型
		if ext!=".png"&&ext!=".jpg"&&ext!=".gif"&&ext!=".jpeg"{
			app.FailWithMessage(e.GetMsg(e.IMAGE_TYPE_ERROR),e.IMAGE_TYPE_ERROR,c)
			return
		}

		fileName,result := GetAvatarName(ext)
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
func GetAvatarName(ext string) (path string,result bool) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	datetime := time.Now().Format("200612")

	// 上传路径加上时间日期
	dstPath := PATH + DIR + datetime

	// 文件的相对路径
	dstFile := dstPath + DIR + datetime + strconv.Itoa(random.Int()) +  strings.Replace(ext, DIR, POS, 1)
	fmt.Println(dstFile)

	// 检测下这个文件夹是否存在 如果没有就新建这个文件夹
	_, err := os.Stat(dstPath)
	res := os.IsNotExist(err)
	if res == true {
		err := os.MkdirAll(dstPath, os.ModePerm)
		if err != nil {
			return "",false
		}
	}
	return dstFile,true
}