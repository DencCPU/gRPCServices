package memory

import (
	"context"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// Управление работой рынков
func (s *Storage) AccessControl(ctx context.Context, timeout time.Duration) string {

	//Добавление названия рынков в слайс
	var markets = make([]int, 0, len(s.date))
	for key := range s.date {
		markets = append(markets, key)
	}

	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return "Market management has been completed"

		case <-ticker.C:

			d := rand.Intn(3)
			switch d {

			case 0: //Market lock

				n := rand.Intn(len(markets))

				s.mu.Lock()
				key := markets[n]
				if s.date[key].Enable != false {
					s.date[key].Enable = false
					s.mu.Unlock()
					continue
				}
				s.mu.Unlock()

			case 1: //Delete market
				n := rand.Intn(len(markets))
				s.mu.Lock()
				key := markets[n]
				if s.date[key].Enable != false {
					s.date[key].Enable = false
					delete_at := time.Now().Local()
					s.date[key].DeleteAt = &delete_at

					s.mu.Unlock()
					continue
				}
				s.mu.Unlock()

			case 2:

				n := rand.Intn(len(markets))

				s.mu.Lock()
				key := markets[n]
				if s.date[key].Enable == false {
					s.date[key].Enable = true
					s.date[key].DeleteAt = nil

					s.mu.Unlock()
					continue
				}
				s.mu.Unlock()

			}
		}
	}
}
