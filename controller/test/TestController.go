package test

import (
	"fmt"
	"gin-blog/controller/upload"
	"gin-blog/models/resource"
	"gin-blog/pkg/app"
	"gin-blog/pkg/e"
	"gin-blog/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"path"
	"strings"
)

// TestPort 测试接口
func TestPort(c *gin.Context)  {
	userInfo := util.GetUserInfo(c)
	fmt.Println(userInfo.(map[string]interface{})["is_black"])
	data := make(map[string] interface{})
	app.OkWithCodeData("测试成功",data,0,c)
	return
}

// TestThumb 测试缩略图效果
func TestThumb(c *gin.Context)  {
	name,_   := c.GetQuery("type")
	f, err   := c.FormFile(name)
	if err != nil {
		app.Fail(c)
		return
	}
	ext :=strings.ToLower(path.Ext(f.Filename))
	if ext!=".png"&&ext!=".jpg"&&ext!=".gif"&&ext!=".jpeg"{
		app.FailWithMessage(e.GetMsg(e.IMAGE_TYPE_ERROR),e.IMAGE_TYPE_ERROR,c)
		return
	}
	originPath := "resource/public/pic/" + name
	fileName,_,result := upload.GetNewFileName(ext,originPath)
	newOriginPath := "resource/public/pic/thumb"
	newFileName,_,result := upload.GetNewFileName(ext,newOriginPath)
	if result != true {
		app.Fail(c)
		return
	}
	err = c.SaveUploadedFile(f,fileName)
	_, err = util.MakeThumbnail(fileName,newFileName)
	if err != nil {
		return
	}
	app.OK(c)
	return
}

// TestPortSecond 测试缩略图效果
func TestPortSecond(c *gin.Context)  {
	photoList := resource.GetAboutPhotos(42)
	var r resource.PhotosList
	err := c.ShouldBindBodyWith(&r, binding.JSON)
	if err != nil {
		return 
	}
	var  oldArr []interface{}
	var  newArr []interface{}
	var  newInsert []interface{}
	var  delArr []interface{}
	for _,value := range photoList{
		oldArr = append(oldArr,value.Id)
	}

	for _,value := range r.PhotosArr{
		if value.Id > 0 {
			newArr = append(newArr,value.Id)
		} else {
			newInsert = append(newInsert,value)
		}
	}
	// 求出需要删除的图片
	for _,value := range oldArr{
		if !util.InArray(newArr,value) {
			delArr = append(delArr,value)
		}
	}
	// 执行删除操作
	if len(delArr) > 0 {
		resource.DelPhotosList(delArr)
	}

	// 执行新增操作
	if len(newInsert) > 0 {
		resource.BatchInsertPhoto(newInsert)
	}
}