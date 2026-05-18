package usecase

import (
	"context"

	userservicedto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/user_service"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

func (s *Service) Validation(ctx context.Context, accessToken string) (userservicedto.Output, error) {
	ctx, span := s.tracer.Start(ctx, "Validation access token")
	defer span.End()

	span.SetAttributes(
		attribute.String("access token:", accessToken),
	)

	output, err := s.userClient.Validation(ctx, accessToken)
	if err != nil {
		s.logger.Error("error validation access token",
			zap.String("spanID:", span.SpanContext().SpanID().String()),
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "validation failed")
		return userservicedto.Output{}, err
	}

	s.logger.Info("Validation completed successfully",
		zap.String("spanID:", span.SpanContext().SpanID().String()),
	)

	span.SetStatus(codes.Ok, "validation successfuly")

	return output, nil
}
