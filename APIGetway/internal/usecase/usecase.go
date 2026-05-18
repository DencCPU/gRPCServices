package usecase

import (
	"context"

	orderdto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/order_service"
	spotservicedto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/spot_service"
	"github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/tokens"
	userservicedto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/user_service"
	orderdomain "github.com/DencCPU/gRPCServices/APIGetway/internal/domain/order"
	userdomain "github.com/DencCPU/gRPCServices/APIGetway/internal/domain/user"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type UserClient interface {
	RegistrationUser(ctx context.Context, newUser userdomain.User) (tokens.PairToken, error)
	Validation(ctx context.Context, accessToken string) (userservicedto.Output, error)
	UpdateAccessToken(ctx context.Context, accessToken, refreshToken string) (tokens.PairToken, error)
	AuthenticationUser(ctx context.Context, email, password string) (tokens.PairToken, error)
}

type OrderClient interface {
	CreateNewOrder(ctx context.Context, order orderdomain.OrderInfo) (orderdto.Output, error)
	GetStatus(ctx context.Context, input orderdto.GetInput) (orderdto.GetOutput, error)
	GetStreamStatus(ctx context.Context, input orderdto.GetInput, msgChan chan orderdto.StreamOutput) error
}

type SpotClient interface {
	ViewEnableMarkets(ctx context.Context, input spotservicedto.Input) ([]spotservicedto.Output, error)
}
type Service struct {
	userClient  UserClient
	orderClient OrderClient
	spotClient  SpotClient
	logger      *zap.Logger
	tracer      trace.Tracer
}

func NewService(userClient UserClient, orderClient OrderClient, spotClient SpotClient, logger *zap.Logger, tracer trace.Tracer) *Service {
	return &Service{userClient, orderClient, spotClient, logger, tracer}
}
