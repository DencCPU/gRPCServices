package orderdomain

import "github.com/shopspring/decimal"

type UserRole int32

const (
	USER_ROLE_UNSPECIFIED  UserRole = 0
	USER_ROLE_BASIC_USER   UserRole = 1
	USER_ROLE_PREMIUM_USER UserRole = 2
)

type OrderType int64

const (
	ORDER_TYPE_UNSPECIFIED OrderType = 0
	ORDER_TYPE_NORMAL      OrderType = 1
	ORDER_TYPE_EXPRESS     OrderType = 2
)

const (
	StatusCreated    = "created"
	StatusProcessing = "processing"
	StatusComplited  = "complited"
	StatusCancel     = "cancel"
)

type Order struct {
	UserId         string
	MarketId       string
	OrderType      OrderType
	Price          decimal.Decimal
	Quantity       int64
	Status         string
	UserRole       UserRole
	IdempotencyKey string
}

type Key struct {
	UserId  string
	OrderId string
}

type Market struct {
	ID   string
	Name string
}

type OrderInfo struct {
	OrderType OrderType
	UserId    string
	OrderId   string
}

type ReceivedOrderInfo struct {
	Status     string
	OrderId    string
	Price      decimal.Decimal
	Quantity   int64
	MarketName string
}
