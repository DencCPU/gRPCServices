package notify

import (
	"context"
	"sync"
	"time"

	orderdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/order"
)

func (s *StatusStorage) UpdateStatusSubs(ctx context.Context, key orderdomain.Key) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		var laststatus string
		ticker := time.NewTicker(s.TikerInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return

			case <-ticker.C:
				status := s.GetStatus(key)
				if laststatus != status {

					for _, ch := range s.Subs[key] {

						select {
						case ch <- status:
						default:
						}
					}

					if status == orderdomain.StatusComplited {
						return
					}

					laststatus = status
				}
			}
		}

	}()

	go func() {
		wg.Wait()
		for _, ch := range s.Subs[key] {
			close(ch)
		}
	}()
}
