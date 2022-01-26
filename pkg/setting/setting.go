package setting

import (
	"fmt"
	"github.com/go-ini/ini"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	Cfg *ini.File

	RunMode string

	HTTPPort int
	ReadTimeout time.Duration
	WriteTimeout time.Duration
	PageSize int
	JwtSecret string
)

func init()  {
	var err error
	path := "conf/app.ini"
	Cfg, err = ini.Load(path)
	if err != nil {
		absolutelyConfFilePath := GetAbsolutelyPath() + path
		fmt.Println(absolutelyConfFilePath)
		Cfg, err = ini.Load(absolutelyConfFilePath)
		if err != nil {
			log.Fatalf("Fail to parse '%v': %v", absolutelyConfFilePath,err)
		}
	}

	LoadBase()
	LoadServer()
	LoadApp()
}

func LoadBase() {
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
}


func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
	HTTPPort = 8000
	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTimeout =  time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}
func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}
	JwtSecret = sec.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")
	PageSize = sec.Key("PAGE_SIZE").MustInt(10)
}

func GetAbsolutelyPath() string {
	file, err := exec.LookPath(os.Args[0])
	path, err := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	path = path[:index]
	if err != nil {
		log.Fatal("get path error")
	}
	return path + `/`
}