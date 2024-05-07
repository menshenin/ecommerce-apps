package cart

import (
	"context"
	"errors"

	"route256.ozon.ru/project/cart/internal/model"
)

// AddItemRequest Запрос на добавление товара
type AddItemRequest struct {
	UserID int64  `validate:"required"`
	SkuID  int64  `validate:"required,has_product"`
	Count  uint16 `validate:"required"`
}

// AddItem Добавление товара в корзину
func (c *Service) AddItem(ctx context.Context, request AddItemRequest) (err error) {
	err = c.validate.StructCtx(ctx, &request)
	if err != nil {
		return
	}
	count, err := c.lomsClient.AvailableCount(ctx, model.SKU(request.SkuID))
	if err != nil {
		return err
	}
	if count < int32(request.Count) {
		return model.ErrInsufficientStocks
	}
	cart, err := c.cartRepo.GetByUserID(ctx, model.UserID(request.UserID))
	if err != nil {
		if errors.Is(err, model.ErrCartNotFound) {
			cart, err = c.cartRepo.Create(ctx, model.UserID(request.UserID))
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return cart.AddItemBySKU(model.SKU(request.SkuID), request.Count)
}
