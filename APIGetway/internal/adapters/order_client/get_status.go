package orderclient

import (
	"context"
	"errors"
	"fmt"

	orderdto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/order_service"
	"github.com/DencCPU/gRPCServices/Protobuf/gen/order_service"
	"github.com/shopspring/decimal"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Client) GetStatus(ctx context.Context, input orderdto.GetInput) (orderdto.GetOutput, error) {
	result, err := c.breaker.Execute(func() (interface{}, error) {
		req := order_service.GetOrderReq{
			OrderId: input.OrderId,
			UserId:  input.UserId,
		}
		resp, err := c.GetOrderStatus(ctx, &req)
		if err != nil {
			return orderdto.GetOutput{}, err
		}

		if resp.OrderId == "" || resp.OrderStatus == "" {
			return orderdto.GetOutput{}, errors.New("incorrect response from the server")
		}
		return resp, err
	})

	if err != nil {
		if errors.Is(err, gobreaker.ErrOpenState) {
			return orderdto.GetOutput{}, status.Errorf(codes.Unavailable, "service is temporarily unavailable")
		}
		return orderdto.GetOutput{}, err
	}

	resp, ok := result.(*order_service.GetOrderResp)
	if !ok {
		return orderdto.GetOutput{}, fmt.Errorf("Inappropriate result type:%T", result)
	}
	uints := decimal.NewFromInt(resp.Price.Units)
	nanos := decimal.NewFromInt32(resp.Price.Nanos).Shift(-9)

	p := uints.Add(nanos)
	price := p.String()

	output := orderdto.GetOutput{
		Status:     resp.OrderStatus,
		OrderId:    resp.OrderId,
		Price:      price,
		Quantity:   resp.Quantity,
		MarketName: resp.MarketName,
	}
	return output, nil
}
