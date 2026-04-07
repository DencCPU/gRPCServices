package usecase

import (
	"context"

	orderdto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/order_service"
	spotservicedto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/spot_service"
	"github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/tokens"
	userservicedto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/user_service"
	"github.com/DencCPU/gRPCServices/APIGetway/internal/domain/order"
	userdomain "github.com/DencCPU/gRPCServices/APIGetway/internal/domain/user"
	"go.uber.org/zap"
)

type UserClient interface {
	RegistrationUser(ctx context.Context, newUser userdomain.User) (tokens.PairToken, error)
	Validation(ctx context.Context, accessToken string) (userservicedto.Output, error)
	UpdateAccessToken(ctx context.Context, accessToken, refreshToken string) (tokens.PairToken, error)
}

type OrderClient interface {
	CreateNewOrder(ctx context.Context, order order.OrderInfo) (orderdto.Output, error)
	GetStatus(ctx context.Context, input orderdto.GetInput) (orderdto.GetOutput, error)
	GetStreamStatus(ctx context.Context, input orderdto.GetInput, msgChan chan orderdto.StreamOutput) error
}

type SpotClient interface {
	ViewEnableMarkets(ctx context.Context, role string) ([]spotservicedto.Output, error)
}
type Service struct {
	user_client  UserClient
	order_client OrderClient
	spot_client  SpotClient
	logger       *zap.Logger
}

func NewService(userClient UserClient, orderClient OrderClient, spotClient SpotClient, logger *zap.Logger) *Service {
	return &Service{userClient, orderClient, spotClient, logger}
}
