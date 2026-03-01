package notify

import "Academy/gRPCServices/OrderService/internal/domain/order"

// Получение актуального статуса заказа
func (s *StatusStorage) GetStatus(key order.Key) string {
	s.mu.Lock()
	status := s.Status[key]
	s.mu.Unlock()
	return status
}
