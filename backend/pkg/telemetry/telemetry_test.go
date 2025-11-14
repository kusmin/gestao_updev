package telemetry

import (
	"context"
	"testing"

	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
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
