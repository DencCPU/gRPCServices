package orderclient

import (
	"context"
	"fmt"
	"io"

	orderdto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/order_service"
	"github.com/DencCPU/gRPCServices/Protobuf/gen/order_service"
)

func (c *Client) GetStreamStatus(ctx context.Context, input orderdto.GetInput, msgChan chan orderdto.StreamOutput) error {
	req := order_service.StreamOrderUpdateReq{
		OrderId: input.Order_id,
		UserId:  input.User_id,
	}

	stream, err := c.StreamOrderUpdate(ctx, &req)
	if err != nil {
		return fmt.Errorf("failed to create stream: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			resp, err := stream.Recv()
			if err == io.EOF {

				return nil
			}

			msg := orderdto.StreamOutput{
				Order_status: resp.OrderStatus,
				Update_time:  resp.UpdateStatusTime.AsTime(),
			}

			select {
			case msgChan <- msg:

			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}
