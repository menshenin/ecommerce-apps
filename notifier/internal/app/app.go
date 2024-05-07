// Package app Приложение
package app

import (
	"context"
	"log/slog"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.opentelemetry.io/otel/trace"
	"route256.ozon.ru/project/notifier/internal/consumer"
)

// Config Конфиг приложения
type Config struct {
	KafkaAddr      string
	KafkaTopic     string
	ConsumerGroup  string
	TracerProvider trace.TracerProvider
	Logger         *slog.Logger
}

// App Приложение
type App struct {
	cons *consumer.Consumer
}

// New Конструктор
func New(conf Config) (*App, error) {
	kafkaClient, err := kgo.NewClient(
		kgo.SeedBrokers(conf.KafkaAddr),
		kgo.ConsumerGroup(conf.ConsumerGroup),
		kgo.ConsumeTopics(conf.KafkaTopic))

	if err != nil {
		return nil, err
	}

	return &App{
		cons: consumer.New(kafkaClient, conf.Logger, conf.TracerProvider),
	}, nil
}

// Run Запуск приложения
func (a *App) Run(ctx context.Context) {
	a.cons.Run(ctx)
}
