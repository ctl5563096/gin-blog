package blog

import (
	"encoding/json"
	"fmt"
	"gin-blog/pkg/util"
	jsonTime "gin-blog/pkg/util/json"
)

// 自定义表名
var tableName = "t_go_article"

type Article struct {
	Title     string `json:"title"`
	Summary   string `json:"summary"`
	Content   string `json:"content"`
	ViewCount int	 `json:"view"`
	Author    string `json:"author"`
}

type IndexArticle struct {
	Id     	  int    `json:"id"`
	Title     string `json:"title"`
	Summary   string `json:"summary"`
	Content   string `json:"content"`
	Cover	  string `json:"cover"`
	ViewCount int	 `json:"view"`
	Like 	  int	 `json:"like"`
	Author    string `json:"author"`
	CreatedAt jsonTime.JSONTime `json:"created_at"`
}

type ArticleSearchList struct {
	KeyWords string `json:"keywords"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Order    string `json:"order"`
}

type ArticleList struct {
	Id		  int 	 	`json:"id"`
	Title     string 	`json:"title"`
	ViewCount int	 	`json:"view_count"`
	IsShow	  int 	 	`json:"is_show"`
	Like 	  int	 	`json:"like"`
	Author    string 	`json:"author"`
	Sort 	  int 		`json:"sort"`
	CreatedAt jsonTime.JSONTime `json:"created_at"`
}

type ArticleDetail struct {
	Id		  int 	 	`json:"id"`
	Title     string 	`json:"title"`
	Summary   string	`json:"summary"`
	IsShow	  int 	 	`json:"is_show"`
	Cover 	  string	`json:"cover"`
	Author    string 	`json:"author"`
	Content   string 	`json:"content"`
	ViewCount int	 	`json:"view_count"`
	Like 	  int	 	`json:"like"`
	CreatedAt jsonTime.JSONTime `json:"created_at"`
}

// EditArticleStruct 修改结构体
type EditArticleStruct struct {
	Id 		 int	`json:"id" validate:"required"`
	Title 	 string `json:"title" validate:"required"`
	Summary  string `json:"summary" validate:"required"`
	Content  string `json:"content" validate:"required"`
	Cover    string `json:"cover"`
	Author   string `json:"author"`
	IsShow   uint  	`json:"is_show"`
}

// GetIndexArticle 获取首页的文章
func GetIndexArticle() ([]*IndexArticle,error) {
	var res [] *IndexArticle
	db := db.Table(tableName).Where("is_show = ?",1).Where("is_delete = ?",0).Order("created_at desc").Limit(5).Find(&res)
	if db.Error != nil {
		util.WriteLog("mysql_error",2,db.Error.Error())
		return nil,db.Error
	}
	return res,nil
}

// CreatArticle 创建新文章
func CreatArticle(params interface{}) (id int ,err error) {
	err = db.Table(tableName).Create(params).Error
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

// SearchArticle 有keywords就搜索 没有正常查找
func SearchArticle(params *ArticleSearchList) ([]*ArticleList ,error)  {
	var res [] *ArticleList
	dbLast := db.Table(tableName).Select("id,title,`view` as view_count,is_show,`like`,author,created_at")
	if params.KeyWords != "" {
		keywords := "%" + params.KeyWords + "%"
		dbLast = dbLast.Where(" is_delete = ? and title LIKE ? ", 0,keywords)
	}else {
		dbLast = dbLast.Where("is_delete = ? " ,0)
	}
	dbLast.Order("updated_at " + params.Order).Limit(params.PageSize).Offset((params.Page - 1) * params.PageSize).Find(&res)
	if dbLast.Error != nil {
		fmt.Println(dbLast.Error.Error())
		return nil,dbLast.Error
	}
	return res,nil
}

// SearchArticleCount 有keywords就搜索 没有正常查找
func SearchArticleCount(params *ArticleSearchList,ch chan int) {
	var res int
	dbLast := db.Table(tableName)
	if params.KeyWords != "" {
		keywords := "%" + params.KeyWords + "%"
		dbLast = dbLast.Where(" is_delete = ? and title LIKE ? ", 0,keywords)
	}else {
		dbLast = dbLast.Where("is_delete = ? " ,0)
	}
	dbLast.Count(&res)
	ch <-res
	close(ch)
}

// SearchArticleCountA 有keywords就搜索 没有正常查找
func SearchArticleCountA(params *ArticleSearchList) int{
	var res int
	dbLast := db.Table(tableName)
	if params.KeyWords != "" {
		keywords := "%" + params.KeyWords + "%"
		dbLast = dbLast.Where(" is_delete = ? and title LIKE ? ", 0,keywords)
	}else {
		dbLast = dbLast.Where("is_delete = ? " ,0)
	}
	dbLast.Count(&res)
	return res
}

// GetDetail 获取文章详情
func GetDetail(id int) ArticleDetail {
	var r ArticleDetail
	db.Table(tableName).Select("id,title,`view` as view_count,is_show,`like`,author,created_at,cover,content,summary").Where("id = ? " ,id).Where("is_delete = ?" , 0).First(&r)
	return r
}

// EditDetail 获取文章详情
func EditDetail(params *EditArticleStruct) bool {
	db.Table(tableName).Select("id").Where("id = ?" ,params.Id).Updates(params)
	if db.Error != nil {
		util.WriteLog("update_rule_err",4,db.Error.Error())
		return false
	}
	return true
}