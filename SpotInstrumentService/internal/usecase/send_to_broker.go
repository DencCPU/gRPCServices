package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

func (s *SpotService) SendToBroker(ctx context.Context, wg *sync.WaitGroup, getMarketsInterval time.Duration, relayInterval time.Duration) {

	wg.Add(1)
	go func() {
		s.logger.Info("Kafka producer started")
		defer wg.Done()

		ticker := time.NewTicker(getMarketsInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				markets := s.storage.GetAllMarkets()
				if len(markets) > 0 {
					s.outbox.Add(markets)
				}
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(relayInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				pendingEvents := s.outbox.GetPending()
				if len(pendingEvents) == 0 {
					continue
				}
				for _, el := range pendingEvents {
					markets := el.Markets
					fmt.Println(markets)
				}

				for _, event := range pendingEvents {

					err := s.kafka.Send(ctx, event)

					if err != nil {
						s.logger.Error("failed to send event to Kafka",
							zap.String("event_id", event.BasedEvent.EventId),
							zap.Error(err))
						continue
					}

					s.logger.Info("the massege sent to kafka:",
						zap.String("eventID:", event.BasedEvent.EventId),
					)

					err = s.outbox.Remove(event.BasedEvent.EventId)
					if err != nil {
						s.logger.Error("failed to remove event from outbox",
							zap.String("event_id", event.BasedEvent.EventId),
							zap.Error(err))
						continue
					}
					s.logger.Info("event sent and removed from outbox",
						zap.String("event_id", event.BasedEvent.EventId))
				}
			}
		}
	}()
}
