package usecase

import (
	"context"

	orderdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/order"
	"go.uber.org/zap"
)

func (o *OrderService) GetStatus(ctx context.Context, key orderdomain.Key) (string, error) {

	ctx, span := o.tracer.Start(ctx, "Get status")
	defer span.End()
	status, err := o.GetOrderState(ctx, key)
	if err != nil {
		o.logger.Error("error receiving order status:",
			zap.Error(err),
		)
		return "", err
	}
	return status, nil
}
