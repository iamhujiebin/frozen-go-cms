package route

import (
	"frozen-go-cms/hilo-common/mycontext"
	"frozen-go-cms/resp"
	"github.com/gin-gonic/gin"
)

/**
 * 主要是解决错误的统一处理
 */
type HandlerFunc func(c *gin.Context) (*mycontext.MyContext, error)

// 对错误进行处理，
func wrapper(handler HandlerFunc) func(c *gin.Context) {
	return func(c *gin.Context) {
		var err error
		var myContext *mycontext.MyContext
		myContext, err = handler(c)
		if err != nil {
			reqUri := c.Request.RequestURI
			method := c.Request.Method

			switch err.(type) {
			default:
				// 注意这里，如果是原生的error, 可能打印不出来,使用errors.Wrap配合%+v可以打印堆栈信息,建议上游使用
				myContext.Log.Errorf("request err -> url:%v, method:%v, err :%+v\n", reqUri, method, err)
				resp.ResponseErrWithString(c, err.Error())
			}
		}
		return
	}
}
