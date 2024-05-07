// Package event События
package event

import (
	"time"

	"route256.ozon.ru/project/loms/internal/model"
)

// Type Тип события
type Type string

// ExtraData Дополнительные данные для события
type ExtraData map[string]any

// Типы событий заказа
const (
	// OrderCreated Заказ создан
	OrderCreated Type = "order-created"

	// OrderFailed Создание заказа завершилось с ошибкой
	OrderFailed Type = "order-failed"

	// OrderAwaitingPayment Заказ ждёт оплаты
	OrderAwaitingPayment Type = "order-awaiting-payment"

	// OrderPayed  Заказ оплачен
	OrderPayed Type = "order-payed"

	// OrderCancelled Заказ отменён
	OrderCancelled Type = "order-cancelled"
)

// OrderEvent Событие заказа
type OrderEvent struct {
	OrderID   model.OrderID `json:"order_id"`
	EventType Type          `json:"event_type"`
	Timestamp time.Time     `json:"timestamp"`
	Extra     ExtraData     `json:"extra"`
}

// CreateOrderEvent Создание события заказа
func CreateOrderEvent(orderID model.OrderID, eventType Type, extra ExtraData) OrderEvent {
	return OrderEvent{
		OrderID:   orderID,
		EventType: eventType,
		Timestamp: time.Now(),
		Extra:     extra,
	}
}
