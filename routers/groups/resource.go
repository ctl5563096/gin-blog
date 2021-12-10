package groups

import (
	"gin-blog/controller/resource"
	"github.com/gin-gonic/gin"
)

// ResourceBaseRouter /** 资源路由 **/
func ResourceBaseRouter(Router *gin.RouterGroup) {
	// v1版接口
	apiRouterV1 := Router.Group("/v1/resource")
	{
		// 新增接口
		apiRouterV1.POST("/hot/resource", resource.CreateNewHotRecord)
		// 获取后台列表
		apiRouterV1.GET("/hot", resource.GetList)
		// 更新热门资源的置顶状态
		apiRouterV1.PUT("/hot/updateTopStatus", resource.SetTopStatus)
		// 更新热门资源的删除状态
		apiRouterV1.DELETE("/hot", resource.DelResource)
		// 根据资源类型获取资源列表
		apiRouterV1.GET("/byType", resource.GetResourceByType)
	}
}
