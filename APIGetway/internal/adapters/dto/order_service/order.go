package orderdto

import (
	"time"
)

type GetInput struct {
	UserId  string `json:"user_id"`
	OrderId string `json:"order_id"`
}

type GetOutput struct {
	Status     string
	OrderId    string
	Price      string
	Quantity   int64
	MarketName string
}

type StreamOutput struct {
	OrderStatus string
	UpdateTime  time.Time
}

type Output struct {
	OrderId     string
	OrderStatus string
}

func NewOutput(orderId, orderStatus string) Output {
	return Output{orderId, orderStatus}
}
