package mylogrus

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"time"
)

const logDir = "/var/log/hilo/"

var filenamePrefix string
var MyLog = logrus.New()

func Info(v interface{}) {
	MyLog.Info("")
}

func init() {
	filenamePrefix = logDir + filepath.Base(os.Args[0]) + "."
	// stderr日志重定向
	MyLog.SetOutput(os.Stdout)
	RewriteStderrFile()

	MyLog.SetFormatter(&logrus.TextFormatter{
		ForceQuote:      false,
		DisableQuote:    true,
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true,
	})
	hook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: getLevelWrite(logrus.DebugLevel),
		logrus.InfoLevel:  getLevelWrite(logrus.InfoLevel),
		logrus.WarnLevel:  getLevelWrite(logrus.WarnLevel),
		logrus.ErrorLevel: getLevelWrite(logrus.ErrorLevel),
		logrus.FatalLevel: getLevelWrite(logrus.FatalLevel),
		logrus.PanicLevel: getLevelWrite(logrus.PanicLevel),
	}, &logrus.TextFormatter{ForceQuote: false, DisableQuote: true, TimestampFormat: time.RFC3339Nano})
	MyLog.AddHook(hook)
	MyLog.SetLevel(logrus.InfoLevel)
	MyLog.SetReportCaller(true)
}

func GetInfoLog() io.Writer {
	return getLevelWrite(logrus.InfoLevel)
}

func getLevelWrite(level logrus.Level) io.Writer {
	var name string
	switch level {
	case logrus.DebugLevel:
		name = "debug.log"
	case logrus.InfoLevel:
		name = "info.log"
	case logrus.WarnLevel:
		name = "warn.log"
	case logrus.ErrorLevel:
		name = "error.log"
	case logrus.FatalLevel:
		name = "fatal.log"
	case logrus.PanicLevel:
		name = "panic.log"
	}
	name = filenamePrefix + name
	writer, err := rotatelogs.New(
		name+".%Y%m%d%H",
		rotatelogs.WithLinkName(name),          // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),  // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Hour), // 日志切割时间间隔
	)
	if err != nil {
		MyLog.Fatal("Failed to create log file:", err.Error())
	}
	return writer
}

func GetSqlLog() io.Writer {
	var name string = "sql.log"
	name = filenamePrefix + name
	writer, err := rotatelogs.New(
		name+".%Y%m%d%H",
		rotatelogs.WithLinkName(name),          // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),  // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Hour), // 日志切割时间间隔
	)
	if err != nil {
		MyLog.Fatal("Failed to create log file:", err.Error())
	}
	return writer
}
