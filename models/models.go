package models

import (
	"fmt"
	"gin-blog/models/blog"
	"gin-blog/models/oss"
	"gin-blog/models/resource"
	"gin-blog/models/rule"
	"gin-blog/models/system"
	"gin-blog/models/user"
	"github.com/jinzhu/gorm"
	"strings"
)

const TablePrefix = "ci_"

// InitModel 初始化各个链接池
func InitModel()  {
	fmt.Println(">开始初始化各数据库连接池...")

	// 分别初始化各个库的链接池
	user.Init()
	// 初始化新blog数据库
	blog.Init()
	// 初始化权限数据库
	rule.Init()
	// 初始化资源数据库
	oss.Init()
	// 初始化系统数据库
	system.Init()
	// 初始化资源的数据库
	resource.Init()

	// 设置表名【注意所有数据库链接都会通用这个方法】
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		if strings.Contains(defaultTableName, ".") {
			return defaultTableName
		} else {
			return TablePrefix + defaultTableName
		}
	}
	fmt.Println(">>>初始化数据库连接池完成")
}

// CloseDB 关闭各数据库链接
func CloseDB() {
	user.CloseDB()
}