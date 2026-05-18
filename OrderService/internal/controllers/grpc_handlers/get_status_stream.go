package orderhandlers

import (
	"io"
	"time"

	orderdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/order"
	order "github.com/DencCPU/gRPCServices/Protobuf/gen/order_service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handlers) StreamOrderUpdate(req *order.StreamOrderUpdateReq, stream order.OrderService_StreamOrderUpdateServer) error {
	//Validation request
	if err := req.Validate(); err != nil {
		return err
	}

	key := orderdomain.Key{
		OrderId: req.OrderId,
		UserId:  req.UserId,
	}
	stateCh, err := h.Service.StreamGetState(stream.Context(), key)
	if err != nil {
		return err
	}

	for {
		select {
		case status, ok := <-stateCh:
			if !ok {
				return io.EOF
			}

			update_time := timestamppb.New(time.Now())
			stream.Send(&order.StreamOrderUpdateResp{OrderStatus: status, UpdateStatusTime: update_time})

		case <-stream.Context().Done():
			return nil
		}
	}
}
