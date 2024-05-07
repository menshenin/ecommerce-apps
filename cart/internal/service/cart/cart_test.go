package cart

import (
	"context"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"route256.ozon.ru/project/cart/internal/model"
	"route256.ozon.ru/project/cart/internal/service/cart/mocks"
)

func TestCart_AddItem(t *testing.T) {
	t.Parallel()
	t.Run("add item to existed cart", func(t *testing.T) {
		cart := &model.Cart{
			UserID: model.UserID(555),
			AddItemBySKU: func(sku model.SKU, count uint16) error {
				assert.EqualValues(t, sku, 123)
				assert.EqualValues(t, count, 3)
				return nil
			},
		}
		lomsClient := mocks.NewLomsClientMock(t).AvailableCountMock.
			Expect(context.Background(), model.SKU(123)).
			Return(5, nil)
		cartRepo := mocks.NewCartsRepositoryMock(t).
			GetByUserIDMock.
			Expect(context.Background(), model.UserID(555)).
			Return(cart, nil)

		itemRepo := mocks.NewItemRepositoryMock(t).GetItemsBySKUMock.
			Expect(context.Background(), model.SKU(123)).
			Return(map[model.SKU]*model.Item{
				123: {
					SKU:   123,
					Name:  "test",
					Price: 10,
				}}, nil)

		c, err := New(cartRepo, itemRepo, lomsClient)
		assert.NoError(t, err)
		_ = c.AddItem(context.Background(), AddItemRequest{
			UserID: 555,
			SkuID:  123,
			Count:  3,
		})
	})

	t.Run("add item to not existed cart", func(t *testing.T) {
		cart := &model.Cart{
			UserID: model.UserID(555),
			AddItemBySKU: func(sku model.SKU, count uint16) error {
				assert.EqualValues(t, sku, 123)
				assert.EqualValues(t, count, 3)
				return nil
			},
		}
		cartRepo := mocks.NewCartsRepositoryMock(t).
			GetByUserIDMock.
			Expect(context.Background(), model.UserID(555)).
			Return(nil, model.ErrCartNotFound).
			CreateMock.
			Expect(context.Background(), model.UserID(555)).
			Return(cart, nil)

		itemRepo := mocks.NewItemRepositoryMock(t).GetItemsBySKUMock.
			Expect(context.Background(), model.SKU(123)).
			Return(map[model.SKU]*model.Item{
				123: {
					SKU:   123,
					Name:  "test",
					Price: 10,
				}}, nil)
		lomsClient := mocks.NewLomsClientMock(t).AvailableCountMock.
			Expect(context.Background(), model.SKU(123)).
			Return(10, nil)

		c, err := New(cartRepo, itemRepo, lomsClient)
		assert.NoError(t, err)
		assert.NoError(t, c.AddItem(context.Background(), AddItemRequest{
			UserID: 555,
			SkuID:  123,
			Count:  3,
		}))
	})

	t.Run("add item, not existed item", func(t *testing.T) {
		itemRepo := mocks.NewItemRepositoryMock(t).GetItemsBySKUMock.
			Expect(context.Background(), model.SKU(123)).
			Return(nil, nil)

		c, err := New(nil, itemRepo, nil)
		assert.NoError(t, err)
		err = c.AddItem(context.Background(), AddItemRequest{
			UserID: 555,
			SkuID:  123,
			Count:  3,
		})
		var vadlidationErr validator.ValidationErrors
		assert.ErrorAs(t, err, &vadlidationErr)
	})

	t.Run("add item, insufficient stocks", func(t *testing.T) {
		lomsClient := mocks.NewLomsClientMock(t).AvailableCountMock.
			Expect(context.Background(), model.SKU(123)).
			Return(5, nil)

		itemRepo := mocks.NewItemRepositoryMock(t).GetItemsBySKUMock.
			Expect(context.Background(), model.SKU(123)).
			Return(map[model.SKU]*model.Item{
				123: {
					SKU:   123,
					Name:  "test",
					Price: 10,
				}}, nil)

		c, err := New(nil, itemRepo, lomsClient)
		assert.NoError(t, err)
		err = c.AddItem(context.Background(), AddItemRequest{
			UserID: 555,
			SkuID:  123,
			Count:  10,
		})
		assert.ErrorIs(t, err, model.ErrInsufficientStocks)
	})
}

