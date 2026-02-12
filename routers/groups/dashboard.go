package groups

import (
	"gin-blog/controller/dashboard"
	token "gin-blog/middleware"
	"github.com/gin-gonic/gin"
)

// DashboardBaseRouter 注册面板相关路由
func DashboardBaseRouter(Router *gin.RouterGroup) {
	dashboardGroup := Router.Group("/v1/dashboard").Use(token.BeforeBusiness())
	{
		dashboardGroup.GET("/stats", dashboard.GetStats)
	}
}
