package usecase

import (
	"context"

	orderdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/order"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

func (o *OrderService) GetStatus(ctx context.Context, key orderdomain.Key) (orderdomain.ReceivedOrderInfo, error) {

	ctx, span := o.tracer.Start(ctx, "get order status")
	defer span.End()
	span.SetAttributes(
		attribute.String("userID", key.OrderId),
		attribute.String("OrderID", key.OrderId),
	)

	orderInfo, err := o.storage.GetOrderState(ctx, key)
	if err != nil {
		o.logger.Error("error receiving order status:",
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "get status failed")
		return orderdomain.ReceivedOrderInfo{}, err
	}

	o.logger.Info("Order status received:",
		zap.String("UserID:", key.UserId),
		zap.String("OrderID", key.OrderId),
	)

	span.SetStatus(codes.Ok, "get status successfuly")
	return orderInfo, nil
}
