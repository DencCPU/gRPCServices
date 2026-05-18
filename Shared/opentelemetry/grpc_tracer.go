package opentelemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func NewTracerProviderGrpc(ctx context.Context, serviceName, collectorHost, collectorPort string, percentage int) (*sdktrace.TracerProvider, error) {

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(collectorHost+":"+collectorPort),
	)
	if err != nil {
		return nil, err
	}
	// fmt.Println(host + ":" + port)
	fraction := float64(percentage) / 100
	sampler := sdktrace.TraceIDRatioBased(fraction) //spans processing percentage

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter), //Processing spans in batches
		sdktrace.WithSampler(sampler),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName(serviceName),
			),
		),
	)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	otel.SetTracerProvider(tp)
	return tp, nil
}
