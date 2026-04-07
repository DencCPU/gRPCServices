package memory

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	domainmarket "github.com/DencCPU/gRPCServices/SpotInstrumentService/internal/domain/market"
	"github.com/google/uuid"
)

// Добавление маркетов
func (s *Storage) AddMarkets(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file:%w", err)
	}
	re := regexp.MustCompile(`([a-zA-Z0-9]+)+\s?([a-zA-Z0-9]+)?`)
	markets := re.FindAllString(string(file), -1)

	//Заполнение storage
	if len(markets) == 0 {
		return errors.New("market list is empty")
	}

	for _, m := range markets {
		id := uuid.New().String()
		s.date[m] = &domainmarket.Market{ID: id, Name: m, Enable: true, Delete_at: nil}
	}

	return nil
}
