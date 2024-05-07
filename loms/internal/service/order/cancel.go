package order

import (
	"context"
	"fmt"

	"route256.ozon.ru/project/loms/internal/event"
	"route256.ozon.ru/project/loms/internal/model"
)

// CancelRequest Запрос на отмену заказа
type CancelRequest struct {
	OrderID model.OrderID
}

// Cancel Отмена заказа
func (s *Service) Cancel(ctx context.Context, request CancelRequest) (err error) {
	order, err := s.ordersRepo.GetByID(ctx, request.OrderID)
	if err != nil {
		return
	}
	status, err := order.GetStatus(ctx)
	if err != nil {
		return
	}
	if status == model.OrderStatusPayed {
		return fmt.Errorf("order already payed")
	}
	stock, err := s.stocksRepo.GetCurrentStock()
	if err != nil {
		return
	}
	items, err := order.GetItems(ctx)
	if err != nil {
		return
	}
	err = stock.CancelReserve(ctx, items)
	if err != nil {
		return
	}

	err = order.UpdateStatus(ctx, model.OrderStatusCancelled)
	if err != nil {
		return err
	}
	s.eventProducer.Produce(ctx, event.CreateOrderEvent(order.ID, event.OrderCancelled, nil))

	return nil
}
