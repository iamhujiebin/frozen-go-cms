package redisCli

import (
	"context"
	"frozen-go-cms/common/mylogrus"
	"frozen-go-cms/common/resource/config"
	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         config.GetConfigRedis().REDIS_HOST,
		Password:     config.GetConfigRedis().REDIS_PASSWORD, // no password set
		DB:           0,                                      // use default DB
		PoolSize:     20,
		MinIdleConns: 20,
	})
	mylogrus.MyLog.Infoln(config.GetConfigRedis().REDIS_HOST)
	mylogrus.MyLog.Infoln(config.GetConfigRedis().REDIS_PASSWORD)
	pong, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		mylogrus.MyLog.Warn(err)
		mylogrus.MyLog.Fatal("redis db0 connect fail")
	} else {
		mylogrus.MyLog.Info("redis db0 connection success - ", pong)
	}
}

func GetRedis() *redis.Client {
	return RedisClient
}
