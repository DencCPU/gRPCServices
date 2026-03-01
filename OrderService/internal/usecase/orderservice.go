package usecase

import (
	"Academy/gRPCServices/OrderService/internal/domain/order"
	"context"
)

type Storage interface {
	AddOrderStorage(context.Context, order.Order, []int64) (string, string, error) //Добавление нового заказа в хранилище
	GetOrderState(context.Context, order.Key) (string, error)                      //Получение статуса заказа
	ControlOrder(orderType string, user_id int64, orderID string) chan string
}

type MarketsService interface {
	GetEnableMarkets(context.Context) ([]int64, error) //Получение списка доступных рынков
}

type Notify interface {
	AddNewState(user_id int64, orderId string, statCh chan string)
	GetStatus(key order.Key) string
	AddNewSub(key order.Key) (chan string, func())
	GetNumbersSubsChan(key order.Key) int
	UpdateStatusSubs(key order.Key)
}

type OrderService struct {
	Storage
	MarketsService
	Notify
}

func NewOrderServ(in_memory Storage, markets_service MarketsService, notify Notify) *OrderService {
	return &OrderService{in_memory, markets_service, notify}
}
