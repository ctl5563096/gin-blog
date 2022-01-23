package resource

import (
	"fmt"
	"gin-blog/models/resource"
	"gin-blog/pkg/app"
	"gin-blog/pkg/util"
	valid "gin-blog/vaild"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"strconv"
)

// CreateNewMusic 新建音乐
func CreateNewMusic(c *gin.Context) {
	var (
		r        resource.CreateMusicData
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

	id,_ := resource.CreateMusicRecord(&r)
	if id <= 0  {
		app.Fail(c)
		return
	}
	app.OkWithData(id,c)
	return
}

// GetMusicList 获取音乐列表
func GetMusicList(c *gin.Context)  {
	var page, _ = 	strconv.Atoi(c.DefaultQuery("page","1"))
	var pageSize, _ = 	strconv.Atoi(c.DefaultQuery("pageSize","10"))
	var keywords= 	c.DefaultQuery("keywords","")
	var order= 	c.DefaultQuery("order","desc")

	if !valid.Page(page,pageSize) {
		app.Fail(c)
		return
	}

	// 获取列表
	res,err := resource.GetMusicList(page,pageSize,keywords,order)
	// 获取总数
	count,_ := resource.GetMusicCount(keywords)
	dataReturn := make(map[string] interface{})
	if err != nil {
		app.Fail(c)
		return
	}
	dataReturn["list"] = res
	dataReturn["count"] = count
	app.OkWithData(dataReturn,c)
}

// DeleteMusicRecord 删除音乐记录
func DeleteMusicRecord(c *gin.Context)  {
	var id, _ = 	strconv.Atoi(c.DefaultQuery("id","0"))

	if id <= 0 {
		app.Fail(c)
		return
	}
	res := resource.DeleteMusicRecord(id)
	fmt.Println(res)
	if !res {
		app.Fail(c)
		return
	}
	app.OK(c)
	return
}

// GetMusicDetailBackend 后台获取音频详情
func GetMusicDetailBackend(c *gin.Context)  {
	var id, _ = 	strconv.Atoi(c.DefaultQuery("id","0"))
	if id <= 0 {
		app.Fail(c)
		return
	}
	// 获取详情
	result,err := resource.GetAudioDetail(id)
	if err != nil {
		app.Fail(c)
		return
	}
	app.OkWithData(result,c)
	return
}

// UpdateAudio 更新音频
func UpdateAudio(c *gin.Context)  {
	var id, _ = 	strconv.Atoi(c.DefaultQuery("id","0"))
	if id <= 0 {
		app.Fail(c)
		return
	}
	var (
		r        resource.CreateMusicData
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
	// 获取详情
	res := resource.UpdateAudioDetail(id,r)
	fmt.Println(res)
	fmt.Println(11111)
	if !res {
		app.Fail(c)
		return
	}
	app.OK(c)
	return
}

