package groups

import (
	"gin-blog/routers/api/login"
	"github.com/gin-gonic/gin"

)

// LoginBaseRouter /** 登陆基本接口 **/
func LoginBaseRouter(Router *gin.RouterGroup) {
	// v1版接口
	apiRouterV1 := Router.Group("login/v1")
	{
		// 获取主业务接口
		apiRouterV1.POST("login", login.LoginBackend)
	}
}
