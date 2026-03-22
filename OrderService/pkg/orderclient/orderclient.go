package orderclient

import (
	orderconfig "Academy/gRPCServices/OrderService/config"
	orderAPI "Academy/gRPCServices/Protobuf/gen/order"

	"Academy/gRPCServices/Shared/config"
	"fmt"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type Client struct {
	orderAPI.OrderServiceClient
}

// УБРАТЬ В ОТДЕЛЬНЮ ПАПКУ
const (
	globalPathToEnv = "."    //Дериктория, в которой находиться общий файл env
	envFile         = ".env" //Название файла .env
	configType      = "yaml" //Тип конфиг файла (yaml)

	pathToLocalEnv = "PATH_ORDERSERVICE_CONFIG_ENV" //Переменная окружения, в которой лежит путь в локально env файлу
	pathToConfig   = "ORDER_CONFIG_PATH"            //Переменная окружения, в которой лежит путь к конфигу
)

func NewClient() (*Client, error) {
	loader := config.NewConfigLoader(globalPathToEnv, envFile, configType, pathToLocalEnv, pathToConfig)
	cfg, err := config.NewConfig[orderconfig.Config](loader)
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

	client := orderAPI.NewOrderServiceClient(conn)
	return &Client{client}, nil
}
