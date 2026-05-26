package outboxdomain

import (
	"time"

	domainusers "github.com/DencCPU/gRPCServices/SpotInstrumentService/internal/domain/users"
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
	UserAccess domainusers.UserRole `json:"user_access"`
}
