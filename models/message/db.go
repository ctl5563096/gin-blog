package message

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"os"
	"time"
)

var db *gorm.DB

const (
	DataBase     = "goBlog"
	maxOpenConns = 1
	maxIdleConns = 1
	maxLifetime  = 28
)

func Init() {
	var err error
	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_POST"),
		DataBase))

	if err != nil {
		log.Fatalf("Message models.Init err: %v", err)
	}

	db.SingularTable(true)
	db.DB().SetMaxIdleConns(maxIdleConns)
	db.DB().SetMaxOpenConns(maxOpenConns)
	db.DB().SetConnMaxLifetime(time.Second * maxLifetime)

	if os.Getenv("RUNMODE") == "debug" {
		db.LogMode(true)
	}
}

func GetDB() *gorm.DB {
	return db
}
