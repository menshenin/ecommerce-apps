// Package stock Реализация репозитория с информацией по складам
package stock

import (
	"context"
	"sync"

	"route256.ozon.ru/project/loms/internal/model"
)

// Repository Репозиторий
type Repository struct {
	stock Stock
}

// New Конструктор
func New() *Repository {
	return &Repository{stock: Stock{
		mutex:   &sync.RWMutex{},
		storage: make(map[model.SKU]model.StockItem),
	}}
}

// GetCurrentStock Получение текущего склада
func (r *Repository) GetCurrentStock() (*model.Stock, error) {
	return &model.Stock{
		Reserve: func(_ context.Context, items []model.OrderItem) error {
			for _, item := range items {
				if err := r.stock.Reserve(item.SKU, item.Count); err != nil {
					return err
				}
			}
			return nil
		},
		AvailableCount: func(_ context.Context, sku model.SKU) (int32, error) {
			return r.stock.AvailableCount(sku)
		},
		CancelReserve: func(_ context.Context, items []model.OrderItem) error {
			for _, item := range items {
				if err := r.stock.CancelReserve(item.SKU, item.Count); err != nil {
					return err
				}
			}
			return nil
		},
		WriteOff: func(_ context.Context, items []model.OrderItem) error {
			for _, item := range items {
				if err := r.stock.WriteOff(item.SKU, item.Count); err != nil {
					return err
				}
			}
			return nil
		},
		Load: func(_ context.Context, items []model.StockItem) error {
			r.stock.Load(items)
			return nil
		},
	}, nil
}
