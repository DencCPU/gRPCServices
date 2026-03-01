package orderclient

import (
	serverconfig "Academy/gRPCServices/OrderService/config/server"
	orderAPI "Academy/gRPCServices/Protobuf/gen/order"
	"fmt"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type Client struct {
	orderAPI.OrderServiceClient
}

func NewClient() (*Client, error) {
	cfg, err := serverconfig.NewConfig()
	if err != nil {
		return nil, err
	}

	host := cfg.Server.Host
	port := cfg.Server.Port
	dsn := fmt.Sprintf("%s:%d", host, port)

	conn, err := grpc.NewClient(dsn,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		return nil, err
	}

	client := orderAPI.NewOrderServiceClient(conn)
	return &Client{client}, nil
}
