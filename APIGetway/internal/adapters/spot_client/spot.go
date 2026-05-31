package spotclient

import (
	"time"

	"github.com/DencCPU/gRPCServices/Protobuf/gen/spot_service"
	"github.com/DencCPU/gRPCServices/SpotInstrumentService/pkg/spotclient"
	"github.com/sony/gobreaker"
)

type Client struct {
	spot_service.SpotInstrumentServiceClient
	breaker           *gobreaker.CircuitBreaker
	ConnectionTimeout time.Duration
}

func NewClient(breaker *gobreaker.CircuitBreaker, connTimeout time.Duration) (*Client, error) {
	client, err := spotclient.NewClient()
	if err != nil {
		return nil, err
	}
	return &Client{client, breaker, connTimeout}, nil
}
