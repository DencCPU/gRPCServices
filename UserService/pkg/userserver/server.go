package userserver

import (
	"net"

	"github.com/DencCPU/gRPCServices/Shared/interceptors"
	userconfig "github.com/DencCPU/gRPCServices/UserService/config"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

type Server struct {
	*grpc.Server
	Listener net.Listener
}

func NewServer(cfg userconfig.Server, logger *zap.Logger) (*Server, error) {

	host := cfg.Host
	port := cfg.Port
	lis, err := net.Listen(cfg.Network, host+":"+port)
	if err != nil {
		return nil, err
	}
	limiter := rate.NewLimiter(rate.Limit(cfg.RequestPerSecondLimit), 1)
	//Interceptors
	newServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.RateLimiter(logger, limiter),
			interceptors.UnaryPanicRecoveryInterceptor(logger),
			interceptors.XRequestID,
			interceptors.LoggerInterseptor(logger),
		),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	return &Server{newServer, lis}, nil
}
