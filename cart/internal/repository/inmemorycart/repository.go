package inmemorycart

import (
	"context"
	"sync"

	"route256.ozon.ru/project/cart/internal/model"
)

// Repository in-memory репозиторий для хранения корзин
type Repository struct {
	mutex   *sync.RWMutex
	storage map[model.UserID]*Cart
}

// New Конструктор
func New() *Repository {
	return &Repository{
		mutex:   &sync.RWMutex{},
		storage: make(map[model.UserID]*Cart),
	}
}

// GetByUserID Получение корзины
func (r *Repository) GetByUserID(_ context.Context, userID model.UserID) (*model.Cart, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if cart, ok := r.storage[userID]; ok {
		return hydrate(cart), nil
	}
	return nil, model.ErrCartNotFound
}

// Create Создание корзины
func (r *Repository) Create(_ context.Context, userID model.UserID) (*model.Cart, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.storage[userID] = NewCart(userID)

	return hydrate(r.storage[userID]), nil
}

func hydrate(cart *Cart) *model.Cart {
	return &model.Cart{
		UserID:          cart.UserID,
		GetItems:        cart.GetItems,
		AddItemBySKU:    cart.AddItemBySKU,
		DeleteItemBySKU: cart.DeleteItemBySKU,
		Clear:           cart.Clear,
	}
}
