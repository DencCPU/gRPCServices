package spotservice

import (
	"time"

	orderconfig "github.com/DencCPU/gRPCServices/OrderService/config"
	spot "github.com/DencCPU/gRPCServices/Protobuf/gen/spot_service"
	servererror "github.com/DencCPU/gRPCServices/Shared/validation/server_error"
	"github.com/DencCPU/gRPCServices/SpotInstrumentService/pkg/spotclient"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

type Client struct {
	spot.SpotInstrumentServiceClient
	breaker *gobreaker.CircuitBreaker
}

func NewClient(cfg orderconfig.BreakerSetting, logger *zap.Logger) (*Client, error) {
	client, err := spotclient.NewClient()
	if err != nil {
		return nil, err
	}
	setting := gobreaker.Settings{
		Name:        cfg.Name,                                  //Название брейкера
		MaxRequests: cfg.MaxRequests,                           //Максимально кол-во запросов, пропускаемых в полуоткрытом режиме(Half-open)
		Interval:    time.Duration(cfg.Interval) * time.Second, //Период сброса статистики подсчета неудачных запросов в закрытом режиме(Close). Время в секундах.
		Timeout:     time.Duration(cfg.Timeout) * time.Second,  //Время прибывания брейкера в открытом состоянии, перед переходов в Half-open.
		ReadyToTrip: func(counts gobreaker.Counts) bool { //Функция, определяющая условие перехода из Close в Open
			return counts.ConsecutiveFailures > cfg.MaxFailRequest
		},
		OnStateChange: func(name string, from, to gobreaker.State) { //Функция для работы с логированием
			if to == gobreaker.StateOpen {
				logger.Info("Брейкер перешел в состояние Open")
			}
			if to == gobreaker.StateHalfOpen {
				logger.Info("Брейкер перешел в состояние Half-open")
			}
			if to == gobreaker.StateClosed {
				logger.Info("Брейкер перешел в состояние Close")
			}
		},
		IsSuccessful: func(err error) bool { //Функция, определяющая, какие ошибки учитываются
			return servererror.ServerError(err)
		},
	}
	breaker := gobreaker.NewCircuitBreaker(setting)
	return &Client{client, breaker}, nil
}
