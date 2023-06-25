//go:build !windows
// +build !windows

package mylogrus

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
)

var stdErrFileHandler *os.File

func RewriteStderrFile() {
	filename := logDir + filepath.Base(os.Args[0]) + ".stderr.log"
	//if runtime.GOOS == "darwin" { // mac本地调试
	//	filename = "./log/hilo/" + filepath.Base(os.Args[0]) + ".stderr.log"
	//}
	if exits, _ := pathExists(filename); exits {
		os.Rename(filename, filename+"_"+time.Now().Format("20060102150405"))
	}

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	stdErrFileHandler = file //把文件句柄保存到全局变量，避免被GC回收

	if err = syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd())); err != nil {
		fmt.Println(err)
		return
	}
	// 内存回收前关闭文件描述符
	runtime.SetFinalizer(stdErrFileHandler, func(fd *os.File) {
		fd.Close()
	})

	return
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
