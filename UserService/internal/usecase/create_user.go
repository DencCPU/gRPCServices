package usecase

import (
	"context"

	tokensdto "github.com/DencCPU/gRPCServices/UserService/internal/adapters/dto/tokens"
	domainuser "github.com/DencCPU/gRPCServices/UserService/internal/domain/user"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

func (s *Service) CreateUser(ctx context.Context, user domainuser.User) (tokensdto.PairToken, error) {

	//Add user to database and generate refresh token
	ctx, span := s.tracer.Start(ctx, "Add user:")
	defer span.End()
	span.SetAttributes(
		attribute.String("name", user.Name),
		attribute.String("email", user.Email),
		attribute.String("password", user.Password),
		attribute.String("role:", user.Role),
	)

	user_id, refreshToken, err := s.storage.AddUser(ctx, user)
	if err != nil {
		s.logger.Error("error adding user to database:",
			zap.String("spanID:", span.SpanContext().SpanID().String()),
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "add user failed")
		return tokensdto.PairToken{}, err
	}

	s.logger.Info("adding user succefully:",
		zap.String("spanID:", span.SpanContext().SpanID().String()),
	)

	//Create accses token
	span.AddEvent("create accsess token")
	accsesToken, ttl, err := s.jwt.CreateAccessToken(user_id, user.Email, user.Role)
	if err != nil {
		s.logger.Error("jwt generation error:",
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "create access token failed")
		return tokensdto.PairToken{}, err
	}

	pairToken := tokensdto.NewPairToken(accsesToken, refreshToken, ttl)
	s.logger.Info("token creation succeful:",
		zap.String("spanID:", span.SpanContext().SpanID().String()),
	)
	span.SetStatus(codes.Ok, "create user successfuly")

	return pairToken, nil
}
