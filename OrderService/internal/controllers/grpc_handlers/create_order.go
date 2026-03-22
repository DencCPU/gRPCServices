package orderhandlers

import (
	"Academy/gRPCServices/OrderService/internal/domain/order"
	orderAPI "Academy/gRPCServices/Protobuf/gen/order"
	"context"
	"fmt"

	"github.com/shopspring/decimal"
)

func (h *Handlers) CreateOrder(ctx context.Context, req *orderAPI.CreateReq) (*orderAPI.CreateResp, error) {

	price, err := decimal.NewFromString(req.Price)
	if err != nil {
		return nil, fmt.Errorf("Неправильно указанная суммма:%w", err)
	}
	newOrder := order.Order{
		User_id:    req.UserId,
		Market_id:  req.MarketId,
		Order_type: req.OrderType,
		Price:      price,
		Quantity:   req.Quantity,
	}

	orderID, status, err := h.Service.CreateOrder(ctx, newOrder)
	if err != nil {
		return &orderAPI.CreateResp{}, err
	}
	resp := orderAPI.CreateResp{OrderId: orderID, Status: status}
	return &resp, nil
}
