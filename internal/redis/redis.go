package redis

import (
	"github.com/go-redis/redis"
	"github.com/sagoresarker/redis-pubsub-example-eda/config"
)

type Redis struct {
	RedisClient *redis.Client
}

func NewRedis(config config.Config) Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost + ":" + config.RedisPort,
		Password: config.RedisPassword,
		DB:       0,
	})

	return Redis{
		RedisClient: client,
	}

}
