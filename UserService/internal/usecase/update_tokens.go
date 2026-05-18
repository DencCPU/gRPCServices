package usecase

import (
	"context"

	tokensdto "github.com/DencCPU/gRPCServices/UserService/internal/adapters/dto/tokens"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

func (s *Service) UpdateTokens(ctx context.Context, inputTokens tokensdto.InputTokens) (tokensdto.PairToken, error) {

	//Update refresh token
	ctx, span := s.tracer.Start(ctx, "Update refresh token:")
	defer span.End()
	span.SetAttributes(
		attribute.String("access token:", inputTokens.AccsesToken),
		attribute.String("refresh token:", inputTokens.RefreshToken),
	)
	refreshToken, err := s.storage.UpdateRefreshToken(ctx, inputTokens.RefreshToken)
	if err != nil {
		s.logger.Error("refresh token update error:",
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "update refresh token failed")
		return tokensdto.PairToken{}, err
	}

	s.logger.Info("refresh token update",
		zap.String("spanID:", span.SpanContext().SpanID().String()),
	)

	//Update accses token
	ctx, span = s.tracer.Start(ctx, "Update access token:")
	defer span.End()
	span.AddEvent("update accses token")
	accessToken, ttl, err := s.jwt.UpdateAccessToken(inputTokens.AccsesToken)
	if err != nil {
		s.logger.Error("access token update error:",
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "update access token failed")
		return tokensdto.PairToken{}, err
	}

	s.logger.Info("access token update",
		zap.String("spanID:", span.SpanContext().SpanID().String()),
	)

	span.SetStatus(codes.Ok, "update tokens successfuly")

	pairToken := tokensdto.NewPairToken(accessToken, refreshToken, ttl)
	return pairToken, nil
}
