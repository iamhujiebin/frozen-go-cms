package mysql

import (
	"context"
	"fmt"
	. "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

func MyNew(writer Writer, config Config) Interface {
	var (
		infoStr      = "%s[info] "
		warnStr      = "%s[warn] "
		errStr       = "%s[error] "
		traceStr     = "%s[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s[%.3fms] [rows:%v] %s"
	)

	//if config.Colorful {
	//	infoStr = Green + "%s\n" + Reset + Green + "[info] " + Reset
	//	warnStr = BlueBold + "%s\n" + Reset + Magenta + "[warn] " + Reset
	//	errStr = Magenta + "%s\n" + Reset + Red + "[error] " + Reset
	//	traceStr = Green + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
	//	traceWarnStr = Green + "%s " + Yellow + "%s\n" + Reset + RedBold + "[%.3fms] " + Yellow + "[rows:%v]" + Magenta + " %s" + Reset
	//	traceErrStr = RedBold + "%s " + MagentaBold + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
	//}
	myTraceStr := " traceId:%v userId:%v"
	infoStr += myTraceStr
	warnStr += myTraceStr
	errStr += myTraceStr
	traceStr += myTraceStr
	traceWarnStr += myTraceStr
	traceErrStr += myTraceStr

	return &myLogger{
		Writer:       writer,
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

type myLogger struct {
	Writer
	Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// LogMode log mode
func (l *myLogger) LogMode(level LogLevel) Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l myLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= Info {
		l.Printf(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l myLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= Warn {
		l.Printf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l myLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= Error {
		l.Printf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (l myLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	traceId, userId := ctx.Value("traceId"), ctx.Value("userId")
	if l.LogLevel > Silent {
		elapsed := time.Since(begin)
		switch {
		case err != nil && l.LogLevel >= Error:
			sql, rows := fc()
			if rows == -1 {
				l.Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql, traceId, userId)
			} else {
				l.Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql, traceId, userId)
			}
		case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
			if rows == -1 {
				l.Printf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql, traceId, userId)
			} else {
				l.Printf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql, traceId, userId)
			}
		case l.LogLevel == Info:
			sql, rows := fc()
			if rows == -1 {
				l.Printf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql, traceId, userId)
			} else {
				l.Printf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql, traceId, userId)
			}
		}
	}
}
