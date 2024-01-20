package req

import (
	"frozen-go-cms/common/mycontext"
	"frozen-go-cms/common/resource/mysql"
	"frozen-go-cms/myerr/bizerr"
	"github.com/gin-gonic/gin"
	"reflect"
)

func GetUserId(c *gin.Context) (mysql.ID, error) {
	if userIdStr, ok := c.Keys[mycontext.USERID]; ok {
		userId := userIdStr.(uint64)
		return userId, nil
	}
	return 0, bizerr.ParaMissing
}

func GetNonEmptyFields(config interface{}, tagName string) map[string]interface{} {
	result := make(map[string]interface{})

	t := reflect.TypeOf(config)
	v := reflect.ValueOf(config)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// 判断字段是否为空
		if value.IsNil() {
			continue
		}

		// 获取tag为form的注释
		tag := field.Tag.Get(tagName)
		if tag != "" {
			result[tag] = value.Interface()
		}
	}

	return result
}
