package groups

import (
	"gin-blog/controller/resource"
	token "gin-blog/middleware"
	"github.com/gin-gonic/gin"
)

// ResourceBaseRouter /** 资源路由 **/
func ResourceBaseRouter(Router *gin.RouterGroup) {
	// v1版接口
	apiRouterV1 := Router.Group("/v1/resource")
	{
		// 前端路由接口
		apiRouterV1.GET("/music/list", resource.GetMusicList)
		apiRouterV1.GET("/photos/list", resource.GetPhotoList)
		apiRouterV1.GET("/ByTag",resource.GetResourceCodeArticle)
		apiRouterV1.Use(token.BeforeBusiness())
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
		// 新增音乐资源
		apiRouterV1.POST("/music/created", resource.CreateNewMusic)
		// 删除音乐资源
		apiRouterV1.DELETE("/music/delete", resource.DeleteMusicRecord)
		// 获取音乐资源详情
		apiRouterV1.GET("/music/detail", resource.GetMusicDetailBackend)
		// 修改音乐资源
		apiRouterV1.PUT("/music/detail", resource.UpdateAudio)
		// 创建图片资源
		apiRouterV1.POST("/photo/created",resource.CreatePhoto)
		// 删除图片资源
		apiRouterV1.DELETE("/photo/delete",resource.DeletePhotoRecord)
		// 修改图片资源
		apiRouterV1.PUT("/photo/detail", resource.UpdateAudio)
		// 获取资源详情
		apiRouterV1.GET("/photo/detail", resource.GetMusicDetailBackend)
	}
}
