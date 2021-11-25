package groups

import (
	"gin-blog/controller/upload"
	"github.com/gin-gonic/gin"
)

// UploadBaseRouter /** 上传文件的基本接口 **/
func UploadBaseRouter(Router *gin.RouterGroup) {
	// v1版接口 上传类需要验证对应的token
	apiRouterV1 := Router.Group("/v1/uploads")
	{
		// 上传头像
		apiRouterV1.POST("avatar", upload.Avatar)
		// 上传图片
		apiRouterV1.POST("pic", upload.Photo)
	}
}
