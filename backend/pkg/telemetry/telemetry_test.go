package telemetry

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"github.com/stretchr/testify/assert"
)

func TestBuildResourceIgnoresInvalidEnvAttributes(t *testing.T) {
	t.Setenv("OTEL_RESOURCE_ATTRIBUTES", "invalid pattern")

	res, err := buildResource(context.Background(), Config{ServiceName: "test-svc"})
	if err != nil {
		t.Fatalf("expected buildResource not to fail, got %v", err)
	}

	if res == nil {
		t.Fatalf("expected buildResource to return a resource instance")
	}

	foundService := false
	for _, attr := range res.Attributes() {
		if attr.Key == semconv.ServiceNameKey && attr.Value.AsString() == "test-svc" {
			foundService = true
			break
		}
	}

	if !foundService {
		t.Fatalf("expected service.name attribute to be present")
	}
}

func TestBuildResourceMergesEnvAttributes(t *testing.T) {
	t.Setenv("OTEL_RESOURCE_ATTRIBUTES", "team=qa,service.version=1.2.3")

	res, err := buildResource(context.Background(), Config{ServiceName: "test-svc"})
	if err != nil {
		t.Fatalf("expected buildResource not to fail, got %v", err)
	}

	attrs := res.Attributes()
	teamAttr := attribute.Key("team")
	serviceVersion := semconv.ServiceVersionKey

	assertAttr := func(key attribute.Key, expected string) {
		for _, attr := range attrs {
			if attr.Key == key && attr.Value.AsString() == expected {
				return
			}
		}
		t.Fatalf("expected attribute %s=%s to be present", key, expected)
	}

	assertAttr(teamAttr, "qa")
	assertAttr(serviceVersion, "1.2.3")
}

func TestTelemetryAccessors(t *testing.T) {
	tel := &Telemetry{}

	// Test TracerProvider
	tp := tel.TracerProvider()
	assert.Nil(t, tp) // Should be nil if not initialized

	// Test MeterProvider
	mp := tel.MeterProvider()
	assert.Nil(t, mp) // Should be nil if not initialized

	// Test MetricsHandler
	mh := tel.MetricsHandler()
	assert.Nil(t, mh) // Should be nil if not initialized

	// Initialize Telemetry with dummy providers
	tel.tracerProvider = sdktrace.NewTracerProvider()
	tel.meterProvider = metric.NewMeterProvider()
	// promExporter is needed for MetricsHandler
	exporter, _ := prometheus.New(prometheus.WithoutScopeInfo())
	tel.promExporter = exporter

	tp = tel.TracerProvider()
	assert.NotNil(t, tp)

	mp = tel.MeterProvider()
	assert.NotNil(t, mp)

	mh = tel.MetricsHandler()
	assert.NotNil(t, mh)
}

func TestParseHeaders(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want map[string]string
	}{
		{
			name: "empty string",
			raw:  "",
			want: map[string]string{},
		},
		{
			name: "single header",
			raw:  "key=value",
			want: map[string]string{"key": "value"},
		},
		{
			name: "multiple headers",
			raw:  "key1=value1, key2=value2",
			want: map[string]string{"key1": "value1", "key2": "value2"},
		},
		{
			name: "headers with spaces",
			raw:  " key1 = value1 , key2 = value2 ",
			want: map[string]string{"key1": "value1", "key2": "value2"},
		},
		{
			name: "invalid header format",
			raw:  "key1=value1, invalid-header, key2=value2",
			want: map[string]string{"key1": "value1", "key2": "value2"},
		},
		{
			name: "header with multiple equals signs",
			raw:  "key=value=with=equals",
			want: map[string]string{"key": "value=with=equals"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseHeaders(tt.raw)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInitPrometheusProvider(t *testing.T) {
	res := resource.NewSchemaless(semconv.ServiceName("test-service"))
	exporter, meterProvider, err := initPrometheusProvider(res)
	assert.NoError(t, err)
	assert.NotNil(t, exporter)
	assert.NotNil(t, meterProvider)
}

func TestInitTracerProvider(t *testing.T) {
	res := resource.NewSchemaless(semconv.ServiceName("test-service"))

	// Test with tracing disabled
	tp, err := initTracerProvider(context.Background(), res, Config{EnableTracing: false})
	assert.NoError(t, err)
	assert.NotNil(t, tp)

	// Test with tracing enabled and OTLP endpoint
	tp, err = initTracerProvider(context.Background(), res, Config{
		EnableTracing: true,
		OTLPEndpoint:  "http://localhost:4318",
		OTLPHeaders:   "key=value",
		OTLPInsecure:  false,
	})
	assert.NoError(t, err)
	assert.NotNil(t, tp)

	// Test with tracing enabled, OTLP endpoint and insecure
	tp, err = initTracerProvider(context.Background(), res, Config{
		EnableTracing: true,
		OTLPEndpoint:  "http://localhost:4318",
		OTLPHeaders:   "key=value",
		OTLPInsecure:  true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, tp)
}

func TestInit(t *testing.T) {
	ctx := context.Background()
	cfg := Config{
		ServiceName:    "test-service",
		ServiceVersion: "1.0.0",
		Environment:    "test",
		EnableTracing:  true,
		OTLPEndpoint:   "http://localhost:4318",
		OTLPHeaders:    "key=value",
		OTLPInsecure:   false,
	}

	telemetry, err := Init(ctx, cfg)
	assert.NoError(t, err)
	assert.NotNil(t, telemetry)
	assert.NotNil(t, telemetry.tracerProvider)
	assert.NotNil(t, telemetry.meterProvider)
	assert.NotNil(t, telemetry.promExporter)

	// Test Init with tracing disabled
	cfg.EnableTracing = false
	telemetry, err = Init(ctx, cfg)
	assert.NoError(t, err)
	assert.NotNil(t, telemetry)
	assert.NotNil(t, telemetry.tracerProvider)
	assert.NotNil(t, telemetry.meterProvider)
	assert.NotNil(t, telemetry.promExporter)

	// Test Init with invalid OTLP endpoint (should not return error from Init directly, but from initTracerProvider)
	cfg.EnableTracing = true
	cfg.OTLPEndpoint = "invalid-endpoint"
	telemetry, err = Init(ctx, cfg)
	assert.NoError(t, err) // Init itself should not return error for invalid endpoint, initTracerProvider handles it
	assert.NotNil(t, telemetry)
}

func TestShutdown(t *testing.T) {
	ctx := context.Background()

	// Test with nil Telemetry
	var nilTel *Telemetry
	err := nilTel.Shutdown(ctx)
	assert.NoError(t, err)

	// Test with initialized Telemetry
	cfg := Config{
		ServiceName:    "test-service",
		ServiceVersion: "1.0.0",
		Environment:    "test",
		EnableTracing:  true,
		OTLPEndpoint:   "http://localhost:4318",
		OTLPHeaders:    "key=value",
		OTLPInsecure:   false,
	}
	telemetry, err := Init(ctx, cfg)
	assert.NoError(t, err)
	assert.NotNil(t, telemetry)

	err = telemetry.Shutdown(ctx)
	assert.NoError(t, err)
}
