// Package pgxstats Предоставляет OpenTelemetry метрики и трассировку для pgx
package pgxstats

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const sqlOperationUnknown = "UNKNOWN"

type queryDataKey struct{}

type queryData struct {
	time time.Time
	op   string
}

// PgxStats Метрики и трассировщик
type PgxStats struct {
	tracer       trace.Tracer
	queryCounter metric.Int64Counter
	queryTime    metric.Float64Histogram
}

// New Конструктор
func New(meterProvider metric.MeterProvider, tracerProvider trace.TracerProvider) (pgxStat *PgxStats, err error) {
	pgxStat = &PgxStats{
		tracer: tracerProvider.Tracer("pgxstat"),
	}
	meter := meterProvider.Meter("pgxstat")
	pgxStat.queryCounter, err = meter.Int64Counter("query_count",
		metric.WithUnit("{count}"),
		metric.WithDescription("Total number of queries"))
	if err != nil {
		return
	}
	pgxStat.queryTime, err = meter.Float64Histogram("response_time",
		metric.WithUnit("ms"),
		metric.WithDescription("Query time in milliseconds"))

	return
}

// TraceQueryStart Начало выполнения запроса
func (p *PgxStats) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	queryOpName := sqlOperationName(data.SQL)
	p.queryCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("operation_name", queryOpName)))

	ctx, _ = p.tracer.Start(ctx, "query "+data.SQL, trace.WithSpanKind(trace.SpanKindClient))

	return context.WithValue(ctx, queryDataKey{}, queryData{
		time: time.Now(),
		op:   queryOpName,
	})
}

// TraceQueryEnd Конец выполнения запроса
func (p *PgxStats) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	qData := ctx.Value(queryDataKey{}).(queryData)
	attributes := []attribute.KeyValue{
		attribute.String("operation_name", qData.op),
	}
	if data.Err != nil {
		attributes = append(attributes, attribute.String("error", data.Err.Error()))
	}
	p.queryTime.Record(ctx,
		float64(time.Since(qData.time))/float64(time.Millisecond),
		metric.WithAttributes(attributes...))

	span := trace.SpanFromContext(ctx)
	span.End()
}

func sqlOperationName(sql string) string {
	parts := strings.Fields(sql)
	if len(parts) == 0 {
		return sqlOperationUnknown
	}

	return strings.ToUpper(parts[0])
}
