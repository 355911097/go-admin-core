package config

import (
	"github.com/redis/go-redis/v9"
)

var _redisCluster *redis.ClusterClient

// GetRedisClusterClient 获取 Redis Cluster 客户端
func GetRedisClusterClient() *redis.ClusterClient {
	return _redisCluster
}

// SetRedisClusterClient 设置 Redis Cluster 客户端
func SetRedisClusterClient(c *redis.ClusterClient) {
	if _redisCluster != nil && _redisCluster != c {
		_redisCluster.Close()
	}
	_redisCluster = c
}

// GetRedisClusterOptions 将配置转换为 ClusterClient 所需的选项
func (e RedisConnectOptions) GetRedisClusterOptions() (*redis.ClusterOptions, error) {
	r := &redis.ClusterOptions{
		Addrs:      e.Addrs,
		Username:   e.Username,
		Password:   e.Password,
		MaxRetries: e.MaxRetries,
		PoolSize:   e.PoolSize,
	}
	var err error
	r.TLSConfig, err = getTLS(e.Tls)
	return r, err
}
func GetRedis() redis.UniversalClient {
	if _redis != nil {
		return _redis
	} else {
		return _redisCluster
	}
}
