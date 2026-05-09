package util

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const logRetentionDays = 30

var (
	logCleanupMu   sync.Mutex
	lastLogCleanup string
)

// WriteLog /** 系统级别日志打印到文件目录 **/
func WriteLog(fileNameInput string, level int, message string) {
	var prefix string
	absPath := GetAbsolutelyPath()
	logDir := absPath + "runtime/logs"
	fileName := logDir + "/" + logLevelName(level) + "-" + ReturnCurrentTime("first") + ".log"

	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Printf("mkdir log dir failed: %s, err: %v", logDir, err)
		return
	}
	cleanupExpiredLogs(logDir)

	// 如果要追加写入文件则注意以免被覆盖
	logFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	// 打开文件失败
	if err != nil {
		log.Printf("open log file failed: %s, err: %v", fileName, err)
		return
	}

	// 关闭文件句柄 防止系统打开文件过多
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			fmt.Println("关闭文件失败!!!")
		}
	}(logFile)
	switch level {
	case 1:
		prefix = "[info]"
		break
	case 2:
		prefix = "[info][mysql]"
		break
	case 3:
		prefix = "[info][work]"
		break
	case 4:
		prefix = "[error]"
		break
	case 5:
		prefix = "[warning]"
		break
	default:
		prefix = "[error]"
	}
	// 构造日志对象 并且写入日志
	debugLog := log.New(logFile, prefix, log.Llongfile)
	debugLog.SetPrefix(prefix)
	debugLog.Println("[" + fileNameInput + "] " + message)
}

func logLevelName(level int) string {
	switch level {
	case 4:
		return "error"
	case 5:
		return "warning"
	default:
		return "info"
	}
}

func cleanupExpiredLogs(logDir string) {
	today := ReturnCurrentTime("first")

	logCleanupMu.Lock()
	if lastLogCleanup == today {
		logCleanupMu.Unlock()
		return
	}
	lastLogCleanup = today
	logCleanupMu.Unlock()

	cutoff := time.Now().AddDate(0, 0, -logRetentionDays)
	entries, err := os.ReadDir(logDir)
	if err != nil {
		log.Printf("read log dir failed: %v", err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".log" {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			log.Printf("read log file info failed: %s, err: %v", entry.Name(), err)
			continue
		}
		if info.ModTime().After(cutoff) {
			continue
		}

		path := filepath.Join(logDir, entry.Name())
		if err := os.Remove(path); err != nil {
			log.Printf("remove expired log failed: %s, err: %v", path, err)
		}
	}
}
