package resource

import (
	"encoding/json"
	"fmt"
)

const PhotosTableName = "t_go_photos"

// PhotosData 创建图片结构体
type PhotosData struct {
	Id 		  int    `json:"id"`
	Title     string `json:"title" validate:"required"`
	Summary   string `json:"summary"`
	Url 	  string `json:"url" validate:"required"`
	Thumb	  string `json:"thumb"`
	FileName  string `json:"file_name" validate:"required"`
	IsTop     int 	 `json:"is_top"`
}

// CreatPhotoRecord 创建记录
func CreatPhotoRecord(params *PhotosData) (int,error) {
	err := db.Table(PhotosTableName).Create(params).Error
	if err != nil {
		fmt.Println()
		return 0,err
	}
	m := make(map[string]interface{})
	j,_ := json.Marshal(params)
	err = json.Unmarshal(j, &m)
	if err != nil {
		return 0, err
	}
	recordId := int(m["id"].(float64))
	return recordId,err
}

// GetPhotoList 获取图片列表
func GetPhotoList(page int, pageSize int,keyword string,order string)([]PhotosData,error)  {
	var r []PhotosData
	dbq := db.Table(PhotosTableName).Select("id,title,summary,url,thumb")
	if keyword != "" {
		keywords := "%" + keyword + "%"
		dbq = dbq.Where(" is_delete = ? and title LIKE ? ", 1,keywords)
	} else {
		dbq = dbq.Where(" is_delete = 1")
	}
	dbq.Order("updated_at " + order).Limit(pageSize).Offset((page - 1) * pageSize).Find(&r)
	if dbq.Error != nil {
		fmt.Println(dbq.Error.Error())
		return nil,dbq.Error
	}
	return r,nil
}

// GetPhotoCount 获取总条数
func GetPhotoCount(keywords string) (int,error) {
	var count int
	dbq := db.Table(PhotosTableName)
	if keywords != "" {
		keywords := "%" + keywords + "%"
		dbq = dbq.Where(" is_delete = ? and title LIKE ? ", 1,keywords)
	} else {
		dbq = dbq.Where(" is_delete = 1")
	}
	dbq.Count(&count)
	return count,nil
}

// DeletePhotoRecord 删除图片记录
func DeletePhotoRecord(id int) bool {
	values := map[string]interface{}{
		"is_delete": 2,
	}
	err :=db.Table(PhotosTableName).Where("id = ?" ,id).Updates(values)
	if err.Error != nil {
		return false
	}
	return true
}