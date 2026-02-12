package dashboard

import (
	"gin-blog/models/dashboard"
	"gin-blog/pkg/app"
	"github.com/gin-gonic/gin"
)

// GetStats 获取面板统计数据
func GetStats(c *gin.Context) {
	stats := dashboard.GetDashboardStats()
	app.OkWithData(stats, c)
}
