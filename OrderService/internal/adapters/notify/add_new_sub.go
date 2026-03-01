package notify

import "Academy/gRPCServices/OrderService/internal/domain/order"

// Добавление нового подписчика для получения актуального статуса
func (s *StatusStorage) AddNewSub(key order.Key) (chan string, func()) {
	ch := make(chan string, 1) // буфер 1, чтобы генератор не блокировался
	s.mu.Lock()
	s.Subs[key] = append(s.Subs[key], ch)
	s.mu.Unlock()

	
	unsubscribe := func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		subs := s.Subs[key]
		for i, c := range subs {
			if c == ch {
				s.Subs[key] = append(subs[:i], subs[i+1:]...)
				break
			}
		}
		close(ch)
	}

	return ch, unsubscribe
}
