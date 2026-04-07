package orderclient

import (
	"context"
	"errors"

	orderdto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/order_service"
	"github.com/DencCPU/gRPCServices/APIGetway/internal/domain/order"
	"github.com/DencCPU/gRPCServices/Protobuf/gen/order_service"
)

func (c *Client) CreateNewOrder(ctx context.Context, order order.OrderInfo) (orderdto.Output, error) {
	req := order_service.CreateOrderReq{}
	switch order.Order_type {
	case "normal":
		req = order_service.CreateOrderReq{
			UserId:    order.User_id,
			MarketId:  order.Market_id,
			OrderType: order_service.OrderType_ORDER_TYPE_NORMAL,
			Price:     order.Price,
			Quantity:  order.Quantity,
		}
	case "express":
		req = order_service.CreateOrderReq{
			UserId:    order.User_id,
			MarketId:  order.Market_id,
			OrderType: order_service.OrderType_ORDER_TYPE_EXPRESS,
			Price:     order.Price,
			Quantity:  order.Quantity,
		}
	default:
		return orderdto.Output{}, errors.New("incorrect oreder type")
	}
	resp, err := c.CreateOrder(ctx, &req)
	if err != nil {
		return orderdto.Output{}, err
	}
	output := orderdto.NewOutput(resp.OrderId, resp.OrderStatus)
	return output, nil
}
