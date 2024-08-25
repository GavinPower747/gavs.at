package observability

import (
	"context"
	"errors"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func SetupOTelSDK(ctx context.Context, serviceName string) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			if shutErr := fn(ctx); shutErr != nil {
				err = errors.Join(err, shutErr)
			}
		}
		return err
	}

	handleErr := func(inErr error) {
		if inErr != nil {
			err = errors.Join(inErr, shutdown(ctx))
		}
	}

	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(os.Getenv("GIT_COMMIT")),
		attribute.String("build.commit", os.Getenv("GIT_COMMIT")),
		attribute.String("build.branch", os.Getenv("GIT_BRANCH")),
		attribute.String("build.date", os.Getenv("BUILD_DATE")),
		semconv.TelemetrySDKLanguageGo,
	)

	otel.SetTextMapPropagator(newPropagator())

	tracerProvider, err := newTraceProvider(ctx, resource)
	if err != nil {
		handleErr(err)
		return
	}

	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	return shutdown, nil
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider(ctx context.Context, resource *resource.Resource) (*trace.TracerProvider, error) {
	exporter, err := otlptracehttp.New(ctx)

	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource),
	)
	return tp, nil
}
