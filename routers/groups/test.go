package groups

import (
	"gin-blog/controller/test"
	"github.com/gin-gonic/gin"
)

// TestRouter /** 测试接口 **/
func TestRouter(Router *gin.RouterGroup) {
	// v1版接口
	apiRouterV1 := Router.Group("v1/test")
	{
		// 获取主业务接口
		apiRouterV1.POST("testSecond", test.TestPortSecond)
		// 获取主业务接口
		apiRouterV1.GET("test", test.TestPort)
		// 测试图片缩略图
		apiRouterV1.GET("thumb", test.TestThumb)
	}

}
