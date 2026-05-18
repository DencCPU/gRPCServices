package memory

import (
	domainmarket "github.com/DencCPU/gRPCServices/SpotInstrumentService/internal/domain/market"
	domainusers "github.com/DencCPU/gRPCServices/SpotInstrumentService/internal/domain/users"
)

// Получение доступных рынков
func (s *Storage) GetEnableMarkets(input domainusers.Input) ([]*domainmarket.Market, string) {
	var size int

	if len(s.date) < input.PageSize || input.PageSize == 0 {
		size = len(s.date)
	} else {
		size = input.PageSize
	}

	var (
		enableMarkets []*domainmarket.Market
		pageToken     string
		i             int
	)

	if input.PageToken == "" {
		for i = 0; i < size; i++ {
			s.mu.RLock()
			key := s.keys[i]
			if (s.date[key].DeleteAt == nil || s.date[key].Enable == true) && s.date[key].UserAccess == input.UserRole {
				enableMarkets = append(enableMarkets, s.date[key])
			}
			s.mu.RUnlock()
		}
	} else {
		for s.keys[i] != input.PageToken {
			i++
		}
		i++
		for input.PageSize != 0 && i < len(s.keys) {
			s.mu.RLock()
			key := s.keys[i]
			if (s.date[key].DeleteAt == nil || s.date[key].Enable == true) && s.date[key].UserAccess == input.UserRole {
				enableMarkets = append(enableMarkets, s.date[key])
			}
			input.PageSize--
			s.mu.RUnlock()
		}
	}

	pageToken = s.keys[i-1]
	return enableMarkets, pageToken
}
