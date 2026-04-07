package usecase

import (
	"context"

	orderdto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/order_service"
	"go.uber.org/zap"
)

func (s *Service) GetStreamStatus(ctx context.Context, input orderdto.GetInput, msgChan chan orderdto.StreamOutput) error {
	err := s.order_client.GetStreamStatus(ctx, input, msgChan)
	if err != nil {
		s.logger.Error("streaming status error:",
			zap.Error(err),
		)
		return err
	}
	return nil
}
