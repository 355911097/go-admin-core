package config

import (
	"github.com/355911097/go-admin-core/redisqueue"
	"github.com/355911097/go-admin-core/storage"
	"github.com/355911097/go-admin-core/storage/queue"
	"time"
)

type Queue struct {
	Redis  *QueueRedis
	Memory *QueueMemory
	NSQ    *QueueNSQ `json:"nsq" yaml:"nsq"`
}

type QueueRedis struct {
	RedisConnectOptions
	Producer *redisqueue.ProducerOptions
	Consumer *redisqueue.ConsumerOptions
}

type QueueMemory struct {
	PoolSize uint
}

type QueueNSQ struct {
	NSQOptions
	ChannelPrefix string
}

var QueueConfig = new(Queue)

// Empty 空设置
func (e Queue) Empty() bool {
	return e.Memory == nil && e.Redis == nil && e.NSQ == nil
}

// Setup 启用顺序 redis > 其他 > memory
func (e Queue) Setup() (storage.AdapterQueue, error) {
	if e.Redis != nil {
		e.Redis.Consumer.ReclaimInterval = e.Redis.Consumer.ReclaimInterval * time.Second
		e.Redis.Consumer.BlockingTimeout = e.Redis.Consumer.BlockingTimeout * time.Second
		e.Redis.Consumer.VisibilityTimeout = e.Redis.Consumer.VisibilityTimeout * time.Second
		client := GetRedis()
		if client == nil {
			if e.Redis.Addr != "" {
				options, err := e.Redis.RedisConnectOptions.GetRedisOptions()
				if err != nil {
					return nil, err
				}
				c := redis.NewClient(options)
				_redis = c
			} else {
				options, err := e.Redis.RedisConnectOptions.GetRedisClusterOptions()
				if err != nil {
					return nil, err
				}
				c := redis.NewClusterClient(options)
				_redisCluster = c
			}
		}
		e.Redis.Producer.RedisClient = client
		e.Redis.Consumer.RedisClient = client
		return queue.NewRedis(e.Redis.Producer, e.Redis.Consumer)
	}
	if e.NSQ != nil {
		cfg, err := e.NSQ.GetNSQOptions()
		if err != nil {
			return nil, err
		}
		return queue.NewNSQ(e.NSQ.Addresses, cfg, e.NSQ.ChannelPrefix)
	}
	return queue.NewMemory(e.Memory.PoolSize), nil
}
