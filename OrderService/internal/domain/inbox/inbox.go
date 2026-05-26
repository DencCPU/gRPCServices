package inboxdomain

import (
	"time"

	orderdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/order"
)

type MarketEvent struct {
	BasedEvent BasedEvent `json:"based_event"`
	Markets    []Market   `json:"markets"`
}

type BasedEvent struct {
	EventId    string    `json:"event_id"`
	EventType  string    `json:"event_type"`
	OccurredAt time.Time `json:"occured_at"`
}

type Market struct {
	MarketId   string               `json:"market_id"`
	MarketName string               `json:"market_name"`
	UserAccess orderdomain.UserRole `json:"user_access"`
}

type KafkaMessage struct {
	Key       string
	Value     MarketEvent `json:"value"`
	Offset    int64       `json:"offset"`
	Partition int         `json:"partition"`
}
