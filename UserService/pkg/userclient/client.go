package userclient

import (
	"context"
	"fmt"
	"time"

	user "github.com/DencCPU/gRPCServices/Protobuf/gen/user_service"
	"github.com/DencCPU/gRPCServices/Shared/config"
	entryuserservice "github.com/DencCPU/gRPCServices/Shared/enter_points/entry_user_service"
	userconfig "github.com/DencCPU/gRPCServices/UserService/config"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type Client struct {
	user.UserServiceClient
}

func NewClient() (*Client, error) {
	loader := config.NewConfigLoader(
		entryuserservice.GlobalPathToEnv,
		entryuserservice.EnvFile,
		entryuserservice.ConfigType,
		entryuserservice.PathToLocalEnv,
		entryuserservice.PathToConfig,
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cfg, err := config.NewConfig[userconfig.Config](ctx, loader)
	if err != nil {

		return nil, fmt.Errorf("Ошибка получения конфига:%w", err)

	}

	host := cfg.Server.Host
	port := cfg.Server.Port
	dsn := fmt.Sprintf("%s:%s", host, port)

	conn, err := grpc.NewClient(
		dsn,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		return nil, err
	}

	client := user.NewUserServiceClient(conn)
	return &Client{client}, nil
}
