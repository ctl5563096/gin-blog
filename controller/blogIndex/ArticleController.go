package blogIndex

import (
	"gin-blog/models/blog"
	"gin-blog/pkg/app"
	"gin-blog/pkg/e"
	"gin-blog/pkg/util"
	valid "gin-blog/vaild"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"math"
	"strconv"
)

// AddArticleStruct 新增结构体
type AddArticleStruct struct {
	ID 		 int	`json:"id"`
	Title 	 string `json:"title" validate:"required"`
	Summary  string `json:"summary" validate:"required"`
	Content  string `json:"content" validate:"required"`
	Cover    string `json:"cover"`
	Author   string `json:"author"`
	IsShow   uint  	`json:"is_show"`
}

// GetIndexBlog 首页最新文章
func GetIndexBlog(c *gin.Context)  {
	data := make(map[string] interface{})
	res,err := blog.GetIndexArticle()
	if err != nil {
		app.FailWithMessage("获取文章首页失败!",1,c)
		return
	}
	data["list"] = res
	app.OkWithData(data,c)
	return
}

// AddArticle 添加文章
func AddArticle(c *gin.Context)  {
	var r AddArticleStruct
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
					util.WriteLog("user_business_error",4,z)
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

	// 展示默认作者名称
	r.Author = "shy"

	_,_ = blog.CreatArticle(&r)
	app.OK(c)
	return
}

// GetArticleList 获取文章列表
func GetArticleList(c *gin.Context)  {
	var page, _ = 	strconv.Atoi(c.DefaultQuery("page","1"))
	var pageSize, _ = 	strconv.Atoi(c.DefaultQuery("pageSize","10"))
	var keywords= 	c.DefaultQuery("keywords","")
	var order= 	c.DefaultQuery("order","desc")
	dataReturn := make(map[string] interface{})
	var searchStruct blog.ArticleSearchList
	searchStruct.Page = page
	searchStruct.PageSize = pageSize
	searchStruct.KeyWords = keywords
	searchStruct.Order	= order
	res,err := blog.SearchArticle(&searchStruct)
	count := blog.SearchArticleCountA(&searchStruct)
	//var count int
	//ch := make(chan int,1)
	//// 这里去执行下协程去查询下总数 使用管道的话必须在协程里面关闭close(ch)管道 否则程序不会终端
	//err = app.GoroutineNotPanic(func() (err error) {
	//	blog.SearchArticleCount(&searchStruct,ch)
	//	count = <-ch
	//	return
	//})
	if page <= 0 {
		app.FailWithMessage(e.GetMsg(e.PARAMS_ERROR),e.PARAMS_ERROR,c)
		return
	}
	if err != nil{
		app.FailWithMessage(e.GetMsg(e.FIND_ARTICLE_ERROR),e.FIND_ARTICLE_ERROR,c)
		return
	}
	for key,value := range res{
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
	dataReturn["list"] = res
	dataReturn["count"] = count
	app.OkWithData(dataReturn,c)
	return
}

// GetArticleDetail 获取文章详情
func GetArticleDetail(c *gin.Context)  {
	var id, _ = 	strconv.Atoi(c.DefaultQuery("id","0"))
	if id <= 0 {
		app.FailWithMessage(e.GetMsg(e.MISS_PARAMS),e.MISS_PARAMS,c)
		return
	}
	res := blog.GetDetail(id)
	app.OkWithData(res,c)
	return
}

// EditArticleDetail 修改文章详情
func EditArticleDetail(c *gin.Context)  {
	var r blog.EditArticleStruct
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
					util.WriteLog("user_business_error",4,z)
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

	// 执行修改的逻辑
	res := blog.EditDetail(&r)
	if !res {
		app.Fail(c)
		return
	}
	app.OK(c)
	return
}