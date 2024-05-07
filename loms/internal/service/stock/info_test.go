package stock

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"route256.ozon.ru/project/loms/internal/model"
	"route256.ozon.ru/project/loms/internal/service/stock/mocks"
)

func TestService_Info(t *testing.T) {
	stocksRepo := mocks.NewStocksRepositoryMock(t)
	stocksRepo.GetCurrentStockMock.Expect().Return(&model.Stock{
		AvailableCount: func(_ context.Context, sku model.SKU) (int32, error) {
			assert.Equal(t, sku, model.SKU(123))
			return 5, nil
		},
	}, nil)
	service := New(stocksRepo)
	resp, err := service.Info(context.Background(), InfoRequest{Sku: model.SKU(123)})
	assert.NoError(t, err)
	assert.Equal(t, resp, &InfoResponse{Count: 5})
}
