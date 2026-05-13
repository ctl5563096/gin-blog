package ledger

import (
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

const (
	TransactionTableName         = "t_family_transaction"
	BudgetTableName              = "t_family_budget"
	UserTableName                = "t_family_user"
	AccountTableName             = "t_family_account"
	MemberTableName              = "t_family_account_member"
	SavingTableName              = "t_family_saving"
	SavingFlowTableName          = "t_family_saving_flow"
	CategoryConfigTableName      = "t_family_transaction_category_config"
	MonthlySavingConfigTableName = "t_family_monthly_saving_config"

	AccountTypePersonal = "personal"
	AccountTypeFamily   = "family"
	BudgetScopeAccount  = "account"
	BudgetScopeMember   = "member"
	SavingScopeFamily   = "family"
	SavingScopePersonal = "personal"
)

type User struct {
	Id        int       `json:"id" gorm:"primary_key"`
	OpenId    string    `json:"openid" gorm:"column:openid"`
	UnionId   string    `json:"unionid" gorm:"column:unionid"`
	NickName  string    `json:"nickName" gorm:"column:nick_name"`
	AvatarUrl string    `json:"avatarUrl" gorm:"column:avatar_url"`
	LastLogin time.Time `json:"lastLogin" gorm:"column:last_login"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

type Transaction struct {
	Id                 int       `json:"id" gorm:"primary_key"`
	AccountId          int       `json:"accountId" gorm:"column:account_id"`
	OwnerId            string    `json:"ownerId" gorm:"column:owner_id"`
	Type               string    `json:"type"`
	Amount             float64   `json:"amount"`
	Category           string    `json:"category"`
	Description        string    `json:"description"`
	Note               string    `json:"note"`
	Date               string    `json:"date"`
	PaymentMethod      string    `json:"paymentMethod" gorm:"column:payment_method"`
	MemberName         string    `json:"memberName" gorm:"column:member_name"`
	SavingBalanceAfter float64   `json:"savingBalanceAfter" gorm:"-"`
	IsDelete           int       `json:"isDelete" gorm:"column:is_delete"`
	CreatedAt          time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt          time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

type SavingFlow struct {
	Id            int       `json:"id" gorm:"primary_key"`
	AccountId     int       `json:"accountId" gorm:"column:account_id"`
	TransactionId int       `json:"transactionId" gorm:"column:transaction_id"`
	Scope         string    `json:"scope"`
	UserOpenId    string    `json:"userOpenid" gorm:"column:user_openid"`
	Action        string    `json:"action"`
	DeltaAmount   float64   `json:"deltaAmount" gorm:"column:delta_amount"`
	BalanceAfter  float64   `json:"balanceAfter" gorm:"column:balance_after"`
	RecordType    string    `json:"recordType" gorm:"column:record_type"`
	CategoryName  string    `json:"categoryName" gorm:"column:category_name"`
	Description   string    `json:"description"`
	IsDelete      int       `json:"isDelete" gorm:"column:is_delete"`
	CreatedAt     time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

type TransactionCategoryConfig struct {
	Id            int       `json:"id" gorm:"primary_key"`
	AccountId     int       `json:"accountId" gorm:"column:account_id"`
	CreatorOpenId string    `json:"creatorOpenid" gorm:"column:creator_openid"`
	CategoryName  string    `json:"categoryName" gorm:"column:category_name"`
	Amount        float64   `json:"amount"`
	Type          string    `json:"type"`
	Source        string    `json:"source"`
	IsDelete      int       `json:"isDelete" gorm:"column:is_delete"`
	CreatedAt     time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

type MonthlySavingConfig struct {
	Id            int       `json:"id" gorm:"primary_key"`
	AccountId     int       `json:"accountId" gorm:"column:account_id"`
	CreatorOpenId string    `json:"creatorOpenid" gorm:"column:creator_openid"`
	CategoryName  string    `json:"categoryName" gorm:"column:category_name"`
	Amount        float64   `json:"amount"`
	Type          string    `json:"type"`
	Source        string    `json:"source"`
	IsDelete      int       `json:"isDelete" gorm:"column:is_delete"`
	CreatedAt     time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

type CreateTransaction struct {
	AccountId                 int     `json:"accountId" validate:"gte=0"`
	OwnerId                   string  `json:"ownerId" validate:"max=64"`
	Type                      string  `json:"type" validate:"required,oneof=expense income"`
	Amount                    float64 `json:"amount" validate:"required,gt=0"`
	Category                  string  `json:"category" validate:"required,max=32"`
	Description               string  `json:"description" validate:"max=80"`
	Note                      string  `json:"note" validate:"max=80"`
	Date                      string  `json:"date" validate:"required,len=10"`
	PaymentMethod             string  `json:"paymentMethod" validate:"required,max=32"`
	MemberName                string  `json:"memberName" validate:"required,max=32"`
	SaveAsCategoryConfig      bool    `json:"saveAsCategoryConfig"`
	SaveAsMonthlySavingConfig bool    `json:"saveAsMonthlySavingConfig"`
	ImportedSavingBalance     float64 `json:"importedSavingBalance"`
	HasImportedSavingBalance  bool    `json:"hasImportedSavingBalance"`
}

type ImportTransactionResult struct {
	Total         int     `json:"total"`
	Success       int     `json:"success"`
	Skipped       int     `json:"skipped"`
	BalanceSynced bool    `json:"balanceSynced"`
	FinalBalance  float64 `json:"finalBalance"`
}

type UpdateTransaction struct {
	Id            int     `json:"id" validate:"required,gt=0"`
	AccountId     int     `json:"accountId" validate:"gte=0"`
	OwnerId       string  `json:"ownerId" validate:"max=64"`
	Type          string  `json:"type" validate:"required,oneof=expense income"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	Category      string  `json:"category" validate:"required,max=32"`
	Description   string  `json:"description" validate:"max=80"`
	Note          string  `json:"note" validate:"max=80"`
	Date          string  `json:"date" validate:"required,len=10"`
	PaymentMethod string  `json:"paymentMethod" validate:"required,max=32"`
	MemberName    string  `json:"memberName" validate:"required,max=32"`
}

type Budget struct {
	Id        int       `json:"id" gorm:"primary_key"`
	OwnerId   string    `json:"ownerId" gorm:"column:owner_id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

type SaveBudget struct {
	AccountId int     `json:"accountId" validate:"gte=0"`
	OwnerId   string  `json:"ownerId" validate:"max=64"`
	Scope     string  `json:"scope" validate:"omitempty,oneof=account member"`
	Amount    float64 `json:"amount" validate:"gte=0"`
}

type UpdateProfile struct {
	OwnerId   string `json:"ownerId" validate:"max=64"`
	NickName  string `json:"nickName" validate:"max=64"`
	AvatarUrl string `json:"avatarUrl" validate:"max=512"`
}

type Summary struct {
	Income       float64 `json:"income"`
	Expense      float64 `json:"expense"`
	Balance      float64 `json:"balance"`
	TodayExpense float64 `json:"todayExpense"`
	Count        int     `json:"count"`
}

type Account struct {
	Id           int       `json:"id" gorm:"primary_key"`
	AccountType  string    `json:"accountType" gorm:"column:account_type"`
	Name         string    `json:"name"`
	OwnerOpenId  string    `json:"ownerOpenid" gorm:"column:owner_openid"`
	BudgetAmount float64   `json:"budgetAmount" gorm:"column:budget_amount"`
	IsDelete     int       `json:"isDelete" gorm:"column:is_delete"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

type AccountMember struct {
	Id                   int       `json:"id" gorm:"primary_key"`
	AccountId            int       `json:"accountId" gorm:"column:account_id"`
	UserOpenId           string    `json:"userOpenid" gorm:"column:user_openid"`
	Role                 string    `json:"role"`
	PersonalBudgetAmount float64   `json:"personalBudgetAmount" gorm:"column:personal_budget_amount"`
	DisplayName          string    `json:"displayName" gorm:"column:display_name"`
	Status               int       `json:"status"`
	CreatedAt            time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt            time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

type AccountDTO struct {
	Id                   int     `json:"id" gorm:"column:id"`
	AccountType          string  `json:"accountType" gorm:"column:account_type"`
	Name                 string  `json:"name" gorm:"column:name"`
	OwnerOpenId          string  `json:"ownerOpenid" gorm:"column:owner_openid"`
	BudgetAmount         float64 `json:"budgetAmount" gorm:"column:budget_amount"`
	Role                 string  `json:"role" gorm:"column:role"`
	PersonalBudgetAmount float64 `json:"personalBudgetAmount" gorm:"column:personal_budget_amount"`
}

type CreateAccount struct {
	AccountType                string  `json:"accountType" validate:"required,oneof=personal family"`
	Name                       string  `json:"name" validate:"max=32"`
	BudgetAmount               float64 `json:"budgetAmount" validate:"gte=0"`
	PersonalBudgetAmount       float64 `json:"personalBudgetAmount" validate:"gte=0"`
	InitialSavingAmount        float64 `json:"initialSavingAmount" validate:"gte=0"`
	MonthlySavingIncomeAmount  float64 `json:"monthlySavingIncomeAmount" validate:"gte=0"`
	MonthlySavingExpenseAmount float64 `json:"monthlySavingExpenseAmount" validate:"gte=0"`
}

type BudgetInfo struct {
	AccountId            int     `json:"accountId"`
	AccountType          string  `json:"accountType"`
	Amount               float64 `json:"amount"`
	PersonalAmount       float64 `json:"personalAmount"`
	PersonalBudgetAmount float64 `json:"personalBudgetAmount"`
}

type Saving struct {
	Id            int       `json:"id" gorm:"primary_key"`
	AccountId     int       `json:"accountId" gorm:"column:account_id"`
	Scope         string    `json:"scope"`
	UserOpenId    string    `json:"userOpenid" gorm:"column:user_openid"`
	TargetAmount  float64   `json:"targetAmount" gorm:"column:target_amount"`
	CurrentAmount float64   `json:"currentAmount" gorm:"column:current_amount"`
	Note          string    `json:"note"`
	IsDelete      int       `json:"isDelete" gorm:"column:is_delete"`
	CreatedAt     time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

type SavingStatus struct {
	Scope         string  `json:"scope"`
	Title         string  `json:"title"`
	TargetAmount  float64 `json:"targetAmount"`
	CurrentAmount float64 `json:"currentAmount"`
	Percent       int     `json:"percent"`
	Note          string  `json:"note"`
}

type SaveSaving struct {
	AccountId     int     `json:"accountId" validate:"gte=0"`
	OwnerId       string  `json:"ownerId" validate:"max=64"`
	Scope         string  `json:"scope" validate:"required,oneof=family personal"`
	TargetAmount  float64 `json:"targetAmount" validate:"gte=0"`
	CurrentAmount float64 `json:"currentAmount" validate:"gte=0"`
	Note          string  `json:"note" validate:"max=80"`
}

type TransactionCategoryConfigDTO struct {
	Id            int     `json:"id"`
	CategoryName  string  `json:"categoryName"`
	Amount        float64 `json:"amount"`
	Type          string  `json:"type"`
	Source        string  `json:"source"`
	AccountId     int     `json:"accountId"`
	CreatorOpenId string  `json:"creatorOpenid"`
}

type MonthlySavingConfigDTO struct {
	Id            int     `json:"id"`
	CategoryName  string  `json:"categoryName"`
	Amount        float64 `json:"amount"`
	Type          string  `json:"type"`
	Source        string  `json:"source"`
	AccountId     int     `json:"accountId"`
	CreatorOpenId string  `json:"creatorOpenid"`
}

type UpdateTransactionCategoryConfig struct {
	Id           int     `json:"id" validate:"required,gt=0"`
	CategoryName string  `json:"categoryName" validate:"required,max=32"`
	Amount       float64 `json:"amount" validate:"required,gte=0"`
	Type         string  `json:"type" validate:"required,oneof=expense income"`
	Source       string  `json:"source" validate:"max=80"`
}

type UpdateMonthlySavingConfig struct {
	Id           int     `json:"id" validate:"required,gt=0"`
	CategoryName string  `json:"categoryName" validate:"required,max=32"`
	Amount       float64 `json:"amount" validate:"required,gte=0"`
	Type         string  `json:"type" validate:"required,oneof=expense income"`
	Source       string  `json:"source" validate:"max=80"`
}

func GetUser(openid string) (User, error) {
	var user User
	err := db.Table(UserTableName).Where("openid = ?", openid).First(&user).Error
	return user, err
}

func UpdateUserProfile(params *UpdateProfile) (User, error) {
	values := map[string]interface{}{
		"nick_name":  params.NickName,
		"avatar_url": params.AvatarUrl,
		"updated_at": time.Now(),
	}
	err := db.Table(UserTableName).Where("openid = ?", params.OwnerId).Updates(values).Error
	if err != nil {
		return User{}, err
	}
	return GetUser(params.OwnerId)
}

func UpsertUser(openid string, unionid string) (User, error) {
	var user User
	err := db.Table(UserTableName).Where("openid = ?", openid).First(&user).Error
	now := time.Now()
	if err == nil {
		err = db.Table(UserTableName).Where("openid = ?", openid).Updates(map[string]interface{}{
			"unionid":    unionid,
			"last_login": now,
			"updated_at": now,
		}).Error
		if err != nil {
			return user, err
		}
		user.UnionId = unionid
		user.LastLogin = now
		user.UpdatedAt = now
		return user, nil
	}
	if err != gorm.ErrRecordNotFound {
		return user, err
	}

	user = User{
		OpenId:    openid,
		UnionId:   unionid,
		LastLogin: now,
	}
	err = db.Table(UserTableName).Create(&user).Error
	return user, err
}

func GetUserAccounts(openid string) ([]AccountDTO, error) {
	var list []AccountDTO
	err := db.Table(AccountTableName+" a").
		Select("a.id, a.account_type, a.name, a.owner_openid, a.budget_amount, m.role, m.personal_budget_amount").
		Joins("INNER JOIN "+MemberTableName+" m ON m.account_id = a.id").
		Where("m.user_openid = ? AND m.status = ? AND a.is_delete = ?", openid, 1, 1).
		Order("a.id ASC").
		Scan(&list).Error
	return list, err
}

func GetAccountForUser(openid string, accountId int) (AccountDTO, error) {
	var account AccountDTO
	err := db.Table(AccountTableName+" a").
		Select("a.id, a.account_type, a.name, a.owner_openid, a.budget_amount, m.role, m.personal_budget_amount").
		Joins("INNER JOIN "+MemberTableName+" m ON m.account_id = a.id").
		Where("a.id = ? AND m.user_openid = ? AND m.status = ? AND a.is_delete = ?", accountId, openid, 1, 1).
		Scan(&account).Error
	if err == nil && account.Id == 0 {
		err = gorm.ErrRecordNotFound
	}
	return account, err
}

func CreateAccountRecord(openid string, params *CreateAccount) (AccountDTO, error) {
	user, _ := GetUser(openid)
	name := params.Name
	if name == "" {
		if params.AccountType == AccountTypeFamily {
			name = "家庭账本"
		} else {
			name = "个人账本"
		}
	}

	account := Account{
		AccountType:  params.AccountType,
		Name:         name,
		OwnerOpenId:  openid,
		BudgetAmount: params.BudgetAmount,
		IsDelete:     1,
	}
	if params.AccountType == AccountTypePersonal {
		account.BudgetAmount = 0
	}
	tx := db.Begin()
	if err := tx.Table(AccountTableName).Create(&account).Error; err != nil {
		tx.Rollback()
		return AccountDTO{}, err
	}

	displayName := user.NickName
	if displayName == "" {
		displayName = "微信用户"
	}
	member := AccountMember{
		AccountId:            account.Id,
		UserOpenId:           openid,
		Role:                 "owner",
		PersonalBudgetAmount: params.PersonalBudgetAmount,
		DisplayName:          displayName,
		Status:               1,
	}
	if err := tx.Table(MemberTableName).Create(&member).Error; err != nil {
		tx.Rollback()
		return AccountDTO{}, err
	}

	saving := buildInitialSaving(account.Id, account.AccountType, openid, params.InitialSavingAmount)
	if err := tx.Table(SavingTableName).Create(&saving).Error; err != nil {
		tx.Rollback()
		return AccountDTO{}, err
	}
	if params.MonthlySavingIncomeAmount > 0 {
		if err := createMonthlySavingConfig(tx, account.Id, openid, "income", params.MonthlySavingIncomeAmount, "每月固定储蓄收入", "账本初始化"); err != nil {
			tx.Rollback()
			return AccountDTO{}, err
		}
	}
	if params.MonthlySavingExpenseAmount > 0 {
		if err := createMonthlySavingConfig(tx, account.Id, openid, "expense", params.MonthlySavingExpenseAmount, "每月固定储蓄支出", "账本初始化"); err != nil {
			tx.Rollback()
			return AccountDTO{}, err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return AccountDTO{}, err
	}
	return GetAccountForUser(openid, account.Id)
}

func CreateRecord(params *CreateTransaction) (*Transaction, error) {
	record := &Transaction{
		AccountId:     params.AccountId,
		OwnerId:       params.OwnerId,
		Type:          params.Type,
		Amount:        params.Amount,
		Category:      params.Category,
		Description:   params.Description,
		Note:          params.Note,
		Date:          params.Date,
		PaymentMethod: params.PaymentMethod,
		MemberName:    params.MemberName,
		IsDelete:      1,
	}
	tx := db.Begin()
	if err := tx.Table(TransactionTableName).Create(record).Error; err != nil {
		tx.Rollback()
		return record, err
	}
	balanceAfter, err := applySavingDelta(tx, params.AccountId, params.OwnerId, signedTransactionAmount(params.Type, params.Amount), record.Id, "create", params.Type, params.Category, params.Description)
	if err != nil {
		tx.Rollback()
		return record, err
	}
	record.SavingBalanceAfter = balanceAfter
	if params.SaveAsCategoryConfig {
		if err := upsertTransactionCategoryConfig(tx, params); err != nil {
			tx.Rollback()
			return record, err
		}
	}
	if params.SaveAsMonthlySavingConfig {
		if err := upsertMonthlySavingConfig(tx, params); err != nil {
			tx.Rollback()
			return record, err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return record, err
	}
	return record, nil
}

func CreateRecords(params []CreateTransaction) (int, error) {
	tx := db.Begin()
	success := 0
	for _, item := range params {
		record := &Transaction{
			AccountId:     item.AccountId,
			OwnerId:       item.OwnerId,
			Type:          item.Type,
			Amount:        item.Amount,
			Category:      item.Category,
			Description:   item.Description,
			Note:          item.Note,
			Date:          item.Date,
			PaymentMethod: item.PaymentMethod,
			MemberName:    item.MemberName,
			IsDelete:      1,
		}
		if err := tx.Table(TransactionTableName).Create(record).Error; err != nil {
			tx.Rollback()
			return success, err
		}
		if item.HasImportedSavingBalance {
			_, err := applySavingAbsoluteBalance(tx, item.AccountId, item.OwnerId, item.ImportedSavingBalance, record.Id, "import", item.Type, item.Category, item.Description)
			if err != nil {
				tx.Rollback()
				return success, err
			}
		} else {
			_, err := applySavingDelta(tx, item.AccountId, item.OwnerId, signedTransactionAmount(item.Type, item.Amount), record.Id, "import", item.Type, item.Category, item.Description)
			if err != nil {
				tx.Rollback()
				return success, err
			}
		}
		if item.SaveAsCategoryConfig {
			if err := upsertTransactionCategoryConfig(tx, &item); err != nil {
				tx.Rollback()
				return success, err
			}
		}
		if item.SaveAsMonthlySavingConfig {
			if err := upsertMonthlySavingConfig(tx, &item); err != nil {
				tx.Rollback()
				return success, err
			}
		}
		success++
	}
	if err := tx.Commit().Error; err != nil {
		return success, err
	}
	return success, nil
}

func UpdateRecord(params *UpdateTransaction) error {
	var oldRecord Transaction
	if err := db.Table(TransactionTableName).
		Where("id = ? AND account_id = ? AND is_delete = ?", params.Id, params.AccountId, 1).
		First(&oldRecord).Error; err != nil {
		return err
	}
	values := map[string]interface{}{
		"type":           params.Type,
		"amount":         params.Amount,
		"category":       params.Category,
		"description":    params.Description,
		"note":           params.Note,
		"date":           params.Date,
		"payment_method": params.PaymentMethod,
		"member_name":    params.MemberName,
		"updated_at":     time.Now(),
	}
	tx := db.Begin()
	if err := tx.Table(TransactionTableName).
		Where("id = ? AND account_id = ? AND is_delete = ?", params.Id, params.AccountId, 1).
		Updates(values).Error; err != nil {
		tx.Rollback()
		return err
	}
	oldAmount := signedTransactionAmount(oldRecord.Type, oldRecord.Amount)
	newAmount := signedTransactionAmount(params.Type, params.Amount)
	if _, err := applySavingDelta(tx, params.AccountId, params.OwnerId, newAmount-oldAmount, params.Id, "update", params.Type, params.Category, params.Description); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func DeleteRecord(id int, accountId int) error {
	var oldRecord Transaction
	if err := db.Table(TransactionTableName).
		Where("id = ? AND account_id = ? AND is_delete = ?", id, accountId, 1).
		First(&oldRecord).Error; err != nil {
		return err
	}
	tx := db.Begin()
	if err := tx.Table(TransactionTableName).
		Where("id = ? AND account_id = ? AND is_delete = ?", id, accountId, 1).
		Update("is_delete", 2).Error; err != nil {
		tx.Rollback()
		return err
	}
	if _, err := applySavingDelta(tx, accountId, oldRecord.OwnerId, -signedTransactionAmount(oldRecord.Type, oldRecord.Amount), oldRecord.Id, "delete", oldRecord.Type, oldRecord.Category, oldRecord.Description); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func GetList(accountId int, month string, recordType string) ([]Transaction, error) {
	var list []Transaction
	dbq := db.Table(TransactionTableName+" t").
		Select("t.*, sf.balance_after AS saving_balance_after").
		Joins("LEFT JOIN ("+
			"SELECT f1.transaction_id, f1.balance_after "+
			"FROM "+SavingFlowTableName+" f1 "+
			"INNER JOIN ("+
			"SELECT transaction_id, MAX(id) AS max_id FROM "+SavingFlowTableName+" WHERE is_delete = 1 GROUP BY transaction_id"+
			") f2 ON f1.id = f2.max_id"+
			") sf ON sf.transaction_id = t.id").
		Where("t.account_id = ? AND t.is_delete = ?", accountId, 1)
	if month != "" {
		dbq = dbq.Where("t.date LIKE ?", month+"%")
	}
	if recordType != "" && recordType != "all" {
		dbq = dbq.Where("t.type = ?", recordType)
	}
	err := dbq.Order("t.date desc,t.id desc").Find(&list).Error
	return list, err
}

func GetSummary(accountId int, month string, today string) (Summary, error) {
	list, err := GetList(accountId, month, "all")
	if err != nil {
		return Summary{}, err
	}

	result := Summary{Count: len(list)}
	for _, item := range list {
		if item.Type == "income" {
			result.Income += item.Amount
		}
		if item.Type == "expense" {
			result.Expense += item.Amount
			if item.Date == today {
				result.TodayExpense += item.Amount
			}
		}
	}
	result.Balance = result.Income - result.Expense
	return result, nil
}

func GetBudget(accountId int, openid string) (BudgetInfo, error) {
	account, err := GetAccountForUser(openid, accountId)
	if err != nil {
		return BudgetInfo{}, err
	}
	return BudgetInfo{
		AccountId:            account.Id,
		AccountType:          account.AccountType,
		Amount:               getAccountBudgetAmount(account),
		PersonalAmount:       account.PersonalBudgetAmount,
		PersonalBudgetAmount: account.PersonalBudgetAmount,
	}, nil
}

func getAccountBudgetAmount(account AccountDTO) float64 {
	if account.AccountType == AccountTypePersonal {
		return account.PersonalBudgetAmount
	}
	return account.BudgetAmount
}

func SaveBudgetRecord(params *SaveBudget) error {
	scope := params.Scope
	if scope == "" {
		scope = BudgetScopeAccount
	}
	values := map[string]interface{}{
		"updated_at": time.Now(),
	}
	if scope == BudgetScopeMember {
		values["personal_budget_amount"] = params.Amount
		return db.Table(MemberTableName).
			Where("account_id = ? AND user_openid = ? AND status = ?", params.AccountId, params.OwnerId, 1).
			Updates(values).Error
	}
	account, err := GetAccountForUser(params.OwnerId, params.AccountId)
	if err != nil {
		return err
	}
	if account.AccountType == AccountTypePersonal {
		return SaveBudgetRecord(&SaveBudget{
			AccountId: params.AccountId,
			OwnerId:   params.OwnerId,
			Scope:     BudgetScopeMember,
			Amount:    params.Amount,
		})
	}
	values["budget_amount"] = params.Amount
	return db.Table(AccountTableName).
		Where("id = ? AND is_delete = ?", params.AccountId, 1).
		Updates(values).Error
}

func GetSavings(account AccountDTO, openid string) ([]SavingStatus, error) {
	scope := primarySavingScope(account.AccountType)
	userOpenId := savingUserOpenId(scope, openid)
	var record Saving
	err := db.Table(SavingTableName).
		Where("account_id = ? AND scope = ? AND user_openid = ? AND is_delete = ?", account.Id, scope, userOpenId, 1).
		First(&record).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return []SavingStatus{buildSavingStatus(scope, record)}, nil
}

func SaveSavingRecord(account AccountDTO, params *SaveSaving) error {
	scope := primarySavingScope(account.AccountType)
	if params.Scope != scope {
		return gorm.ErrRecordNotFound
	}
	userOpenId := savingUserOpenId(scope, params.OwnerId)
	now := time.Now()
	values := map[string]interface{}{
		"target_amount":  params.TargetAmount,
		"current_amount": params.CurrentAmount,
		"note":           params.Note,
		"updated_at":     now,
	}

	var existed Saving
	err := db.Table(SavingTableName).
		Where("account_id = ? AND scope = ? AND user_openid = ? AND is_delete = ?", account.Id, scope, userOpenId, 1).
		First(&existed).Error
	if err == nil {
		return db.Table(SavingTableName).
			Where("id = ?", existed.Id).
			Updates(values).Error
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	record := Saving{
		AccountId:     account.Id,
		Scope:         scope,
		UserOpenId:    userOpenId,
		TargetAmount:  params.TargetAmount,
		CurrentAmount: params.CurrentAmount,
		Note:          params.Note,
		IsDelete:      1,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	return db.Table(SavingTableName).Create(&record).Error
}

func SyncSavingCurrentAmount(account AccountDTO, openid string, amount float64, note string) error {
	scope := primarySavingScope(account.AccountType)
	userOpenId := savingUserOpenId(scope, openid)
	now := time.Now()

	var existed Saving
	err := db.Table(SavingTableName).
		Where("account_id = ? AND scope = ? AND user_openid = ? AND is_delete = ?", account.Id, scope, userOpenId, 1).
		First(&existed).Error
	if err == nil {
		values := map[string]interface{}{
			"current_amount": amount,
			"updated_at":     now,
		}
		if note != "" {
			values["note"] = note
		}
		return db.Table(SavingTableName).
			Where("id = ?", existed.Id).
			Updates(values).Error
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	record := buildInitialSaving(account.Id, account.AccountType, openid, amount)
	if note != "" {
		record.Note = note
	}
	return db.Table(SavingTableName).Create(&record).Error
}

func buildInitialSaving(accountId int, accountType string, openid string, amount float64) Saving {
	scope := primarySavingScope(accountType)
	now := time.Now()
	return Saving{
		AccountId:     accountId,
		Scope:         scope,
		UserOpenId:    savingUserOpenId(scope, openid),
		TargetAmount:  0,
		CurrentAmount: amount,
		Note:          "初始储蓄",
		IsDelete:      1,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func primarySavingScope(accountType string) string {
	if accountType == AccountTypeFamily {
		return SavingScopeFamily
	}
	return SavingScopePersonal
}

func savingUserOpenId(scope string, openid string) string {
	if scope == SavingScopePersonal {
		return openid
	}
	return ""
}

func signedTransactionAmount(recordType string, amount float64) float64 {
	if recordType == "income" {
		return amount
	}
	return -amount
}

func applySavingDelta(tx *gorm.DB, accountId int, openid string, delta float64, transactionId int, action string, recordType string, categoryName string, description string) (float64, error) {
	if delta == 0 {
		return currentSavingBalance(tx, accountId, openid)
	}
	var account Account
	if err := tx.Table(AccountTableName).
		Where("id = ? AND is_delete = ?", accountId, 1).
		First(&account).Error; err != nil {
		return 0, err
	}
	scope := primarySavingScope(account.AccountType)
	userOpenId := savingUserOpenId(scope, openid)

	var saving Saving
	err := tx.Table(SavingTableName).
		Where("account_id = ? AND scope = ? AND user_openid = ? AND is_delete = ?", accountId, scope, userOpenId, 1).
		First(&saving).Error
	if err == gorm.ErrRecordNotFound {
		record := buildInitialSaving(accountId, account.AccountType, openid, 0)
		if err = tx.Table(SavingTableName).Create(&record).Error; err != nil {
			return 0, err
		}
		saving = record
	} else if err != nil {
		return 0, err
	}
	balanceAfter := saving.CurrentAmount + delta
	if err = tx.Table(SavingTableName).
		Where("id = ?", saving.Id).
		Updates(map[string]interface{}{
			"current_amount": balanceAfter,
			"updated_at":     time.Now(),
		}).Error; err != nil {
		return 0, err
	}
	flow := SavingFlow{
		AccountId:     accountId,
		TransactionId: transactionId,
		Scope:         scope,
		UserOpenId:    userOpenId,
		Action:        action,
		DeltaAmount:   delta,
		BalanceAfter:  balanceAfter,
		RecordType:    recordType,
		CategoryName:  limitString(categoryName, 32),
		Description:   limitString(description, 80),
		IsDelete:      1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err = tx.Table(SavingFlowTableName).Create(&flow).Error; err != nil {
		return 0, err
	}
	return balanceAfter, nil
}

func applySavingAbsoluteBalance(tx *gorm.DB, accountId int, openid string, balanceAfter float64, transactionId int, action string, recordType string, categoryName string, description string) (float64, error) {
	var account Account
	if err := tx.Table(AccountTableName).
		Where("id = ? AND is_delete = ?", accountId, 1).
		First(&account).Error; err != nil {
		return 0, err
	}
	scope := primarySavingScope(account.AccountType)
	userOpenId := savingUserOpenId(scope, openid)

	var saving Saving
	err := tx.Table(SavingTableName).
		Where("account_id = ? AND scope = ? AND user_openid = ? AND is_delete = ?", accountId, scope, userOpenId, 1).
		First(&saving).Error
	if err == gorm.ErrRecordNotFound {
		record := buildInitialSaving(accountId, account.AccountType, openid, 0)
		if err = tx.Table(SavingTableName).Create(&record).Error; err != nil {
			return 0, err
		}
		saving = record
	} else if err != nil {
		return 0, err
	}

	delta := balanceAfter - saving.CurrentAmount
	if err = tx.Table(SavingTableName).
		Where("id = ?", saving.Id).
		Updates(map[string]interface{}{
			"current_amount": balanceAfter,
			"updated_at":     time.Now(),
		}).Error; err != nil {
		return 0, err
	}
	flow := SavingFlow{
		AccountId:     accountId,
		TransactionId: transactionId,
		Scope:         scope,
		UserOpenId:    userOpenId,
		Action:        action,
		DeltaAmount:   delta,
		BalanceAfter:  balanceAfter,
		RecordType:    recordType,
		CategoryName:  limitString(categoryName, 32),
		Description:   limitString(description, 80),
		IsDelete:      1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err = tx.Table(SavingFlowTableName).Create(&flow).Error; err != nil {
		return 0, err
	}
	return balanceAfter, nil
}

func currentSavingBalance(tx *gorm.DB, accountId int, openid string) (float64, error) {
	var account Account
	if err := tx.Table(AccountTableName).
		Where("id = ? AND is_delete = ?", accountId, 1).
		First(&account).Error; err != nil {
		return 0, err
	}
	scope := primarySavingScope(account.AccountType)
	userOpenId := savingUserOpenId(scope, openid)
	var saving Saving
	err := tx.Table(SavingTableName).
		Where("account_id = ? AND scope = ? AND user_openid = ? AND is_delete = ?", accountId, scope, userOpenId, 1).
		First(&saving).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return saving.CurrentAmount, nil
}

func buildSavingStatus(scope string, record Saving) SavingStatus {
	title := "账本储蓄"
	if scope == SavingScopeFamily {
		title = "家庭储蓄"
	} else if scope == SavingScopePersonal {
		title = "个人储蓄"
	}
	percent := 0
	if record.TargetAmount > 0 {
		percent = int(record.CurrentAmount / record.TargetAmount * 100)
	}
	return SavingStatus{
		Scope:         scope,
		Title:         title,
		TargetAmount:  record.TargetAmount,
		CurrentAmount: record.CurrentAmount,
		Percent:       percent,
		Note:          record.Note,
	}
}

func GetTransactionCategoryConfigs(accountId int, openid string) ([]TransactionCategoryConfigDTO, error) {
	var list []TransactionCategoryConfigDTO
	err := db.Table(CategoryConfigTableName).
		Select("id, category_name, amount, type, source, account_id, creator_openid").
		Where("account_id = ? AND creator_openid = ? AND is_delete = ?", accountId, openid, 1).
		Order("updated_at desc, id desc").
		Scan(&list).Error
	return list, err
}

func GetMonthlySavingConfigs(accountId int, openid string) ([]MonthlySavingConfigDTO, error) {
	var list []MonthlySavingConfigDTO
	err := db.Table(MonthlySavingConfigTableName).
		Select("id, category_name, amount, type, source, account_id, creator_openid").
		Where("account_id = ? AND creator_openid = ? AND is_delete = ?", accountId, openid, 1).
		Order("updated_at desc, id desc").
		Scan(&list).Error
	return list, err
}

func UpdateTransactionCategoryConfigRecord(accountId int, openid string, params *UpdateTransactionCategoryConfig) error {
	values := map[string]interface{}{
		"category_name": limitString(params.CategoryName, 32),
		"amount":        params.Amount,
		"type":          params.Type,
		"source":        limitString(params.Source, 80),
		"updated_at":    time.Now(),
	}
	return db.Table(CategoryConfigTableName).
		Where("id = ? AND account_id = ? AND creator_openid = ? AND is_delete = ?", params.Id, accountId, openid, 1).
		Updates(values).Error
}

func DeleteTransactionCategoryConfigRecord(id int, accountId int, openid string) error {
	return db.Table(CategoryConfigTableName).
		Where("id = ? AND account_id = ? AND creator_openid = ? AND is_delete = ?", id, accountId, openid, 1).
		Updates(map[string]interface{}{
			"is_delete":  2,
			"updated_at": time.Now(),
		}).Error
}

func UpdateMonthlySavingConfigRecord(accountId int, openid string, params *UpdateMonthlySavingConfig) error {
	values := map[string]interface{}{
		"category_name": limitString(params.CategoryName, 32),
		"amount":        params.Amount,
		"type":          params.Type,
		"source":        limitString(params.Source, 80),
		"updated_at":    time.Now(),
	}
	return db.Table(MonthlySavingConfigTableName).
		Where("id = ? AND account_id = ? AND creator_openid = ? AND is_delete = ?", params.Id, accountId, openid, 1).
		Updates(values).Error
}

func DeleteMonthlySavingConfigRecord(id int, accountId int, openid string) error {
	return db.Table(MonthlySavingConfigTableName).
		Where("id = ? AND account_id = ? AND creator_openid = ? AND is_delete = ?", id, accountId, openid, 1).
		Updates(map[string]interface{}{
			"is_delete":  2,
			"updated_at": time.Now(),
		}).Error
}

func upsertTransactionCategoryConfig(tx *gorm.DB, params *CreateTransaction) error {
	now := time.Now()
	values := map[string]interface{}{
		"amount":     params.Amount,
		"type":       params.Type,
		"source":     params.Description,
		"updated_at": now,
		"is_delete":  1,
	}

	var existed TransactionCategoryConfig
	err := tx.Table(CategoryConfigTableName).
		Where("account_id = ? AND creator_openid = ? AND category_name = ? AND type = ?", params.AccountId, params.OwnerId, params.Category, params.Type).
		First(&existed).Error
	if err == nil {
		return tx.Table(CategoryConfigTableName).
			Where("id = ?", existed.Id).
			Updates(values).Error
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	record := TransactionCategoryConfig{
		AccountId:     params.AccountId,
		CreatorOpenId: params.OwnerId,
		CategoryName:  params.Category,
		Amount:        params.Amount,
		Type:          params.Type,
		Source:        params.Description,
		IsDelete:      1,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	return tx.Table(CategoryConfigTableName).Create(&record).Error
}

func upsertMonthlySavingConfig(tx *gorm.DB, params *CreateTransaction) error {
	now := time.Now()
	values := map[string]interface{}{
		"amount":     params.Amount,
		"source":     params.Description,
		"updated_at": now,
		"is_delete":  1,
	}

	var existed MonthlySavingConfig
	err := tx.Table(MonthlySavingConfigTableName).
		Where("account_id = ? AND creator_openid = ? AND category_name = ? AND type = ?", params.AccountId, params.OwnerId, params.Category, params.Type).
		First(&existed).Error
	if err == nil {
		return tx.Table(MonthlySavingConfigTableName).
			Where("id = ?", existed.Id).
			Updates(values).Error
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}
	return createMonthlySavingConfig(tx, params.AccountId, params.OwnerId, params.Type, params.Amount, params.Category, params.Description)
}

func createMonthlySavingConfig(tx *gorm.DB, accountId int, ownerId string, recordType string, amount float64, categoryName string, source string) error {
	now := time.Now()
	record := MonthlySavingConfig{
		AccountId:     accountId,
		CreatorOpenId: ownerId,
		CategoryName:  limitString(categoryName, 32),
		Amount:        amount,
		Type:          recordType,
		Source:        limitString(source, 80),
		IsDelete:      1,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	return tx.Table(MonthlySavingConfigTableName).Create(&record).Error
}

func limitString(value string, max int) string {
	runes := []rune(strings.TrimSpace(value))
	if len(runes) == 0 {
		return ""
	}
	if len(runes) <= max {
		return string(runes)
	}
	return string(runes[:max])
}
