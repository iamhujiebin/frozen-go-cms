package domain

import (
	"frozen-go-cms/common/mycontext"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type CtxAndDb struct {
	Db *gorm.DB
	*mycontext.MyContext
	Redis *redis.Client
}
