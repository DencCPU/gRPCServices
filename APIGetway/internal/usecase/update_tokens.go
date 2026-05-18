package usecase

import (
	"context"

	"github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/tokens"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

func (s *Service) UpdateTokens(ctx context.Context, accessToken, refreshToken string) (tokens.PairToken, error) {
	ctx, span := s.tracer.Start(ctx, "Update token")
	defer span.End()

	span.SetAttributes(
		attribute.String("access token:", accessToken),
		attribute.String("refresh token:", refreshToken),
	)
	pairToken, err := s.userClient.UpdateAccessToken(ctx, accessToken, refreshToken)
	if err != nil {
		s.logger.Error("tokens refresh error",
			zap.String("spanID", span.SpanContext().SpanID().String()),
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "update tokens failed")
		return tokens.PairToken{}, err
	}

	s.logger.Info("tokens has been update",
		zap.String("spanID:", span.SpanContext().SpanID().String()),
	)

	span.SetStatus(codes.Ok, "update tokens successfuly")
	return pairToken, nil
}
