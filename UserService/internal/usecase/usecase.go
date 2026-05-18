package usecase

import (
	"context"
	"time"

	postgresdto "github.com/DencCPU/gRPCServices/UserService/internal/adapters/dto/postgres"
	"github.com/DencCPU/gRPCServices/UserService/internal/adapters/dto/user"
	domainuser "github.com/DencCPU/gRPCServices/UserService/internal/domain/user"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Storage interface {
	AddUser(ctx context.Context, user domainuser.User) (userID string, refreshToken string, err error)
	Authentication(ctx context.Context, email, password string) (postgresdto.AuthUser, error)
	UpdatePassword(ctx context.Context, email string, password string) (err error)
	UpdateRefreshToken(ctx context.Context, token string) (string, error)
	UpdateExpireAt(ctx context.Context, userId string) (string, error)
}
type JWT interface {
	CreateAccessToken(userID, email, role string) (accessToken string, expireAt time.Time, err error)
	UpdateAccessToken(token string) (newToken string, expireAt time.Time, err error)
	Validation(accessToken string) (user user.Output, err error)
}

type Service struct {
	storage Storage
	jwt     JWT
	logger  *zap.Logger
	tracer  trace.Tracer
}

func NewService(storage Storage, logger *zap.Logger, jwt JWT, tracer trace.Tracer) *Service {
	return &Service{storage, jwt, logger, tracer}
}
