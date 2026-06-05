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
	)

	if input.PageToken == "" {

		switch input.UserRole {
		case domainusers.USER_ROLE_PREMIUM_USER:
			for i := 0; i < len(s.date) && size > 0; i++ {
				s.mu.RLock()
				if s.date[i].DeleteAt == nil || s.date[i].Enable == true {
					enableMarkets = append(enableMarkets, s.date[i])
					size--
				}
				s.mu.RUnlock()
			}

		case domainusers.USER_ROLE_BASIC_USER:
			for i := 0; i < len(s.date) && size > 0; i++ {
				s.mu.RLock()
				if (s.date[i].DeleteAt == nil || s.date[i].Enable == true) && input.UserRole == s.date[i].UserAccess {
					enableMarkets = append(enableMarkets, s.date[i])
					size--
				}
				s.mu.RUnlock()
			}
		}

	} else {

		startIdx := -1
		for i := 0; i < len(s.date); i++ {
			s.mu.RLock()
			if s.date[i].ID == input.PageToken {
				startIdx = i + 1
				s.mu.RUnlock()
				break
			}
			s.mu.RUnlock()
		}

		if startIdx == -1 {
			return nil, ""
		}

		remaining := size

		switch input.UserRole {
		case domainusers.USER_ROLE_PREMIUM_USER:
			for i := startIdx; remaining > 0 && i < len(s.date); i++ {
				s.mu.RLock()
				if s.date[i].DeleteAt == nil || s.date[i].Enable == true {
					enableMarkets = append(enableMarkets, s.date[i])
					remaining--
				}
				s.mu.RUnlock()
			}

		case domainusers.USER_ROLE_BASIC_USER:
			for i := startIdx; remaining > 0 && i < len(s.date); i++ {
				s.mu.RLock()
				if (s.date[i].DeleteAt == nil || s.date[i].Enable == true) && s.date[i].UserAccess == input.UserRole {
					enableMarkets = append(enableMarkets, s.date[i])
					remaining--
				}
				s.mu.RUnlock()
			}
		}
	}
	if len(enableMarkets) > 0 {
		lastMarket := enableMarkets[len(enableMarkets)-1]
		pageToken = lastMarket.ID
	}

	return enableMarkets, pageToken
}
