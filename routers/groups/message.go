package groups

import (
	"gin-blog/controller/message"
	token "gin-blog/middleware"
	"github.com/gin-gonic/gin"
)

// MessageBaseRouter 前台留言接口
func MessageBaseRouter(Router *gin.RouterGroup) {
	apiRouterV1 := Router.Group("/v1/message")
	{
		apiRouterV1.POST("/create", message.Create)
	}

	apiRouterAdmin := Router.Group("/v1/message").Use(token.BeforeBusiness())
	{
		apiRouterAdmin.GET("/list", message.GetList)
	}
}
