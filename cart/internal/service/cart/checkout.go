package cart

import (
	"context"

	"route256.ozon.ru/project/cart/internal/model"
)

// CheckoutRequest Запрос на формление заказа
type CheckoutRequest struct {
	UserID int64 `validate:"required"`
}

// CheckoutResponse Ответ
type CheckoutResponse struct {
	OrderID model.OrderID
}

// Checkout Оформление заказа
func (c *Service) Checkout(ctx context.Context, request CheckoutRequest) (resp *CheckoutResponse, err error) {
	err = c.validate.Struct(&request)
	if err != nil {
		return
	}
	cart, err := c.cartRepo.GetByUserID(ctx, model.UserID(request.UserID))
	if err != nil {
		return
	}
	orderID, err := c.lomsClient.Checkout(ctx, cart)
	if err != nil {
		return
	}
	err = cart.Clear()
	if err != nil {
		return
	}
	return &CheckoutResponse{OrderID: orderID}, nil
}
