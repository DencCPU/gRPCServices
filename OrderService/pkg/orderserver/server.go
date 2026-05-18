package orderserver

import (
	"fmt"
	"net"

	orderconfig "github.com/DencCPU/gRPCServices/OrderService/config"
	"github.com/DencCPU/gRPCServices/Shared/interceptors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

type Server struct {
	*grpc.Server
	Listener net.Listener
}

func New(cfg orderconfig.Server, logger *zap.Logger) (*Server, error) {

	host := cfg.Host
	port := cfg.Port
	dsn := fmt.Sprintf("%s:%d", host, port)
	lis, err := net.Listen(cfg.Network, dsn)
	if err != nil {
		return nil, err
	}

	limiter := rate.NewLimiter(rate.Limit(cfg.RequestPerSecondLimit), 1)

	newServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.RateLimiter(logger, limiter),
			interceptors.UnaryPanicRecoveryInterceptor(logger),
			interceptors.XRequestID,
			interceptors.LoggerInterseptor(logger),
		), grpc.StatsHandler(otelgrpc.NewServerHandler()))

	return &Server{newServer, lis}, nil
}
