package orderserver

import (
	serverconfig "Academy/gRPCServices/OrderService/config/server"
	"Academy/gRPCServices/SpotInstrumentService/pkg/interseptors"
	"fmt"
	"net"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type Server struct {
	*grpc.Server
	Listener net.Listener
}

func New() (*Server, error) {

	cfg, err := serverconfig.NewConfig()
	if err != nil {
		return nil, err
	}

	host := cfg.Server.Host
	port := cfg.Server.Port
	dsn := fmt.Sprintf("%s:%d", host, port)
	lis, err := net.Listen(cfg.Server.Network, dsn)
	if err != nil {
		return nil, err
	}

	newServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interseptors.UnaryPanicRecoveryInterceptor,
			interseptors.XRequestID,
			interseptors.LoggerInterseptor,
		), grpc.StatsHandler(otelgrpc.NewServerHandler()))

	return &Server{newServer, lis}, nil
}
