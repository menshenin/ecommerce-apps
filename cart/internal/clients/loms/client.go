// Package loms Структуры для работы с сервисом LOMS
package loms

import (
	"context"

	"google.golang.org/grpc"
	"route256.ozon.ru/project/cart/internal/model"
	lomspb "route256.ozon.ru/project/cart/internal/pkg/pb/loms"
)

// GrpcClient Интерфейс grpc-клиента для работы с сервисом LOMS
type GrpcClient interface {
	CreateOrder(ctx context.Context, in *lomspb.CreateOrderRequest, opts ...grpc.CallOption) (*lomspb.CreateOrderResponse, error)
	StocksInfo(ctx context.Context, in *lomspb.StocksInfoRequest, opts ...grpc.CallOption) (*lomspb.StocksInfoResponse, error)
}

// Client Клиент для работы с сервисом LOMS
type Client struct {
	grpcClient GrpcClient
}

// AvailableCount Получение доступного количества товара
func (c *Client) AvailableCount(ctx context.Context, sku model.SKU) (int32, error) {
	resp, err := c.grpcClient.StocksInfo(ctx, &lomspb.StocksInfoRequest{Sku: int64(sku)})
	if err != nil {
		return 0, err
	}
	return resp.Count, nil
}

// Checkout Оформления заказа
func (c *Client) Checkout(ctx context.Context, cart *model.Cart) (orderID model.OrderID, err error) {
	cartItems, err := cart.GetItems()
	if err != nil {
		return
	}
	items := make([]*lomspb.Item, len(cartItems))
	for i, cartItem := range cartItems {
		items[i] = &lomspb.Item{
			Sku:   uint32(cartItem.SKU),
			Count: uint32(cartItem.Count),
		}
	}
	request := &lomspb.CreateOrderRequest{
		UserId: int64(cart.UserID),
		Items:  items,
	}

	resp, err := c.grpcClient.CreateOrder(ctx, request)
	if err != nil {
		return
	}

	return model.OrderID(resp.GetOrderId()), nil
}

// New Конструктор
func New(lomsGrpcClient GrpcClient) *Client {
	return &Client{
		grpcClient: lomsGrpcClient,
	}
}
