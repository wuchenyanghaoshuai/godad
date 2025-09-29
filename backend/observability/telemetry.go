package observability

import (
	"context"
	"log"
	"sync/atomic"
	"time"

	"godad-backend/config"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var tracingEnabled atomic.Bool

// InitTracing configures OpenTelemetry tracing if exporter endpoint is provided.
func InitTracing(ctx context.Context, cfg *config.Config) func(context.Context) error {
	endpoint := cfg.Observability.OTLPExporterEndpoint
	if endpoint == "" {
		tracingEnabled.Store(false)
		return func(context.Context) error { return nil }
	}

	opts := []otlptracegrpc.Option{otlptracegrpc.WithEndpoint(endpoint)}
	if cfg.Observability.OTLPInsecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		log.Printf("OpenTelemetry exporter 初始化失败: %v", err)
		tracingEnabled.Store(false)
		return func(context.Context) error { return nil }
	}

	res, err := resource.Merge(resource.Default(), resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("godad-backend"),
		attribute.String("service.environment", cfg.Server.Environment),
	))
	if err != nil {
		res = resource.Default()
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter,
			tracesdk.WithBatchTimeout(5*time.Second),
			tracesdk.WithMaxExportBatchSize(256),
		),
		tracesdk.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{},
	))

	tracingEnabled.Store(true)
	log.Println("OpenTelemetry tracing 已启用")

	return func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			return err
		}
		return exporter.Shutdown(ctx)
	}
}

// TracingEnabled 返回是否启用了 tracing。
func TracingEnabled() bool {
	return tracingEnabled.Load()
}

// GinTraceMiddleware returns an OTEL middleware when tracing is enabled.
func GinTraceMiddleware(serviceName string) gin.HandlerFunc {
	if !TracingEnabled() {
		return nil
	}
	return otelgin.Middleware(serviceName)
}
