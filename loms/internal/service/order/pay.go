package order

import (
	"context"
	"fmt"

	"route256.ozon.ru/project/loms/internal/event"
	"route256.ozon.ru/project/loms/internal/model"
)

// PayRequest Запрос на оплату заказа
type PayRequest struct {
	OrderID model.OrderID
}

// Pay Оплата заказа
func (s *Service) Pay(ctx context.Context, request PayRequest) (err error) {
	order, err := s.ordersRepo.GetByID(ctx, request.OrderID)
	if err != nil {
		return
	}
	status, err := order.GetStatus(ctx)
	if err != nil {
		return
	}
	if status != model.OrderStatusAwaitingPayment {
		return fmt.Errorf("order has wrong status: %s", status.String())
	}

	stock, err := s.stocksRepo.GetCurrentStock()
	if err != nil {
		return
	}
	items, err := order.GetItems(ctx)
	if err != nil {
		return
	}
	err = stock.WriteOff(ctx, items)
	if err != nil {
		return
	}

	err = order.UpdateStatus(ctx, model.OrderStatusPayed)
	if err != nil {
		return err
	}
	s.eventProducer.Produce(ctx, event.CreateOrderEvent(order.ID, event.OrderPayed, nil))

	return nil
}
