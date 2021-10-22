package user

import (
	"gin-blog/models/user"
	"gin-blog/pkg/app"
	"gin-blog/pkg/e"
	"gin-blog/pkg/util"
	valid "gin-blog/vaild"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-playground/validator/v10"
	"strconv"
)

// AdminUser 用户结构体
type AdminUser struct {
	ID 		 int	`json:"id"`
	UserName string `json:"user_name" validate:"required,lt=11,gt=6"`
	Password string `json:"password"  validate:"required"`
	PhoneNum string `json:"phone_num" validate:"required,len=11"`
	Email	 string `json:"email" validate:"email"`
	Sex	 	 int 	`json:"sex"`
	Avatar	 string `json:"avatar"`
}

// UpdateUserInfo 更新用户结构体
type UpdateUserInfo struct {
	ID 		 int	`json:"id" validate:"required"`
	UserName string `json:"user_name" validate:"required,lt=11,gt=6"`
	PhoneNum string `json:"phone_num" validate:"required,len=11"`
	Email	 string `json:"email" validate:"email"`
	Sex	 	 int 	`json:"sex"`
	Avatar	 string `json:"avatar"`
}

// CreateUser 创建用户
func CreateUser(c *gin.Context)  {
	var r AdminUser
	var errStr string
	var errorMap map[string][]string
	err := c.ShouldBind(&r)
	validate := validator.New()
	err = validate.Struct(r)
	if err != nil {
		switch err.(type) {
			case validator.ValidationErrors:
				errorMap = valid.Translate(err)
				//循环遍历Map 只返回第一个错误信息
				for _,v:= range errorMap{
					for _,z := range v{
						util.WriteLog("user_business_error",4,z)
						app.FailWithMessage(z,4,c)
						return
					}
				}
		default:
			errStr = "未知错误"
		}
		app.FailWithMessage(errStr,1,c)
		return
	}

	// 插入用户逻辑 新建用户先统一使用默认的头像 性别为未知
	avatar 	 :="/resource/public/pic/avatar.jpeg"
	r.Avatar = avatar
	r.Sex 	 = 3

	_,_ = user.CreatUser(&r)
	app.OK(c)
	return
}

//UpdateUser 更新用户信息
func UpdateUser(c *gin.Context)  {
	var r UpdateUserInfo
	err := c.ShouldBind(&r)
	validate := validator.New()
	err = validate.Struct(r)
	if err != nil {
		util.WriteLog("user_business_error",4,err.Error())
		app.FailWithMessage(err.Error(),4,c)
		return
	}

	// 查询是否存在该用户
	res := user.GetUserById(r.ID)
	if !res {
		app.FailWithMessage(e.GetMsg(e.USER_NOT_EXIST),e.USER_NOT_EXIST,c)
		return
	}
	// 执行更新
	user.UpdateUser(&r)
	if r.ID < 1 {
		app.FailWithMessage(e.GetMsg(e.DELETE_ERROR),e.DELETE_ERROR,c)
		return
	}
	app.OK(c)
	return
}

// OpenUser 启用用户
func OpenUser(c *gin.Context)  {
	userInfo := util.GetUserInfo(c)
	// 当前发起操作的用户权限
	if userInfo.(map[string]interface{})["role"].(float64) < 1{
		app.FailWithMessage(e.GetMsg(e.OPERATION_ONT_PERMITTED),e.OPERATION_ONT_PERMITTED,c)
		return
	}
	operateId         := c.Query("id")
	// 转整形
	intOperateId, _   := strconv.Atoi(operateId)
	if  intOperateId < 1{
		app.FailWithMessage(e.GetMsg(e.MISS_PARAMS),e.MISS_PARAMS,c)
		return
	}
	id := user.OpenUser(intOperateId)
	if id < 1 {
		app.FailWithMessage(e.GetMsg(e.UPDATE_ERROR),e.UPDATE_ERROR,c)
		return
	}
	app.OK(c)
	return
}

// BlackUser 对用户进行拉黑操作
func BlackUser(c *gin.Context)  {
	userInfo := util.GetUserInfo(c)
	// 当前发起操作的用户权限
	if userInfo.(map[string]interface{})["role"].(float64) < 1{
		app.FailWithMessage(e.GetMsg(e.OPERATION_ONT_PERMITTED),e.OPERATION_ONT_PERMITTED,c)
		return
	}
	operateId         := c.Query("id")
	status         	  := c.Query("status")
	// 转整形
	intOperateId, _   := strconv.Atoi(operateId)
	intStatus, _   	  := strconv.Atoi(status)
	if  intOperateId < 1{
		app.FailWithMessage(e.GetMsg(e.MISS_PARAMS),e.MISS_PARAMS,c)
		return
	}
	id := user.ChangeBlackStatus(intOperateId,intStatus)
	if id < 1 {
		app.FailWithMessage(e.GetMsg(e.UPDATE_ERROR),e.UPDATE_ERROR,c)
		return
	}
	app.OK(c)
	return
}