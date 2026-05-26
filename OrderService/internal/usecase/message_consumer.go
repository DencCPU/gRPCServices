package usecase

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

func (o *OrderService) MessageConsumer(ctx context.Context, wg *sync.WaitGroup, errChan chan error) {
	o.logger.Info("starting markets message consumer")

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				o.logger.Info("markets consumer stopped")
				return

			default:
				msg, err := o.kafka.ReadMessage(ctx)
				if err != nil {
					if err == context.Canceled {
						return
					}
					o.logger.Error("failed to read message", zap.Error(err))
					errChan <- fmt.Errorf("read message error: %w", err)
					continue
				}

				if o.inbox.CheckEvent(msg.Value.BasedEvent.EventId) {
					o.logger.Debug("duplicate message skipped",
						zap.String("event_id", msg.Value.BasedEvent.EventId))

					err := o.kafka.Commit(ctx, msg)
					if err != nil {
						o.logger.Error("failed to commit offset", zap.Error(err))
					}
					continue
				}

				newMarkets := o.inbox.Add(msg)
				o.logger.Info("processing markets update",
					zap.String("event_id", msg.Value.BasedEvent.EventId),
					zap.Int("new_markets_count", len(newMarkets)))

				o.storage.UpdateMarketCache(newMarkets)
				o.storage.RemoveIrrelevantMarkets()

				err = o.kafka.Commit(ctx, msg)
				if err != nil {
					o.logger.Error("failed to commit offset", zap.Error(err))
					errChan <- fmt.Errorf("commit error: %w", err)
					continue
				}

				o.logger.Info("markets update processed successfully",
					zap.String("event_id", msg.Value.BasedEvent.EventId))
			}
		}
	}()
}
