// Package main Основное приложение
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"route256.ozon.ru/project/notifier/internal/app"
	"route256.ozon.ru/project/notifier/internal/pkg/tracelog"
	"route256.ozon.ru/project/notifier/internal/pkg/tracerprovider"
)

const (
	kafkaAddrEnv          = "KAFKA_ADDR"
	kafkaTopicEnv         = "KAFKA_TOPIC"
	kafkaConsumerGroupEnv = "CONSUMER_GROUP"

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

	tracerProvider, err := tracerprovider.NewJaegerTracerProvider(
		os.Getenv(jaegerAddrEnd),
		os.Getenv(jaegerServiceNameEnv),
		os.Getenv(jaegerEnvironmentEnv),
	)

	checkErr(logger, err)

	conf := app.Config{
		KafkaAddr:      os.Getenv(kafkaAddrEnv),
		KafkaTopic:     os.Getenv(kafkaTopicEnv),
		ConsumerGroup:  os.Getenv(kafkaConsumerGroupEnv),
		TracerProvider: tracerProvider,
		Logger:         logger,
	}

	a, err := app.New(conf)
	checkErr(logger, err)
	done := make(chan struct{})
	go func() {
		a.Run(ctx)
		done <- struct{}{}
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), stopTimeout)
	go func() {
		defer cancel()
		<-done
		checkErr(logger, tracerProvider.Shutdown(ctx))
	}()

	<-ctx.Done()
}

func checkErr(logger *slog.Logger, err error) {
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
