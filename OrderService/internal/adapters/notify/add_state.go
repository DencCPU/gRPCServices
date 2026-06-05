package notify

import (
	"fmt"

	orderdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/order"
)

// Add a new status
func (s *StatusStorage) AddNewState(userId string, orderId string, statCh chan string) {
	key := orderdomain.Key{UserId: userId, OrderId: orderId}
	for state := range statCh {
		s.mu.Lock()
		fmt.Println("NOTIFY ORDER STATUS:", state)
		s.Status[key] = state
		s.mu.Unlock()
	}
}
