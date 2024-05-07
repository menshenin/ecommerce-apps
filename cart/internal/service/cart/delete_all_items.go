package cart

import (
	"context"

	"route256.ozon.ru/project/cart/internal/model"
)

// DeleteAllItemsRequest Запрос на удаление всех товаров
type DeleteAllItemsRequest struct {
	UserID int64 `validate:"required"`
}

// DeleteAllItems Удаление всех товаров их корзины
func (c *Service) DeleteAllItems(ctx context.Context, request DeleteAllItemsRequest) (err error) {
	err = c.validate.Struct(&request)
	if err != nil {
		return
	}
	cart, err := c.cartRepo.GetByUserID(ctx, model.UserID(request.UserID))
	if err != nil {
		return
	}

	return cart.Clear()
}
