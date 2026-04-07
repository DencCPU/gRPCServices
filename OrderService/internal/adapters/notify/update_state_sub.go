package notify

import (
	"sync"

	orderdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/order"
)

func (s *StatusStorage) UpdateStatusSubs(key orderdomain.Key) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		var laststatus string

		for {
			status := s.GetStatus(key)
			if laststatus != status {

				// Рассылаем всем подписчикам
				for _, ch := range s.Subs[key] {
					select {
					case ch <- status:
					default:
					}
				}
			} else {
				for _, ch := range s.Subs[key] {
					close(ch)
				}
				return
			}
			laststatus = status
			//Поставить тикер
		}
	}()

	go func() {
		wg.Wait()
	}()
}
