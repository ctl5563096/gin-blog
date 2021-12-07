package resource

import (
	_interface "gin-blog/interface"
	"gin-blog/models/resource"
	"gin-blog/models/system"
	"gin-blog/pkg/util"

	//"gin-blog/models/system"

	//"gin-blog/models/system"
	"gin-blog/pkg/app"
	"gin-blog/pkg/e"
	//"gin-blog/pkg/util"
	valid "gin-blog/vaild"
	"github.com/gin-gonic/gin"
	"math"
	"strconv"
)

// CreateHotRecord  新建热门资源
func CreateHotRecord(c *gin.Context)  {
	var r resource.CreateData
	err := c.ShouldBind(&r)
	err = valid.RequestData(&r)
	if err != nil {
		app.Fail(c)
		return
	}

	res := resource.CreatedRecord(r)
	if !res {
		app.Fail(c)
		return
	}
	app.OK(c)
	return
}

// GetList 获取列表
func GetList(c *gin.Context)  {
	var page, _ 	= 	strconv.Atoi(c.DefaultQuery("page","1"))
	var pageSize, _ = 	strconv.Atoi(c.DefaultQuery("pageSize","10"))
	var keywords	= 	c.DefaultQuery("keywords","")
	var order		= 	c.DefaultQuery("order","desc")
	if page <= 0 {
		app.FailWithMessage(e.GetMsg(e.PARAMS_ERROR),e.PARAMS_ERROR,c)
		return
	}
	dataReturn := make(map[string] interface{})
	var searchStruct resource.SearchStruct
	searchStruct.Page 	  = page
	searchStruct.PageSize = pageSize
	searchStruct.Keywords = keywords
	searchStruct.Order	  = order

	list,_  := resource.GetList(searchStruct)
	count := resource.GetCount(searchStruct)

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

	dataReturn["list"] = list
	dataReturn["count"] = count

	app.OkWithData(dataReturn,c)
	return
}

// SetTopStatus 设置为置顶状态
func SetTopStatus(c *gin.Context)  {
	var id, _ 	= 	strconv.Atoi(c.DefaultQuery("id","0"))
	var status, _ 	= 	strconv.Atoi(c.DefaultQuery("status","0"))
	if (status != 1 && status != 2) || id <= 0  {
		app.FailWithMessage(e.GetMsg(e.MISS_PARAMS),e.MISS_PARAMS,c)
		return
	}
	result := resource.SetTopStatus(id,status,"is_top")
	if !result {
		app.Fail(c)
		return
	}
	app.OK(c)
	return
}

// DelResource 删除热门资源
func DelResource(c *gin.Context)  {
	var id, _ 	= 	strconv.Atoi(c.DefaultQuery("id","0"))
	if  id <= 0  {
		app.FailWithMessage(e.GetMsg(e.MISS_PARAMS),e.MISS_PARAMS,c)
		return
	}
	var status = 2
	result := resource.SetTopStatus(id,status,"is_delete")
	if !result {
		app.Fail(c)
		return
	}
	app.OK(c)
	return
}

// GetResourceByType 根据资源类型获取资源
func GetResourceByType(c *gin.Context)  {
	var resourceType, _ 	= 	strconv.Atoi(c.DefaultQuery("type","0"))
	var page, _ 	= 	strconv.Atoi(c.DefaultQuery("page","0"))
	// 获取总共的资源
	res := system.GetTagsTypeList("resourceType")
	var  hyStack []int
	for _,v := range res{
		hyStack = append(hyStack,v.ParamValue)
	}
	// 如果资源不在参数表返回false
	if !util.InArrayHelper(resourceType,hyStack) {
		app.Fail(c)
		return
	}

	factory := new(_interface.Factory)
	p 	:= factory.Generate(resourceType)
	if p == nil {
		app.Fail(c)
		return

	}
	data := p.GetList(page,10)
	app.OkWithData(data,c)
}
