package userhandlers

import (
	"context"

	"github.com/DencCPU/gRPCServices/Protobuf/gen/user_service"
	tokensdto "github.com/DencCPU/gRPCServices/UserService/internal/adapters/dto/tokens"
	"github.com/DencCPU/gRPCServices/UserService/internal/adapters/dto/user"
	domainuser "github.com/DencCPU/gRPCServices/UserService/internal/domain/user"
)

type Service interface {
	CreateUser(context.Context, domainuser.User) (tokensdto.PairToken, error)
	UpdateTokens(context.Context, tokensdto.InputTokens) (tokensdto.PairToken, error)
	Validation(ctx context.Context, accessToken string) (user.Output, error)
}

type Handler struct {
	user_service.UnimplementedUserServiceServer
	Service
}

func NewHandler(service Service) *Handler {
	return &Handler{Service: service}
}
