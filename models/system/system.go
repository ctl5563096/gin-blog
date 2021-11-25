package system

import (
	"fmt"
	"gin-blog/pkg/util"
)

type Params struct {
	Id 			uint `json:"id"`
	Name string `json:"name" validate:"required"`
	Code string `json:"code" validate:"required"`
	IsEnabled int `json:"is_enabled"`
	ParamName string `json:"param_name" validate:"required"`
	ParamValue int `json:"param_value" validate:"required"`
	Weight int `json:"weight"`
}

type ParamsList struct {
	Id 			uint 	`json:"id"`
	Name 		string  `json:"name"`
	Code 		string 	`json:"code"`
	ParamName 	string  `json:"param_name"`
	ParamValue 	int 	`json:"param_value"`
	Sort 	  	int 	`json:"sort"`
	Weight 		int 	`json:"weight"`
	IsEnabled   int 	`json:"is_enabled"`
}

type DelParam struct {
	Id int `json:"id"`
}

type EditStruct struct {
	Id 			 int 		`json:"id" validate:"required"`
	Name 	 	 string 	`json:"name" validate:"required"`
	Code 	 	 string 	`json:"code" validate:"required"`
	ParamName	 string 	`json:"param_name" validate:"required"`
	ParamValue 	 int 		`json:"param_value" validate:"required"`
	Weight 	 	 int 		`json:"weight"`
	IsEnabled	 int 		`json:"is_enabled" validate:"required"`
}

// 自定义表名
var tableName = "t_go_parameter"

// CreateRecord 新建记录
func CreateRecord(params *Params) bool {
	db.Table(tableName).Create(params)
	if db.Error != nil {
		return false
	}
	return true
}

// GetList 获取参数列表
func GetList(page int, pageSize int,keywords string,order string) ([]*ParamsList,error) {
	var r []*ParamsList
	dbReturn := db.Table(tableName)
	if keywords != "" {
		keywords := "%" + keywords + "%"
		dbReturn = dbReturn.Where(" name LIKE ? or param_name LIKE ? ", keywords,keywords)
	}
	dbReturn.Order("updated_at " + order).Limit(pageSize).Offset((page - 1) * pageSize).Find(&r)
	if dbReturn.Error != nil {
		fmt.Println(dbReturn.Error.Error())
		return nil,dbReturn.Error
	}
	return r,nil
}

// GetTotal 获取总数
func GetTotal(keywords string) int {
	var count int
	dbReturn := db.Table(tableName)
	if keywords != "" {
		keywords := "%" + keywords + "%"
		dbReturn = dbReturn.Where(" name LIKE ? or param_name LIKE ? ", keywords,keywords)
	}
	dbReturn.Count(&count)
	return count
}

// DeleteParam 删除参数
func DeleteParam(id int) error  {
	var r DelParam
	db.Table(tableName).Where("id = ?", id).Delete(&r)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

// UpdateDetail 修改权限
func UpdateDetail(editStruct *EditStruct)  bool{
	data := make(map[string]interface{})
	data["name"] = editStruct.Name
	data["code"] = editStruct.Code
	data["is_enabled"] = editStruct.IsEnabled
	data["param_name"] = editStruct.ParamName
	data["param_value"] = editStruct.ParamValue
	data["weight"] = editStruct.Weight
	db.Table(tableName).Select("id").Where("id = ?" ,editStruct.Id).Updates(data)
	if db.Error != nil {
		util.WriteLog("update_rule_err",4,db.Error.Error())
		return false
	}
	return true
}

// GetParamDetail 获取参数详情
func GetParamDetail(id int) ParamsList {
	var r ParamsList
	db.Table(tableName).Where("id = ?", id).First(&r)
	return r
}