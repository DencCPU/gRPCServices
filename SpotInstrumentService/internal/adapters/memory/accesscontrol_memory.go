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
func (s *Storage) AccessControl(ctx context.Context) string {

	//Добавление названия рынков в слайс
	var markets = make([]string, 0, len(s.date))
	for key := range s.date {
		markets = append(markets, key)
	}

	for {
		select {
		case <-ctx.Done():
			// fmt.Println("Время жизни цикла истекло")
			return "Market management has been completed"
		default:
			d := rand.Intn(3) //Случайная блокировка или удаление рынка
			switch d {

			case 0: //Блокировка доступа случайного маркета

				n := rand.Intn(len(markets))

				s.mu.Lock()
				key := markets[n]
				if s.date[key].Enable != false { //Проверка доступа к рынку
					s.date[key].Enable = false
					s.mu.Unlock()
					break
				}
				s.mu.Unlock()

			case 1: //Удаление случайного маркета с рынка
				n := rand.Intn(len(markets))

				s.mu.Lock()
				key := markets[n]
				if s.date[key].Enable != false { //Проверка доступа к рынку
					s.date[key].Enable = false
					delete_at := time.Now().Local()
					s.date[key].Delete_at = &delete_at
					s.mu.Unlock()
					break
				}
				s.mu.Unlock()

			case 2: //Востановление доступа к маркету на рынке

				n := rand.Intn(len(markets))

				s.mu.Lock()
				key := markets[n]
				if s.date[key].Enable == false {
					s.date[key].Enable = true
					s.date[key].Delete_at = nil
					s.mu.Unlock()
					break
				}
				s.mu.Unlock()

			}
		}
		time.Sleep(1 * time.Second)
	}
}
