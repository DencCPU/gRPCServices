package notify

import orderdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/order"

// Получение актуального статуса заказа
func (s *StatusStorage) GetStatus(key orderdomain.Key) string {
	s.mu.RLock()
	status := s.Status[key]
	s.mu.RUnlock()
	if status == "" {
		return "created"
	}
	return status
}
