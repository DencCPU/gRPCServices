package usecase

import (
	"context"

	inboxdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/inbox"
	orderdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/order"
	"go.opentelemetry.io/otel/trace"

	"go.uber.org/zap"
)

type Storage interface {
	AddOrderStorage(ctx context.Context, newOrder orderdomain.Order) (orderID string, orderStatus string, err error) //Добавление нового заказа в хранилище
	GetOrderState(ctx context.Context, key orderdomain.Key) (orderInfo orderdomain.ReceivedOrderInfo, err error)
	IdempotencyCheck(idepotencyKey string) bool

	UpdateMarketCache(markets []orderdomain.Market)
	RemoveIrrelevantMarkets()
}

type MarketsService interface {
	GetEnableMarkets(ctx context.Context, userID string, userRole orderdomain.UserRole) ([]orderdomain.Market, error)
}

type Notify interface {
	AddNewState(string, string, chan string)
	GetStatus(orderdomain.Key) string
	AddNewSub(orderdomain.Key) chan string
	GetNumbersSubsChan(orderdomain.Key) int
	UpdateStatusSubs(context.Context, orderdomain.Key)
}

type Kafka interface {
	ReadMessage(ctx context.Context) (inboxdomain.KafkaMessage, error)
	Commit(ctx context.Context, msg inboxdomain.KafkaMessage) error
}

type Inbox interface {
	Add(msg inboxdomain.KafkaMessage) []orderdomain.Market
	CheckEvent(eventId string) bool
}

type OrderService struct {
	storage     Storage
	spotService MarketsService
	notify      Notify
	kafka       Kafka
	inbox       Inbox
	logger      *zap.Logger
	tracer      trace.Tracer
}

func NewOrderServ(store Storage, markets_service MarketsService, notify Notify, kafka Kafka, inbox Inbox, logger *zap.Logger, tracer trace.Tracer) *OrderService {
	return &OrderService{
		store,
		markets_service,
		notify,
		kafka,
		inbox,
		logger,
		tracer,
	}
}
