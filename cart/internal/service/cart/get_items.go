package cart

import (
	"context"
	"sort"

	"route256.ozon.ru/project/cart/internal/model"
)

// GetItemsRequest Запрос на получение товаров
type GetItemsRequest struct {
	UserID int64 `validate:"required"`
}

// GetItemsResponse Ответ
type GetItemsResponse struct {
	Items []struct {
		model.Item
		model.CartItem
	}
	TotalPrice uint32
}

// GetItems Получение всех товаров, находящихся в корзине
func (c *Service) GetItems(ctx context.Context, request GetItemsRequest) (resp GetItemsResponse, err error) {
	err = c.validate.Struct(&request)
	if err != nil {
		return
	}

	cart, err := c.cartRepo.GetByUserID(ctx, model.UserID(request.UserID))
	if err != nil {
		return
	}
	items, err := cart.GetItems()
	if err != nil {
		return
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].SKU < items[j].SKU
	})
	resp.Items = make([]struct {
		model.Item
		model.CartItem
	}, len(items))
	skus := make([]model.SKU, len(items))
	for i := range items {
		skus[i] = items[i].SKU
	}
	itemsMap, err := c.itemRepo.GetItemsBySKU(ctx, skus...)
	if err != nil {
		return
	}
	for i := range items {
		item := *itemsMap[items[i].SKU]
		resp.Items[i].CartItem = items[i]
		resp.Items[i].Item = item
		resp.TotalPrice += item.Price * uint32(resp.Items[i].Count)
	}

	return
}
