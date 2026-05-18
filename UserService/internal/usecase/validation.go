package usecase

import (
	"context"
	"errors"

	"github.com/DencCPU/gRPCServices/UserService/internal/adapters/dto/user"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

func (s *Service) Validation(ctx context.Context, accessToken string) (user.Output, error) {
	ctx, span := s.tracer.Start(ctx, "validation token")
	defer span.End()
	span.SetAttributes(
		attribute.String("access token:", accessToken),
	)
	output, err := s.jwt.Validation(accessToken)
	if err != nil {
		s.logger.Error("token validation error:",
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "validation failed")
		return user.Output{}, err
	}

	if output.Role == "" || output.UserId == "" {
		s.logger.Error("data from token was not received")
		span.RecordError(errors.New("data from token was not received"))
		span.SetStatus(codes.Error, "data from token was not received")

		return user.Output{}, errors.New("data from token was not received")
	}

	s.logger.Info("token validation successful",
		zap.String("spanID:", span.SpanContext().SpanID().String()),
	)
	span.SetStatus(codes.Ok, "validation token successfuly")
	return output, nil
}
