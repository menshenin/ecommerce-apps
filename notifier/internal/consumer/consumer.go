// Package consumer Kafka-консьюмер
package consumer

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

// Consumer Консьюмер
type Consumer struct {
	client *kgo.Client
	logger *slog.Logger
	tracer trace.Tracer
}

// New Конструктор
func New(client *kgo.Client, logger *slog.Logger, tracerProvider trace.TracerProvider) *Consumer {
	return &Consumer{
		client: client,
		logger: logger,
		tracer: tracerProvider.Tracer("consumer"),
	}
}

// Run Запуск консьюмера
func (c *Consumer) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			c.client.Close()
			return
		default:
			fetches := c.client.PollFetches(ctx)
			iter := fetches.RecordIter()
			for !iter.Done() {
				c.processMessage(ctx, iter.Next())
			}
		}
	}
}

func (c *Consumer) processMessage(ctx context.Context, record *kgo.Record) {
	span := c.startSpan(ctx, record)
	defer span.End()
	c.logger.InfoContext(ctx, "Message received", "data", string(record.Value))
}

func (c *Consumer) startSpan(ctx context.Context, record *kgo.Record) trace.Span {
	span := trace.Span(noop.Span{})
	var data struct {
		Extra map[string]any `json:"extra"`
	}
	err := json.Unmarshal(record.Value, &data)
	if err != nil {
		c.logger.ErrorContext(ctx, err.Error())
		return span
	}

	if s, ok := data.Extra["trace_id"].(string); ok {
		traceID, err := trace.TraceIDFromHex(s)
		if err != nil {
			c.logger.Error(err.Error())
			return span
		}
		if s, ok := data.Extra["span_id"].(string); ok {
			spanID, err := trace.SpanIDFromHex(s)
			if err != nil {
				c.logger.ErrorContext(ctx, err.Error())
				return span
			}

			ctx = trace.ContextWithRemoteSpanContext(ctx, trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.FlagsSampled,
			}))

			_, span = c.tracer.Start(ctx, "Process message", trace.WithSpanKind(trace.SpanKindConsumer))
		}
	}

	return span
}
