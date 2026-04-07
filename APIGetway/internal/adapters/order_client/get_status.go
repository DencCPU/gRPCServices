package orderclient

import (
	"context"

	orderdto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/order_service"
	"github.com/DencCPU/gRPCServices/Protobuf/gen/order_service"
)

func (c *Client) GetStatus(ctx context.Context, input orderdto.GetInput) (orderdto.GetOutput, error) {
	req := order_service.GetOrderReq{
		OrderId: input.Order_id,
		UserId:  input.User_id,
	}
	resp, err := c.GetOrderStatus(ctx, &req)
	if err != nil {
		return orderdto.GetOutput{}, err
	}
	output := orderdto.GetOutput{
		Order_id:     resp.OrderId,
		Order_status: resp.OrderStatus,
	}
	return output, nil
}
