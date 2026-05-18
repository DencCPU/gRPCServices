package usecase

import (
	"context"

	orderdto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/order_service"
	orderdomain "github.com/DencCPU/gRPCServices/APIGetway/internal/domain/order"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

func (s *Service) CreateOrder(ctx context.Context, order orderdomain.OrderInfo) (orderdto.Output, error) {
	ctx, span := s.tracer.Start(ctx, "Create new order")
	defer span.End()

	span.SetAttributes(
		attribute.String("UserID", order.UserId),
		attribute.String("marketID", order.MarketId),
		attribute.String("orderType", order.OrderType),
		attribute.String("price,", order.Price),
		attribute.Int64("quantity", order.Quantity),
		attribute.Int64("userRole", int64(order.UserRole)),
	)
	output, err := s.orderClient.CreateNewOrder(ctx, order)
	if err != nil {
		s.logger.Error("error creating a new order:",
			zap.String("spanID:", span.SpanContext().SpanID().String()),
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "create new order failed")
		return orderdto.Output{}, err
	}

	s.logger.Info("Order created",
		zap.String("spanID:", span.SpanContext().SpanID().String()),
	)
	span.SetStatus(codes.Ok, "create order successfuly")
	return output, nil
}
