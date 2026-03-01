package usecase

import (
	"Academy/gRPCServices/OrderService/internal/domain/order"
	"context"
	"fmt"
)

func (o *OrderService) CreateOrder(ctx context.Context, newOrder order.Order) (string, string, error) {

	//Получения списка доступных рынков
	marketsID, err := o.GetEnableMarkets(ctx)
	if err != nil {
		return "", "", err
	}
	fmt.Println("Получили рынки")

	//Создание нового заказа
	orderID, status, err := o.AddOrderStorage(ctx, newOrder, marketsID)
	if err != nil {
		return "", "", err
	}
	fmt.Println("Заказ создан")

	//Выполнение заказа
	stateCh := o.ControlOrder(newOrder.Order_type, newOrder.User_id, orderID)

	go o.Notify.AddNewState(newOrder.User_id, orderID, stateCh)

	return orderID, status, nil
}
