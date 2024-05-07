package cart

import (
	"context"

	"route256.ozon.ru/project/cart/internal/model"
)

// DeleteItemRequest Запрос на удаление товара
type DeleteItemRequest struct {
	UserID int64 `validate:"required"`
	SkuID  int64 `validate:"required"`
}

// DeleteItem удаление товара из корзины
func (c *Service) DeleteItem(ctx context.Context, request DeleteItemRequest) (err error) {
	err = c.validate.Struct(&request)
	if err != nil {
		return
	}
	cart, err := c.cartRepo.GetByUserID(ctx, model.UserID(request.UserID))
	if err != nil {
		return
	}

	return cart.DeleteItemBySKU(model.SKU(request.SkuID))
}
