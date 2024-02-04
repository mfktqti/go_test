package utils

import (
	"os"
	"os/exec"
	"path/filepath"
)

// GetProcessName 获取进程的名称
func GetProcessName() string {
	fullFile, _ := exec.LookPath(os.Args[0])
	filename := filepath.Base(fullFile)
	if filename != "" {
		return filename
	}
	return "app"
}
