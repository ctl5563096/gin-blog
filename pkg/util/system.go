package util

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

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