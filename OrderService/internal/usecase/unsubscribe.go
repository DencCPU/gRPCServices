package usecase

import orderdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/order"

func (o *OrderService) Unsubscribe(key orderdomain.Key, ch chan string) {
	o.notify.Unsubscribe(key, ch)
	o.logger.Info("Unsubscribe user")
}
