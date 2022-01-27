package resource

import (
	"gin-blog/models/resource"
	"gin-blog/pkg/app"
	"gin-blog/pkg/util"
	valid "gin-blog/vaild"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"strconv"
)

// CreatePhoto 创建图片资源
func CreatePhoto(c *gin.Context)  {
	var (
		r        resource.PhotosData
		errStr   string
		errorMap map[string][]string
	)
	err := c.ShouldBindBodyWith(&r, binding.JSON)
	validate := validator.New()
	err = validate.Struct(r)
	if err != nil {
		switch err.(type) {
		case validator.ValidationErrors:
			errorMap = valid.Translate(err)
			//循环遍历Map 只返回第一个错误信息
			for _, v := range errorMap {
				for _, z := range v {
					util.WriteLog("user_business_error", 4, z)
					app.FailWithMessage(z, 4, c)
					return
				}
			}
		default:
			errStr = "未知错误"
		}
		app.FailWithMessage(errStr, 1, c)
		return
	}

	id,_ := resource.CreatPhotoRecord(&r)
	if id <= 0  {
		app.Fail(c)
		return
	}
	app.OkWithData(id,c)
	return
}

// GetPhotoList 获取图片列表
func GetPhotoList(c *gin.Context)  {
	var page, _ = 	strconv.Atoi(c.DefaultQuery("page","1"))
	var pageSize, _ = 	strconv.Atoi(c.DefaultQuery("pageSize","10"))
	var keywords= 	c.DefaultQuery("keywords","")
	var order= 	c.DefaultQuery("order","desc")

	if !valid.Page(page,pageSize) {
		app.Fail(c)
		return
	}

	// 获取列表
	res,err := resource.GetPhotoList(page,pageSize,keywords,order)
	// 获取总数
	count,_ := resource.GetPhotoCount(keywords)
	if err != nil {
		app.Fail(c)
		return
	}
	dataReturn := make(map[string] interface{})
	dataReturn["list"] = res
	dataReturn["count"] = count
	app.OkWithData(dataReturn,c)
}

// DeletePhotoRecord 删除图片记录
func DeletePhotoRecord(c *gin.Context)  {
	var id, _ = 	strconv.Atoi(c.DefaultQuery("id","0"))

	if id <= 0 {
		app.Fail(c)
		return
	}
	res := resource.DeletePhotoRecord(id)
	if !res {
		app.Fail(c)
		return
	}
	app.OK(c)
	return
}

// GetPhotoDetailBackend 获取图片资源详情
func GetPhotoDetailBackend(c *gin.Context)  {
	var id, _ = 	strconv.Atoi(c.DefaultQuery("id","0"))
	if id <= 0 {
		app.Fail(c)
		return
	}
	// 获取详情
	result,err := resource.GetPhotoDetail(id)
	if err != nil {
		app.Fail(c)
		return
	}
	app.OkWithData(result,c)
	return
}

// UpdatePhotosResource 更新图片资源
func UpdatePhotosResource(c *gin.Context)  {
	var id, _ = 	strconv.Atoi(c.DefaultQuery("id","0"))
	if id <= 0 {
		app.Fail(c)
		return
	}
	var (
		r        resource.PhotosData
		errStr   string
		errorMap map[string][]string
	)
	err := c.ShouldBindBodyWith(&r, binding.JSON)
	validate := validator.New()
	err = validate.Struct(r)
	if err != nil {
		switch err.(type) {
		case validator.ValidationErrors:
			errorMap = valid.Translate(err)
			//循环遍历Map 只返回第一个错误信息
			for _, v := range errorMap {
				for _, z := range v {
					util.WriteLog("user_business_error", 4, z)
					app.FailWithMessage(z, 4, c)
					return
				}
			}
		default:
			errStr = "未知错误"
		}
		app.FailWithMessage(errStr, 1, c)
		return
	}
	// 更新图片详情
	res := resource.UpdatePhotoDetail(id,r)
	if !res {
		app.Fail(c)
		return
	}
	app.OK(c)
	return
}