// Package model Модели
package model

import (
	"context"
	"errors"
)

// UserID Идентификатор пользователя
type UserID int64

// OrderID Идентификатор заказа
type OrderID int64

// SKU Идентификатор товара
type SKU int64

// OrderStatus Статус заказа
type OrderStatus int

var (
	// ErrOrderNotFound Заказ не найден
	ErrOrderNotFound = errors.New("order not found")

	// ErrSkuNotFound Товар не найден
	ErrSkuNotFound = errors.New("sku not found")

	// ErrReserveMoreThenTotalCount Попытка зарезервировать больше, чем есть на складе
	ErrReserveMoreThenTotalCount = errors.New("reserving more then total count")

	// ErrCancelMoreThenReserved Попытка отменить резерв большего количества товаров, чем зарезервировано
	ErrCancelMoreThenReserved = errors.New("cancel more then reserved")

	// ErrWriteOffMoreThenReserved Попытка списать больше товара, чем зарезервировано
	ErrWriteOffMoreThenReserved = errors.New("write off more then reserved")
)

const (
	// OrderStatusNew Новый заказ
	OrderStatusNew OrderStatus = iota + 1

	// OrderStatusAwaitingPayment Заказ ожидает оплаты
	OrderStatusAwaitingPayment

	// OrderStatusFailed Заказ не смог завершиться
	OrderStatusFailed

	// OrderStatusPayed Заказ оплачен
	OrderStatusPayed

	// OrderStatusCancelled Заказ отменён
	OrderStatusCancelled
)

// String Строковое представление статуса заказа
func (status OrderStatus) String() string {
	switch status {
	case OrderStatusCancelled:
		return "CANCELLED"
	case OrderStatusAwaitingPayment:
		return "AWAITING_PAYMENT"
	case OrderStatusFailed:
		return "FAILED"
	case OrderStatusPayed:
		return "PAYED"
	case OrderStatusNew:
		return "NEW"
	default:
		return "UNKNOWN"
	}
}

// StockItem Товар на складе
type StockItem struct {
	SKU        SKU
	TotalCount int32
	Reserved   int32
}

// Stock Склад
type Stock struct {
	Reserve        func(ctx context.Context, items []OrderItem) error
	AvailableCount func(ctx context.Context, sku SKU) (int32, error)
	CancelReserve  func(ctx context.Context, items []OrderItem) error
	WriteOff       func(ctx context.Context, items []OrderItem) error
	Load           func(ctx context.Context, items []StockItem) error
}

// OrderItem Товар в заказе
type OrderItem struct {
	SKU   SKU
	Count uint16
}

// Order Заказ
type Order struct {
	ID           OrderID
	UserID       UserID
	UpdateStatus func(ctx context.Context, status OrderStatus) error
	GetItems     func(ctx context.Context) ([]OrderItem, error)
	GetStatus    func(ctx context.Context) (OrderStatus, error)
}
