package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewRedis redis模式
func NewRedisCluster(client *redis.ClusterClient, options *redis.ClusterOptions) (*RedisCluster, error) {
	if client == nil {
		client = redis.NewClusterClient(options)
	}
	r := &RedisCluster{
		client: client,
	}
	err := r.connect()
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Redis cache implement
type RedisCluster struct {
	client *redis.ClusterClient
}

func (*RedisCluster) String() string {
	return "redis"
}

// connect connect test
func (r *RedisCluster) connect() error {
	var err error
	_, err = r.client.Ping(context.Background()).Result()
	return err
}

// Get from key
func (r *RedisCluster) Get(key string) (string, error) {
	return r.client.Get(context.Background(), key).Result()
}

// Set value with key and expire time
func (r *RedisCluster) Set(key string, val interface{}, expire int) error {
	return r.client.Set(context.Background(), key, val, time.Duration(expire)*time.Second).Err()
}

// Del delete key in redis
func (r *RedisCluster) Del(key string) error {
	return r.client.Del(context.Background(), key).Err()
}

// HashGet from key
func (r *RedisCluster) HashGet(hk, key string) (string, error) {
	return r.client.HGet(context.Background(), hk, key).Result()
}

// HashDel delete key in specify redis's hashtable
func (r *RedisCluster) HashDel(hk, key string) error {
	return r.client.HDel(context.Background(), hk, key).Err()
}

// Increase
func (r *RedisCluster) Increase(key string) error {
	return r.client.Incr(context.Background(), key).Err()
}

func (r *RedisCluster) Decrease(key string) error {
	return r.client.Decr(context.Background(), key).Err()
}

// Set ttl
func (r *RedisCluster) Expire(key string, dur time.Duration) error {
	return r.client.Expire(context.Background(), key, dur).Err()
}

// GetClient 暴露原生client
func (r *RedisCluster) GetClient() *redis.ClusterClient {
	return r.client
}
