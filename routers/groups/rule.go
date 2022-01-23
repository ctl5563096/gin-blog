package groups

import (
	"gin-blog/controller/rule"
	token "gin-blog/middleware"
	"github.com/gin-gonic/gin"
)

// RuleBaseRouter /** 文章基本接口 **/
func RuleBaseRouter(Router *gin.RouterGroup) {
	// v1版接口
	apiRouterV1 := Router.Group("/v1/rule").Use(token.BeforeBusiness())
	{
		// 获取侧边栏菜单
		apiRouterV1.GET("menu", rule.GetMenu)
		// 根据角色获取角色的权限
		apiRouterV1.GET("getRuleByRole", rule.GetRuleByRoleId)
		// 获取所有角色
		apiRouterV1.GET("role", rule.GetAllRole).Use(token.ValidateRule())
		// 获取所有的权限
		apiRouterV1.GET("rule", rule.GetAllRule)
		// 修改角色的权限
		apiRouterV1.PUT("rule", rule.ChangeRoleRule)
		// 删除权限
		apiRouterV1.DELETE("rule", rule.DelRule)
		// 新增权限
		apiRouterV1.POST("rule", rule.AddRule)
		// 根据权限ID获取权限详情
		apiRouterV1.GET("ruleDetail", rule.GetRuleDetailById)
		// 根据权限ID修改权限详情
		apiRouterV1.PUT("ruleDetail", rule.EditRuleById)
	}
}