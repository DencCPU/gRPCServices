package memory

import (
	"sync"

	domainmarket "github.com/DencCPU/gRPCServices/SpotInstrumentService/internal/domain/market"
	"go.uber.org/zap"
)

// In-memory хранилище для хранения рынков
type Storage struct {
	date   map[int]*domainmarket.Market //Хранилище
	mu     sync.RWMutex
	logger *zap.Logger
}

// Создание нового хранинлища
func NewStorage(logger *zap.Logger) (*Storage, error) {
	s := Storage{
		date:   make(map[int]*domainmarket.Market),
		logger: logger,
	}

	return &s, nil
}
