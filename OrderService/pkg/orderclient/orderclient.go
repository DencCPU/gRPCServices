package orderclient

import (
	"context"
	"fmt"
	"time"

	orderconfig "github.com/DencCPU/gRPCServices/OrderService/config"
	order "github.com/DencCPU/gRPCServices/Protobuf/gen/order_service"
	"github.com/DencCPU/gRPCServices/Shared/config"
	entryorderservice "github.com/DencCPU/gRPCServices/Shared/enter_points/entry_order_service"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type Client struct {
	order.OrderServiceClient
}

func NewClient() (*Client, error) {
	loader := config.NewConfigLoader(
		entryorderservice.GlobalPathToEnv,
		entryorderservice.EnvFile,
		entryorderservice.ConfigType,
		entryorderservice.PathToLocalEnv,
		entryorderservice.PathToConfig,
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cfg, err := config.NewConfig[orderconfig.Config](ctx, loader)
	if err != nil {

		return nil, fmt.Errorf("Ошибка получения конфига:%w", err)

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

	client := order.NewOrderServiceClient(conn)
	return &Client{client}, nil
}
