package util

import (
	"log"
	"os"
)

// WriteLog /** 系统级别日志打印到文件目录 **/
func WriteLog(fileNameInput string,level int,message string )  {
	var prefix string
	var emptyErr = true
	// 文件名
	fileName := "runtime/logs/" + fileNameInput + "-" + ReturnCurrentTime("first") + ".log"

	// 先检查文件是否存在
	if _, err := os.Open(fileName)

	err != nil {
		// 判断文件是否存在
		emptyErr = os.IsExist(err)
	}

	// 如果不存在就循环创建文件夹 存在的话直接打开
	if !emptyErr {
		err := os.MkdirAll("runtime/logs", 0777)
		if err != nil {
			log.Fatalln("mkdir error")
		}
		_,err = os.Create(fileName)
	}

	// 如果要追加写入文件则注意以免被覆盖
	logFile,err := os.OpenFile(fileName,os.O_WRONLY|os.O_CREATE|os.O_APPEND,0777)
	// 打开文件失败
	if err != nil {
		log.Fatalln(err.Error())
	}

	// 关闭文件句柄 防止系统打开文件过多
	defer logFile.Close()
	switch level {
		case 1:
			prefix = "[info]"
			break
		case 2:
			prefix = "[mysql]"
			break
		case 3:
			prefix = "[work]"
			break
		default:
			prefix = "[error]"
	}
	// 构造日志对象 并且写入日志
	debugLog := log.New(logFile,prefix,log.Llongfile)
	debugLog.SetPrefix(prefix)
	debugLog.Println(message)
}
