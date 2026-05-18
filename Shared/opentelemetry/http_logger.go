package opentelemetry

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
)

func NewLoggerProviderHttp(ctx context.Context, serviceName, collectorHost, collectorPort string) (*sdklog.LoggerProvider, error) {

	exporter, err := otlploghttp.New(
		ctx,
		otlploghttp.WithEndpoint(collectorHost+":"+collectorPort),
		otlploghttp.WithURLPath("/v1/logs"),
		otlploghttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	provider := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(
			sdklog.NewBatchProcessor(
				exporter,
			),
		),
	)

	global.SetLoggerProvider(provider)
	return provider, nil
}
