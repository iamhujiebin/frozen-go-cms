package route

import (
	"frozen-go-cms/resp"
	"github.com/gin-gonic/gin"
)

/**
controller层全局异常处理
*/

// 等级最高，为了只为最后有返回值到前端
func ExceptionHandle(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			//打印错误堆栈信息
			resp.ResponseErrWithString(c, r)
			//终止后续接口调用，不加的话recover到异常后，还会继续执行接口里后续代码
			c.Abort()
		}
	}()
	c.Next()
}

// LoggerHandle 日志Handle
func LoggerHandle(c *gin.Context) {
	c.Next()
}
