package main

import (
	"fmt"
	"gin-blog/pkg/cache"
	"gin-blog/pkg/setting"
	"gin-blog/pkg/util"
	"gin-blog/routers"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	defer func() {
		if err := recover(); err != nil {

			fmt.Println("Panic Occured and Recovered in, Error Info: ", err)

		}
	}()
	// 初始化env配置
	var appEnv  = os.Getenv("APP_ENV")
	var envFile = ".env"
	if appEnv != "" {
		envFile = ".env" + "." +appEnv
	}
	envPath :=  util.GetAbsolutelyPath() + `/` +envFile
	// 加载env目录
	err := godotenv.Load(envPath)
	if err != nil {
		err := godotenv.Load(envFile)
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	util.WriteLog("project",3,"项目启动")
	util.WriteLog("project",3,"端口:" + strconv.Itoa(setting.HTTPPort))
	// 初始化路由
	router := routers.InitRouter()
	// 初始化缓存
	cache.Init()
	// 初始化语言包
	//valid.Init()
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}