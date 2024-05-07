// Package loms Рализация grpc сервера LOMS
package loms

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"route256.ozon.ru/project/loms/internal/model"
	lomspb "route256.ozon.ru/project/loms/internal/pkg/pb/loms"
	"route256.ozon.ru/project/loms/internal/service/order"
	"route256.ozon.ru/project/loms/internal/service/stock"
)

// Server Сервера
type Server struct {
	lomspb.UnimplementedLomsServer
	orderService  *order.Service
	stocksService *stock.Service
}

// New Конструктор
func New(orderService *order.Service, stocksService *stock.Service) *Server {
	return &Server{
		orderService:  orderService,
		stocksService: stocksService,
	}
}

// CreateOrder Создание нового заказа
func (s *Server) CreateOrder(ctx context.Context, request *lomspb.CreateOrderRequest) (*lomspb.CreateOrderResponse, error) {
	items := make([]order.ItemCreateRequest, len(request.GetItems()))
	for i, item := range request.GetItems() {
		items[i] = order.ItemCreateRequest{
			Sku:   model.SKU(item.GetSku()),
			Count: uint16(item.GetCount()),
		}
	}
	resp, err := s.orderService.Create(ctx, order.CreateRequest{
		UserID: model.UserID(request.GetUserId()),
		Items:  items,
	})
	if err != nil {
		return nil, err
	}

	return &lomspb.CreateOrderResponse{
		OrderId: int64(resp.OrderID),
	}, nil
}

// OrderInfo Получение информации о заказе
func (s *Server) OrderInfo(ctx context.Context, request *lomspb.OrderInfoRequest) (*lomspb.OrderInfoResponse, error) {
	info, err := s.orderService.Info(ctx, order.InfoRequest{
		OrderID: model.OrderID(request.GetOrderId()),
	})
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return nil, status.Error(codes.NotFound, model.ErrOrderNotFound.Error())
		}
		return nil, err
	}
	orderItems := make([]*lomspb.Item, len(info.Items))
	for i, item := range info.Items {
		orderItems[i] = &lomspb.Item{
			Sku:   uint32(item.SKU),
			Count: uint32(item.Count),
		}
	}

	return &lomspb.OrderInfoResponse{
		UserId: int64(info.UserID),
		Status: convertOrderStatus(info.Status),
		Items:  orderItems,
	}, nil
}

// OrderPay Оплата заказа
func (s *Server) OrderPay(ctx context.Context, request *lomspb.OrderPayRequest) (*emptypb.Empty, error) {
	err := s.orderService.Pay(ctx, order.PayRequest{
		OrderID: model.OrderID(request.GetOrderId()),
	})
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return nil, status.Error(codes.NotFound, model.ErrOrderNotFound.Error())
		}
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// OrderCancel Отмена заказа
func (s *Server) OrderCancel(ctx context.Context, request *lomspb.OrderCancelRequest) (*emptypb.Empty, error) {
	err := s.orderService.Cancel(ctx, order.CancelRequest{
		OrderID: model.OrderID(request.GetOrderId()),
	})

	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return nil, status.Error(codes.NotFound, model.ErrOrderNotFound.Error())
		}
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// StocksInfo Получение информации об остатках на складе
func (s *Server) StocksInfo(ctx context.Context, request *lomspb.StocksInfoRequest) (*lomspb.StocksInfoResponse, error) {
	info, err := s.stocksService.Info(ctx, stock.InfoRequest{Sku: model.SKU(request.GetSku())})
	if err != nil {
		return nil, err
	}

	return &lomspb.StocksInfoResponse{Count: info.Count}, nil
}

func convertOrderStatus(status model.OrderStatus) lomspb.OrderStatus {
	switch status {
	case model.OrderStatusCancelled:
		return lomspb.OrderStatus_ORDER_STATUS_CANCELLED
	case model.OrderStatusAwaitingPayment:
		return lomspb.OrderStatus_ORDER_STATUS_AWAITING_PAYMENT
	case model.OrderStatusFailed:
		return lomspb.OrderStatus_ORDER_STATUS_FAILED
	case model.OrderStatusPayed:
		return lomspb.OrderStatus_ORDER_STATUS_PAYED
	case model.OrderStatusNew:
		return lomspb.OrderStatus_ORDER_STATUS_NEW
	default:
		return lomspb.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
}
