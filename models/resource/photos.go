package resource

import (
	"encoding/json"
	"fmt"
	"gin-blog/pkg/util"
	jsonTime "gin-blog/pkg/util/json"
)

const PhotosTableName = "t_go_photos"
const PhotosListTableName = "t_go_photos_list"

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

type PhotosDetail struct {
	Id 		  	  int    `json:"id"`
	Title     	  string `json:"title" validate:"required"`
	Summary   	  string `json:"summary"`
	Url 	  	  string `json:"url" validate:"required"`
	Thumb	  	  string `json:"thumb"`
	FileName  	  string `json:"file_name" validate:"required"`
	IsTop     	  int 	 `json:"is_top"`
	CreatedAt     jsonTime.JSONTime 	 `json:"created_at"`
}

type PhotosList struct {
	PhotosArr []Photo `json:"changePhotoArr"`
}

type Photo struct {
	Id 		  	  int    `json:"id"`
	Uid           int 	 `json:"uid"`
	Status        string `json:"status"`
	ResourceId    int 	 `json:"resource_id"`
	Url 	  	  string `json:"url"`
	Thumb	  	  string `json:"thumb"`
	FileName  	  string `json:"file_name"`
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

// GetPhotoDetail 获取图片详情
func GetPhotoDetail(id int) (PhotosDetail,error) {
	var r PhotosDetail
	err := db.Table(PhotosTableName).Where("id = ?",id).Where("is_delete = ?", 1).First(&r)
	if err.Error != nil {
		return r,err.Error
	}
	return r,nil
}

// UpdatePhotoDetail 更新图片详情
func UpdatePhotoDetail(id int, data PhotosData) bool  {
	err := db.Table(PhotosTableName).Where("id = ?",id).Update(data)
	if err.Error != nil {
		return false
	}
	return true
}

// BatchInsertPhoto 批量插入角色和权限
func BatchInsertPhoto (data []interface{}) bool {
	sql := util.GetBranchInsertSql(data,PhotosListTableName)
	db.Exec(sql)
	if db.Error != nil {
		return false
	}
	return true
}

// GetAboutPhotos 去查询与之关联的系列的图片
func GetAboutPhotos(resourceId int) []Photo {
	var r []Photo
	db.Table(PhotosListTableName).Select("id,uid,url,thumb,file_name,status,resource_id").
		Where("resource_id = ?",resourceId).
		Where("is_delete = ?",1).Find(&r)
	return r
}

// DelPhotosList 批量更新删除状态
func DelPhotosList(data []interface{}) {
	db.Table(PhotosListTableName).Where("id in (?)", data).Updates(map[string]interface{}{"is_delete": 2})
}