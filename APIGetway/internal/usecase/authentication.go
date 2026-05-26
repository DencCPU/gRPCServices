package usecase

import (
	"context"

	"github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/tokens"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (s *Service) Authentication(ctx context.Context, email, password string) (tokens.PairToken, error) {
	ctx, span := s.tracer.Start(ctx, "Autentication")
	defer span.End()

	span.SetAttributes(
		attribute.String("email", email),
	)

	pairToken, err := s.userClient.AuthenticationUser(ctx, email, password)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "authentication failed")
		return tokens.PairToken{}, err
	}

	span.SetStatus(codes.Ok, "authentication succeccfuly")
	return pairToken, nil
}
