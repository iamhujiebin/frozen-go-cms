package utils

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"frozen-go-cms/common/resource/config"
	"frozen-go-cms/common/resource/mysql"
	"strconv"
	"strings"
	"time"
)

// 去除slice中的重复元素
func UniqueSliceUInt64(sliceIn []uint64) []uint64 {
	sliceOut := make([]uint64, 0, len(sliceIn))
	m := make(map[uint64]struct{}, len(sliceIn))
	for _, i := range sliceIn {
		if _, ok := m[i]; !ok {
			m[i] = struct{}{}
			sliceOut = append(sliceOut, i)
		}
	}
	return sliceOut
}

func ToString(s interface{}) (string, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return "", nil
	}
	return string(b), nil
}

func SliceToMapUInt64(s []uint64) map[uint64]struct{} {
	m := make(map[uint64]struct{}, len(s))
	for _, i := range s {
		m[i] = struct{}{}
	}
	return m
}

func IfLogoutStr(condition bool, trueVal, falseVal string) string {
	if condition {
		return trueVal
	}
	return falseVal
}

func IfLogoutNick(condition bool, code string, nick string) string {
	if condition {
		return "Hilo No." + code
	}
	return nick
}

func IfLogout(logoutTime int64) bool {
	return logoutTime > 0 && time.Now().Unix() > logoutTime
}

func BirthdayToUint64(birthday *mysql.Timestamp) *uint64 {
	if *birthday == 0 {
		return nil
	}
	return (*uint64)(birthday)
}

// 空字符串转成nil
func StrNil(msg string) *string {
	if msg == "" {
		return nil
	}
	return &msg
}

func TypeToUint8(t *mysql.Type) *uint8 {
	if *t == 0 {
		return nil
	} else {
		return (*uint8)(t)
	}
}

func StrToString(str *mysql.Str) *string {
	return (*string)(str)
}

func NumToUint32(num *mysql.Num) *uint32 {
	return (*uint32)(num)
}

func IsInStringList(str string, list []string) bool {
	for _, v := range list {
		if str == v {
			return true
		}
	}
	return false
}

// 缩短url: 去掉AWS前缀，以便aws接口使用
func StripAwsPrefix(url string) string {
	if !strings.HasPrefix(url, config.GetConfigAws().AWS_CDN) {
		return url
	}
	newUrl := url[len(config.GetConfigAws().AWS_CDN):]
	return newUrl
}

// 保留两位小数
func Decimal(value float64) float64 {
	newValue, err := strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	if err != nil {
		return value
	}
	return newValue
}

func GetMD5Str(str string) string {
	md5.New()
	has := md5.Sum([]byte(str))
	return fmt.Sprintf("%X", has)
}
