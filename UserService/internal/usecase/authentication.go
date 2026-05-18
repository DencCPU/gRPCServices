package usecase

import (
	"context"

	tokensdto "github.com/DencCPU/gRPCServices/UserService/internal/adapters/dto/tokens"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

func (s *Service) AuthenticationUser(ctx context.Context, email, password string) (tokensdto.PairToken, error) {
	ctx, span := s.tracer.Start(ctx, "authentication user")
	defer span.End()

	span.SetAttributes(
		attribute.String("email", email),
		attribute.String("password", password),
	)

	authUser, err := s.storage.Authentication(ctx, email, password)
	if err != nil {
		s.logger.Error("user authentication error:",
			zap.String("spanID:", span.SpanContext().SpanID().String()),
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "authentication failed")

		return tokensdto.PairToken{}, err
	}
	s.logger.Info("user confirmed")

	span.AddEvent("start update expireAt")
	ctx, span = s.tracer.Start(ctx, "update expireAt in refresh token")
	defer span.End()
	span.SetAttributes(
		attribute.String("userID", authUser.ID),
	)

	refreshToken, err := s.storage.UpdateExpireAt(ctx, authUser.ID)
	if err != nil {
		s.logger.Error("error update expireAt in refresh token:",
			zap.String("spanID:", span.SpanContext().SpanID().String()),
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "update expireAt failed")
		return tokensdto.PairToken{}, err
	}
	s.logger.Info("update expireAt in refresh token successfully",
		zap.String("spanID:", span.SpanContext().SpanID().String()),
	)

	span.AddEvent("create access token")

	accessToken, expireAt, err := s.jwt.CreateAccessToken(authUser.ID, email, authUser.Role)
	if err != nil {
		s.logger.Error("access token creation error",
			zap.String("spanID:", span.SpanContext().SpanID().String()),
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "create access token failed")
		return tokensdto.PairToken{}, err
	}
	s.logger.Info("authentication successful",
		zap.String("spanID:", span.SpanContext().SpanID().String()),
	)
	span.SetStatus(codes.Ok, "authentication successfully")

	pairToken := tokensdto.NewPairToken(accessToken, refreshToken, expireAt)
	return pairToken, nil
}
