package usecase

import (
	"context"

	orderdto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/order_service"
	"go.uber.org/zap"
)

func (s *Service) GetOrderStatus(ctx context.Context, input orderdto.GetInput) (orderdto.GetOutput, error) {
	output, err := s.order_client.GetStatus(ctx, input)
	if err != nil {
		s.logger.Error("error getting order status:",
			zap.Error(err),
		)
		return orderdto.GetOutput{}, err
	}
	s.logger.Info("order status received")
	return output, nil
}
