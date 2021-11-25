package rule

import (
	"encoding/json"
	"gin-blog/pkg/util"
)

// 自定义表名
var tableNameRule = "t_go_rule"

type AllRule struct {
	Id  		int 		`json:"id"`
	Pid 		int 		`json:"pid"`
	Label 		string  	`json:"label"`
	ChildNode 	[]*AllRule	`json:"childNode"`
}

type DeleteRule struct {
	First 		int `json:"first"`
	Second 		int `json:"second"`
	Third 		int `json:"third"`
}

type AddRuleStruct struct {
	ID 		 	 int		`json:"id"`
	RuleName 	 string 	`json:"rule_name" validate:"required,lt=11,gt=1"`
	IsMenu 	 	 int 		`json:"is_menu"`
	Icon 	 	 string 	`json:"icon"`
	Pid	 	 	 int 		`json:"pid"`
	Controller	 string 	`json:"controller"`
	Action   	 string 	`json:"action"`
	Sort   	  	 int		`json:"sort"`
	Url 		 string  	`json:"url"`
}

type DetailRule struct {
	Id 			 int 		`json:"id"`
	RuleName 	 string 	`json:"rule_name"`
	IsMenu 	 	 int 		`json:"is_menu"`
	Icon 	 	 string 	`json:"icon"`
	Pid	 	 	 int 		`json:"pid"`
	Controller	 string 	`json:"controller"`
	Action   	 string 	`json:"action"`
	Sort   	  	 int		`json:"sort"`
	Url 		 string  	`json:"url"`
	CreatedAt    string		`json:"created_at"`
}

type EditRule struct {
	Id 			 int 		`json:"id" validate:"required"`
	RuleName 	 string 	`json:"rule_name" validate:"required"`
	IsMenu 	 	 int 		`json:"is_menu"`
	Icon 	 	 string 	`json:"icon"`
	Pid	 	 	 int 		`json:"pid"`
	Controller	 string 	`json:"controller"`
	Action   	 string 	`json:"action"`
	Sort   	  	 int		`json:"sort"`
	Url 		 string  	`json:"url"`
	CreatedAt    string		`json:"created_at"`
}

// GetAllRule 获取所有的权限
func GetAllRule() []*AllRule {
	var rs []*AllRule
	db.Table(tableNameRule).Select("id,pid,rule_name as label").Find(&rs)
	return rs
}

// BatchDeleteRule 删除权限以及旗下所有的子权限
func BatchDeleteRule(ruleId int) error {
	var res []DeleteRule
	var sql string
	sql = "select IFNULL(c1.id ,0) as first ,IFNULL(c2.id ,0) as second ,IFNULL(c3.id ,0) as third from t_go_rule as c1 " +
		"left join t_go_rule as c2 on c1.id = c2.pid " +
		"left join t_go_rule as c3 on c2.id = c3.pid " +
		" where c1.id = ?"
	db.Raw(sql,ruleId).Scan(&res)
	var deleteArr = make([]int, 0, len(res) * 3)
	for _,v := range res{
		if v.First > 0 {
			deleteArr = append(deleteArr,v.First)
		}
		if v.Second > 0 {
			deleteArr = append(deleteArr,v.Second)
		}
		if v.Third > 0 {
			deleteArr = append(deleteArr,v.Third)
		}
	}
	newArr := util.RemoveDuplicate(deleteArr)
	db.Table(tableNameRule).Where("id in (?)",newArr).Delete(DeleteData{})
	if db.Error != nil {
		return db.Error
	}
	return nil
}

// CreatRule 创建新的规则
func CreatRule(ruleStruct *AddRuleStruct)  (returnId int ,err error){
	err = db.Table(tableNameRule).Create(ruleStruct).Error
	if err != nil {
		return 0,err
	}
	m := make(map[string]interface{})
	j,_ := json.Marshal(ruleStruct)
	err = json.Unmarshal(j, &m)
	if err != nil {
		return 0, err
	}
	recordId := int(m["id"].(float64))
	return recordId,err
}

// GetRuleById 根据rule权限获取权限详情
func GetRuleById(id int) []*DetailRule {
	var r []*DetailRule
	db.Table(tableNameRule).Where("id = ? " , id).Find(&r)
	return r
}

// UpdateRule 修改权限
func UpdateRule(editStruct *EditRule)  bool{
	data := make(map[string]interface{})
	data["rule_name"] 	= editStruct.RuleName
	data["is_menu"] 	= editStruct.RuleName
	data["icon"] 		= editStruct.Icon
	data["controller"] 	= editStruct.Controller
	data["action"] 		= editStruct.Action
	data["sort"] 		= editStruct.Sort
	db.Table(tableNameRule).Select("id").Where("id = ?" ,editStruct.Id).Updates(data)
	if db.Error != nil {
		util.WriteLog("update_rule_err",4,db.Error.Error())
		return false
	}
	return true
}