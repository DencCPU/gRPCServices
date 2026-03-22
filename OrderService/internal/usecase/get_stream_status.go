package usecase

import (
	"Academy/gRPCServices/OrderService/internal/domain/order"
	"context"
	"fmt"
)

// Получение статуса заказа в стриминге
func (o *OrderService) StreamGetState(ctx context.Context, key order.Key) chan string {

	//Добавление новой подписки для получения статусов
	stateCh := o.Notify.AddNewSub(key)
	fmt.Println("Подписан новый клиент")

	//Получение кол-ва каналов
	quantiryCh := o.GetNumbersSubsChan(key)
	if quantiryCh == 1 {
		o.UpdateStatusSubs(key)
		fmt.Println("Запущена служба")
	}

	return stateCh
}
