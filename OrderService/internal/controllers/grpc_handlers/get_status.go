package orderhandlers

import (
	"context"
	"fmt"

	orderdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/order"
	"github.com/DencCPU/gRPCServices/Protobuf/gen/money"
	order "github.com/DencCPU/gRPCServices/Protobuf/gen/order_service"
	"github.com/shopspring/decimal"
)

func (h *Handlers) GetOrderStatus(ctx context.Context, req *order.GetOrderReq) (*order.GetOrderResp, error) {
	//Валидация запроса
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("uncorrect request format:%w", err)
	}
	//Добавить валидацию запроса
	key := orderdomain.Key{
		OrderId: req.OrderId,
		UserId:  req.UserId,
	}
	orderInfo, err := h.Service.GetStatus(ctx, key)
	if err != nil {
		return &order.GetOrderResp{}, err
	}
	uints := orderInfo.Price.IntPart()
	nanos := orderInfo.Price.Sub(decimal.NewFromInt(uints)).Shift(9).IntPart()
	resp := order.GetOrderResp{
		OrderStatus: orderInfo.Status,
		OrderId:     orderInfo.OrderId,
		Price: &money.Money{
			CurrencyCode: "RUS",
			Units:        uints,
			Nanos:        int32(nanos),
		},
		Quantity:   orderInfo.Quantity,
		MarketName: orderInfo.MarketName,
	}
	return &resp, nil
}
