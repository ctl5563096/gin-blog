package login

import (
	_ "encoding/json"
	"gin-blog/models/user"
	"gin-blog/pkg/app"
	"gin-blog/pkg/e"
	_ "gin-blog/pkg/e"
	_ "gin-blog/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	_ "github.com/gin-gonic/gin/binding"
	_ "github.com/go-playground/validator/v10"
)

// Request 验证规则
type Request struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
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
	//这里生成token
	token,err 	  := GetToken(requestData.Username)
	if err != nil {
		app.FailWithMessage(e.GetMsg(e.ERROR_AUTH_TOKEN),e.ERROR_AUTH_TOKEN,c)
		return
	}
	data["token"] = token
	app.OkWithCodeData("登陆成功",data,e.SUCCESS,c)
	return
}

// GetToken 生成token
func GetToken(userName string) (string,error) {
	return "22222222",nil
}