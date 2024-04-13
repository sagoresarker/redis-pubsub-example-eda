package publisher

import (
	"context"
	"encoding/json"
	"log"

	"github.com/sagoresarker/redis-pubsub-example-eda/internal/redis"
)

type Message struct {
	Channel string
	Data    interface{}
}

type MessagePublisher struct {
	redisClient redis.Redis
}

func NewMessagePublisher(redisClient redis.Redis) *MessagePublisher {
	return &MessagePublisher{redisClient}
}

func (p *MessagePublisher) PublishMessages(ctx context.Context, message Message) {
	serializedMessage, err := json.Marshal(message.Data)
	if err != nil {
		log.Printf("[%s] Failed to serialized message: %v", message.Channel, err)
		return
	}

	err = p.redisClient.RedisClient.Publish(message.Channel, serializedMessage).Err()
	if err != nil {
		log.Printf("[%s] Failed to publish message: %v", message.Channel, err)
	}
}
