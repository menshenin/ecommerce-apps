package dborder

import (
	"context"

	"route256.ozon.ru/project/loms/internal/model"
	"route256.ozon.ru/project/loms/internal/repository/dborder/db"
)

var statusMap = map[model.OrderStatus]db.OrderStatus{
	model.OrderStatusNew:             db.OrderStatusNew,
	model.OrderStatusPayed:           db.OrderStatusPayed,
	model.OrderStatusCancelled:       db.OrderStatusCancelled,
	model.OrderStatusAwaitingPayment: db.OrderStatusAwaitingPayment,
	model.OrderStatusFailed:          db.OrderStatusFailed,
}

func convertToDomainStatus(status db.OrderStatus) model.OrderStatus {
	for domainStatus, dbStatus := range statusMap {
		if status == dbStatus {
			return domainStatus
		}
	}

	return 0
}

func convertToDBStatus(status model.OrderStatus) db.OrderStatus {
	if s, ok := statusMap[status]; ok {
		return s
	}

	return ""
}

func (r *Repository) convertToDomainOrder(order db.Order) *model.Order {
	return &model.Order{
		ID:     model.OrderID(order.ID),
		UserID: model.UserID(order.UserID),
		UpdateStatus: func(ctx context.Context, status model.OrderStatus) error {
			order.Status = convertToDBStatus(status)
			return db.New(r.master).UpdateStatus(ctx, db.UpdateStatusParams{
				ID:     order.ID,
				Status: order.Status,
			})
		},
		GetItems: func(ctx context.Context) ([]model.OrderItem, error) {
			items, err := db.New(r.slave).GetOrderItems(ctx, order.ID)
			if err != nil {
				return nil, err
			}
			models := make([]model.OrderItem, len(items))
			for i, item := range items {
				models[i] = model.OrderItem{
					SKU:   model.SKU(item.Sku),
					Count: uint16(item.Count),
				}
			}

			return models, nil
		},
		GetStatus: func(_ context.Context) (model.OrderStatus, error) {
			return convertToDomainStatus(order.Status), nil
		},
	}
}
