package common

import (
	"fmt"
	"runtime"
)

// FuncName 获取当前运行的方法名称
func FuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return fmt.Sprintf("方法：%s，", f.Name())
}
