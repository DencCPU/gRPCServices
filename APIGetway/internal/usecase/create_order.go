package usecase

import (
	"context"

	orderdto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/order_service"
	"github.com/DencCPU/gRPCServices/APIGetway/internal/domain/order"
)

func (s *Service) CreateOrder(ctx context.Context, order order.OrderInfo) (orderdto.Output, error) {
	output, err := s.order_client.CreateNewOrder(ctx, order)
	if err != nil {
		return orderdto.Output{}, err
	}
	return output, nil
}
