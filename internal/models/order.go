package models

type Order struct {
	ID     string      `json:"id"`
	UserID string      `json:"user_id"`
	Items  []Item      `json:"items"`
	Total  float64     `json:"total"`
	Status OrderStatus `json:"status"`
}

type Item struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
)
