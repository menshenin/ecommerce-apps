// Package producer Продьюсер событий
package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.opentelemetry.io/otel/trace"
	"route256.ozon.ru/project/loms/internal/event"
)

// Producer Продьюсер событий
type Producer struct {
	client *kgo.Client
	logger *slog.Logger
	tracer trace.Tracer
}

// New Конструктор
func New(client *kgo.Client, logger *slog.Logger, traceProvider trace.TracerProvider) *Producer {
	return &Producer{
		client: client,
		logger: logger,
		tracer: traceProvider.Tracer("kafka_producer"),
	}
}

// Produce Отправка событий
func (e *Producer) Produce(ctx context.Context, orderEvent event.OrderEvent) {
	var span trace.Span
	if e.tracer != nil {
		span = e.startSpan(ctx, &orderEvent)
	}

	data, err := json.Marshal(orderEvent)
	if err != nil {
		e.logger.ErrorContext(ctx, "serialize event error", "err", err)
		return
	}

	r := &kgo.Record{
		Key:   []byte(strconv.Itoa(int(orderEvent.OrderID))),
		Value: data,
	}

	e.client.Produce(context.Background(), r, func(record *kgo.Record, err error) {
		if span != nil {
			span.End()
		}
		if err != nil {
			e.logger.ErrorContext(ctx, "error on producing", "err", err, "event", string(record.Value))
		}
	})
}

// Stop Остановка продьюсера
func (e *Producer) Stop() {
	e.client.Close()
}

func (e *Producer) startSpan(ctx context.Context, orderEvent *event.OrderEvent) trace.Span {
	ctx, span := e.tracer.Start(
		ctx,
		fmt.Sprintf("Sending event: %s", orderEvent.EventType),
		trace.WithSpanKind(trace.SpanKindProducer))
	if orderEvent.Extra == nil {
		orderEvent.Extra = make(event.ExtraData, 2)
	}
	orderEvent.Extra["span_id"] = span.SpanContext().SpanID().String()
	orderEvent.Extra["trace_id"] = span.SpanContext().TraceID().String()

	return span
}
