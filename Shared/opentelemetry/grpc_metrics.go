package opentelemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func NewMetricProviderGrpc(ctx context.Context, serverName, collectorHost, collectorPort string, interval time.Duration) (*metric.MeterProvider, error) {
	//Создание экспортера
	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(collectorHost+":"+collectorPort),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	reader := metric.NewPeriodicReader(exporter, metric.WithInterval(interval))

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
		metric.WithReader(reader),
		metric.WithResource(res),
	)

	otel.SetMeterProvider(provider)

	return provider, nil
}
