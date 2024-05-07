package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"route256.ozon.ru/project/loms/internal/event"
	"route256.ozon.ru/project/loms/internal/model"
	"route256.ozon.ru/project/loms/internal/service/order/mocks"
)

func TestService_Cancel(t *testing.T) {
	ordersRepo := mocks.NewOrdersRepositoryMock(t)
	updateStatusCalled := false
	ordersRepo.GetByIDMock.Expect(context.Background(), model.OrderID(123)).
		Return(&model.Order{
			ID:     0,
			UserID: 0,
			UpdateStatus: func(_ context.Context, status model.OrderStatus) error {
				updateStatusCalled = true
				assert.Equal(t, model.OrderStatusCancelled, status)
				return nil
			},
			GetItems: func(_ context.Context) ([]model.OrderItem, error) {
				return []model.OrderItem{
					{
						SKU:   123,
						Count: 10,
					},
				}, nil
			},
			GetStatus: func(_ context.Context) (model.OrderStatus, error) {
				return model.OrderStatusAwaitingPayment, nil
			},
		}, nil)

	cancelReserveCalled := false
	stocksRepo := mocks.NewStocksRepositoryMock(t)
	stocksRepo.GetCurrentStockMock.Expect().Return(&model.Stock{
		CancelReserve: func(_ context.Context, items []model.OrderItem) error {
			cancelReserveCalled = true
			assert.Equal(t, []model.OrderItem{
				{
					SKU:   123,
					Count: 10,
				},
			}, items)
			return nil
		},
	}, nil)
	eventProducer := mocks.NewEventProducerMock(t).ProduceMock.Inspect(func(_ context.Context, orderEvent event.OrderEvent) {
		assert.Equal(t, event.OrderCancelled, orderEvent.EventType)
	}).Return()
	service := New(ordersRepo, stocksRepo, eventProducer)
	err := service.Cancel(context.Background(), CancelRequest{OrderID: model.OrderID(123)})
	assert.NoError(t, err)
	assert.True(t, updateStatusCalled)
	assert.True(t, cancelReserveCalled)
}
