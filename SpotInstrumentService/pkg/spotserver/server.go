package grpcserver

import (
	"net"
	"time"

	"github.com/DencCPU/gRPCServices/Shared/interceptors"
	spotconfig "github.com/DencCPU/gRPCServices/SpotInstrumentService/config"
	redisadapter "github.com/DencCPU/gRPCServices/SpotInstrumentService/internal/adapters/redis"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

// Структура сервера
type Server struct {
	*grpc.Server              //Сервер
	Listener     net.Listener //Настройки подключения
}

// Создание нового сервера
func New(redis *redisadapter.RedisDB, cfg spotconfig.Server, logger *zap.Logger) (*Server, error) {

	//Инициализация интерфеса listener
	host := cfg.Host
	port := cfg.Port
	lis, err := net.Listen(cfg.Network, host+":"+port)
	if err != nil {
		return nil, err
	}

	limiter := rate.NewLimiter(rate.Limit(cfg.RequestPerSecondLimit), 1)

	//Регистрация интерсепторов на сервере
	newServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.RateLimiter(logger, limiter),
			interceptors.UnaryPanicRecoveryInterceptor(logger),
			interceptors.XRequestID,
			interceptors.LoggerInterseptor(logger),
			redisadapter.RedisCacheInterceptor(redis, 10*time.Minute),
		),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	return &Server{newServer, lis}, nil
}
