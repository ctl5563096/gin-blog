package token

import (
	"encoding/json"
	"gin-blog/pkg/app"
	"gin-blog/pkg/cache/mainCache"
	"gin-blog/pkg/e"
	"gin-blog/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

// AdminUser 用户结构体
type AdminUser struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	PhoneNum string `json:"phone_num"`
	IsBlack  int `json:"is_black"`
}

// BeforeBusiness 业务中间件检查token是否有效
func BeforeBusiness() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先从从GET的请求体里面取 最后再从POST里面获取
		token := c.DefaultQuery("token","")
		if  token == ""{
			requestData := make(map[string]interface{}) //注意该结构接受的内容
			err  := c.BindJSON(&requestData)
			if err != nil {
				app.FailWithMessage(e.GetMsg(e.MISS_TOKEN),e.MISS_TOKEN,c)
				c.Abort()
				return
			}
			for k,v := range requestData {
				if k == "token" {
					token = v.(string)
				}
			}
			if token == "" {
				app.FailWithMessage(e.GetMsg(e.MISS_TOKEN),e.MISS_TOKEN,c)
				c.Abort()
				return
			}
		}
		redisConnect := mainCache.MainRedisConn.Get()
		// 使用完连接一定要关闭返回给连接池 不然会造成redis连接数过高
		defer func(redisConnect redis.Conn) {
			err := redisConnect.Close()
			if err != nil {
				util.WriteLog("close_redis_error",4,"close redis connect error,cache:mainCache")
			}
		}(redisConnect)
		// 从缓存里面查询一下是否存在token 如果存在则将对应的信息存进request里面方便后面调用
		res,err := redis.String(redisConnect.Do("GET",token))
		if err != nil{
				app.FailWithMessage(e.GetMsg(e.ERROR_AUTH_CHECK_TOKEN_FAIL),e.ERROR_AUTH_CHECK_TOKEN_FAIL,c)
				c.Abort()
				return
		}
		// 先判断用户的有效性 将用户信息注入到ctx里面 方便后面调用
		userInfo := make(map[string]interface{})
		err = json.Unmarshal([]byte(res), &userInfo)
		if err != nil {
			app.FailWithMessage("系统错误",1,c)
			c.Abort()
			return
		}
		// 由于是接口值 需要强转一下类型
		if int(userInfo["is_black"].(float64)) > 0 {
			app.FailWithMessage(e.GetMsg(e.USER_IN_BLACK),e.USER_IN_BLACK,c)
			c.Abort()
			return
		}
		// 注入ctx上下文里面
		c.Set("userInfo", userInfo)
		// 相当于php lumen return $next($request)
		c.Next()
	}
}