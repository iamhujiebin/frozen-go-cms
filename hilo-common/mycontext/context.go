package mycontext

import (
	"context"
	"frozen-go-cms/hilo-common/mylogrus"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

const (
	TRACEID       = "traceId"
	USERID        = "userId"
	EXTERNAL_ID   = "externalId"
	CODE          = "code"
	NICK          = "nick"
	AVATAR        = "avatar"
	COUNTRY       = "country"
	EXTERNALID1   = "externalId1"
	EXTERNALID2   = "externalId2"
	MGRID         = "mgrId"
	DEVICETYPE    = "deviceType"
	DEVICEVERSION = "deviceVersion"
	APP_VERSION   = "appVersion"
	ACTION_RESULt = "actionResult"
	URL           = "url"
	METHOD        = "method"
	IMEI          = "imei"
	LANGUAGE      = "language"
	TOKEN         = "token"
	CARRIER       = "carrier"  // 运营商,460 开头是中国,如:46001,46007
	TIMEZONE      = "timeZone" // 时区 GMT+8 / GMT+8:00

	InnerEncrypt = "innerEncrypt" // 加密用的key,服务端内部用
)

/**
 * 主要是完成日志打印
 * @param
 * @return
 **/

type MyContext struct {
	context.Context
	Log *logrus.Entry
	Cxt map[string]interface{}
}

func CreateMyContextWith(traceId interface{}) *MyContext {
	cxt := map[string]interface{}{}
	cxt[TRACEID] = traceId
	return CreateMyContext(cxt)
}

func CreateMyContext(ctx map[string]interface{}) *MyContext {
	var traceId string
	if traceIdTemp, ok := ctx[TRACEID]; ok {
		traceId, ok = traceIdTemp.(string)
	} else {
		traceId = strings.Replace(uuid.NewV4().String(), "-", "", -1)
	}

	var userId string
	if userIdTemp, ok := ctx[USERID]; ok {
		userId = strconv.FormatUint(userIdTemp.(uint64), 10)
	}

	var deviceTypeStr string
	if deviceTypeTemp, ok := ctx[DEVICETYPE]; ok {
		deviceTypeStr, ok = deviceTypeTemp.(string)
	}

	var sAppVersion string
	if appVersionTmp, ok := ctx[APP_VERSION]; ok {
		sAppVersion, ok = appVersionTmp.(string)
	}

	var url string
	if urlTmp, ok := ctx[URL]; ok {
		url, ok = urlTmp.(string)
	}

	var method string
	if methodTmp, ok := ctx[METHOD]; ok {
		method, ok = methodTmp.(string)
	}
	_ctx := context.WithValue(context.Background(), "traceId", traceId)
	_ctx = context.WithValue(_ctx, "userId", userId)
	return &MyContext{
		Context: _ctx,
		Log:     CreateContextLog(userId, traceId, deviceTypeStr, sAppVersion, url, method),
		Cxt:     ctx,
	}
}

/**
 * 创建上下文的日志
 **/
func CreateContextLog(userId string, traceId string, deviceType string, deviceVersion string, url string, method string) *logrus.Entry {
	return mylogrus.MyLog.WithFields(logrus.Fields{
		USERID:      userId,
		TRACEID:     traceId,
		DEVICETYPE:  deviceType,
		APP_VERSION: deviceVersion,
		URL:         url,
		METHOD:      method,
	})
}
