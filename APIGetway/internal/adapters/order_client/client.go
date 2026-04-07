package orderclient

import (
	"github.com/DencCPU/gRPCServices/OrderService/pkg/orderclient"
	"github.com/DencCPU/gRPCServices/Protobuf/gen/order_service"
)

type Client struct {
	order_service.OrderServiceClient
}

func NewClient() (*Client, error) {
	client, err := orderclient.NewClient()
	if err != nil {
		return &Client{}, err
	}
	return &Client{client}, nil
}
