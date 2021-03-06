package groups

import (
	"gin-blog/controller/login"
	"github.com/gin-gonic/gin"
)

// LoginBaseRouter /** 登陆基本接口 **/
func LoginBaseRouter(Router *gin.RouterGroup) {
	// v1版接口
	apiRouterV1 := Router.Group("login/v1")
	{
		// 获取主业务接口
		apiRouterV1.POST("login", login.LoginBackend)
		// 检查token是否有效
		apiRouterV1.POST("check-token", login.CheckToken)
		// 登出接口
		apiRouterV1.POST("logout", login.LoginOut)
	}
}
