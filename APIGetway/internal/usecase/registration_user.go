package usecase

import (
	"context"

	"github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/tokens"
	userdomain "github.com/DencCPU/gRPCServices/APIGetway/internal/domain/user"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

func (s *Service) RegistrationUser(ctx context.Context, newUser userdomain.User) (tokens.PairToken, error) {
	ctx, span := s.tracer.Start(ctx, "Registration user")
	defer span.End()

	span.SetAttributes(
		attribute.String("name", newUser.Name),
		attribute.String("email", newUser.Email),
		attribute.String("password", newUser.Password),
		attribute.String("role:", newUser.Role),
	)

	pairToken, err := s.userClient.RegistrationUser(ctx, newUser)
	if err != nil {
		s.logger.Error("error registering a new user",
			zap.String("spanID:", span.SpanContext().SpanID().String()),
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "registration user failed")
		return tokens.PairToken{}, err
	}

	s.logger.Info("User is registred",
		zap.String("spanID:", span.SpanContext().SpanID().String()),
	)

	span.SetStatus(codes.Ok, "registration user successfuly")

	return pairToken, nil
}
