package groups

import (
	"gin-blog/controller/upload"
	token "gin-blog/middleware"
	"github.com/gin-gonic/gin"
)

// UploadBaseRouter /** 上传文件的基本接口 **/
func UploadBaseRouter(Router *gin.RouterGroup) {
	// v1版接口 上传类需要验证对应的token 上传类也需要接入token验证
	apiRouterV1 := Router.Group("/v1/uploads").Use(token.BeforeBusiness())
	{
		// 上传头像
		apiRouterV1.POST("avatar", upload.Avatar)
		// 上传图片
		apiRouterV1.POST("pic", upload.Photo)
		// 上传音频
		apiRouterV1.POST("audio", upload.Music)
	}
}
