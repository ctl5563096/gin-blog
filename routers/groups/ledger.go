package groups

import (
	"gin-blog/controller/ledger"
	token "gin-blog/middleware"
	"github.com/gin-gonic/gin"
)

func LedgerBaseRouter(Router *gin.RouterGroup) {
	authRouterV1 := Router.Group("/v1/ledger/auth")
	{
		authRouterV1.POST("/login", ledger.WxLogin)
	}

	apiRouterV1 := Router.Group("/v1/ledger").Use(token.LedgerAuth())
	{
		apiRouterV1.GET("/accounts", ledger.GetAccounts)
		apiRouterV1.POST("/accounts", ledger.CreateAccount)
		apiRouterV1.GET("/transactions", ledger.GetTransactions)
		apiRouterV1.POST("/transactions", ledger.CreateTransaction)
		apiRouterV1.PUT("/transactions", ledger.UpdateTransaction)
		apiRouterV1.DELETE("/transactions", ledger.DeleteTransaction)
		apiRouterV1.GET("/summary", ledger.GetSummary)
		apiRouterV1.GET("/budget", ledger.GetBudget)
		apiRouterV1.PUT("/budget", ledger.SaveBudget)
		apiRouterV1.GET("/user/profile", ledger.GetProfile)
		apiRouterV1.PUT("/user/profile", ledger.UpdateProfile)
		apiRouterV1.POST("/user/avatar", ledger.UploadAvatar)
	}
}
