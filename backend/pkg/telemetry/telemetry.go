package telemetry

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// Config define parâmetros de inicialização da telemetria.
type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	EnableTracing  bool
	OTLPEndpoint   string
	OTLPHeaders    string
	OTLPInsecure   bool
}

// Telemetry encapsula providers e exportadores.
type Telemetry struct {
	tracerProvider *sdktrace.TracerProvider
	meterProvider  *metric.MeterProvider
	promExporter   *prometheus.Exporter
}

// Init configura tracer/meter providers e exportadores.
func Init(ctx context.Context, cfg Config) (*Telemetry, error) {
	res, err := buildResource(ctx, cfg)
	if err != nil {
		return nil, err
	}

	promExporter, meterProvider, err := initPrometheusProvider(res)
	if err != nil {
		return nil, fmt.Errorf("init prometheus exporter: %w", err)
	}

	tracerProvider, err := initTracerProvider(ctx, res, cfg)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(tracerProvider)
	otel.SetMeterProvider(meterProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &Telemetry{
		tracerProvider: tracerProvider,
		meterProvider:  meterProvider,
		promExporter:   promExporter,
	}, nil
}

// Shutdown encerra providers e exportadores.
func (t *Telemetry) Shutdown(ctx context.Context) error {
	var err error
	if t == nil {
		return nil
	}
	if t.meterProvider != nil {
		err = errors.Join(err, t.meterProvider.Shutdown(ctx))
	}
	if t.tracerProvider != nil {
		err = errors.Join(err, t.tracerProvider.Shutdown(ctx))
	}
	return err
}

// TracerProvider expõe o provider de tracing.
func (t *Telemetry) TracerProvider() *sdktrace.TracerProvider {
	if t == nil {
		return nil
	}
	return t.tracerProvider
}

// MeterProvider expõe o provider de métricas.
func (t *Telemetry) MeterProvider() *metric.MeterProvider {
	if t == nil {
		return nil
	}
	return t.meterProvider
}

// MetricsHandler retorna o handler HTTP para /metrics.
func (t *Telemetry) MetricsHandler() http.Handler {
	if t == nil || t.promExporter == nil {
		return nil
	}
	return promhttp.Handler()
}

func initPrometheusProvider(res *resource.Resource) (*prometheus.Exporter, *metric.MeterProvider, error) {
	exporter, err := prometheus.New(prometheus.WithoutScopeInfo())
	if err != nil {
		return nil, nil, err
	}
	meterProvider := metric.NewMeterProvider(
		metric.WithReader(exporter),
		metric.WithResource(res),
	)
	return exporter, meterProvider, nil
}

func initTracerProvider(ctx context.Context, res *resource.Resource, cfg Config) (*sdktrace.TracerProvider, error) {
	if cfg.EnableTracing && cfg.OTLPEndpoint != "" {
		client := otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(cfg.OTLPEndpoint),
			otlptracehttp.WithURLPath("v1/traces"),
			otlptracehttp.WithHeaders(parseHeaders(cfg.OTLPHeaders)),
		)
		if cfg.OTLPInsecure {
			client = otlptracehttp.NewClient(
				otlptracehttp.WithEndpoint(cfg.OTLPEndpoint),
				otlptracehttp.WithInsecure(),
				otlptracehttp.WithURLPath("v1/traces"),
				otlptracehttp.WithHeaders(parseHeaders(cfg.OTLPHeaders)),
			)
		}
		exporter, err := otlptrace.New(ctx, client)
		if err != nil {
			return nil, fmt.Errorf("init otlp exporter: %w", err)
		}
		return sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(res),
		), nil
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
	), nil
}

func resourceAttributes(cfg Config) []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		semconv.ServiceName(cfg.ServiceName),
		attribute.String("deployment.environment", cfg.Environment),
	}
	if cfg.ServiceVersion != "" {
		attrs = append(attrs, semconv.ServiceVersion(cfg.ServiceVersion))
	}
	return attrs
}

func parseHeaders(raw string) map[string]string {
	headers := map[string]string{}
	if raw == "" {
		return headers
	}
	pairs := strings.Split(raw, ",")
	for _, pair := range pairs {
		parts := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(parts) != 2 {
			continue
		}
		headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return headers
}

func buildResource(ctx context.Context, cfg Config) (*resource.Resource, error) {
	base := resource.NewSchemaless(resourceAttributes(cfg)...)

	envRes, err := resource.New(ctx, resource.WithFromEnv())
	if err != nil {
		otel.Handle(fmt.Errorf("apply otel resource from environment: %w", err))
	} else if envRes != nil {
		if merged, mergeErr := resource.Merge(base, envRes); mergeErr == nil {
			base = merged
		} else {
			otel.Handle(fmt.Errorf("merge otel resource from environment: %w", mergeErr))
		}
	}

	sysRes, err := resource.New(
		ctx,
		resource.WithProcess(),
		resource.WithOS(),
		resource.WithTelemetrySDK(),
	)
	if err != nil {
		return nil, fmt.Errorf("build otel resource: %w", err)
	}

	merged, err := resource.Merge(base, sysRes)
	if err != nil {
		return nil, fmt.Errorf("merge otel resource: %w", err)
	}

	return merged, nil
}
