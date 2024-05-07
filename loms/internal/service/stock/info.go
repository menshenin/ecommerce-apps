package stock

import (
	"context"

	"route256.ozon.ru/project/loms/internal/model"
)

// InfoRequest Запрос на получение данных о товаре
type InfoRequest struct {
	Sku model.SKU ``
}

// InfoResponse Ответ
type InfoResponse struct {
	Count int32
}

// Info Информация о товаре на складе
func (s *Service) Info(ctx context.Context, request InfoRequest) (resp *InfoResponse, err error) {
	err = s.validate.Struct(&request)
	if err != nil {
		return
	}

	stock, err := s.stocksRepo.GetCurrentStock()
	if err != nil {
		return
	}
	count, err := stock.AvailableCount(ctx, request.Sku)
	if err != nil {
		return
	}

	return &InfoResponse{
		Count: count,
	}, nil
}
