package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/sagoresarker/redis-pubsub-example-eda/config"
	"github.com/sagoresarker/redis-pubsub-example-eda/internal/handlers"
	"github.com/sagoresarker/redis-pubsub-example-eda/internal/publisher"
	"github.com/sagoresarker/redis-pubsub-example-eda/internal/redis"
	"github.com/sagoresarker/redis-pubsub-example-eda/internal/subscriber"
)

func main() {
	cfg := config.LoadConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	redisClient := redis.NewRedis(cfg)
	defer redisClient.RedisClient.Close()

	pub := publisher.NewMessagePublisher(redisClient)

	sub := subscriber.NewMessageConsumer(redisClient)

	go sub.ConsumerMessages(ctx, []string{"user.created", "user.updated", "order.created", "order.updated"})

	router := gin.Default()

	userHandler := handlers.NewUserHandler(pub)
	orderHandler := handlers.NewOrderHandler(pub)

	router.POST("/users", userHandler.CreateUser)
	router.PUT("/users/:id", userHandler.UpdateUser)
	router.POST("/orders", orderHandler.CreateOrder)
	router.PUT("/orders/:id/status", orderHandler.UpdateOrderStatus)

	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Wait for interrupt signal to gracefully shut down
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")
}
