package outbox

import (
	"errors"
	"sync"
	"time"

	domainmarket "github.com/DencCPU/gRPCServices/SpotInstrumentService/internal/domain/market"
	outboxdomain "github.com/DencCPU/gRPCServices/SpotInstrumentService/internal/domain/outbox"
	"github.com/google/uuid"
)

type Outbox struct {
	events map[string]*outboxdomain.MarketEvent
	mu     sync.RWMutex
}

func NewOutbox() *Outbox {
	return &Outbox{events: make(map[string]*outboxdomain.MarketEvent)}
}

func (o *Outbox) Add(markets []*domainmarket.Market) string {
	eventId := uuid.New().String()

	marketSlice := make([]outboxdomain.Market, 0, len(markets))

	for _, market := range markets {
		m := outboxdomain.Market{
			MarketId:   market.ID,
			MarketName: market.Name,
			UserAccess: market.UserAccess,
		}
		marketSlice = append(marketSlice, m)
	}

	event := &outboxdomain.MarketEvent{
		BasedEvent: outboxdomain.BasedEvent{
			EventId:    eventId,
			EventType:  "View enable markets",
			OccurredAt: time.Now(),
		},
		Markets: marketSlice,
	}
	o.mu.Lock()
	o.events[eventId] = event
	o.mu.Unlock()

	return eventId
}

func (o *Outbox) GetPending() []*outboxdomain.MarketEvent {
	o.mu.RLock()
	defer o.mu.RUnlock()

	result := make([]*outboxdomain.MarketEvent, 0, len(o.events))
	for _, event := range o.events {
		result = append(result, event)
	}
	return result
}

func (o *Outbox) Remove(eventId string) error {
	o.mu.Lock()
	if _, exist := o.events[eventId]; !exist {
		return errors.New("event not found")
	}
	delete(o.events, eventId)
	o.mu.Unlock()
	return nil
}

func (o *Outbox) Get(eventId string) (*outboxdomain.MarketEvent, error) {
	o.mu.RLock()
	if _, exist := o.events[eventId]; !exist {
		return nil, errors.New("event not found")
	}
	event := o.events[eventId]
	o.mu.RUnlock()

	return event, nil
}
