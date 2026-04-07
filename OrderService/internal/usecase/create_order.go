package usecase

import (
	"context"
	"fmt"

	orderdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/order"
	"go.uber.org/zap"
)

func (o *OrderService) CreateOrder(ctx context.Context, newOrder orderdomain.Order) (string, string, error) {
	//Получения списка доступных рынков
	ctx, span := o.tracer.Start(ctx, "Get enable markets")
	defer span.End()

	markets, err := o.GetEnableMarkets(ctx)
	if err != nil {
		o.logger.Error("Error getting available markets:",
			zap.Error(err),
		)
		return "", "", err
	}
	o.logger.Info("Received a list of available markets",
		zap.String("get enable markets span ID:", span.SpanContext().TraceID().String()),
	)

	//Создание нового заказа
	ctx, span = o.tracer.Start(ctx, "AddOrder")
	defer span.End()
	orderID, status, err := o.AddOrderStorage(ctx, newOrder, markets)
	if err != nil {
		return "", "", err
	}
	o.logger.Info("Order created:",
		zap.String("OrderID:", orderID),
		zap.String("add order span ID:", span.SpanContext().TraceID().String()),
	)

	//Выполнение заказа
	fmt.Println("The order has been processed")
	stateCh := o.ControlOrder(newOrder.Order_type, newOrder.User_id, orderID)

	go o.Notify.AddNewState(newOrder.User_id, orderID, stateCh)

	return orderID, status, nil
}
