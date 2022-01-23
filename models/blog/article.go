package blog

import (
	"encoding/json"
	"fmt"
	"gin-blog/common"
	"gin-blog/models/system"
	"gin-blog/pkg/util"
	jsonTime "gin-blog/pkg/util/json"
	"strconv"
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

type Tag struct {
	Id 	  int 	 `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
	Css   string `json:"css"`
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
	Tags      []Tag  `json:"tags"`
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

type CodeArticleList struct {
	 Id		int `json:"id"`
	 Title  string `json:"title"`
	 Author string `json:"author"`
	 CreatedAt jsonTime.JSONTime `json:"created_at"`
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

// BatchInsertTags 批量插入tag
func BatchInsertTags(data []interface{}) bool {
	sql := util.GetBranchInsertSql(data,"t_go_resource_tags_relation")
	fmt.Println(sql)
	db.Exec(sql)
	if db.Error != nil {
		return false
	}
	return true
}

// EditArticleTags 修改文章的tags 先删除后删除
func EditArticleTags(id int,originList []int,newList []interface{}) bool {
	tx := db.Begin()
	// 执行删除
	tx.Table("t_go_resource_tags_relation").Where("resource_type = ? and resource_id = ? and code = ? and param_value in (?)" , "article",id,"articleType",originList).Delete(&system.DelParam{})
	if tx.Error != nil {
		// 回滚事务
		tx.Rollback()
		return false
	}

	// 执行新增
	sql := util.GetBranchInsertSql(newList,"t_go_resource_tags_relation")
	util.WriteLog("exec_origin_sql" , 2,sql)
	tx.Exec(sql)
	if tx.Error != nil {
		// 回滚事务
		tx.Rollback()
		return false
	}
	tx.Commit()
	if tx.Error != nil {
		util.WriteLog("commit_error" , 4,"提交事务失败文章ID" + string(rune(id)))
	}
	return true
}

// GetIndexArticleTags 获取首页文章的tags
func GetIndexArticleTags(idData []int) []Tag {
	var  sql string
	var  r 	 []Tag
	sql = "SELECT a.resource_id as id ,a.param_value as value,b.param_name as name,a.tag_style as css from t_go_resource_tags_relation as " +
		"a LEFT JOIN t_go_parameter as b on a.param_value = b.param_value WHERE a.resource_type = " +
		"'article' AND a.code = 'articleType' AND b.code = 'articleType' AND a.resource_id IN (?) ORDER BY a.created_at asc"
	db.Raw(sql,idData).Scan(&r)
	return r
}

// GetChooseList 获取文章选择列表
func GetChooseList(params common.GetListParams,result *[]common.ReturnGetList )  {
	db.Table(tableName).Select("id,title,cover").Where("is_show = ? and is_delete = ?",1,0).Limit(params.PageSize).Offset((params.Page - 1) * params.PageSize).Find(&result)
}

// SearchArticleByTag 根据
func SearchArticleByTag(page int,pageSize int) []CodeArticleList{
	// 查询tag relation表的关系
	var (
		sql string
		r []CodeArticleList
	)
	page = page - 1
	sql = "SELECT a.id,a.title,a.created_at,a.author FROM t_go_article AS a " +
		"LEFT JOIN (SELECT resource_id FROM t_go_resource_tags_relation WHERE resource_type = 'articleType' " +
		"AND CODE = 'articleType' AND 'param_value' = 14 GROUP BY resource_id ) AS b ON a.id = b.resource_id WHERE " +
		"a.is_show = 1 AND a.is_delete = 0 ORDER BY a.created_at DESC limit " + strconv.Itoa(page * pageSize) + "," +  strconv.Itoa(pageSize)
	db.Raw(sql).Scan(&r)
	return r
}

// SearchArticleByTagCount 计算总数
func SearchArticleByTagCount() int {
	var total int
	sql := "SELECT count(a.id) as total FROM t_go_article AS a " +
		"LEFT JOIN (SELECT resource_id FROM t_go_resource_tags_relation WHERE resource_type = 'articleType' " +
		"AND CODE = 'articleType' AND 'param_value' = 14 GROUP BY resource_id ) AS b ON a.id = b.resource_id WHERE " +
		"a.is_show = 1 AND a.is_delete = 0 "
	db.Raw(sql).Count(&total)
	return total
}