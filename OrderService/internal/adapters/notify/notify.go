package notify

import (
	"Academy/gRPCServices/OrderService/internal/domain/order"
	"sync"
)

type StatusStorage struct {
	Status map[order.Key]string
	Subs   map[order.Key][]chan string
	mu     sync.Mutex
}

func NewStatStorage() *StatusStorage {
	return &StatusStorage{Status: make(map[order.Key]string), Subs: make(map[order.Key][]chan string)}
}

func (s *StatusStorage) GetNumbersSubsChan(key order.Key) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	l := len(s.Subs[key])
	return l
}
