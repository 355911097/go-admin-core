package locker

import (
	"time"

	"github.com/355911097/go-admin-core/redislock"
)

// NewRedis 初始化locker
func NewRedis(c *redis.UniversalClient) *Redis {
	return &Redis{
		client: c,
	}
}

type Redis struct {
	client *redis.UniversalClient
	mutex  *redislock.Client
}

func (Redis) String() string {
	return "redis"
}

func (r *Redis) Lock(key string, ttl int64, options *redislock.Options) (*redislock.Lock, error) {
	if r.mutex == nil {
		r.mutex = redislock.New(r.client)
	}
	return r.mutex.Obtain(key, time.Duration(ttl)*time.Second, options)
}
