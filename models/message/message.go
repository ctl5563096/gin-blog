package message

const TableName = "t_go_message"

type CreateMessage struct {
	Id        int    `json:"id"`
	Email     string `json:"email" validate:"required,email"`
	Content   string `json:"content" validate:"required,min=5,max=2000"`
	Ip        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	IsDelete  int    `json:"is_delete"`
}

type MessageListItem struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	Content   string `json:"content"`
	Ip        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	CreatedAt string `json:"created_at"`
}

func CreateRecord(params *CreateMessage) error {
	if params.IsDelete == 0 {
		params.IsDelete = 1
	}
	return db.Table(TableName).Create(params).Error
}

func GetList(page int, pageSize int, keywords string) ([]MessageListItem, error) {
	var list []MessageListItem
	dbq := db.Table(TableName).Select("id,email,content,ip,user_agent,created_at").Where("is_delete = ?", 1)
	if keywords != "" {
		like := "%" + keywords + "%"
		dbq = dbq.Where("email LIKE ? OR content LIKE ?", like, like)
	}
	err := dbq.Order("id desc").Limit(pageSize).Offset((page - 1) * pageSize).Find(&list).Error
	return list, err
}

func GetCount(keywords string) (int, error) {
	var count int
	dbq := db.Table(TableName).Where("is_delete = ?", 1)
	if keywords != "" {
		like := "%" + keywords + "%"
		dbq = dbq.Where("email LIKE ? OR content LIKE ?", like, like)
	}
	err := dbq.Count(&count).Error
	return count, err
}
