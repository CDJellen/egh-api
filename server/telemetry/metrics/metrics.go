package metrics

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric/global"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func NewGrpcMetricsProvider(ctx context.Context, serviceName string, metricAddr string, freq int) (provider *sdkmetric.MeterProvider, err error) {
	exporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithEndpoint(metricAddr))
	if err != nil {
		log.Printf("failed to obtain grpc trace: %+v", err)
		return nil, err
	}

	// define applicable resource tags
	resources := resource.NewWithAttributes(semconv.SchemaURL, attribute.String("tracer address", metricAddr), attribute.String("service name", "egh-api"))

	// define a new periodic reader
	interval := sdkmetric.WithInterval(time.Duration(freq) * time.Second)
	pReader := sdkmetric.NewPeriodicReader(exporter, interval)

	// define a new metrics provider
	provider = sdkmetric.NewMeterProvider(sdkmetric.WithResource(resources), sdkmetric.WithReader(pReader))

	global.SetMeterProvider(provider)

	return provider, nil

}
