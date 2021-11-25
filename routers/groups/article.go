package groups

import (
	"gin-blog/controller/blogIndex"
	"github.com/gin-gonic/gin"
)

// ArticleBaseRouter /** 文章基本接口 **/
func ArticleBaseRouter(Router *gin.RouterGroup) {
	// v1版接口
	apiRouterV1 := Router.Group("/v1/article")
	{
		// 获取获取博客首页文章
		apiRouterV1.GET("index", blogIndex.GetIndexBlog)
		// 下面的路由都需要验证token
		//apiRouterV1.Use(token.BeforeBusiness())
		// 创建新的博客文章
		apiRouterV1.POST("create", blogIndex.AddArticle)
		// 获取文章列表
		apiRouterV1.GET("list",blogIndex.GetArticleList)
		// 获取文章详情
		apiRouterV1.GET("detail",blogIndex.GetArticleDetail)
		// 修改文章详情
		apiRouterV1.PUT("detail",blogIndex.EditArticleDetail)
	}
}