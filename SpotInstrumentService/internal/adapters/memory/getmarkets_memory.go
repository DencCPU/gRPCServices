package memory

import (
	"fmt"

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

	fmt.Println("USERROLE:", input.UserRole)

	if input.PageToken == "" {

		switch input.UserRole {
		case domainusers.USER_ROLE_PREMIUM_USER:
			fmt.Println("PREMIUM ROLE")
			for i = 0; i < size; i++ {
				s.mu.RLock()
				if s.date[i].DeleteAt == nil || s.date[i].Enable == true {
					fmt.Println("USERACCESS:", s.date[i].UserAccess)
					enableMarkets = append(enableMarkets, s.date[i])
				}
				s.mu.RUnlock()
			}
		case domainusers.USER_ROLE_BASIC_USER:
			fmt.Println("BASIC ROLE")
			for i = 0; i < size; i++ {
				s.mu.RLock()
				if (s.date[i].DeleteAt == nil || s.date[i].Enable == true) && input.UserRole == s.date[i].UserAccess {
					fmt.Println("USERACCESS:", s.date[i].UserAccess)
					enableMarkets = append(enableMarkets, s.date[i])
				}
				s.mu.RUnlock()
			}
		}

	} else {

		for s.date[i].ID != input.PageToken {
			i++
		}
		i++

		switch input.UserRole {
		case domainusers.USER_ROLE_PREMIUM_USER:
			for input.PageSize != 0 && i < len(s.date) {
				s.mu.RLock()
				if s.date[i].DeleteAt == nil || s.date[i].Enable == true {
					fmt.Println("USERACCESS:", s.date[i].UserAccess)
					enableMarkets = append(enableMarkets, s.date[i])
				}
				input.PageSize--
				s.mu.RUnlock()
			}
		case domainusers.USER_ROLE_BASIC_USER:
			for input.PageSize != 0 && i < len(s.date) {
				s.mu.RLock()
				fmt.Println("USERACCESS:", s.date[i].UserAccess)
				if (s.date[i].DeleteAt == nil || s.date[i].Enable == true) && s.date[i].UserAccess == input.UserRole {
					enableMarkets = append(enableMarkets, s.date[i])
				}
				input.PageSize--
				s.mu.RUnlock()
			}
		}

	}

	pageToken = s.date[i-1].ID
	return enableMarkets, pageToken
}
