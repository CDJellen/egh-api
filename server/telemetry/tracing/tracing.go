package tracing

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func NewGrpcTraceProvider(ctx context.Context, serviceName string, traceAddr string, sampling float64) (provider *trace.TracerProvider, err error) {
	batcher, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(traceAddr))
	if err != nil {
		log.Printf("failed to obtain grpc trace: %+v", err)
		return nil, err
	}

	// define applicable resource tags
	resources := resource.NewWithAttributes(semconv.SchemaURL, attribute.String("tracer address", traceAddr), attribute.String("service name", "egh-api"))

	// define sampling rate
	sampler := trace.ParentBased(trace.TraceIDRatioBased(sampling))

	// define a new tracing provider
	provider = trace.NewTracerProvider(trace.WithBatcher(batcher), trace.WithResource(resources), trace.WithSampler(sampler))

	// define a new tracing propagator
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{})

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagator)

	return provider, nil

}
