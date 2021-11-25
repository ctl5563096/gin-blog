package user

import (
	"encoding/json"
	"gin-blog/pkg/util"
)

// 自定义表名
var tableName = "ci_admin_user"

type AdminUser struct {
	Id		 int 	`json:"id"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	PhoneNum string `json:"phone_num"`
	IsBlack  int `json:"is_black"`
	Role 	 int `json:"role"`
}

type FindUser struct {
	Id int `json:"id"`
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

// CreatUser 创建用户
func CreatUser(params interface{}) (id int ,err error) {
	err = db.Create(params).Error
	if err != nil {
		return 0,err
	}
	m := make(map[string]interface{})
	j,_ := json.Marshal(params)
	err = json.Unmarshal(j, &m)
	if err != nil {
		return 0, err
	}
	recordId := int(m["id"].(float64))
	return recordId,err
}

// GetUserById 根据用户id判断用户是否存在
func GetUserById(id int) bool {
	var result FindUser
	db.Table(tableName).Where("id = ?" ,id).First(&result)
	if result.Id < 1  {
		return false
	}
	return true
}

//UpdateUser 更新用户信息
func UpdateUser(params interface{}){
	var user FindUser
	db.Table(tableName).Model(&user).Save(params)
}

//OpenUser 更新用户信息
func OpenUser(id int) (res int){
	var user FindUser
	db.Table(tableName).Model(&user).Where("id = ?" , id).Update("is_use",1)
	return user.Id
}

//ChangeBlackStatus 更新用户的黑名单状态
func ChangeBlackStatus(id int,status int) (res int){
	var user FindUser
	if status< 0 || status > 1{
		return 0
	}
	db.Table(tableName).Model(&user).Where("id = ?" , id).Update("is_use",status)
	return user.Id
}