// Package order Сервис для работы с заказами
package order

import (
	"context"

	"route256.ozon.ru/project/loms/internal/event"
	"route256.ozon.ru/project/loms/internal/model"
)

// EventProducer Эмиттер событий заказа
//
//go:generate minimock -i EventProducer -s "_mock.go" -o ./mocks/
type EventProducer interface {
	Produce(ctx context.Context, event event.OrderEvent)
}

// OrdersRepository Репозиторий заказов
//
//go:generate minimock -i OrdersRepository -s "_mock.go" -o ./mocks/
type OrdersRepository interface {
	GetByID(ctx context.Context, id model.OrderID) (*model.Order, error)
	Create(ctx context.Context, userID model.UserID, items []model.OrderItem) (*model.Order, error)
}

// StocksRepository Репозиторий складов
//
//go:generate minimock -i StocksRepository -s "_mock.go" -o ./mocks/
type StocksRepository interface {
	GetCurrentStock() (*model.Stock, error)
}

// Service Сервис
type Service struct {
	eventProducer EventProducer
	ordersRepo    OrdersRepository
	stocksRepo    StocksRepository
}

// New Конструктор
func New(ordersRepo OrdersRepository, stocksRepo StocksRepository, eventProducer EventProducer) *Service {
	return &Service{
		eventProducer: eventProducer,
		ordersRepo:    ordersRepo,
		stocksRepo:    stocksRepo,
	}
}
