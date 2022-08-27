package CommenDb

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	tr "go.opentelemetry.io/otel/trace"
)

// tracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func tracerProvider(service string, version string) (*trace.TracerProvider, error) {
	// Create the Jaeger exporter

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://jaeger:14268/api/traces")))
	if err != nil {
		return nil, err
	}
	tp := trace.NewTracerProvider(
		// Always be sure to batch in production.
		trace.WithBatcher(exp),
		// Record information about this application in a Resource.
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			semconv.ServiceVersionKey.String(version),
		)),
	)
	return tp, nil
}
func (a *Api) CreateNewTracer(ServiceVersion string) tr.Tracer {
	tp, err := tracerProvider(a.Service, ServiceVersion)
	if err != nil {
		println(err)
		return nil
	}
	otel.SetTracerProvider(tp)
	newTrace := tp.Tracer(a.Service)
	return newTrace
}
func (a *Api) CreateSpanWithTraceId() (tr.Span, context.Context) {
	var traceId tr.TraceID
	traceId, err := tr.TraceIDFromHex(a.TraceId)
	if err != nil {
		fmt.Println("trace id is not valid : ", err.Error())
	}
	var spanContextConfig tr.SpanContextConfig
	spanContextConfig.TraceID = traceId
	spanContextConfig.TraceFlags = 01
	spanContextConfig.Remote = false
	spanContext := tr.NewSpanContext(spanContextConfig)
	requestContext := context.Background()
	requestContext = tr.ContextWithSpanContext(requestContext, spanContext)
	tracer := a.CreateNewTracer(a.ServiceVersion)
	ctx, span := tracer.Start(requestContext, a.Route)
	return span, ctx
}
func (a *Api) CreateContextWithTraceIdAndSpanId() context.Context {
	var traceId tr.TraceID
	traceId, err := tr.TraceIDFromHex(a.TraceId)
	if err != nil {
		fmt.Println("trace id is not valid : ", err.Error())
	}
	var spanID tr.SpanID
	spanID, err = tr.SpanIDFromHex(a.SpanId)
	if err != nil {
		fmt.Println("span id is not valid : ", err.Error())
	}
	var spanContextConfig tr.SpanContextConfig
	spanContextConfig.TraceID = traceId
	spanContextConfig.SpanID = spanID
	spanContextConfig.TraceFlags = 01
	spanContextConfig.Remote = false
	spanContext := tr.NewSpanContext(spanContextConfig)
	requestContext := context.Background()
	requestContext = tr.ContextWithSpanContext(requestContext, spanContext)
	return requestContext
}
func (a *Api) CreateChildSpanWithTraceIdAndSpanId() (tr.Span, context.Context) {
	var traceId tr.TraceID
	traceId, err := tr.TraceIDFromHex(a.TraceId)
	if err != nil {
		fmt.Println("trace id is not valid : ", err.Error())
	}
	var spanID tr.SpanID
	spanID, err = tr.SpanIDFromHex(a.SpanId)
	if err != nil {
		fmt.Println("span id is not valid : ", err.Error())
	}
	var spanContextConfig tr.SpanContextConfig
	spanContextConfig.TraceID = traceId
	spanContextConfig.SpanID = spanID
	spanContextConfig.TraceFlags = 01
	spanContextConfig.Remote = false
	spanContext := tr.NewSpanContext(spanContextConfig)
	requestContext := context.Background()
	requestContext = tr.ContextWithSpanContext(requestContext, spanContext)
	tracer := a.CreateNewTracer(a.ServiceVersion)
	ctx, span := tracer.Start(requestContext, a.SpanId)
	return span, ctx
}
