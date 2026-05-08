package ledger

import (
	"github.com/jinzhu/gorm"
	"time"
)

const (
	TransactionTableName = "t_family_transaction"
	BudgetTableName      = "t_family_budget"
	UserTableName        = "t_family_user"
	AccountTableName     = "t_family_account"
	MemberTableName      = "t_family_account_member"

	AccountTypePersonal = "personal"
	AccountTypeFamily   = "family"
	BudgetScopeAccount  = "account"
	BudgetScopeMember   = "member"
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
	Id            int       `json:"id" gorm:"primary_key"`
	AccountId     int       `json:"accountId" gorm:"column:account_id"`
	OwnerId       string    `json:"ownerId" gorm:"column:owner_id"`
	Type          string    `json:"type"`
	Amount        float64   `json:"amount"`
	Category      string    `json:"category"`
	Note          string    `json:"note"`
	Date          string    `json:"date"`
	PaymentMethod string    `json:"paymentMethod" gorm:"column:payment_method"`
	MemberName    string    `json:"memberName" gorm:"column:member_name"`
	IsDelete      int       `json:"isDelete" gorm:"column:is_delete"`
	CreatedAt     time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

type CreateTransaction struct {
	AccountId     int     `json:"accountId" validate:"gte=0"`
	OwnerId       string  `json:"ownerId" validate:"max=64"`
	Type          string  `json:"type" validate:"required,oneof=expense income"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	Category      string  `json:"category" validate:"required,max=32"`
	Note          string  `json:"note" validate:"max=80"`
	Date          string  `json:"date" validate:"required,len=10"`
	PaymentMethod string  `json:"paymentMethod" validate:"required,max=32"`
	MemberName    string  `json:"memberName" validate:"required,max=32"`
}

type UpdateTransaction struct {
	Id            int     `json:"id" validate:"required,gt=0"`
	AccountId     int     `json:"accountId" validate:"gte=0"`
	OwnerId       string  `json:"ownerId" validate:"max=64"`
	Type          string  `json:"type" validate:"required,oneof=expense income"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	Category      string  `json:"category" validate:"required,max=32"`
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
	AccountType          string  `json:"accountType" validate:"required,oneof=personal family"`
	Name                 string  `json:"name" validate:"max=32"`
	BudgetAmount         float64 `json:"budgetAmount" validate:"gte=0"`
	PersonalBudgetAmount float64 `json:"personalBudgetAmount" validate:"gte=0"`
}

type BudgetInfo struct {
	AccountId            int     `json:"accountId"`
	AccountType          string  `json:"accountType"`
	Amount               float64 `json:"amount"`
	PersonalAmount       float64 `json:"personalAmount"`
	PersonalBudgetAmount float64 `json:"personalBudgetAmount"`
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
	tx.Commit()
	return GetAccountForUser(openid, account.Id)
}

func CreateRecord(params *CreateTransaction) (*Transaction, error) {
	record := &Transaction{
		AccountId:     params.AccountId,
		OwnerId:       params.OwnerId,
		Type:          params.Type,
		Amount:        params.Amount,
		Category:      params.Category,
		Note:          params.Note,
		Date:          params.Date,
		PaymentMethod: params.PaymentMethod,
		MemberName:    params.MemberName,
		IsDelete:      1,
	}
	err := db.Table(TransactionTableName).Create(record).Error
	return record, err
}

func UpdateRecord(params *UpdateTransaction) error {
	values := map[string]interface{}{
		"type":           params.Type,
		"amount":         params.Amount,
		"category":       params.Category,
		"note":           params.Note,
		"date":           params.Date,
		"payment_method": params.PaymentMethod,
		"member_name":    params.MemberName,
		"updated_at":     time.Now(),
	}
	return db.Table(TransactionTableName).
		Where("id = ? AND account_id = ? AND is_delete = ?", params.Id, params.AccountId, 1).
		Updates(values).Error
}

func DeleteRecord(id int, accountId int) error {
	return db.Table(TransactionTableName).
		Where("id = ? AND account_id = ? AND is_delete = ?", id, accountId, 1).
		Update("is_delete", 2).Error
}

func GetList(accountId int, month string, recordType string) ([]Transaction, error) {
	var list []Transaction
	dbq := db.Table(TransactionTableName).Where("account_id = ? AND is_delete = ?", accountId, 1)
	if month != "" {
		dbq = dbq.Where("date LIKE ?", month+"%")
	}
	if recordType != "" && recordType != "all" {
		dbq = dbq.Where("type = ?", recordType)
	}
	err := dbq.Order("date desc,id desc").Find(&list).Error
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
