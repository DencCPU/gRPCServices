package inbox

import (
	"sync"

	inboxdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/inbox"
	orderdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/order"
)

type Inbox struct {
	cacheEvent map[string]inboxdomain.MarketEvent
	mu         sync.RWMutex
}

func NewInbox() *Inbox {
	return &Inbox{cacheEvent: make(map[string]inboxdomain.MarketEvent)}
}

func (i *Inbox) Add(msg inboxdomain.KafkaMessage) []orderdomain.Market {
	key := msg.Value.BasedEvent.EventId
	i.mu.Lock()
	i.cacheEvent[key] = msg.Value
	eventMarkets := i.cacheEvent[key].Markets
	i.mu.Unlock()

	markets := make([]orderdomain.Market, 0, len(eventMarkets))

	for _, market := range eventMarkets {
		m := orderdomain.Market{
			MarketId:   market.MarketId,
			MarketName: market.MarketName,
			// UserAccess: market.UserAccess,
		}
		markets = append(markets, m)
	}
	return markets
}

func (i *Inbox) CheckEvent(eventId string) bool {
	i.mu.RLock()
	if _, exist := i.cacheEvent[eventId]; exist {
		return true
	}
	i.mu.RUnlock()
	return false
}
