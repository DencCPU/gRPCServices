package usecase

import (
	"context"

	"github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/tokens"
)

func (s *Service) UpdateTokens(ctx context.Context, accessToken, refreshToken string) (tokens.PairToken, error) {
	pairToken, err := s.user_client.UpdateAccessToken(ctx, accessToken, refreshToken)
	if err != nil {
		return tokens.PairToken{}, err
	}
	return pairToken, nil
}
