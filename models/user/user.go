package user

import (
	"gin-blog/pkg/util"
)

type AdminUser struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	PhoneNum string `json:"phone_num"`
	IsBlack int `json:"is_black"`
}

// GetUsers 获取用户列表
func GetUsers(pageNum int, pageSize int, maps interface {}) ([]*AdminUser,error) {
	var res  []*AdminUser
	db := db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&res)
	if db.Error != nil {
		util.WriteLog("mysql_error",2,db.Error.Error())
		return nil,db.Error
	}
	return res,nil
}

// GetUserPassWordByUserName 根据用户名搜索用户
func GetUserPassWordByUserName(maps interface {}) ([]*AdminUser,error) {
	var res  []*AdminUser
	db := db.Where(maps).First(&res)
	if db.Error != nil {
		util.WriteLog("mysql_error",2,db.Error.Error())
		return nil,db.Error
	}
	return res,db.Error
}