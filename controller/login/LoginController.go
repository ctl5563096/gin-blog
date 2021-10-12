package login

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	_ "encoding/json"
	"gin-blog/models/user"
	"gin-blog/pkg/app"
	"gin-blog/pkg/cache/mainCache"
	"gin-blog/pkg/e"
	_ "gin-blog/pkg/e"
	"gin-blog/pkg/util"
	_ "gin-blog/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	_ "github.com/gin-gonic/gin/binding"
	_ "github.com/go-playground/validator/v10"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

// Request 验证规则
type Request struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AdminUser 用户结构体
type AdminUser struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	PhoneNum string `json:"phone_num"`
	IsBlack  int `json:"is_black"`
}

// LoginBackend /** 后台登陆 **/
func LoginBackend(c *gin.Context)  {
	var requestData Request
	data := make(map[string] interface{})
	err  := c.ShouldBindBodyWith(&requestData,binding.JSON)
	if err != nil {
		app.LoseWithParameter(err.Error(),c)
		return
	}
	// 利用username去查询密码再做对比
	find := make(map[string]interface{})
	find["user_name"] = requestData.Username
	find["is_use"]	 = 1
	find["is_black"] = 0
	res,err := user.GetUserPassWordByUserName(find)
	if err != nil {
		app.Fail(c)
		return
	}
	data["user_info"]  = res
	// 判断密码是否相等
	if res !=nil && res[0].Password != requestData.Password{
		app.FailWithMessage(e.GetMsg(e.PASSWORD_ERROR),e.PASSWORD_ERROR,c)
		return
	}
	// 先判断切片长度
	if len(res) > 0 {
		token, err := GetTokenByAdminUser((*AdminUser)(res[0]))
		//这里生成token
		//token,err 	  := GetToken(requestData.Username)
		if err != nil {
			app.FailWithMessage(e.GetMsg(e.ERROR_AUTH_TOKEN), e.ERROR_AUTH_TOKEN, c)
			return
		}
		data["token"] = token
		app.OkWithCodeData("登陆成功", data, e.SUCCESS, c)
	}else {
		app.FailWithMessage("系统错误",1,c)
	}
	return
}

// GetToken 生成token
func GetToken(userName string) (string,error) {
	// token生成算法 使用缓存
	redisConnect := mainCache.MainRedisConn.Get()
	// 使用完连接一定要关闭返回给连接池 不然会造成redis连接数过高
	defer func(redisConnect redis.Conn) {
		err := redisConnect.Close()
		if err != nil {
			util.WriteLog("close_redis_error",4,"close redis connect error,cache:mainCache")
		}
	}(redisConnect)
	h := md5.New()
	h.Write([]byte((strconv.FormatInt(time.Now().Unix(), 10)) + "_" + userName))
	// md5加密下token
	token := hex.EncodeToString(h.Sum(nil))
	_, err := redisConnect.Do("SET",token,"test_token","ex",86400)
	if err != nil {
		util.WriteLog("get_token_error",4,"get_token_error:" + userName + "。result：" + err.Error())
		return "", err
	}
	return token,nil
}

// CheckToken 检查token是否过期
func CheckToken(c *gin.Context)  {
	data := make(map[string] interface{})
	token := c.DefaultQuery("token","")
	if  token == ""{
		app.FailWithMessage(e.GetMsg(e.ERROR_AUTH_CHECK_TOKEN_FAIL),e.ERROR_AUTH_CHECK_TOKEN_FAIL,c)
		return
	}
	app.OkWithCodeData("登陆成功",data,e.SUCCESS,c)
	return
}

// GetTokenByAdminUser 生成token 暂时用较为简单的md5加密
func GetTokenByAdminUser(userInfo *AdminUser) (string,error) {
	// token生成算法 使用缓存
	redisConnect := mainCache.MainRedisConn.Get()
	// 使用完连接一定要关闭返回给连接池 不然会造成redis连接数过高
	defer func(redisConnect redis.Conn) {
		err := redisConnect.Close()
		if err != nil {
			util.WriteLog("close_redis_error",4,"close redis connect error,cache:mainCache")
		}
	}(redisConnect)
	h := md5.New()
	h.Write([]byte(userInfo.PhoneNum + "_" + userInfo.UserName + "_" + (strconv.FormatInt(time.Now().Unix(), 10))))
	userJson, err := json.Marshal(userInfo)
	token  := hex.EncodeToString(h.Sum(nil))
	_, err = redisConnect.Do("SET", token, userJson, "ex", 86400)
	if err != nil {
		util.WriteLog("get_token_error",4,"get_token_error:" + userInfo.UserName + "_" + userInfo.PhoneNum + "。result：" + err.Error())
		return "", err
	}
	return token,nil
}

// LoginOut 登出接口
func LoginOut(c *gin.Context)  {
	requestData := make(map[string]interface{}) //注意该结构接受的内容
	err  := c.BindJSON(&requestData)
	if err != nil {
		app.FailWithMessage(e.GetMsg(e.MISS_PARAMS),e.MISS_PARAMS,c)
		return
	}
	// token生成算法 使用缓存
	redisConnect := mainCache.MainRedisConn.Get()
	// 使用完连接一定要关闭返回给连接池 不然会造成redis连接数过高
	defer func(redisConnect redis.Conn) {
		err := redisConnect.Close()
		if err != nil {
			util.WriteLog("close_redis_error",4,"close redis connect error,cache:mainCache")
		}
	}(redisConnect)
	// 删除token
	_, err = redisConnect.Do("DEL",requestData["token"])
	if err != nil {
		app.FailWithMessage(e.GetMsg(e.SYSTEM_ERROR),e.SYSTEM_ERROR,c)
		return
	}
	app.OK(c)
	return
}