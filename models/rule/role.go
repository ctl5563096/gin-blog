package rule

// 自定义表名
var tableNameRole = "t_go_role"

type RoleResponse struct {
	Id 		 int 	`json:"id"`
	RoleName string `json:"role"`
}

// GetAllRole 获取所有的角色
func GetAllRole() []*RoleResponse {
	var rs []*RoleResponse
	db.Table(tableNameRole).Where("is_enabled = ? ",1).Find(&rs)
	return rs
}
