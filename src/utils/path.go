package utils

import (
	"os"
)

// 检查文件是否存在
func PathExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// 获取一个文件信息
func PathInfo(path string) (os.FileInfo, error) {
	return os.Stat(path)
}
