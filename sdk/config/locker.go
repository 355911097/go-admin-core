package config

import (
	"github.com/355911097/go-admin-core/storage"
	"github.com/355911097/go-admin-core/storage/locker"
)

var LockerConfig = new(Locker)

type Locker struct {
	Redis *RedisConnectOptions
}

// Empty 空设置
func (e Locker) Empty() bool {
	return e.Redis == nil
}

// Setup 启用顺序 redis > 其他 > memory
func (e Locker) Setup() (storage.AdapterLocker, error) {
	if e.Redis != nil {
		client := GetRedis()
		if client == nil {
			if e.Redis.Addr != "" {
				options, err := e.Redis.GetRedisOptions()
				if err != nil {
					return nil, err
				}
				c := redis.NewClient(options)
				_redis = c
			} else {
				options, err := e.Redis.GetRedisClusterOptions()
				if err != nil {
					return nil, err
				}
				c := redis.NewClusterClient(options)
				_redisCluster = c
			}
		}
		return locker.NewRedis(&client), nil
	}
	return nil, nil
}
