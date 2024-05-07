// Package main LOMS Сервис
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"google.golang.org/grpc"
	"route256.ozon.ru/project/loms/internal/app"
	"route256.ozon.ru/project/loms/internal/pkg/grpcmw"
	"route256.ozon.ru/project/loms/internal/pkg/pgxstats"
	"route256.ozon.ru/project/loms/internal/pkg/tracelog"
	"route256.ozon.ru/project/loms/internal/pkg/tracerprovider"
)

const (
	listenAddr      = ":50001"
	httpGatewayAddr = ":8080"
	stopTimeout     = 2 * time.Second

	pgMasterDSNEnv = "PG_MASTER_DSN"
	pgSlaveDSNEnv  = "PG_SLAVE_DSN"

	kafkaAddrEnv  = "KAFKA_ADDR"
	kafkaTopicEnv = "KAFKA_TOPIC"

	jaegerAddrEnd        = "JAEGER_ADDR"
	jaegerServiceNameEnv = "JAEGER_SERVICE_NAME"
	jaegerEnvironmentEnv = "JAEGER_ENVIRONMENT"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	logger := slog.New(
		tracelog.WithTracePropagation(
			slog.NewJSONHandler(os.Stdout, nil)))

	metricsExporter, err := prometheus.New()
	checkErr(logger, err)

	meterProvider := metric.NewMeterProvider(metric.WithReader(metricsExporter))
	tracerProvider, err := tracerprovider.NewJaegerTracerProvider(
		os.Getenv(jaegerAddrEnd),
		os.Getenv(jaegerServiceNameEnv),
		os.Getenv(jaegerEnvironmentEnv),
	)
	checkErr(logger, err)

	pgxStat, err := pgxstats.New(meterProvider, tracerProvider)
	checkErr(logger, err)

	slaveConfig, err := pgx.ParseConfig(os.Getenv(pgSlaveDSNEnv))
	checkErr(logger, err)

	slaveConfig.Tracer = pgxStat

	slaveConnect, err := pgx.ConnectConfig(ctx, slaveConfig)
	checkErr(logger, err)

	masterConfig, err := pgx.ParseConfig(os.Getenv(pgMasterDSNEnv))
	checkErr(logger, err)
	masterConfig.Tracer = pgxStat

	masterConnect, err := pgx.ConnectConfig(ctx, masterConfig)
	checkErr(logger, err)

	config := app.Config{
		ListenAddr:      listenAddr,
		HTTPGatewayAddr: httpGatewayAddr,
		GrpcServerOpts: []grpc.ServerOption{
			grpc.ChainUnaryInterceptor(grpcmw.Validate, grpcmw.Logged(slog.Default())),
			grpc.StatsHandler(
				otelgrpc.NewServerHandler(
					otelgrpc.WithMeterProvider(meterProvider),
					otelgrpc.WithTracerProvider(tracerProvider),
					otelgrpc.WithPropagators(propagation.NewCompositeTextMapPropagator(
						propagation.TraceContext{}, propagation.Baggage{},
					)))),
		},
		MasterConnect: masterConnect,
		SlaveConnect:  slaveConnect,
		KafkaAddr:     os.Getenv(kafkaAddrEnv),
		KafkaTopic:    os.Getenv(kafkaTopicEnv),
		TraceProvider: tracerProvider,
		MeterProvider: meterProvider,
	}

	a, err := app.New(config)
	checkErr(logger, err)

	go func() {
		if appErr := a.Run(); appErr != nil {
			checkErr(logger, err)
		}
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), stopTimeout)
	go func() {
		defer cancel()
		checkErr(logger, a.Stop(ctx))
	}()

	<-ctx.Done()
}

func checkErr(logger *slog.Logger, err error) {
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
