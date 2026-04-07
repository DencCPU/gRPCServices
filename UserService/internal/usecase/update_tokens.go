package usecase

import (
	"context"

	tokensdto "github.com/DencCPU/gRPCServices/UserService/internal/adapters/dto/tokens"
	"go.uber.org/zap"
)

func (s *Service) UpdateTokens(ctx context.Context, inputTokens tokensdto.InputTokens) (tokensdto.PairToken, error) {

	//Update refresh token
	ctx, span := s.tracer.Start(ctx, "Update refresh token:")
	defer span.End()
	refreshToken, err := s.Storage.UpdateRefreshToken(ctx, inputTokens.RefreshToken)
	if err != nil {
		s.logger.Error("refresh token update error:",
			zap.Error(err),
		)
		return tokensdto.PairToken{}, err
	}
	s.logger.Info("refresh token update")

	//Update accses token
	ctx, span = s.tracer.Start(ctx, "Update accses token:")
	defer span.End()
	accessToken, ttl, err := s.JWT.UpdateAccessToken(inputTokens.AccsesToken)
	if err != nil {
		s.logger.Error("accsess token update error:",
			zap.Error(err),
		)
		return tokensdto.PairToken{}, err
	}
	s.logger.Info("access token update")

	pairToken := tokensdto.NewPairToken(accessToken, refreshToken, ttl)
	return pairToken, nil
}
