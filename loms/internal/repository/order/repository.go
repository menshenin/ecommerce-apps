// Package order Реализация репозитория заказов
package order

import (
	"context"
	"sync"

	"route256.ozon.ru/project/loms/internal/model"
)

// Order Заказ
type Order struct {
	id     model.OrderID
	userID model.UserID
	status model.OrderStatus
	items  []model.OrderItem
}

// Repository Репозиторий
type Repository struct {
	mu      *sync.RWMutex
	storage map[model.OrderID]*Order
}

// GetByID Получение заказа по ID
func (r *Repository) GetByID(_ context.Context, id model.OrderID) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if order, ok := r.storage[id]; ok {
		return hydrate(order), nil
	}
	return nil, model.ErrOrderNotFound
}

// Create Создание заказа
func (r *Repository) Create(_ context.Context, userID model.UserID, items []model.OrderItem) (*model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := r.generateID()
	order := &Order{
		id:     id,
		userID: userID,
		status: model.OrderStatusNew,
		items:  items,
	}
	r.storage[id] = order

	return hydrate(order), nil
}

func (r *Repository) generateID() model.OrderID {
	if len(r.storage) == 0 {
		return 1
	}
	return model.OrderID(len(r.storage) + 1)
}

// New Конструктор
func New() *Repository {
	return &Repository{
		mu:      &sync.RWMutex{},
		storage: make(map[model.OrderID]*Order),
	}
}

func hydrate(o *Order) *model.Order {
	return &model.Order{
		ID:     o.id,
		UserID: o.userID,
		UpdateStatus: func(_ context.Context, status model.OrderStatus) error {
			o.status = status
			return nil
		},
		GetItems: func(_ context.Context) ([]model.OrderItem, error) {
			return o.items, nil
		},
		GetStatus: func(_ context.Context) (model.OrderStatus, error) {
			return o.status, nil
		},
	}
}
