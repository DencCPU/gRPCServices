package orderclient

import (
	"context"
	"errors"
	"fmt"
	"io"

	orderdto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/order_service"
	"github.com/DencCPU/gRPCServices/Protobuf/gen/order_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Client) GetStreamStatus(ctx context.Context, input orderdto.GetInput, msgChan chan orderdto.StreamOutput) error {
	req := &order_service.StreamOrderUpdateReq{
		OrderId: input.OrderId,
		UserId:  input.UserId,
	}

	stream, err := c.StreamOrderUpdate(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create stream: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			resp, err := stream.Recv()

			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
				if status.Code(err) == codes.Canceled {
					return nil
				}
				return fmt.Errorf("stream recv error: %w", err)
			}

			msg := orderdto.StreamOutput{
				OrderStatus: resp.OrderStatus,
				UpdateTime:  resp.UpdateStatusTime.AsTime(),
			}

			select {
			case msgChan <- msg:
			case <-ctx.Done():
				return nil
			}
		}
	}
}
