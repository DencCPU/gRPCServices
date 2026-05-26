package memory

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"time"

	domainmarket "github.com/DencCPU/gRPCServices/SpotInstrumentService/internal/domain/market"
	domainusers "github.com/DencCPU/gRPCServices/SpotInstrumentService/internal/domain/users"
	"github.com/google/uuid"
)

// Добавление маркетов
func (s *Storage) AddMarkets(path string) error {

	//Read file with markets name
	file, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file:%w", err)
	}

	//Create regular expression
	re := regexp.MustCompile(`([a-zA-Z0-9]+)+\s?([a-zA-Z0-9]+)?`)
	markets := re.FindAllString(string(file), -1)

	if len(markets) == 0 {
		return errors.New("market list is empty")
	}

	rand.Seed(time.Now().Unix())
	for i, name := range markets {
		id := uuid.New().String()

		var userAccess domainusers.UserRole
		switch rand.Intn(2) {
		case 0:
			userAccess = domainusers.USER_ROLE_BASIC_USER
		default:
			userAccess = domainusers.USER_ROLE_PREMIUM_USER
		}

		s.date[i] = &domainmarket.Market{ID: id, Name: name, Enable: true, DeleteAt: nil, UserAccess: userAccess}
	}

	return nil
}
