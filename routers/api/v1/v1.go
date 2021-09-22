package v1

import (
	"gin-blog/models/user"
	"gin-blog/pkg/e"
	"gin-blog/pkg/util"
	"github.com/Unknwon/com"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetUser 获取用户
func GetUser(c *gin.Context) {
	keywords := c.DefaultQuery("keyword","")
	username  := c.Query("username")
	pageSize,_ := com.StrTo(c.DefaultQuery("page","10")).Int()

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if username != "" {
		maps["username"] = username
	}

	if keywords != "" {
		maps["keywords"] = keywords
	}
	// 给个默认值
	isUse := 1
	// 状态码返回默认值
	code := e.SUCCESS

	if arg := c.Query("state"); arg != "" {
		isUse = com.StrTo(arg).MustInt()
		maps["is_use"] = isUse
	}

	data["lists"] = user.GetUsers(util.GetPage(c),pageSize,maps)
	data["total"] = user.GetUsersTotal(maps)

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : data,
	})
}

// AddTag 新增文章标签
func AddTag(c *gin.Context) {
	data := make(map[string]interface{})
	c.JSON(http.StatusOK, gin.H{
		"code" : 200,
		"msg" : e.GetMsg(200),
		"data" : data,
	})
}

// EditTag 修改文章标签
func EditTag(c *gin.Context) {
}

// DeleteTag 删除文章标签
func DeleteTag(c *gin.Context) {
}