func TestCart_GetItems(t *testing.T) {
	cart := &model.Cart{
		UserID: model.UserID(555),
		GetItems: func() ([]model.CartItem, error) {
			return []model.CartItem{
				{
					SKU:   200,
					Count: 3,
				},
				{
					SKU:   100,
					Count: 1,
				},
			}, nil
		},
	}
	cartRepo := mocks.NewCartsRepositoryMock(t).
		GetByUserIDMock.
		Expect(context.Background(), model.UserID(555)).
		Return(cart, nil)

	itemRepo := mocks.NewItemRepositoryMock(t).
		GetItemsBySKUMock.
		Expect(context.Background(), model.SKU(100), model.SKU(200)).
		Return(map[model.SKU]*model.Item{
			200: {
				SKU:   200,
				Name:  "test 1",
				Price: 10,
			},
			100: {
				SKU:   100,
				Name:  "test 2",
				Price: 3,
			}}, nil)

	c, err := New(cartRepo, itemRepo, nil)
	assert.NoError(t, err)
	resp, err := c.GetItems(context.Background(), GetItemsRequest{UserID: 555})
	assert.NoError(t, err)
	assert.Equal(t, GetItemsResponse{
		Items: []struct {
			model.Item
			model.CartItem
		}{
			{
				Item: model.Item{
					SKU:   100,
					Name:  "test 2",
					Price: 3,
				},
				CartItem: model.CartItem{
					SKU:   100,
					Count: 1,
				},
			},
			{
				Item: model.Item{
					SKU:   200,
					Name:  "test 1",
					Price: 10,
				},
				CartItem: model.CartItem{
					SKU:   200,
					Count: 3,
				},
			},
		},
		TotalPrice: 33,
	}, resp)
}

func TestCart_DeleteItem(t *testing.T) {
	cart := &model.Cart{
		UserID: model.UserID(555),
		DeleteItemBySKU: func(sku model.SKU) error {
			assert.Equal(t, sku, model.SKU(123))
			return nil
		},
	}
	cartRepo := mocks.NewCartsRepositoryMock(t).
		GetByUserIDMock.
		Expect(context.Background(), model.UserID(555)).
		Return(cart, nil)

	c, err := New(cartRepo, nil, nil)
	assert.NoError(t, err)
	assert.NoError(t, c.DeleteItem(context.Background(), DeleteItemRequest{
		UserID: 555,
		SkuID:  123,
	}))
}

func TestCart_DeleteAllItems(t *testing.T) {
	clearCalls := false
	cart := &model.Cart{
		UserID: model.UserID(555),
		Clear: func() error {
			clearCalls = true
			return nil
		},
	}
	cartRepo := mocks.NewCartsRepositoryMock(t).
		GetByUserIDMock.
		Expect(context.Background(), model.UserID(555)).
		Return(cart, nil)

	c, err := New(cartRepo, nil, nil)
	assert.NoError(t, err)
	assert.NoError(t, c.DeleteAllItems(context.Background(), DeleteAllItemsRequest{UserID: 555}))
	assert.True(t, clearCalls)
}

func TestService_Checkout(t *testing.T) {
	cart := &model.Cart{
		GetItems: func() ([]model.CartItem, error) {
			return []model.CartItem{
				{
					SKU:   12,
					Count: 13,
				},
			}, nil
		},
		Clear: func() error {
			return nil
		},
	}

	lomsClient := mocks.NewLomsClientMock(t).CheckoutMock.
		Expect(context.Background(), cart).
		Return(model.OrderID(1), nil)

	cartRepo := mocks.NewCartsRepositoryMock(t).
		GetByUserIDMock.
		Expect(context.Background(), model.UserID(555)).
		Return(cart, nil)

	c, err := New(cartRepo, nil, lomsClient)
	assert.NoError(t, err)
	resp, err := c.Checkout(context.Background(), CheckoutRequest{UserID: 555})
	assert.NoError(t, err)
	assert.Equal(t, &CheckoutResponse{OrderID: 1}, resp)
}
