package myerr

import (
	"fmt"
	"git.hilo.cn/hilo-common/mylogrus"
	"github.com/pkg/errors"
	"strconv"
)

//业务错误
type BusinessError struct {
	code    uint16
	message string
	err     error
	data    BusinessData
}

func (err *BusinessError) Error() string {
	return err.err.Error()
}

func (err *BusinessError) GetErr() error {
	return err.err
}

func (err *BusinessError) GetCode() uint16 {
	return err.code
}

func (err *BusinessError) GetMsg() string {
	return err.message
}

func (err *BusinessError) GetData() BusinessData {
	return err.data
}

var codes = map[uint16]string{}

//定义必须是明确的。不可以修改，字段等同于翻译中要替换的字符
type BusinessData struct {
	//剩余秒
	Second int `json:"second"`
	//所需数量
	Num       int    `json:"num"`
	Code      string `json:"code"`
	Timestamp int64  `json:"timestamp"`
	//官网充值地址
	CheckOutUrl string `json:"checkOutUrl"`
}

func NewBusiness(err *BusinessError) *BusinessError {
	return &BusinessError{
		code:    err.code,
		message: err.message,
		err:     err.err,
		data:    err.data,
	}
}

func NewBusinessCode(code uint16, msg string, data BusinessData) *BusinessError {
	if _, ok := codes[code]; ok {
		mylogrus.MyLog.Error(fmt.Sprintf("错误码 %d 已经存在，请更换一个", code))
		return nil
	}
	codes[code] = msg
	return &BusinessError{
		code:    code,
		message: msg,
		err:     errors.New("{code:" + strconv.Itoa(int(code)) + ",message:" + msg + "}"),
		data:    data,
	}
}

func NewBusinessCodeNoCheck(code uint16, msg string, data BusinessData) *BusinessError {
	return &BusinessError{
		code:    code,
		message: msg,
		err:     errors.New("{code:" + strconv.Itoa(int(code)) + ",message:" + msg + "}"),
		data:    data,
	}
}