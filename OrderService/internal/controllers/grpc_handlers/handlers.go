package orderhandlers

import (
	"context"

	orderdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/order"
	"github.com/DencCPU/gRPCServices/OrderService/internal/usecase"
	order "github.com/DencCPU/gRPCServices/Protobuf/gen/order_service"
)

type Service interface {
	CreateOrder(ctx context.Context, newOrder orderdomain.Order) (orderID string, orderStatus string, err error)
	GetStatus(ctx context.Context, key orderdomain.Key) (orderInfo orderdomain.ReceivedOrderInfo, err error)
	StreamGetState(ctx context.Context, key orderdomain.Key) (stateChan chan string, err error)
}

type Handlers struct {
	order.UnimplementedOrderServiceServer
	Service Service
}

func NewHandlers(orderService *usecase.OrderService) *Handlers {
	return &Handlers{Service: orderService}
}
