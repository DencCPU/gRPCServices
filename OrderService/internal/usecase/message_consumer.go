package usecase

import (
	"context"
	"sync"

	"go.uber.org/zap"
)

func (o *OrderService) MessageConsumer(ctx context.Context, wg *sync.WaitGroup) {
	o.logger.Info("starting markets message consumer")

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():

				return

			default:
				msg, err := o.kafka.ReadMessage(ctx)
				if err != nil {
					o.logger.Error("read kafka message error:",
						zap.Error(err),
					)
				}

				if o.inbox.CheckEvent(msg.Value.BasedEvent.EventId) {
					o.logger.Debug("duplicate message skipped",
						zap.String("eventID", msg.Value.BasedEvent.EventId))

					err := o.kafka.Commit(ctx, msg)
					if err != nil {
						o.logger.Error("failed to commit offset", zap.Error(err))
					}
					continue
				}

				newMarkets := o.inbox.Add(msg)
				o.logger.Info("processing markets update",
					zap.String("eventID", msg.Value.BasedEvent.EventId),
					zap.Int("new_markets_count", len(newMarkets)))

				o.storage.UpdateMarketCache(newMarkets)

				o.storage.RemoveIrrelevantMarkets()

				err = o.kafka.Commit(ctx, msg)
				if err != nil {
					o.logger.Error("failed to commit offset", zap.Error(err))
					continue
				}

				o.logger.Info("markets update processed successfully",
					zap.String("eventID", msg.Value.BasedEvent.EventId))

			}
		}

	}()
}
