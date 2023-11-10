package utils

import (
	"github.com/jinzhu/now"
	"time"
)

const DEFAULT_LANG = "en"
const DATETIME_FORMAT = "2006-01-02 15:04:05"
const COMPACT_DATE_FORMAT = "20060102"
const DATE_FORMAT = "2006-01-02"
const COMPACT_MONTH_FORMAT = "200601"
const MONTH_FORMAT = "2006-01"
const DAY_SECOND = 24 * 60 * 60

// 获取沙特时区
var CstZone = time.FixedZone("CST", 3*3600) // 东三

func GetZeroTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func GetDayEndTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999, t.Location())
}

// 开始结束
func DayStartEnd(date time.Time) (start, end time.Time) {
	start = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end = time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999, date.Location())
	return
}

// 取最近的一个星期n
func GetLastDayOfWeek(t time.Time, n time.Weekday) time.Time {
	weekDay := t.Weekday()
	// 校正日期
	if weekDay < n {
		weekDay = 7 - n + weekDay
	} else {
		weekDay = weekDay - n
	}
	return t.AddDate(0, 0, -(int(weekDay)))
}

func GetMonday(t time.Time) time.Time {
	return GetLastDayOfWeek(t, time.Monday)
}

// 增加年/月
// 因为golang原生的Time.AddDate增加月份的时候有bug
func AddDate(t time.Time, years int, months int) time.Time {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()

	// firstDayOfMonthAfterAddDate: years 年，months 月后的 那个月份的1号
	firstDayOfMonthAfterAddDate := time.Date(year+years, month+time.Month(months), 1,
		hour, min, sec, t.Nanosecond(), t.Location())
	// firstDayOfMonthAfterAddDate 月份的最后一天
	lastDay := now.New(firstDayOfMonthAfterAddDate).EndOfMonth().Day()

	// 如果 t 的天 > lastDay，则设置为lastDay
	// 如：t 为 2020-03-31 12:00:00 +0800，增加1个月，为4月31号
	// 但是4月没有31号，则设置为4月最后一天lastDay（30号）
	if day > lastDay {
		day = lastDay
	}

	return time.Date(year+years, month+time.Month(months), day,
		hour, min, sec, t.Nanosecond(), t.Location())
}

// 当天结束剩余秒数
func DayRemainSecond(date time.Time) int64 {
	return time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999, date.Location()).Unix() - date.Unix()
}

func GetLastMonthStart(t time.Time) time.Time {
	thisMonthStart := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	lastDay := thisMonthStart.AddDate(0, 0, -1)
	lastMonthStart := time.Date(lastDay.Year(), lastDay.Month(), 1, 0, 0, 0, 0, t.Location())
	return lastMonthStart
}
