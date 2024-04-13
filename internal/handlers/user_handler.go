package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sagoresarker/redis-pubsub-example-eda/internal/models"
	"github.com/sagoresarker/redis-pubsub-example-eda/internal/publisher"
)

type UserHandler struct {
	publisher *publisher.MessagePublisher
}

func NewUserHandler(pub *publisher.MessagePublisher) *UserHandler {
	return &UserHandler{
		publisher: pub,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.publisher.PublishMessages(c.Request.Context(), publisher.Message{
		Channel: "user.created",
		Data:    user,
	})

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.publisher.PublishMessages(c.Request.Context(), publisher.Message{
		Channel: "user.updated",
		Data:    map[string]interface{}{"id": id, "user": user},
	})

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
