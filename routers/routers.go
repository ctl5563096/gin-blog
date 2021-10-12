package routers

import (
	"fmt"
	"gin-blog/pkg/e"
	"gin-blog/pkg/setting"
	"gin-blog/routers/groups"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"time"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	// 部分中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// 鉴权中间件 只对需要鉴权的路由进行拦截
	// 健康检测接口
	r.GET("/health", func(c *gin.Context) {
		c.String(e.SUCCESS, "health - " + fmt.Sprint(time.Now().Unix()))
	})
	gin.SetMode(setting.RunMode)
	group := r.Group("")
	groups.LoginBaseRouter(group)
	groups.TestRouter(group)
	ginpprof.Wrap(r)
	return r
}