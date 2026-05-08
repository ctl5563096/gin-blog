package message

import (
	"gin-blog/models/message"
	"gin-blog/pkg/app"
	"gin-blog/pkg/e"
	"gin-blog/pkg/util"
	valid "gin-blog/vaild"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"strconv"
)

// Create 前台新增留言
func Create(c *gin.Context) {
	var r message.CreateMessage
	var errStr string
	var errorMap map[string][]string

	err := c.ShouldBindBodyWith(&r, binding.JSON)
	validate := validator.New()
	err = validate.Struct(r)
	if err != nil {
		switch err.(type) {
		case validator.ValidationErrors:
			errorMap = valid.Translate(err)
			for _, v := range errorMap {
				for _, z := range v {
					util.WriteLog("message_business_error", 4, z)
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

	r.Ip = c.ClientIP()
	r.UserAgent = c.Request.UserAgent()
	if err = message.CreateRecord(&r); err != nil {
		app.FailWithMessage("留言保存失败", 1, c)
		return
	}
	app.OkWithMessage("留言提交成功", c)
}

// GetList 后台留言列表
func GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keywords := c.DefaultQuery("keywords", "")

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	list, err := message.GetList(page, pageSize, keywords)
	if err != nil {
		app.FailWithMessage(e.GetMsg(e.ERROR), e.ERROR, c)
		return
	}
	count, err := message.GetCount(keywords)
	if err != nil {
		app.FailWithMessage(e.GetMsg(e.ERROR), e.ERROR, c)
		return
	}

	result := make(map[string]interface{})
	result["list"] = list
	result["count"] = count
	app.OkWithData(result, c)
}
