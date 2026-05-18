package breaker

import (
	"time"

	servererror "github.com/DencCPU/gRPCServices/Shared/validation/server_error"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

type Params struct {
	Name           string
	MaxRequest     uint32
	Interval       time.Duration
	Timeout        time.Duration
	MaxFailRequest uint32
}

func NewBreaker(logger *zap.Logger, params Params) *gobreaker.CircuitBreaker {
	setting := gobreaker.Settings{
		Name:        params.Name,       //Breaker name
		MaxRequests: params.MaxRequest, //Maximum number of requests allowed in half-open mode
		Interval:    params.Interval,   //The period for resetting the statistics for counting failed requests in private mode. Time in seconds.
		Timeout:     params.Timeout,    //The time the breaker remains in the open state before transitioning to Half-open.
		ReadyToTrip: func(counts gobreaker.Counts) bool { //A function that determines the condition for transitioning from Close to Open
			return counts.ConsecutiveFailures > params.MaxFailRequest
		},
		OnStateChange: func(name string, from, to gobreaker.State) { //Function for working with logging
			if to == gobreaker.StateOpen {
				logger.Info("Breaker went into a state Open")
			}
			if to == gobreaker.StateHalfOpen {
				logger.Info("Breaker went into a state Half-open")
			}
			if to == gobreaker.StateClosed {
				logger.Info("Breaker went into a state Close")
			}
		},
		IsSuccessful: func(err error) bool { //A function that determines which errors are taken into account
			return servererror.ServerError(err)
		},
	}
	breaker := gobreaker.NewCircuitBreaker(setting)
	return breaker
}
