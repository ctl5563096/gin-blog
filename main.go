package main

import (
	"fmt"
	"gin-blog/models"
	"gin-blog/pkg/setting"
	"gin-blog/pkg/util"
	"gin-blog/routers"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	// 初始化env配置
	var appEnv = os.Getenv("APP_ENV")
	var envFile = ".env"
	if appEnv != "" {
		envFile = ".env" + "." +appEnv
	}
	// 加载env目录
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	util.WriteLog("test",1,"项目启动" +  util.ReturnCurrentTime("second"))
	// 初始化路由
	router := routers.InitRouter()
	// 初始化模型链接池
	models.Init()
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}