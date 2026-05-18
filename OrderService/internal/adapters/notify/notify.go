package notify

import (
	"sync"
	"time"

	orderdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/order"
)

type StatusStorage struct {
	Status        map[orderdomain.Key]string
	Subs          map[orderdomain.Key][]chan string
	TikerInterval time.Duration
	mu            sync.RWMutex
}

func NewStatStorage(interval time.Duration) *StatusStorage {
	return &StatusStorage{
		Status:        make(map[orderdomain.Key]string),
		Subs:          make(map[orderdomain.Key][]chan string),
		TikerInterval: interval,
	}

}

func (s *StatusStorage) GetNumbersSubsChan(key orderdomain.Key) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	l := len(s.Subs[key])
	return l
}
