package groups

import (
	"gin-blog/controller/upload"
	token "gin-blog/middleware"
	"github.com/gin-gonic/gin"
)

// UploadBaseRouter /** 上传文件的基本接口 **/
func UploadBaseRouter(Router *gin.RouterGroup) {
	// v1版接口 上传类需要验证对应的token
	apiRouterV1 := Router.Group("/v1/uploads").Use(token.BeforeBusiness())
	{
		// 创建用户
		apiRouterV1.POST("avatar", upload.Avatar)
	}
}
