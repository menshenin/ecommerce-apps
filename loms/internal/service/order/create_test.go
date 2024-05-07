package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"route256.ozon.ru/project/loms/internal/event"
	"route256.ozon.ru/project/loms/internal/model"
	"route256.ozon.ru/project/loms/internal/service/order/mocks"
)

func TestService_Create(t *testing.T) {
	reserveCalls := false
	testItems := []model.OrderItem{
		{
			SKU:   model.SKU(1),
			Count: 10,
		},
		{
			SKU:   model.SKU(2),
			Count: 20,
		},
	}
	stocksRepo := mocks.NewStocksRepositoryMock(t)
	stocksRepo.GetCurrentStockMock.Expect().Return(&model.Stock{
		Reserve: func(_ context.Context, items []model.OrderItem) error {
			reserveCalls = true
			assert.Equal(t, testItems, items)
			return nil
		},
	}, nil)

	ordersRepo := mocks.NewOrdersRepositoryMock(t)
	ordersRepo.CreateMock.Expect(context.Background(), model.UserID(123), testItems).Return(&model.Order{
		ID: model.OrderID(1),
		UpdateStatus: func(_ context.Context, status model.OrderStatus) error {
			assert.Equal(t, status, model.OrderStatusAwaitingPayment)
			return nil
		},
	}, nil)

	callNum := 0
	eventProducer := mocks.NewEventProducerMock(t).ProduceMock.Inspect(func(_ context.Context, orderEvent event.OrderEvent) {
		callNum++
		switch callNum {
		case 1:
			assert.Equal(t, orderEvent.EventType, event.OrderCreated)
		case 2:
			assert.Equal(t, orderEvent.EventType, event.OrderAwaitingPayment)
		}
	}).Return()

	serv := New(ordersRepo, stocksRepo, eventProducer)
	resp, err := serv.Create(context.Background(), CreateRequest{
		UserID: model.UserID(123),
		Items: []ItemCreateRequest{
			{
				Sku:   model.SKU(1),
				Count: 10,
			},
			{
				Sku:   model.SKU(2),
				Count: 20,
			},
		},
	})

	assert.Nil(t, err)
	assert.Equal(t, model.OrderID(1), resp.OrderID)
	assert.True(t, reserveCalls)
}
