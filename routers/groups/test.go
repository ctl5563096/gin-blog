package groups

import (
	"gin-blog/controller/test"
	token "gin-blog/middleware"
	"github.com/gin-gonic/gin"
)

// TestRouter /** 测试接口 **/
func TestRouter(Router *gin.RouterGroup) {
	// v1版接口
	apiRouterV1 := Router.Group("v1/test")
	apiRouterV1.Use(token.BeforeBusiness())
	{
		// 获取主业务接口
		apiRouterV1.GET("test", test.Test)
	}

}
