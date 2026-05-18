package usecase

import (
	"context"

	orderdto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/order_service"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

func (s *Service) GetOrderStatus(ctx context.Context, input orderdto.GetInput) (orderdto.GetOutput, error) {
	ctx, span := s.tracer.Start(ctx, "Get order status")
	defer span.End()
	span.SetAttributes(
		attribute.String("userID", input.OrderId),
		attribute.String("OrderID", input.OrderId),
	)

	output, err := s.orderClient.GetStatus(ctx, input)
	if err != nil {
		s.logger.Error("error getting order status:",
			zap.String("spanID", span.SpanContext().SpanID().String()),
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "get order status failed")
		return orderdto.GetOutput{}, err
	}

	s.logger.Info("order status received",
		zap.String("spanID", span.SpanContext().SpanID().String()),
	)

	span.SetStatus(codes.Ok, "get status successfuly")

	return output, nil
}
