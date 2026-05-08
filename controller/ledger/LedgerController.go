package ledger

import (
	tokenMiddleware "gin-blog/middleware"
	"gin-blog/models/ledger"
	"gin-blog/pkg/app"
	valid "gin-blog/vaild"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"strconv"
	"time"
)

func CreateTransaction(c *gin.Context) {
	var r ledger.CreateTransaction
	if !bindAndValidate(c, &r) {
		return
	}
	r.OwnerId = getLedgerOwnerId(c)
	account, ok := getLedgerAccount(c, r.OwnerId)
	if !ok {
		return
	}
	r.AccountId = account.Id

	record, err := ledger.CreateRecord(&r)
	if err != nil {
		app.FailWithMessage("账单保存失败", 1, c)
		return
	}
	app.OkWithData(record, c)
}

func UpdateTransaction(c *gin.Context) {
	var r ledger.UpdateTransaction
	if !bindAndValidate(c, &r) {
		return
	}
	r.OwnerId = getLedgerOwnerId(c)
	account, ok := getLedgerAccount(c, r.OwnerId)
	if !ok {
		return
	}
	r.AccountId = account.Id

	if err := ledger.UpdateRecord(&r); err != nil {
		app.FailWithMessage("账单更新失败", 1, c)
		return
	}
	app.OkWithMessage("账单更新成功", c)
}

func DeleteTransaction(c *gin.Context) {
	id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	ownerId := getLedgerOwnerId(c)
	account, ok := getLedgerAccount(c, ownerId)
	if !ok {
		return
	}
	if id <= 0 || ownerId == "" {
		app.FailWithParameter("缺少账单 id", c)
		return
	}

	if err := ledger.DeleteRecord(id, account.Id); err != nil {
		app.FailWithMessage("账单删除失败", 1, c)
		return
	}
	app.OkWithMessage("账单删除成功", c)
}

func GetTransactions(c *gin.Context) {
	ownerId := getLedgerOwnerId(c)
	if ownerId == "" {
		app.MissToken(c)
		return
	}
	account, ok := getLedgerAccount(c, ownerId)
	if !ok {
		return
	}

	list, err := ledger.GetList(account.Id, c.DefaultQuery("month", ""), c.DefaultQuery("type", "all"))
	if err != nil {
		app.FailWithMessage("账单列表获取失败", 1, c)
		return
	}
	app.OkWithData(list, c)
}

func GetSummary(c *gin.Context) {
	ownerId := getLedgerOwnerId(c)
	if ownerId == "" {
		app.MissToken(c)
		return
	}
	account, ok := getLedgerAccount(c, ownerId)
	if !ok {
		return
	}

	month := c.DefaultQuery("month", time.Now().Format("2006-01"))
	today := c.DefaultQuery("today", time.Now().Format("2006-01-02"))
	summary, err := ledger.GetSummary(account.Id, month, today)
	if err != nil {
		app.FailWithMessage("账单汇总获取失败", 1, c)
		return
	}
	app.OkWithData(summary, c)
}

func GetBudget(c *gin.Context) {
	ownerId := getLedgerOwnerId(c)
	if ownerId == "" {
		app.MissToken(c)
		return
	}
	account, ok := getLedgerAccount(c, ownerId)
	if !ok {
		return
	}

	budget, err := ledger.GetBudget(account.Id, ownerId)
	if err != nil {
		app.FailWithMessage("预算获取失败", 1, c)
		return
	}
	app.OkWithData(budget, c)
}

func SaveBudget(c *gin.Context) {
	var r ledger.SaveBudget
	if !bindAndValidate(c, &r) {
		return
	}
	r.OwnerId = getLedgerOwnerId(c)
	account, ok := getLedgerAccount(c, r.OwnerId)
	if !ok {
		return
	}
	r.AccountId = account.Id

	if err := ledger.SaveBudgetRecord(&r); err != nil {
		app.FailWithMessage("预算保存失败", 1, c)
		return
	}
	app.OkWithMessage("预算保存成功", c)
}

func GetAccounts(c *gin.Context) {
	ownerId := getLedgerOwnerId(c)
	if ownerId == "" {
		app.MissToken(c)
		return
	}
	list, err := ledger.GetUserAccounts(ownerId)
	if err != nil {
		app.FailWithMessage("账本列表获取失败", 1, c)
		return
	}
	app.OkWithData(list, c)
}

func CreateAccount(c *gin.Context) {
	var r ledger.CreateAccount
	if !bindAndValidate(c, &r) {
		return
	}
	ownerId := getLedgerOwnerId(c)
	if ownerId == "" {
		app.MissToken(c)
		return
	}
	account, err := ledger.CreateAccountRecord(ownerId, &r)
	if err != nil {
		app.FailWithMessage("账本创建失败", 1, c)
		return
	}
	app.OkWithData(account, c)
}

func getLedgerOwnerId(c *gin.Context) string {
	ownerId, ok := c.Get(tokenMiddleware.LedgerOwnerContextKey)
	if !ok {
		return ""
	}
	ownerIdText, ok := ownerId.(string)
	if !ok {
		return ""
	}
	return ownerIdText
}

func getLedgerAccount(c *gin.Context, ownerId string) (ledger.AccountDTO, bool) {
	if ownerId == "" {
		app.MissToken(c)
		return ledger.AccountDTO{}, false
	}
	accountId, _ := strconv.Atoi(c.GetHeader("X-Ledger-Account-Id"))
	if accountId <= 0 {
		accountId, _ = strconv.Atoi(c.DefaultQuery("accountId", "0"))
	}
	if accountId <= 0 {
		app.FailWithMessage("请选择账本类型", 4, c)
		return ledger.AccountDTO{}, false
	}
	account, err := ledger.GetAccountForUser(ownerId, accountId)
	if err != nil {
		app.FailWithMessage("账本不存在或无权限", 4, c)
		return ledger.AccountDTO{}, false
	}
	return account, true
}

func bindAndValidate(c *gin.Context, params interface{}) bool {
	if err := c.ShouldBindBodyWith(params, binding.JSON); err != nil {
		app.FailWithParameter("请求参数格式错误", c)
		return false
	}

	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		switch err.(type) {
		case validator.ValidationErrors:
			errorMap := valid.Translate(err)
			for _, errors := range errorMap {
				for _, message := range errors {
					app.FailWithMessage(message, 4, c)
					return false
				}
			}
		default:
			app.FailWithMessage("参数校验失败", 1, c)
			return false
		}
		return false
	}
	return true
}
