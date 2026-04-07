package opentelemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func NewTrace(ctx context.Context, serviceName, host, port string) (*sdktrace.TracerProvider, error) {

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(host+":"+port),
	)
	if err != nil {
		return nil, err
	}
	sampler := sdktrace.TraceIDRatioBased(0.3) //Процент выборки для обработки спанов 30%

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter), //Обработка спанов пачками
		sdktrace.WithSampler(sampler),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName(serviceName),
			),
		),
	)

	otel.SetTracerProvider(tp)
	return tp, nil
}
