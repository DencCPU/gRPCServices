package memory

import domainmarket "github.com/DencCPU/gRPCServices/SpotInstrumentService/internal/domain/market"

func (s *Storage) GetAllMarkets() []*domainmarket.Market {
	markets := make([]*domainmarket.Market, 0, len(s.date))

	for _, m := range s.date {
		s.mu.RLock()
		markets = append(markets, m)
		s.mu.RUnlock()
	}
	return markets
}
