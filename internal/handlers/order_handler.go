package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sagoresarker/redis-pubsub-example-eda/internal/models"
	"github.com/sagoresarker/redis-pubsub-example-eda/internal/publisher"
)

type OrderHandler struct {
	publisher *publisher.MessagePublisher
}

func NewOrderHandler(pub *publisher.MessagePublisher) *OrderHandler {
	return &OrderHandler{
		publisher: pub,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.publisher.PublishMessages(c.Request.Context(), publisher.Message{
		Channel: "order.created",
		Data:    order,
	})

	c.JSON(http.StatusCreated, gin.H{"message": "Order created successfully"})
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var status models.OrderStatus
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.publisher.PublishMessages(c.Request.Context(), publisher.Message{
		Channel: "order.updated",
		Data:    map[string]interface{}{"id": id, "status": status},
	})
	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}
