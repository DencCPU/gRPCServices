package orderhandlers

import (
	"time"

	orderdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/order"
	order "github.com/DencCPU/gRPCServices/Protobuf/gen/order_service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handlers) StreamOrderUpdate(req *order.StreamOrderUpdateReq, stream order.OrderService_StreamOrderUpdateServer) error {
	//Валидация запроса
	if err := req.Validate(); err != nil {
		return nil
	}
	key := orderdomain.Key{
		Order_id: req.OrderId,
		User_id:  req.UserId,
	}
	stateCh := h.Service.StreamGetState(stream.Context(), key)

	for {
		select {
		case status, ok := <-stateCh:
			if !ok {
				return nil
			}
			update_time := timestamppb.New(time.Now())
			stream.Send(&order.StreamOrderUpdateResp{OrderStatus: status, UpdateStatusTime: update_time})
		case <-stream.Context().Done():
			return nil
		}
	}
}
