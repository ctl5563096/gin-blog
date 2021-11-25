package token

import (
	//"fmt"
	"github.com/gin-gonic/gin"
)

// ValidateRule 权限验证 将在business里面获取role 然后获取该角色对应的权限进行验证
func ValidateRule() gin.HandlerFunc {
	return func(c *gin.Context) {
		//strArr := strings.Split(c.HandlerNames()[len(c.HandlerNames()) - 1], `/`)
		//routerApi := strArr[len(strArr) - 1]
		//ruleArr := strings.Split(routerApi, `.`)
		// 这里查表检查下是否有具有这个权限
	}
}
