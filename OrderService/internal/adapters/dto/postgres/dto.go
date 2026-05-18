package postgresdto

import (
	"time"

	orderdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/order"
	"github.com/shopspring/decimal"
)

type OrderDTO struct {
	RefUserId   int
	RefMarketId int
	OrderType   string
	Price       decimal.Decimal
	Quantity    int
	Status      string
	RefOrderId  int
	CreatedAt   time.Time
}

type UsersDTO struct {
	UserId    string
	CreatedAt time.Time
}

type MarketDTO struct {
	MarketId   string
	MarketName string
	CreatedAt  time.Time
}

type Oreders_idDTO struct {
	OrderId   string
	CreatedAt time.Time
}

func CreatOrderDTO(order orderdomain.Order) *OrderDTO {
	input := OrderDTO{
		Price:    order.Price,
		Quantity: int(order.Quantity),
	}
	switch order.OrderType {
	case orderdomain.ORDER_TYPE_NORMAL:
		input.OrderType = "normal"
	case orderdomain.ORDER_TYPE_EXPRESS:
		input.OrderType = "express"
	}

	return &input
}

func CreateUserDTO(userId string) *UsersDTO {
	return &UsersDTO{UserId: userId}
}

func CreateMarketDTO(marketId, marketName string) *MarketDTO {
	return &MarketDTO{MarketId: marketId, MarketName: marketName}
}

func CreateOrders_idDTO(orderId string) *Oreders_idDTO {
	return &Oreders_idDTO{OrderId: orderId}
}
