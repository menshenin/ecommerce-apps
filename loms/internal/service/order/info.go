package order

import (
	"context"

	"route256.ozon.ru/project/loms/internal/model"
)

// InfoRequest Запрос информации о заказе
type InfoRequest struct {
	OrderID model.OrderID
}

// InfoResponse Ответ
type InfoResponse struct {
	UserID model.UserID
	Status model.OrderStatus
	Items  []model.OrderItem
}

// Info Получение информации о заказе
func (s *Service) Info(ctx context.Context, request InfoRequest) (resp *InfoResponse, err error) {
	order, err := s.ordersRepo.GetByID(ctx, request.OrderID)
	if err != nil {
		return
	}

	status, err := order.GetStatus(ctx)
	if err != nil {
		return
	}

	items, err := order.GetItems(ctx)
	if err != nil {
		return
	}
	return &InfoResponse{
		UserID: order.UserID,
		Status: status,
		Items:  items,
	}, nil
}
