package app

import (
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// Config Конфиг
type Config struct {
	ListenAddr      string
	HTTPGatewayAddr string
	SwaggerJSONPath string
	GrpcServerOpts  []grpc.ServerOption
	MasterConnect   *pgx.Conn
	SlaveConnect    *pgx.Conn
	KafkaAddr       string
	KafkaTopic      string
	TraceProvider   trace.TracerProvider
	MeterProvider   metric.MeterProvider
}
