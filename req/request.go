package req

import (
	"frozen-go-cms/myerr/bizerr"
	"git.hilo.cn/hilo-common/mycontext"
	"git.hilo.cn/hilo-common/resource/mysql"
	"github.com/gin-gonic/gin"
)

func GetUserId(c *gin.Context) (mysql.ID, error) {
	if userIdStr, ok := c.Keys[mycontext.USERID]; ok {
		userId := userIdStr.(uint64)
		return userId, nil
	}
	return 0, bizerr.ParaMissing
}
