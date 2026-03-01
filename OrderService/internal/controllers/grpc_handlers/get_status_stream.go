package orderhandlers

import (
	"Academy/gRPCServices/OrderService/internal/domain/order"
	orderAPI "Academy/gRPCServices/Protobuf/gen/order"
)

func (h *Handlers) StreamOrderUpdate(req *orderAPI.GetReq, stream orderAPI.OrderService_StreamOrderUpdateServer) error {
	key := order.Key{
		Order_id: req.OrderId,
		User_id:  req.UserId,
	}
	stateCh, unsubscribe := h.Service.StreamGetState(stream.Context(), key)
	defer unsubscribe()

	for {
		select {
		case status, ok := <-stateCh:
			if !ok {
				return nil
			}
			stream.Send(&orderAPI.GetResp{OrderStatus: status})
		case <-stream.Context().Done():
			return nil
		}
	}
}
