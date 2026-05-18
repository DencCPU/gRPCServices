package orderdomain

type UserRole int32

const (
	USER_ROLE_UNSPECIFIED  UserRole = 0
	USER_ROLE_BASIC_USER   UserRole = 1
	USER_ROLE_PREMIUM_USER UserRole = 2
)

type OrderInfo struct {
	UserId    string `json:"user_id"`
	MarketId  string `json:"market_id"`
	OrderType string `json:"order_type"`
	Price     string `json:"price"`
	Quantity  int64  `json:"quantity"`
	UserRole  UserRole
}
