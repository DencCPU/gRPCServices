package spotclient

import (
	"context"
	"time"

	spot "github.com/DencCPU/gRPCServices/Protobuf/gen/spot_service"
	"github.com/DencCPU/gRPCServices/Shared/config"
	entryspotservice "github.com/DencCPU/gRPCServices/Shared/enter_points/entry_spot_service"
	spotconfig "github.com/DencCPU/gRPCServices/SpotInstrumentService/config"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type Client struct {
	spot.SpotInstrumentServiceClient
}

func NewClient() (*Client, error) {

	loader := config.NewConfigLoader(
		entryspotservice.GlobalPathToEnv,
		entryspotservice.EnvFile,
		entryspotservice.ConfigType,
		entryspotservice.PathToLocalEnv,
		entryspotservice.PathToConfig,
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cfg, err := config.NewConfig[spotconfig.Config](ctx, loader)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.NewClient(cfg.Server.Host+":"+cfg.Server.Port, grpc.WithInsecure(), grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		return nil, err
	}

	client := spot.NewSpotInstrumentServiceClient(conn)
	return &Client{client}, nil
}
