package rule

import (
	"gin-blog/pkg/util"
)

// 自定义表名
var tableNameRuleRole = "t_go_role_rule"

type RoleRelation struct {
	Rule int `json:"id"`
}

type InsertData struct {
	Rule int `json:"rule"`
	Role int `json:"role"`
}

type DeleteData struct {
	Id int `json:"id"`
}

// GetRuleByRoleId 根据角色id去获取权限
func GetRuleByRoleId(roleId int) []*RoleRelation {
	var res []*RoleRelation
	db.Table(tableNameRuleRole).Where("role = ?" ,roleId).Find(&res)
	return res
}

// BatchInsert 批量插入角色和权限
func BatchInsert(data []interface{}) {
	sql := util.GetBranchInsertSql(data,tableNameRuleRole)
	db.Exec(sql)
}

// BatchDelete 批量删除
func BatchDelete(data []interface{},role int)  {
	db.Table(tableNameRuleRole).Where("role = ? and rule in (?)",role , data).Delete(DeleteData{})
}