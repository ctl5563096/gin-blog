package user


type AdminUser struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	PhoneNum string `json:"phone_num"`
	IsBlack int `json:"is_black"`
}


func GetUsers(pageNum int, pageSize int, maps interface {}) (tags []AdminUser) {
	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)
	return
}

func GetUsersTotal(maps interface {}) (count int64){
	db.Model(&AdminUser{}).Where(maps).Count(&count)
	return
}

// GetFirst 获取第一条记录
func GetFirst()  (res[]AdminUser){
	db.First(&res)
	return
}