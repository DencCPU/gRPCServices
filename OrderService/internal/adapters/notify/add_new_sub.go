package notify

import orderdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/order"

// Add a new sub
func (s *StatusStorage) AddNewSub(key orderdomain.Key) chan string {
	ch := make(chan string, 10)
	s.mu.Lock()
	s.Subs[key] = append(s.Subs[key], ch)
	s.mu.Unlock()

	return ch
}

func (s *StatusStorage) Unsubscribe(key orderdomain.Key, ch chan string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	subs := s.Subs[key]
	for i, c := range subs {
		if c == ch {
			s.Subs[key] = append(subs[:i], subs[i+1:]...)
			break
		}
	}
}
