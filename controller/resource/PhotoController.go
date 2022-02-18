package resource

import (
	"gin-blog/models/resource"
	"gin-blog/pkg/app"
	"gin-blog/pkg/e"
	"gin-blog/pkg/util"
	valid "gin-blog/vaild"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"strconv"
)

type PhotosList struct {
	NewPhotosArr []resource.Photo `json:"changePhotoArr"`
}
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
	// 这里处理系列图片
	res := SetNewsPhotosList(c,id)
	if !res {
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
	// 获取关联的系列图片
	photoList := resource.GetAboutPhotos(id)
	if err != nil {
		app.Fail(c)
		return
	}
	dataReturn := make(map[string] interface{})
	dataReturn["changePhotoArr"] = photoList
	dataReturn["info"] = result
	app.OkWithData(dataReturn,c)
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
		list     PhotosList
		isChange int
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
	isChange,_ = strconv.Atoi(c.DefaultQuery("isChange","2"))
	// 有变动就需要对比然后插入
	if isChange == 1 {
		err = c.ShouldBindBodyWith(&list, binding.JSON)
		// 这里去更新系列图片
		UpdatePhotoList(list,id)
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


//GetPhotoDetailFront 前端获取详情
func GetPhotoDetailFront(c *gin.Context)  {
	var id, _ = 	strconv.Atoi(c.DefaultQuery("id","0"))
	if id <= 0 {
		app.Fail(c)
		return
	}
	// 获取详情
	result,err := resource.GetPhotoDetail(id)
	// 获取关联的系列图片
	photoList := resource.GetAboutPhotos(id)
	dataReturn := make(map[string] interface{})
	dataReturn["changePhotoArr"] = photoList
	dataReturn["info"] = result
	if err != nil {
		app.Fail(c)
		return
	}
	app.OkWithData(dataReturn,c)
	return
}

//SetNewsPhotosList 保存系列图片
func SetNewsPhotosList(c *gin.Context, id int) bool {
	var r resource.PhotosList
	err := c.ShouldBindBodyWith(&r,binding.JSON)
	if err != nil {
		app.Fail(c)
		return false
	}
	if id <= 0 {
		app.FailWithMessage(e.GetMsg(e.MISS_PARAMS),e.MISS_PARAMS,c)
		return false
	}
	var  newArr []interface{}
	for _,value := range r.PhotosArr{
		value.ResourceId = id
		newArr = append(newArr,value)
	}
	res := resource.BatchInsertPhoto(newArr)
	return  res
}

// UpdatePhotoList 更新关联图片
func UpdatePhotoList(list PhotosList, id int)  {
	// 先去获取旧的
	photoList := resource.GetAboutPhotos(id)

	var  oldArr []interface{}
	var  newArr []interface{}
	var  newInsert []interface{}
	var  delArr []interface{}
	for _,value := range photoList{
		oldArr = append(oldArr,value.Id)
	}

	for _,value := range list.NewPhotosArr {
		if value.Id > 0 {
			newArr = append(newArr,value.Id)
		} else {
			value.ResourceId = id
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
