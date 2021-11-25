package rule

import "gin-blog/pkg/util"

// 自定义表名
var tableName = "t_go_rule"


type MenuList struct {
	Id     	  	int    	 `json:"id"`
	RuleName  	string 	 `json:"rule_name"`
	Pid       	int    	 `json:"pid"`
	Url       	string 	 `json:"url"`
	Icon 	  	string 	 `json:"icon"`
	Sort      	int    	 `json:"sort"`
	ChildNode   []*MenuList `json:"childNode"`
}

// GetBackendMenu 获取菜单
func GetBackendMenu() ([]*MenuList,error) {
	var res [] *MenuList
	db.Table(tableName).Where("is_menu = ?" ,2).Order("sort desc").Find(&res)
	if db.Error != nil {
		util.WriteLog("mysql_error",2,db.Error.Error())
		return nil,db.Error
	}
	return res,nil
}
