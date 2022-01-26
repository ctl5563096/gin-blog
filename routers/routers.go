package routers

import (
	"fmt"
	token "gin-blog/middleware"
	"gin-blog/pkg/e"
	"gin-blog/pkg/setting"
	"gin-blog/routers/groups"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	// 部分中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// 全局异常捕获器
	r.Use(token.Recover)
	// 鉴权中间件 只对需要鉴权的路由进行拦截
	// 健康检测接口
	r.GET("/health", func(c *gin.Context) {
		c.String(e.SUCCESS, "health - " + fmt.Sprint(time.Now().Unix()))
	})
	// 静态资源文件访问
	r.StaticFS("/resource", http.Dir("resource/"))
	gin.SetMode(setting.RunMode)
	group := r.Group("")
	// 注册基础路由
	groups.LoginBaseRouter(group)
	groups.TestRouter(group)
	groups.UserBaseRouter(group)
	groups.UploadBaseRouter(group)
	groups.ArticleBaseRouter(group)
	groups.RuleBaseRouter(group)
	groups.SystemBaseRouter(group)
	groups.ResourceBaseRouter(group)
	ginpprof.Wrap(r)
	// 默认404路由为非法请求
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": -1, "msg": "非法请求"})
	})
	return r
}