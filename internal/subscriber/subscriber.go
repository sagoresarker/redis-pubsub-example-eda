package subscriber

import (
	"context"
	"encoding/json"
	"log"

	redisMain "github.com/go-redis/redis"
	"github.com/sagoresarker/redis-pubsub-example-eda/internal/redis"
)

type MessageConsumer struct {
	redisClient  redis.Redis
	subscription *redisMain.PubSub
}

func NewMessageConsumer(redis redis.Redis) *MessageConsumer {
	return &MessageConsumer{
		redisClient: redis,
	}
}

func (c *MessageConsumer) ConsumerMessages(ctx context.Context, channels []string) {
	for _, channel := range channels {
		go c.handleCustomType1Logic(ctx, channel)
	}
}

func (c *MessageConsumer) handleCustomType1Logic(ctx context.Context, channel string) {
	consumerCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Printf("[%s] Consumer started listening...\n", channel)

	c.subscription = c.redisClient.RedisClient.Subscribe(channel)
	defer c.subscription.Close()

	messageChannel := c.subscription.Channel()

	for {
		select {
		case <-consumerCtx.Done():
			log.Printf("[%s] Consumer stopped listening...\n", channel)
			return
		case msg := <-messageChannel:
			var messageData interface{}
			err := json.Unmarshal([]byte(msg.Payload), &messageData)
			if err != nil {
				log.Printf("[%s] Failed to deserialize message: %v", channel, err)
				continue
			}

			switch channel {
			case "user.created":
				if userData, ok := messageData.(map[string]interface{}); ok {
					log.Printf("[%s] User created: %+v\n", channel, userData)
				} else {
					log.Printf("[%s] Unsupported message format: %+v\n", channel, messageData)
				}
			case "user.updated":
				if userData, ok := messageData.(map[string]interface{}); ok {
					log.Printf("[%s] User updated: ID=%v, User=%+v\n", channel, userData["id"], userData["user"])
				} else {
					log.Printf("[%s] Unsupported message format: %+v\n", channel, messageData)
				}
			case "order.created":
				if orderData, ok := messageData.(map[string]interface{}); ok {
					log.Printf("[%s] Order created: %+v\n", channel, orderData)
				} else {
					log.Printf("[%s] Unsupported message format: %+v\n", channel, messageData)
				}
			case "order.updated":
				if orderData, ok := messageData.(map[string]interface{}); ok {
					log.Printf("[%s] Order updated: ID=%v, Status=%v\n", channel, orderData["id"], orderData["status"])
				} else {
					log.Printf("[%s] Unsupported message format: %+v\n", channel, messageData)
				}
			default:
				log.Printf("[%s] Unsupported message type: %+v\n", channel, messageData)
			}
		}
	}
}
