package usecase

import (
	"context"
	"fmt"

	orderdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/order"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

func (o *OrderService) CreateOrder(ctx context.Context, newOrder orderdomain.Order) (string, string, error) {
	//Idempotency check
	ctx, span := o.tracer.Start(ctx, "Create order")
	defer span.End()

	span.SetAttributes(
		attribute.String("UserID", newOrder.UserId),
		attribute.String("marketID", newOrder.MarketId),
		attribute.Int64("orderType", int64(newOrder.OrderType)),
		attribute.Float64("price,", newOrder.Price.InexactFloat64()),
		attribute.Int64("quantity", newOrder.Quantity),
		attribute.Int64("userRole", int64(newOrder.UserRole)),
		attribute.String("idempotrncyKey", newOrder.IdempotencyKey),
	)
	span.AddEvent("IdempotencyCheck")

	idempotency := o.storage.IdempotencyCheck(newOrder.IdempotencyKey)
	if !idempotency {
		o.logger.Warn("Attempt to create an order using an existing idempotency key",
			zap.String("Key:", newOrder.IdempotencyKey),
		)
		span.SetStatus(codes.Error, "trying to re-create the order")
		return "", "", fmt.Errorf("trying to re-create the order")
	}

	o.logger.Info("The idempotency key has been verified",
		zap.String("Key:", newOrder.IdempotencyKey),
	)

	//Get marketss list
	span.AddEvent("get enable markets")
	ctx, span = o.tracer.Start(ctx, "enable markets")
	defer span.End()

	markets, err := o.spotService.GetEnableMarkets(ctx, newOrder.UserId, newOrder.UserRole)
	if err != nil {
		o.logger.Error("error getting available markets:",
			zap.String("spanID:", span.SpanContext().SpanID().String()),
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "get enable markets failed")
		return "", "", err
	}

	o.logger.Info("Received a list of available markets",
		zap.String("spanID:", span.SpanContext().SpanID().String()),
	)

	span.AddEvent("add order to storage")
	ctx, span = o.tracer.Start(ctx, "add order")
	defer span.End()

	orderID, status, err := o.storage.AddOrderStorage(ctx, newOrder, markets)
	if err != nil {
		o.logger.Error("Error add new order to storage:",
			zap.String("spanID:", span.SpanContext().SpanID().String()),
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "add order to storage failed")
		return "", "", err
	}
	o.logger.Info("Order created:",
		zap.String("OrderID:", orderID),
		zap.String("spanID:", span.SpanContext().SpanID().String()),
	)
	span.SetStatus(codes.Ok, "create order successfuly")

	//Order fulfillment
	o.logger.Info("The order has been processed")

	return orderID, status, nil
}
