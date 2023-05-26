package middleware

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func NewTracing(ctx context.Context, userID int64) context.Context {
	tp := otel.GetTracerProvider()
	tracer := tp.Tracer("")
	ctx, span, _ := startTrace(ctx, tracer, "11111", userID)
	defer span.End()

	return ctx
}

func startTrace(ctx context.Context, tracer trace.Tracer, spanName string, userID int64) (context.Context, trace.Span, string) {
	ctx, span, traceID := newTrace(ctx, tracer, spanName)
	span.SetAttributes(
		attribute.String("trace_id", traceID),
		attribute.Int64("user_id", userID),
	)

	return ctx, span, traceID
}

// NewTrace запускает новую трассировку
func newTrace(parentCtx context.Context,
	tracer trace.Tracer, spanName string) (ctx context.Context, span trace.Span, debugID string) {
	ctx, span = trace.SpanFromContext(parentCtx).TracerProvider().Tracer("").Start(parentCtx, spanName)
	if !span.IsRecording() {
		ctx, span = tracer.Start(parentCtx, spanName)
	}
	traceID := span.SpanContext().TraceID().String()

	return ctx, span, traceID
}
