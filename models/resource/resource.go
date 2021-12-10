package resource

import (
	"fmt"
	"gin-blog/common"
)

type SearchStruct struct {
	Page 	 int 	`json:"page"`
	PageSize int 	`json:"pageSize"`
	Keywords string `json:"keywords"`
	Order 	 string `json:"order"`
}

type CreateData struct {
	Id  		 uint 		`json:"id"`
	ResourceType uint 		`json:"resource_type" validate:"required"`
	ResourceId	 uint 		`json:"resource_id" validate:"required"`
	IsTop		 *uint 		`json:"is_top" gorm:"default:'1'"`
	Title		 string 	`json:"title" validate:"required"`
	Cover		 string 	`json:"cover" validate:"required"`
	Sort 	  	 int 		`json:"sort" gorm:"-"`
}

type Data struct {
	Id  		 uint 		`json:"id"`
	ResourceType uint 		`json:"resource_type" validate:"required"`
	ResourceId	 uint 		`json:"resource_id" validate:"required"`
}

type Resource interface {
	GetResource(msg string) []interface{}
}

const  TableName = "t_go_hot_content"

// GetResource 获取资源列表
func GetResource(tableName string) []Data {
	var r []Data
	db.Table(tableName).Select("id,title.cover").Where("is_deleted = ?",1).Find(&r)
	return r
}

// CreatedRecord 新建记录
func CreatedRecord(params CreateData) bool {
	// 先检查这个资源是否被删除过 如果是就恢复并且更新 否则就是插入
	var r common.Select
	db.Table(TableName).Select("*").Where("resource_id = ? and resource_type = ? and is_delete = ?",params.ResourceId,params.ResourceType,2).First(&r)
	// 去更新
	if r.Id > 0  {
		db.Table(TableName).Where("id = ?",r.Id).Update(map[string]interface{}{"title": params.Title,"cover":params.Cover,"is_top":params.IsTop,"is_delete":1})
		if db.Error != nil {
			return false
		}
		return true
	}
	db.Table(TableName).Create(&params)
	if db.Error != nil || params.Id <= 0{
		return false
	}
	return true
}

// GetList 获取列表
func GetList(params SearchStruct) ([]*CreateData,error) {
	var res [] *CreateData
	dbLast := db.Table(TableName).Select("id,resource_type,resource_id,is_top,cover,title")
	if params.Keywords != "" {
		keywords := "%" + params.Keywords + "%"
		dbLast = dbLast.Where(" is_delete = ? and title LIKE ? ", 1,keywords)
	}else {
		dbLast = dbLast.Where("is_delete = ? " ,1)
	}
	dbLast.Order("updated_at " + params.Order).Limit(params.PageSize).Offset((params.Page - 1) * params.PageSize).Find(&res)
	if dbLast.Error != nil {
		fmt.Println(dbLast.Error.Error())
		return nil,dbLast.Error
	}
	return res,nil
}

// GetCount 获取资源总数
func GetCount(params SearchStruct) int {
	var res int
	dbLast := db.Table(TableName)
	if params.Keywords != "" {
		keywords := "%" + params.Keywords + "%"
		dbLast = dbLast.Where(" is_delete = ? and title LIKE ? ", 1,keywords)
	}else {
		dbLast = dbLast.Where("is_delete = ? " ,1)
	}
	dbLast.Count(&res)
	return res
}

// SetTopStatus 更新资源的状态
func SetTopStatus(id int,status int,filed string) bool {
	db.Table(TableName).Where("id = ?",id).Updates(map[string]interface{}{filed: status})
	if db.Error != nil {
		return false
	}
	return true
}