package blog

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"os"
	"time"
)
// 数据库实例
var db *gorm.DB

type Model struct {
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}



const (
	// DataBase 默认链接库，有些模型里面需要设置库的
	DataBase = "goBlog"
	// 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	maxOpenConns = 1
	// 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	maxIdleConns = 1
	// 可以重用连接的最长时间[5分钟先]
	// 防止closing bad idle connection: EOF（ 在 MySQL Server 主动断开连接之前，MySQL Client 的连接池中的连接被关闭掉），具体值要问DBA
	// 数据库端设置生存时间30s
	maxLifetime = 28
)

// Init 初始化mysql链接
func Init()  {
	//dns := "chentulin:A5563096z@tcp(120.78.13.233:3306)/swoft?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_POST"),
		DataBase))

	if err != nil {
		log.Fatalf("User models.Init err: %v", err)
	}

	// 允许复数
	db.SingularTable(true)
	// 具体配置看一下上述const描述
	db.DB().SetMaxIdleConns(maxIdleConns)
	db.DB().SetMaxOpenConns(maxOpenConns)
	db.DB().SetConnMaxLifetime(time.Second * maxLifetime)
	// 强行重写表前缀
	gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
		return "t_go_" + defaultTableName
	}
	// 判断有没有这个表 如果没有的话就去创建
	//if !db.HasTable(&Article{}) {
	//	fmt.Println(">开始创建文章表...")
	//	db.Set("gorm:table_options", "ENGINE=InnoDB;").CreateTable(&Article{}).AddIndex("index_title","title")//创建表时指存在引擎
	//	fmt.Println("<创建文章表成功...")
	//}

	// 日志[生产必须关闭！]
	if os.Getenv("RUNMODE") == "debug" {
		db.LogMode(true)
	}
}

// CloseDB closes database connection (unnecessary)
func CloseDB() {
	defer func(db *gorm.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("User models.CloseDB err: %v", err)
		}
	}(db)
}


// updateTimeStampForUpdateCallback will set `UpdatedAt` when updating 更新的钩子自动更新update_time
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		err := scope.SetColumn("UpdatedAt", time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			return
		}
	}
}

// GetDB 获取数据库连接实例
func GetDB() *gorm.DB {
	return db
}


