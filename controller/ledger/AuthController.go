package ledger

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	tokenMiddleware "gin-blog/middleware"
	"gin-blog/models/ledger"
	"gin-blog/pkg/app"
	"gin-blog/pkg/cache/mainCache"
	"gin-blog/pkg/e"
	"gin-blog/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gomodule/redigo/redis"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type WxLoginRequest struct {
	Code string `json:"code" binding:"required"`
}

type WxCodeSession struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func WxLogin(c *gin.Context) {
	var r WxLoginRequest
	if err := c.ShouldBindBodyWith(&r, binding.JSON); err != nil || r.Code == "" {
		app.FailWithParameter("缺少微信登录 code", c)
		return
	}

	wxSession, err := code2Session(r.Code)
	if err != nil {
		app.FailWithMessage(err.Error(), e.ERROR_AUTH_CHECK_TOKEN_FAIL, c)
		return
	}

	user, err := ledger.UpsertUser(wxSession.OpenId, wxSession.UnionId)
	if err != nil {
		app.FailWithMessage("用户登录信息保存失败", e.ERROR, c)
		return
	}

	tokenValue, err := createLedgerToken(wxSession)
	if err != nil {
		app.FailWithMessage(e.GetMsg(e.ERROR_AUTH_TOKEN), e.ERROR_AUTH_TOKEN, c)
		return
	}

	accounts, err := ledger.GetUserAccounts(wxSession.OpenId)
	if err != nil {
		app.FailWithMessage("账本信息获取失败", e.ERROR, c)
		return
	}

	app.OkWithData(gin.H{
		"token":     tokenValue,
		"openid":    wxSession.OpenId,
		"unionid":   wxSession.UnionId,
		"user":      user,
		"accounts":  accounts,
		"needSetup": len(accounts) == 0,
		"expiresIn": 86400,
	}, c)
}

func code2Session(code string) (WxCodeSession, error) {
	appId := os.Getenv("WECHAT_MINIAPP_APPID")
	appSecret := os.Getenv("WECHAT_MINIAPP_SECRET")
	if appId == "" {
		appId = os.Getenv("WX_APPID")
	}
	if appSecret == "" {
		appSecret = os.Getenv("WX_SECRET")
	}
	if appId == "" || appSecret == "" {
		return WxCodeSession{}, fmt.Errorf("缺少微信小程序 appid 或 secret 配置")
	}

	requestUrl := "https://api.weixin.qq.com/sns/jscode2session?appid=" + url.QueryEscape(appId) +
		"&secret=" + url.QueryEscape(appSecret) +
		"&js_code=" + url.QueryEscape(code) +
		"&grant_type=authorization_code"

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(requestUrl)
	if err != nil {
		return WxCodeSession{}, fmt.Errorf("微信登录请求失败")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return WxCodeSession{}, fmt.Errorf("微信登录响应读取失败")
	}

	var wxSession WxCodeSession
	if err = json.Unmarshal(body, &wxSession); err != nil {
		return WxCodeSession{}, fmt.Errorf("微信登录响应解析失败")
	}
	if wxSession.ErrCode != 0 {
		return WxCodeSession{}, fmt.Errorf("微信登录失败：%s", wxSession.ErrMsg)
	}
	if wxSession.OpenId == "" {
		return WxCodeSession{}, fmt.Errorf("微信登录未返回 openid")
	}
	return wxSession, nil
}

func createLedgerToken(wxSession WxCodeSession) (string, error) {
	redisConnect := mainCache.MainRedisConn.Get()
	defer redisConnect.Close()
	if err := redisConnect.Err(); err != nil {
		util.WriteLog("ledger-auth", 5, fmt.Sprintf("redis get connection failed: %v", err))
		return "", err
	}

	h := md5.New()
	h.Write([]byte(wxSession.OpenId + "_" + wxSession.SessionKey + "_" + fmt.Sprint(time.Now().UnixNano())))
	tokenValue := hex.EncodeToString(h.Sum(nil))
	util.WriteLog("ledger-auth", 1, fmt.Sprintf("create token start openid=%s tokenPrefix=%s", maskOpenId(wxSession.OpenId), tokenValue[:8]))

	session := tokenMiddleware.LedgerSession{
		OpenId:  wxSession.OpenId,
		UnionId: wxSession.UnionId,
	}
	sessionJson, err := json.Marshal(session)
	if err != nil {
		util.WriteLog("ledger-auth", 4, fmt.Sprintf("marshal ledger session failed openid=%s err=%v", maskOpenId(wxSession.OpenId), err))
		return "", err
	}

	_, err = redisConnect.Do("SET", "ledger:token:"+tokenValue, sessionJson, "EX", 86400)
	if err != nil && err != redis.ErrNil {
		util.WriteLog("ledger-auth", 5, fmt.Sprintf("redis set token failed openid=%s tokenPrefix=%s err=%v", maskOpenId(wxSession.OpenId), tokenValue[:8], err))
		return "", err
	}
	util.WriteLog("ledger-auth", 1, fmt.Sprintf("redis set token success openid=%s tokenPrefix=%s", maskOpenId(wxSession.OpenId), tokenValue[:8]))
	return tokenValue, nil
}

func maskOpenId(openid string) string {
	if len(openid) <= 8 {
		return "***"
	}
	return openid[:4] + "***" + openid[len(openid)-4:]
}
