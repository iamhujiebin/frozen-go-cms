//go:build windows
// +build windows

package mylogrus

import (
	"os"
	"path/filepath"
	"time"
)

func RewriteStderrFile() {
	filename := logDir + filepath.Base(os.Args[0]) + ".stderr.log"
	if exits, _ := pathExists(filename); exits {
		os.Rename(filename, filename+"_"+time.Now().Format("20060102150405"))
	}

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	MyLog.Errorf("stderr log in:%v,err:%v", file, err)
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
