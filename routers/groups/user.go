package groups

import (
	"gin-blog/controller/user"
	token "gin-blog/middleware"
	"github.com/gin-gonic/gin"
)

// UserBaseRouter /** 用户基本接口 **/
func UserBaseRouter(Router *gin.RouterGroup) {
	// v1版接口
	apiRouterV1 := Router.Group("/v1/user")
	{
		// 创建用户
		apiRouterV1.POST("create", user.CreateUser).Use(token.BeforeBusiness())
		// 开启用户使用 这里需要验证token
		apiRouterV1.GET("open", user.OpenUser)
		// 更新用户 这里需要验证token
		apiRouterV1.PUT("update", user.UpdateUser)
		// 对用户的黑名单状态进行操作 这里需要验证token
		apiRouterV1.PUT("black", user.BlackUser)
	}
}
