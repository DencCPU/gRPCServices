package orderclient

import (
	"time"

	"github.com/DencCPU/gRPCServices/OrderService/pkg/orderclient"
	"github.com/DencCPU/gRPCServices/Protobuf/gen/order_service"
	"github.com/sony/gobreaker"
)

type Client struct {
	order_service.OrderServiceClient
	breaker           *gobreaker.CircuitBreaker
	ConnectionTimeout time.Duration
}

func NewClient(breaker *gobreaker.CircuitBreaker, connTimeout time.Duration) (*Client, error) {
	client, err := orderclient.NewClient()
	if err != nil {
		return &Client{}, err
	}
	return &Client{client, breaker, connTimeout}, nil
}
