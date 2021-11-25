package groups

import (
	"gin-blog/controller/system"
	"github.com/gin-gonic/gin"
)

// SystemBaseRouter /** 系统的基本接口 **/
func SystemBaseRouter(Router *gin.RouterGroup) {
	// v1版接口 上传类需要验证对应的token
	apiRouterV1 := Router.Group("/v1/system")
	{
		// 新增参数
		apiRouterV1.POST("param", system.Create)
		// 后台获取参数列表
		apiRouterV1.GET("getParamList", system.GetList)
		// 后台删除指定参数
		apiRouterV1.DELETE("param",system.DelParam)
		// 后台更新参数
		apiRouterV1.PUT("param",system.EditParam)
		// 后台更新参数
		apiRouterV1.GET("param",system.GetDetail)
	}
}
