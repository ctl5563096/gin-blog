package routers

import (
	"gin-blog/pkg/e"
	"fmt"
	"gin-blog/pkg/setting"
	"gin-blog/routers/groups"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"time"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// 健康检测接口
	r.GET("/health", func(c *gin.Context) {
		c.String(e.SUCCESS, "health - " + fmt.Sprint(time.Now().Unix()))
	})
	gin.SetMode(setting.RunMode)
	group := r.Group("")
	groups.LoginBaseRouter(group)
	ginpprof.Wrap(r)
	return r
}