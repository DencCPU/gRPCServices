package userclient

import (
	"time"

	user "github.com/DencCPU/gRPCServices/Protobuf/gen/user_service"
	"github.com/DencCPU/gRPCServices/UserService/pkg/userclient"
	"github.com/sony/gobreaker"
)

type Client struct {
	user.UserServiceClient
	breaker           *gobreaker.CircuitBreaker
	ConnectionTimeout time.Duration
}

func NewClient(breaker *gobreaker.CircuitBreaker, connTimeout time.Duration) (*Client, error) {
	client, err := userclient.NewClient()
	if err != nil {
		return nil, err
	}
	return &Client{client, breaker, connTimeout}, nil
}
