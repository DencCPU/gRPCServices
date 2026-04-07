package usecase

import (
	"context"

	spotservicedto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/spot_service"
)

func (s *Service) ViewEnableMarkets(ctx context.Context, role string) ([]spotservicedto.Output, error) {
	markets, err := s.spot_client.ViewEnableMarkets(ctx, role)
	if err != nil {
		return nil, err
	}
	return markets, nil
}
