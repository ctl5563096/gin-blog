package rule

import (
	"gin-blog/models/rule"
	"gin-blog/pkg/app"
	"gin-blog/pkg/e"
	"gin-blog/pkg/util"
	valid "gin-blog/vaild"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strconv"
)

type Rule struct {
	RuleArr []interface{} `json:"changeRuleArr"`
}

// AddRule 新增规则
func AddRule(c *gin.Context)  {
	var v rule.AddRuleStruct
	var errStr string
	var errorMap map[string][]string
	err := c.ShouldBind(&v)
	validate := validator.New()
	err = validate.Struct(v)
	if err != nil {
		switch err.(type) {
		case validator.ValidationErrors:
			errorMap = valid.Translate(err)
			//循环遍历Map 只返回第一个错误信息
			for _,v:= range errorMap{
				for _,z := range v{
					util.WriteLog("rule_business_error",4,z)
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

	_,_ = rule.CreatRule(&v)
	app.OK(c)
	return
}

// EditRuleById 根据权限id修改权限详情
func EditRuleById(c *gin.Context)  {
	var v rule.EditRule
	var errStr string
	var errorMap map[string][]string
	err := c.ShouldBind(&v)
	validate := validator.New()
	err = validate.Struct(v)
	if err != nil {
		switch err.(type) {
		case validator.ValidationErrors:
			errorMap = valid.Translate(err)
			//循环遍历Map 只返回第一个错误信息
			for _,v:= range errorMap{
				for _,z := range v{
					util.WriteLog("rule_business_error",4,z)
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

	result  	:=	make(map[string] interface{})
	res := rule.UpdateRule(&v)
	result["res"] = res
	app.OkWithData(result,c)
	return
}

// GetRuleDetailById 根据权限id获取权限详情
func GetRuleDetailById(c *gin.Context)  {
	var ruleId, _ = 	strconv.Atoi(c.DefaultQuery("id","0"))
	if ruleId <= 0 {
		app.FailWithMessage(e.GetMsg(e.MISS_PARAMS),e.MISS_PARAMS,c)
		return
	}
	r := rule.GetRuleById(ruleId)
	app.OkWithData(r,c)
	return
}

// GetRuleByRoleId 根据角色id获取角色权限
func GetRuleByRoleId(c *gin.Context)  {
	var role, _ = 	strconv.Atoi(c.DefaultQuery("id","0"))
	if role <= 0 {
		app.FailWithMessage(e.GetMsg(e.MISS_PARAMS),e.MISS_PARAMS,c)
		return
	}
	res  		:= 	rule.GetRuleByRoleId(role)
	app.OkWithData(res,c)
}

// GetAllRole 获取所有角色信息
func GetAllRole(c *gin.Context)  {
	res  		   := 	rule.GetAllRole()
	app.OkWithData(res,c)
	return
}

// GetAllRule 获取所有的权限
func GetAllRule(c *gin.Context)  {
	res  		   := 	rule.GetAllRule()
	// 开始组装树形菜单
	var dataList []*rule.AllRule

	for _,value := range res {
		if value.Pid == 0 {
			sonMenu 		:= getSonRuleMenu(value.Id,res)
			value.ChildNode = sonMenu
			dataList = append(dataList, value)
		}
	}

	app.OkWithData(dataList, c)
	return
}

// ChangeRoleRule 修改角色的权限
func ChangeRoleRule(c *gin.Context)  {
	var role, _ = 	strconv.Atoi(c.DefaultQuery("id","0"))
	// 获取角色修改的权限组 在下面进行比较
	var r Rule
	err := c.ShouldBind(&r)
	if err != nil {
		app.Fail(c)
		return
	}
	if role <= 0 {
		app.FailWithMessage(e.GetMsg(e.MISS_PARAMS),e.MISS_PARAMS,c)
		return
	}
	if role == 1 {
		app.FailWithMessage("超级管理员不允许修改权限，默认拥有全部权限",e.ERROR,c)
		return
	}
	roleRuleRs := rule.GetRuleByRoleId(role)

	// 定义整型不定长切片
	var  oldRule []interface{}
	var  newRule []interface{}
	for _,value := range roleRuleRs{
		oldRule = append(oldRule,value.Rule)
	}
	for _,value := range r.RuleArr{
		newRule = append(newRule,int(value.(float64)))
	}
	oleRuleArr, _ := util.ArrayDiff(newRule, oldRule)
	addRuleArr, _ := util.ArrayDiff(oldRule, newRule)

	if len(oleRuleArr) == 0 &&  len(addRuleArr)  ==0 {
		app.FailWithMessage("请确认你有更改权限的权利！",e.ERROR,c)
		return
	}

	// 新增权限
	if len(addRuleArr) != 0 {
		var insertData []interface{}
		// 批量插入
		for _,v := range addRuleArr {
			var item rule.InsertData
			item.Rule = v.(int)
			item.Role = role
			insertData = append(insertData, item)
		}
		rule.BatchInsert(insertData)
	}

	//删除的权限
	if len(oleRuleArr) != 0 {
		rule.BatchDelete(oleRuleArr,role)
	}

	app.OK(c)
	return
}

// DelRule 删除权限
func DelRule(c *gin.Context)  {
	var ruleRole, _ = 	strconv.Atoi(c.DefaultQuery("id","0"))
	if ruleRole <= 0 {
		app.Fail(c)
		return
	}
	err := rule.BatchDeleteRule(ruleRole)
	if err != nil{
		app.Fail(c)
		return
	}

	app.OK(c)
	return
}

// getSonRuleMenu 获取子权限菜单
func getSonRuleMenu(topId int ,ruleList []*rule.AllRule) []*rule.AllRule {
	var childList []*rule.AllRule
	// 顶级菜单下面的子菜单的
	for _ , value := range ruleList{
		if value.Pid ==  topId{
			childList = append(childList, value)
		}
	}

	// 子菜单求孙菜单的的childNode
	for _ , item := range childList{
		sonMenu 		:= getLastRuleList(item.Id,ruleList)
		item.ChildNode = sonMenu
	}
	return childList
}

// 获取最后一个层级的子菜单
func getLastRuleList(sonId int, ruleList []*rule.AllRule ) []*rule.AllRule {
	var childList []*rule.AllRule
	// 顶级菜单下面的子菜单的
	for _ , value := range ruleList{
		if value.Pid ==  sonId{
			childList = append(childList, value)
		}
	}
	return childList
}