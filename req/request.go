package req

import (
	"frozen-go-cms/hilo-common/mycontext"
	"frozen-go-cms/hilo-common/resource/mysql"
	"frozen-go-cms/myerr/bizerr"
	"github.com/gin-gonic/gin"
)

func GetUserId(c *gin.Context) (mysql.ID, error) {
	if userIdStr, ok := c.Keys[mycontext.USERID]; ok {
		userId := userIdStr.(uint64)
		return userId, nil
	}
	return 0, bizerr.ParaMissing
}
