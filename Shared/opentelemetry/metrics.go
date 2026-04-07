package opentelemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func NewMetricPrometeus(ctx context.Context, serverName string) (*metric.MeterProvider, error) {
	//Создание экспортера
	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	//Информация о сервере
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serverName),
		),
	)
	if err != nil {
		return nil, err
	}

	//Создание провайдера
	provider := metric.NewMeterProvider(
		metric.WithReader(exporter),
		metric.WithResource(res),
	)

	otel.SetMeterProvider(provider)

	return provider, nil
}
