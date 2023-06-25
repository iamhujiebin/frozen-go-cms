package redisCli

import (
	"context"
	"frozen-go-cms/hilo-common/mylogrus"
	"github.com/go-redis/redis/v8"
	"time"
)

//这个用户避免多个服务器并发问题。
func SetNX(key string, value interface{}, expiration time.Duration, callBack func()) {
	flag, err := RedisClient.SetNX(context.Background(), key, value, expiration).Result()
	if err != nil {
		mylogrus.MyLog.Errorf("key:%v lock start setNx err: %v", key, err)
	}
	if !flag {
		mylogrus.MyLog.Infof("key:%v lock setNx has lock", key)
		return
	}
	mylogrus.MyLog.Infof("key:%v lock setNx begin", key)
	callBack()
	//执行结束之后，移除key
	//RedisClient.Del(context.Background(), key)
	mylogrus.MyLog.Infof("key:%v lock setNx end", key)
}

func Lock(key string, expiration time.Duration) bool {
	flag, err := RedisClient.SetNX(context.Background(), key, 1, expiration).Result()
	if err != nil {
		return false
	}
	if !flag {
		return false
	}
	return true
}

func GetCacheInt64(key string) (int64, error) {
	data, err := RedisClient.Get(context.Background(), key).Int64()
	if err != nil && err != redis.Nil {
		return 0, err
	}
	return data, nil
}

func IncrBy(key string, num int64) (int64, error) {
	resNum, err := RedisClient.IncrBy(context.Background(), key, num).Result()
	if err != nil {
		return 0, err
	}
	return resNum, nil
}

func IncrNumExpire(key string, num int64, expiration time.Duration) (int64, error) {
	times, err := IncrBy(key, num)
	if err != nil {
		return 0, err
	}
	ttl, err := RedisClient.TTL(context.Background(), key).Result()
	if err != nil {
		return 0, err
	}
	if ttl == -1 {
		RedisClient.Expire(context.Background(), key, expiration)
	}
	return times, nil
}

func DelCache(key string) error {
	err := RedisClient.Del(context.Background(), key).Err()
	if err != nil {
		mylogrus.MyLog.Errorf("DelCache key:%s, err:%s", key, err)
		return err
	}
	return nil
}
