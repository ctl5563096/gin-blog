package rule

import (
	ruleModel "gin-blog/models/rule"
	"gin-blog/pkg/app"
	"github.com/gin-gonic/gin"
)

// GetMenu 获取菜单
func GetMenu(c *gin.Context){
	data := make(map[string] interface{})
	res,err := ruleModel.GetBackendMenu()
	if err != nil {
		app.FailWithMessage("获取菜单失败!",1,c)
		return
	}
	var menuList []* ruleModel.MenuList
	for _,value := range res {
		if value.Pid == 0 {
			sonMenu 		:= getSonMenu(value.Id,res)
			value.ChildNode = sonMenu
			menuList = append(menuList, value)
		}
	}
	data["list"] = menuList
	app.OkWithData(data,c)
	return
}

// getSonMenu 返回子菜单栏的所有
func getSonMenu(topId int,menuList []*ruleModel.MenuList) []*ruleModel.MenuList  {
	var childList []*ruleModel.MenuList
	// 顶级菜单下面的子菜单的
	for _ , value := range menuList{
		if value.Pid ==  topId{
			childList = append(childList, value)
		}
	}
	// 子菜单求孙菜单的的childNode
	for _ , item := range childList{
		sonMenu 		:= getLastMenu(item.Id,menuList)
		item.ChildNode = sonMenu
	}
	return childList
}

// getLastMenu 返回孙菜单
func getLastMenu(topId int,menuList []*ruleModel.MenuList) []*ruleModel.MenuList  {
	var childList []*ruleModel.MenuList
	// 顶级菜单下面的子菜单的
	for _ , value := range menuList{
		if value.Pid ==  topId{
			childList = append(childList, value)
		}
	}
	return childList
}