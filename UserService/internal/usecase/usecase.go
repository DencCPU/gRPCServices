package usecase

import (
	"context"
	"time"

	"github.com/DencCPU/gRPCServices/UserService/internal/adapters/dto/user"
	domainuser "github.com/DencCPU/gRPCServices/UserService/internal/domain/user"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Storage interface {
	AddUser(context.Context, domainuser.User) (string, string, error)
	UpdatePassword(context.Context, string, string) error
	UpdateRefreshToken(context.Context, string) (string, error)
}
type JWT interface {
	CreateAccessToken(string, string, string) (string, time.Time, error)
	UpdateAccessToken(string) (string, time.Time, error)
	Validation(accessToken string) (user.Output, error)
}

type Service struct {
	Storage
	JWT
	logger *zap.Logger
	tracer trace.Tracer
}

func NewService(storage Storage, logger *zap.Logger, jwt JWT, tracer trace.Tracer) *Service {
	return &Service{storage, jwt, logger, tracer}
}
