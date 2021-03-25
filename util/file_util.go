package util

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// 获取当前脚本的存储路径
func GetFilePath() string {
	str, _ := os.Getwd()
	return str
}

// 获取当前脚本的执行路径
func GetExecPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))

	return path[:index]
}
