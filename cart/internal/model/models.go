// Package model Модели
package model

import "errors"

// UserID Идентификатор пользователя
type UserID int64

// SKU Идентификатор товара
type SKU int64

// OrderID Идентификатор заказа
type OrderID int64

// ErrCartNotFound Корзина не найдена
var ErrCartNotFound = errors.New("cart not found")

// ErrInsufficientStocks Недостаточно товара на складе
var ErrInsufficientStocks = errors.New("insufficient stocks")

// Item Товар
type Item struct {
	SKU   SKU
	Name  string
	Price uint32
}

// CartItem Товар в корзине
type CartItem struct {
	SKU   SKU
	Count uint16
}

// Cart Корзина
type Cart struct {
	UserID          UserID
	GetItems        func() ([]CartItem, error)
	AddItemBySKU    func(sku SKU, count uint16) error
	DeleteItemBySKU func(sku SKU) error
	Clear           func() error
}
