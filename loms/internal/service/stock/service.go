// Package stock Сервис для работы со складами
package stock

import (
	"github.com/go-playground/validator/v10"
	"route256.ozon.ru/project/loms/internal/model"
)

// StocksRepository Репозиторий для получения данных по складам
//
//go:generate minimock -i StocksRepository -s "_mock.go" -o ./mocks/
type StocksRepository interface {
	GetCurrentStock() (*model.Stock, error)
}

// Service Сервис
type Service struct {
	stocksRepo StocksRepository
	validate   *validator.Validate
}

// New Конструктор
func New(stocksRepo StocksRepository) *Service {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return &Service{
		stocksRepo: stocksRepo,
		validate:   validate,
	}
}
