// Package main Основное приложение
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/multierr"
	"route256.ozon.ru/project/cart/internal/app"
	"route256.ozon.ru/project/cart/internal/pkg/tracelog"
	"route256.ozon.ru/project/cart/internal/pkg/tracerprovider"
)

const (
	stopTimeout = 2 * time.Second

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

	conf := app.ReadConfig()
	conf.Logger = logger
	conf.MeterProvider = meterProvider
	conf.TracerProvider = tracerProvider
	a := app.New(conf)
	go func() {
		checkErr(logger, a.Run())
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), stopTimeout)
	go func() {
		defer cancel()
		err := multierr.Combine(a.Stop(ctx), tracerProvider.Shutdown(ctx))
		checkErr(logger, err)
	}()

	<-ctx.Done()
}

func checkErr(logger *slog.Logger, err error) {
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
