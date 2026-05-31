package spotservice

import (
	"time"

	spot "github.com/DencCPU/gRPCServices/Protobuf/gen/spot_service"
	"github.com/DencCPU/gRPCServices/SpotInstrumentService/pkg/spotclient"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

type Client struct {
	spot.SpotInstrumentServiceClient
	breaker           *gobreaker.CircuitBreaker
	ConnectionTimeout time.Duration
}

func NewClient(logger *zap.Logger, breaker *gobreaker.CircuitBreaker, connTimeout time.Duration) (*Client, error) {
	client, err := spotclient.NewClient()
	if err != nil {
		return nil, err
	}
	return &Client{client, breaker, connTimeout}, nil
}
