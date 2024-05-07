// Package inmemorycart Корзина в памяти
package inmemorycart

import (
	"sync"

	"route256.ozon.ru/project/cart/internal/model"
)

// Cart Корзина
type Cart struct {
	*sync.RWMutex
	UserID model.UserID
	items  map[model.SKU]model.CartItem
}

// GetItems Получение товаров из корзины
func (c *Cart) GetItems() ([]model.CartItem, error) {
	c.RLock()
	defer c.RUnlock()
	cartItems := make([]model.CartItem, 0, len(c.items))
	for _, item := range c.items {
		cartItems = append(cartItems, item)
	}

	return cartItems, nil
}

// Clear Очитка корзины
func (c *Cart) Clear() error {
	c.Lock()
	defer c.Unlock()
	clear(c.items)
	return nil
}

// AddItemBySKU Добавление товара
func (c *Cart) AddItemBySKU(sku model.SKU, count uint16) error {
	c.Lock()
	defer c.Unlock()
	c.items[sku] = model.CartItem{
		SKU:   sku,
		Count: c.items[sku].Count + count,
	}
	return nil
}

// DeleteItemBySKU Удаление товара
func (c *Cart) DeleteItemBySKU(sku model.SKU) error {
	c.Lock()
	defer c.Unlock()
	delete(c.items, sku)
	return nil
}

// NewCart Конструктор
func NewCart(userID model.UserID) *Cart {
	return &Cart{
		UserID:  userID,
		RWMutex: &sync.RWMutex{},
		items:   make(map[model.SKU]model.CartItem),
	}
}
