package opentelemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
)

func NewMetricProviderHttp(ctx context.Context, serverName, collectorHost, collectorPort string, interval time.Duration) (*metric.MeterProvider, error) {
	fmt.Println(collectorHost, collectorPort, interval)
	exporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpoint(collectorHost+":"+collectorPort),
		otlpmetrichttp.WithURLPath("/v1/metrics"),
		otlpmetrichttp.WithInsecure(),
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
