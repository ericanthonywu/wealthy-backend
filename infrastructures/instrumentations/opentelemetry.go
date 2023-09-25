package instrumentations

import (
	"log"
	"os"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semConv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func OpenTelemetryExporter(appName, appEnv string) *trace.TracerProvider {
	l := log.New(os.Stdout, "", 0)

	exp, err := JaegerExporter()
	if err != nil {
		l.Fatal(err)
	}

	return trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(ResourceTrack(appName, appEnv)),
	)
}

func ResourceTrack(appName, appEnvironment string) *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semConv.SchemaURL,
			semConv.ServiceName(appName),
			semConv.ServiceVersion("v1"),
			attribute.String(appEnvironment, appEnvironment),
		),
	)
	return r
}
