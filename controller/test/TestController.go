package test

import (
	"fmt"
	"gin-blog/pkg/app"
	"gin-blog/pkg/util"
	"github.com/gin-gonic/gin"
)

// Test 测试接口
func Test(c *gin.Context)  {
	userInfo := util.GetUserInfo(c)
	fmt.Println(userInfo.(map[string]interface{})["is_black"])
	data := make(map[string] interface{})
	app.OkWithCodeData("测试成功",data,0,c)
	return
}