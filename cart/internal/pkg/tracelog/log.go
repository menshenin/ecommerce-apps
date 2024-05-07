package tracelog

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

type logger struct {
	slog.Handler
}

func (l logger) Handle(ctx context.Context, record slog.Record) error {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasSpanID() {
		record.Add("span_id", spanCtx.SpanID().String())
	}
	if spanCtx.HasTraceID() {
		record.Add("trace_id", spanCtx.TraceID().String())
	}

	return l.Handler.Handle(ctx, record)
}

func WithTracePropagation(handler slog.Handler) slog.Handler {
	return logger{handler}
}
