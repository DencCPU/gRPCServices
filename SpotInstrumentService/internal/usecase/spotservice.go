package usecase

import (
	"context"

	domainmarket "github.com/DencCPU/gRPCServices/SpotInstrumentService/internal/domain/market"
	outboxdomain "github.com/DencCPU/gRPCServices/SpotInstrumentService/internal/domain/outbox"
	domainusers "github.com/DencCPU/gRPCServices/SpotInstrumentService/internal/domain/users"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Интерфейс для рынков
type StorageRepo interface {
	GetEnableMarkets(input domainusers.Input) ([]*domainmarket.Market, string) //Получение доступных рынков
	GetAllMarkets() []*domainmarket.Market
}

type Kafka interface {
	Send(ctx context.Context, event *outboxdomain.MarketEvent) error
}

type OutBox interface {
	Add(markets []*domainmarket.Market) string
	GetPending() []*outboxdomain.MarketEvent
	Remove(eventId string) error
}

type SpotService struct {
	storage StorageRepo
	kafka   Kafka
	outbox  OutBox
	logger  *zap.Logger
	tracer  trace.Tracer
}

// Конструктор для SpotInstrument
func NewSpotInstrument(repo StorageRepo, kafka Kafka, outbox OutBox, logger *zap.Logger, tracer trace.Tracer) *SpotService {
	return &SpotService{
		repo,
		kafka,
		outbox,
		logger,
		tracer,
	}
}
