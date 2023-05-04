package resp

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code             uint16      `json:"code"`             // 错误码
	Message          interface{} `json:"message"`          // 消息
	MessageData      interface{} `json:"messageData"`      // 消息详情
	OperationMessage interface{} `json:"operationMessage"` // 操作消息
	Data             interface{} `json:"data"`             // 数据
}

type GameResponse struct {
	RetCode      uint16      `json:"ret_code"`
	RetMsg       string      `json:"ret_msg"`
	SdkErrorCode uint16      `json:"sdk_error_code"`
	Data         interface{} `json:"data"`
}

/**
 * HTTP输出json信息
 * param: *gin.Context
 * param: error        err
 * param: interface{}  data
 */
func ResponseOk(c *gin.Context, data interface{}) {
	// always return http.StatusOK
	response := Response{
		Code:             200,
		Message:          "success",
		OperationMessage: "",
		Data:             data,
	}
	printResponseBody(c, &response)

	c.JSON(http.StatusOK, response)
}

func ResponseWaring(c *gin.Context, err error) {
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	response := Response{
		Code:             300,
		Message:          "warn",
		OperationMessage: msg,
		Data:             nil,
	}
	printResponseBody(c, &response)
	c.JSON(http.StatusOK, response)
}

func ResponseSysError(c *gin.Context, err error) {
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	response := Response{
		Code:             500,
		Message:          "error",
		OperationMessage: msg,
		Data:             nil,
	}
	printResponseBody(c, &response)
	c.JSON(http.StatusOK, response)
}

func ResponseErrWithString(c *gin.Context, msg interface{}) {
	response := Response{
		Code:             500,
		Message:          msg,
		OperationMessage: msg,
		Data:             nil,
	}
	printResponseBody(c, &response)
	c.JSON(http.StatusOK, response)
}

func printResponseBody(c *gin.Context, response interface{}) {
	traceId, _ := c.Get("traceid")
	if _traceId, ok := traceId.(string); ok {
		c.Header("X-Trace-ID", _traceId)
	}

	var userId uint64 = 0
	if strUserId, ok := c.Get("userid"); ok {
		userId = strUserId.(uint64)
	}

	buf, err := json.Marshal(response)
	body := ""
	if len(buf) < 1024 {
		body = string(buf)
	} else {
		body = string(buf[0:1024])
	}

	if err == nil {
		fmt.Printf("request rsp url:%s, traceId:%v, userId:%d, body:%s", c.Request.RequestURI, traceId, userId, body)
	} else {
		fmt.Printf("request rsp url:%s, traceId:%v, userId:%d, body:%s", c.Request.RequestURI, traceId, userId, body)
	}
}
