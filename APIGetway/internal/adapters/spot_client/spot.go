package spotclient

import (
	"github.com/DencCPU/gRPCServices/Protobuf/gen/spot_service"
	"github.com/DencCPU/gRPCServices/SpotInstrumentService/pkg/spotclient"
)

type Client struct {
	spot_service.SpotInstrumentServiceClient
}

func NewClient() (*Client, error) {
	client, err := spotclient.NewClient()
	if err != nil {
		return nil, err
	}
	return &Client{client}, nil
}
