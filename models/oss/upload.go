package oss

type SaveStruct struct {
	Oss string `json:"oss"`
	Local string `json:"local"`
}

// 自定义表名
var tableName = "t_go_local_oss"

// SaveRelation 把oss的和本地的关系存起来
func SaveRelation(oss string, local string) bool {
	db.Table(tableName).Create(SaveStruct{Oss:oss,Local:local})
	if db.Error != nil {
		return false
	}
	return true
}
