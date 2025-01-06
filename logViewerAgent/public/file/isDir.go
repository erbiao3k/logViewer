package file

import "os"

// IsDir 判断是否为目录
func IsDir(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return fi.IsDir()
}
