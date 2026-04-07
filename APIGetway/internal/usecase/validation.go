package usecase

import (
	"context"

	userservicedto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/user_service"
)

func (s *Service) Validation(ctx context.Context, accessToken string) (userservicedto.Output, error) {
	output, err := s.user_client.Validation(ctx, accessToken)
	if err != nil {
		return userservicedto.Output{}, err
	}
	return output, nil
}
