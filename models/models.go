package models

import (
	"fmt"
	"gin-blog/models/user"
	"github.com/jinzhu/gorm"
	"strings"
)

const TablePrefix = "ci_"

// Init 初始化各个链接池
func Init()  {
	fmt.Println(">开始初始化各数据库连接池...")

	// 分别初始化各个库的链接池
	user.Init()

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