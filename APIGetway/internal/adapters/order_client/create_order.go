package orderclient

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	orderdto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/order_service"
	orderdomain "github.com/DencCPU/gRPCServices/APIGetway/internal/domain/order"
	"github.com/DencCPU/gRPCServices/Protobuf/gen/common"
	"github.com/DencCPU/gRPCServices/Protobuf/gen/money"
	"github.com/DencCPU/gRPCServices/Protobuf/gen/order_service"
	"github.com/shopspring/decimal"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Client) CreateNewOrder(ctx context.Context, order orderdomain.OrderInfo) (orderdto.Output, error) {
	d, err := decimal.NewFromString(order.Price)
	if err != nil {
		return orderdto.Output{}, nil
	}
	data := fmt.Sprintf("%s%s%v%s%d",
		order.UserId,
		order.MarketId,
		order.OrderType,
		order.Price,
		order.Quantity,
	)
	hash := sha256.Sum256([]byte(data))

	indempotencyKey := hex.EncodeToString(hash[:])

	uints := d.IntPart()
	nanos := d.Sub(decimal.NewFromInt(uints)).Shift(9).IntPart()

	result, err := c.breaker.Execute(func() (interface{}, error) {
		req := order_service.CreateOrderReq{
			UserId:   order.UserId,
			MarketId: order.MarketId,
			Price: &money.Money{
				CurrencyCode: "RUB",
				Units:        uints,
				Nanos:        int32(nanos),
			},
			Quantity:        order.Quantity,
			UserRole:        common.UserRole(order.UserRole),
			IndempotencyKey: indempotencyKey,
		}
		switch order.OrderType {
		case "normal":
			req.OrderType = order_service.OrderType_ORDER_TYPE_NORMAL
		case "express":
			req.OrderType = order_service.OrderType_ORDER_TYPE_EXPRESS
		default:
			return orderdto.Output{}, errors.New("incorrect oreder type")
		}

		resp, err := c.CreateOrder(ctx, &req)
		if err != nil {
			return orderdto.Output{}, err
		}

		if resp.OrderId == "" || resp.OrderStatus == "" {
			return orderdto.Output{}, errors.New("incorrect response from the server")
		}
		return resp, nil
	})

	if err != nil {
		if errors.Is(err, gobreaker.ErrOpenState) {
			return orderdto.Output{}, status.Errorf(codes.Unavailable, "service is temporarily unavailable")
		}
		return orderdto.Output{}, err
	}

	resp, ok := result.(*order_service.CreateOrderResp)
	if !ok {
		return orderdto.Output{}, fmt.Errorf("Inappropriate result type:%T", result)
	}

	output := orderdto.NewOutput(resp.OrderId, resp.OrderStatus)
	return output, nil
}
