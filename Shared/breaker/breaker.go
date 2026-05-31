package breaker

import (
	"context"
	"time"

	servererror "github.com/DencCPU/gRPCServices/Shared/validation/server_error"
	"github.com/sony/gobreaker"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

type Params struct {
	Name           string
	MaxRequest     uint32
	Interval       time.Duration
	Timeout        time.Duration
	MaxFailRequest uint32
}

func NewBreaker(logger *zap.Logger, params Params, meter metric.Meter) (*gobreaker.CircuitBreaker, error) {

	//Breaker state change metric
	stateChangeCounter, err := meter.Int64Counter("circuit_breaker_state_changes",
		metric.WithDescription("Total number of circuit breaker state changes."), //DescriptionMetric
		metric.WithUnit("{change}"),
	)

	if err != nil {
		return nil, err
	}

	setting := gobreaker.Settings{
		Name:        params.Name,       //Breaker name
		MaxRequests: params.MaxRequest, //Maximum number of requests allowed in half-open mode
		Interval:    params.Interval,   //The period for resetting the statistics for counting failed requests in private mode. Time in seconds.
		Timeout:     params.Timeout,    //The time the breaker remains in the open state before transitioning to Half-open.
		ReadyToTrip: func(counts gobreaker.Counts) bool { //A function that determines the condition for transitioning from Close to Open
			return counts.ConsecutiveFailures >= params.MaxFailRequest
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

			if stateChangeCounter != nil {
				stateChangeCounter.Add(context.Background(), 1,
					metric.WithAttributes(
						attribute.String("breaker_name", name),
						attribute.String("from_state", from.String()),
						attribute.String("to_state", to.String()),
					),
				)
			}
		},
		IsSuccessful: func(err error) bool { //A function that determines which errors are taken into account
			return servererror.ServerError(err)
		},
	}
	breaker := gobreaker.NewCircuitBreaker(setting)
	return breaker, nil
}
