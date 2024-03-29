package route

import (
	"frozen-go-cms/common/mycontext"
	"frozen-go-cms/common/mylogrus"
	"frozen-go-cms/common/resource/config"
	"frozen-go-cms/domain/model/casbin_m"
	"frozen-go-cms/myerr/bizerr"
	"frozen-go-cms/req/jwt"
	"frozen-go-cms/resp"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
	"strings"
	"time"
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

// jwt解密
func JWTApiHandle(c *gin.Context) {
	logger := mylogrus.MyLog.WithField("URL", c.Request.URL).WithField("METHOD", c.Request.Method)
	token := c.GetHeader("Authorization")
	if token == "" {
		logger.Warnf("token err is empty! %v", c.Request.Header)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if len(strings.Split(token, " ")) == 2 {
		token = strings.Split(token, " ")[1]
	} else {
		logger.Warnf("token len is wrong! %v", c.Request.Header)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	claims, err := jwt.ParseToken(token)
	if err != nil {
		logger.Warnf("token parsed err:%v", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	logger = logger.WithField("userId", claims.UserId)
	if time.Now().Unix() > claims.ExpiresAt {
		logger.Warnf("token expire err, now: %d, expiresAt %d", time.Now().Unix(), claims.ExpiresAt)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if claims.Issuer != config.GetConfigJWT().ISSUER_API {
		logger.Warnf("token err issuer:%s, configIssuer %s", claims.Issuer, config.GetConfigJWT().ISSUER_API)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var newToken = token
	// token 连续7天没玩,第八天回来后给新token(线上是30天过期)
	if claims.ExpiresAt-time.Now().Unix() < 86400*7 {
		logger.Infof("token nearly expire err, now:%d,expiresAt:%d", time.Now().Unix(), claims.ExpiresAt)
		newToken, err = jwt.GenerateToken(claims.UserId, claims.Mobile, config.GetConfigJWT().ISSUER_API)
		if err != nil {
			logger.Warnf("token generation failed, err:%v", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}

	c.Set(mycontext.USERID, claims.UserId)

	c.Writer.Header().Add("token", newToken)
	c.Next()
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func CasbinHandle(c *gin.Context) {
	// 获取请求的URI
	obj := c.Request.URL.RequestURI()
	// 获取请求方法
	act := c.Request.Method
	// 获取用户的角色
	sub, _ := c.Get(mycontext.USERID)
	e := casbin_m.Casbin()
	// 判断策略中是否存在
	success, err := e.Enforce(cast.ToString(sub), obj, act)
	if err != nil {
		resp.ResponseBusiness(c, bizerr.AuthFail)
		c.Abort()
		return
	}
	if success {
		c.Next()
	} else {
		resp.ResponseBusiness(c, bizerr.AuthFail)
		c.Abort()
		return
	}
}

func SuperUser(c *gin.Context) {
	userId, _ := c.Get(mycontext.USERID)
	if cast.ToInt(userId) != 1 {
		resp.ResponseBusiness(c, bizerr.AuthFail)
		c.Abort()
	}
}
