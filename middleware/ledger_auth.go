package token

import (
	"encoding/json"
	"gin-blog/pkg/app"
	"gin-blog/pkg/cache/mainCache"
	"gin-blog/pkg/e"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gomodule/redigo/redis"
	"strings"
)

const LedgerOwnerContextKey = "ledgerOwnerId"

type LedgerSession struct {
	OpenId  string `json:"openid"`
	UnionId string `json:"unionid"`
}

func LedgerAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenValue := getLedgerToken(c)
		if tokenValue == "" {
			app.MissToken(c)
			c.Abort()
			return
		}

		redisConnect := mainCache.MainRedisConn.Get()
		defer redisConnect.Close()

		res, err := redis.String(redisConnect.Do("GET", "ledger:token:"+tokenValue))
		if err != nil {
			app.MissToken(c)
			c.Abort()
			return
		}

		var session LedgerSession
		if err = json.Unmarshal([]byte(res), &session); err != nil || session.OpenId == "" {
			app.FailWithMessage(e.GetMsg(e.ERROR_AUTH), e.ERROR_AUTH, c)
			c.Abort()
			return
		}

		c.Set(LedgerOwnerContextKey, session.OpenId)
		c.Next()
	}
}

func getLedgerToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	}
	if authHeader != "" {
		return strings.TrimSpace(authHeader)
	}

	tokenValue := c.DefaultQuery("token", "")
	if tokenValue != "" {
		return tokenValue
	}

	requestData := make(map[string]interface{})
	if err := c.ShouldBindBodyWith(&requestData, binding.JSON); err != nil {
		return ""
	}
	if value, ok := requestData["token"].(string); ok {
		return value
	}
	return ""
}
