package usecase

import (
	"context"

	orderdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/order"
)

// Получение статуса заказа в стриминге
func (o *OrderService) StreamGetState(ctx context.Context, key orderdomain.Key) chan string {

	//Добавление новой подписки для получения статусов
	stateCh := o.Notify.AddNewSub(key)
	o.logger.Info("New client signed")

	//Получение кол-ва каналов
	quantiryCh := o.GetNumbersSubsChan(key)
	if quantiryCh == 1 {
		o.UpdateStatusSubs(key)
		o.logger.Info("Service started")
	}

	return stateCh
}
