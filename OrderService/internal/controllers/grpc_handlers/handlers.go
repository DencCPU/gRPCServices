package orderhandlers

import (
	"Academy/gRPCServices/OrderService/internal/domain/order"
	"Academy/gRPCServices/OrderService/internal/usecase"
	orderAPI "Academy/gRPCServices/Protobuf/gen/order"
	"context"
)

type Service interface {
	CreateOrder(context.Context, order.Order) (string, string, error)
	GetStatus(context.Context, order.Key) (string, error)
	StreamGetState(ctx context.Context, key order.Key) (chan string, func())
}

type Handlers struct {
	orderAPI.UnimplementedOrderServiceServer
	Service Service
}

func NewHandlers(orderService *usecase.OrderService) *Handlers {
	return &Handlers{Service: orderService}
}
