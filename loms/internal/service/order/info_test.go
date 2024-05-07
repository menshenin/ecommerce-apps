package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"route256.ozon.ru/project/loms/internal/model"
	"route256.ozon.ru/project/loms/internal/service/order/mocks"
)

func TestService_Info(t *testing.T) {
	ordersRepo := mocks.NewOrdersRepositoryMock(t)
	ordersRepo.GetByIDMock.
		Expect(context.Background(), model.OrderID(123)).
		Return(&model.Order{
			ID:     123,
			UserID: 11,
			GetItems: func(_ context.Context) ([]model.OrderItem, error) {
				return []model.OrderItem{
					{
						SKU:   1,
						Count: 10,
					},
				}, nil
			},
			GetStatus: func(_ context.Context) (model.OrderStatus, error) {
				return model.OrderStatusAwaitingPayment, nil
			},
		}, nil)

	service := New(ordersRepo, nil, nil)
	resp, err := service.Info(context.Background(), InfoRequest{
		OrderID: model.OrderID(123),
	})
	assert.NoError(t, err)
	assert.Equal(t, &InfoResponse{
		UserID: 11,
		Status: model.OrderStatusAwaitingPayment,
		Items: []model.OrderItem{
			{
				SKU:   1,
				Count: 10,
			},
		},
	}, resp)
}
