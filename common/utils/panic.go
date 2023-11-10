package utils

import (
	"frozen-go-cms/common/mylogrus"
	"runtime/debug"
)

func CheckGoPanic() {
	if r := recover(); r != nil {
		//打印错误堆栈信息
		mylogrus.MyLog.Errorf("ACTION PANIC: %v, stack: %v", r, string(debug.Stack()))
	}
}
