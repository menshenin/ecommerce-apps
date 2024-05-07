package order

import (
	"context"
	"fmt"

	"route256.ozon.ru/project/loms/internal/event"
	"route256.ozon.ru/project/loms/internal/model"
)

// ItemCreateRequest Запрос на создание товара в заказе
type ItemCreateRequest struct {
	Sku   model.SKU
	Count uint16
}

// CreateRequest Создание заказа
type CreateRequest struct {
	UserID model.UserID
	Items  []ItemCreateRequest
}

// CreateResponse Ответ
type CreateResponse struct {
	OrderID model.OrderID
}

// Create Создание заказа
func (s *Service) Create(ctx context.Context, request CreateRequest) (resp *CreateResponse, err error) {
	orderItems := make([]model.OrderItem, len(request.Items))
	for i, item := range request.Items {
		orderItems[i] = model.OrderItem{
			SKU:   item.Sku,
			Count: item.Count,
		}
	}
	order, err := s.ordersRepo.Create(ctx, request.UserID, orderItems)
	if err != nil {
		return
	}
	s.eventProducer.Produce(ctx, event.CreateOrderEvent(order.ID, event.OrderCreated, nil))
	stock, err := s.stocksRepo.GetCurrentStock()
	if err != nil {
		return
	}
	err = stock.Reserve(ctx, orderItems)
	if err != nil {
		updateStatusErr := order.UpdateStatus(ctx, model.OrderStatusFailed)
		if updateStatusErr != nil {
			return nil, fmt.Errorf("update order status error: %v, prev: %v", updateStatusErr, err)
		}
		s.eventProducer.Produce(ctx, event.CreateOrderEvent(order.ID, event.OrderFailed, nil))
	}
	err = order.UpdateStatus(ctx, model.OrderStatusAwaitingPayment)
	if err != nil {
		return
	}
	s.eventProducer.Produce(ctx, event.CreateOrderEvent(order.ID, event.OrderAwaitingPayment, nil))

	return &CreateResponse{
		OrderID: order.ID,
	}, nil
}
