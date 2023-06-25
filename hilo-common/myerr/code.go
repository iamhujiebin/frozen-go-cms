package myerr

import (
	"fmt"
	"frozen-go-cms/hilo-common/mylogrus"
	"github.com/pkg/errors"
	"strconv"
)

//成功
type Success struct {
	code    uint16
	message string
}

func (err *Success) Error() string {
	return err.message
}

//正确的标识符
var success = &Success{code: 200, message: "OK"}

func GetSuccessCode() uint16 {
	return success.code
}

func GetSuccessMsg() string {
	return success.message
}

func GetSuccess() Success {
	return *success
}

//系统错误
type SysError struct {
	code    uint16
	message string
	err     error
}

var sysError = &SysError{
	code:    500,
	message: "",
}

func (err *SysError) Error() string {
	return err.err.Error()
}

func NewSysError(msg string) *SysError {
	return &SysError{
		code:    sysError.code,
		message: msg,
		err:     errors.New("{code:" + strconv.Itoa(int(sysError.code)) + ",message:" + msg + "}"),
	}
}

func NewSysErrorF(format string, args ...interface{}) *SysError {
	return NewSysError(fmt.Sprintf(format, args...))
}

func GetSysErrorCode() uint16 {
	return sysError.code
}

func (sysError *SysError) GetErr() error {
	return sysError.err
}

func (sysError *SysError) GetCode() uint16 {
	return sysError.code
}

func (sysError *SysError) GetMsg() string {
	return sysError.message
}

//警告错误
type WaringError struct {
	code    uint16
	message string
	err     error
}

var waringError = &WaringError{
	code:    300,
	message: "",
}

func GetWaringCode() uint16 {
	return waringError.code
}

func (err *WaringError) Error() string {
	return err.err.Error()
}

func NewWaring(msg string) *WaringError {
	return &WaringError{
		code:    waringError.code,
		message: msg,
		err:     errors.New("{code:" + strconv.Itoa(int(waringError.code)) + ",message:" + msg + "}"),
	}
}

func NewWaringErrorF(format string, args ...interface{}) *WaringError {
	return NewWaring(fmt.Sprintf(format, args...))
}

func (err *WaringError) GetErr() error {
	return err.err
}

func (err *WaringError) GetCode() uint16 {
	return err.code
}

func (err *WaringError) GetMsg() string {
	return err.message
}

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

//包装日志，让日志成堆栈状态
func WrapErrWithStr(err error, msg string) error {
	if h, ok := err.(*BusinessError); ok {
		h.err = errors.Wrap(h.err, msg)
		return h
	} else if h, ok := err.(*WaringError); ok {
		h.err = errors.Wrap(h.err, msg)
		return h
	} else if h, ok := err.(*SysError); ok {
		h.err = errors.Wrap(h.err, msg)
		return h
	} else {
		return errors.Wrap(err, msg)
	}
}

func WrapErr(err error) error {
	if h, ok := err.(*BusinessError); ok {
		h1 := NewBusiness(h)
		h1.err = errors.Wrap(h1.err, "")
		return h1
	} else if h, ok := err.(*WaringError); ok {
		h.err = errors.Wrap(h.err, "")
		return h
	} else if h, ok := err.(*SysError); ok {
		h.err = errors.Wrap(h.err, "")
		return h
	} else {
		return errors.Wrap(err, "")
	}
}
