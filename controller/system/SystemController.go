package system

import (
	"gin-blog/models/system"
	"gin-blog/pkg/app"
	"gin-blog/pkg/e"
	"gin-blog/pkg/util"
	valid "gin-blog/vaild"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"math"
	"strconv"
)

// Create 创建系统参数
func Create(c *gin.Context)  {
	var r system.Params
	var errStr string
	var errorMap map[string][]string
	err := c.ShouldBind(&r)
	validate := validator.New()
	err = validate.Struct(r)
	if err != nil {
		switch err.(type) {
		case validator.ValidationErrors:
			errorMap = valid.Translate(err)
			//循环遍历Map 只返回第一个错误信息
			for _,v:= range errorMap{
				for _,z := range v{
					util.WriteLog("system_business_error",4,z)
					app.FailWithMessage(z,4,c)
					return
				}
			}
		default:
			errStr = "未知错误"
		}
		app.FailWithMessage(errStr,1,c)
		return
	}

	system.CreateRecord(&r)
	app.OK(c)
	return
}

// GetList 获取参数列表
func GetList(c *gin.Context)  {
	var page, _ = 	strconv.Atoi(c.DefaultQuery("page","1"))
	var pageSize, _ = 	strconv.Atoi(c.DefaultQuery("pageSize","10"))
	var keywords= 	c.DefaultQuery("keywords","")
	var order= 	c.DefaultQuery("order","desc")
	dataReturn := make(map[string] interface{})

	// 求列表
	list,_  := system.GetList(page,pageSize,keywords,order)
	// 求总数
	count := system.GetTotal(keywords)

	dataReturn["count"] = count
	dataReturn["list"] = list
	for key,value := range list{
		// 判断下是什么排序再去赋值
		if order == "desc" {
			var pageAll int
			var currentSort int
			pageAll = int(math.Ceil(float64(count) / float64(pageSize)))
			// 如果大于或者总页数就取最后一页
			if  page >= pageAll {
				currentSort = count % pageSize - key
			}else {
				currentSort = (pageAll - page) * pageSize + count % pageSize - key
			}
			value.Sort = currentSort
		}else{
			value.Sort = pageSize * (page - 1) + key + 1
		}
	}
	app.OkWithData(dataReturn,c)
	return
}

// DelParam 删除指定参数
func DelParam(c *gin.Context)  {
	var id, _ = strconv.Atoi(c.DefaultQuery("id","0"))
	if id <= 0 {
		app.FailWithMessage(e.GetMsg(e.PARAMS_ERROR),e.PARAMS_ERROR,c)
		return
	}
	err := system.DeleteParam(id)
	if err != nil {
		app.Fail(c)
		return
	}
	app.OK(c)
	return
}

// EditParam 更新参数
func EditParam(c *gin.Context) {
	var v system.EditStruct
	var errStr string
	var errorMap map[string][]string
	err := c.ShouldBind(&v)
	validate := validator.New()
	err = validate.Struct(v)
	if err != nil {
		switch err.(type) {
		case validator.ValidationErrors:
			errorMap = valid.Translate(err)
			//循环遍历Map 只返回第一个错误信息
			for _, v := range errorMap {
				for _, z := range v {
					util.WriteLog("rule_business_error", 4, z)
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

	result := make(map[string]interface{})
	res := system.UpdateDetail(&v)
	result["res"] = res
	app.OkWithData(result, c)
	return
}

// GetDetail 获取参数详情
func GetDetail(c *gin.Context)  {
	var id, _ = strconv.Atoi(c.DefaultQuery("id","0"))
	if id <= 0 {
		app.FailWithMessage(e.GetMsg(e.PARAMS_ERROR),e.PARAMS_ERROR,c)
		return
	}
	res := system.GetParamDetail(id)
	app.OkWithData(res,c)
	return
}

func GetTags(c *gin.Context)  {
	var code = c.DefaultQuery("code","")
	if code == "" {
		app.FailWithMessage(e.GetMsg(e.PARAMS_ERROR),e.PARAMS_ERROR,c)
		return
	}
	res := system.GetTagsTypeList(code)
	app.OkWithData(res,c)
	return
}