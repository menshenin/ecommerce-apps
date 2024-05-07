// Package stockloader Загрузчик содержимого склада из файла
package stockloader

import (
	"context"
	_ "embed"
	"encoding/json"

	"route256.ozon.ru/project/loms/internal/model"
)

//go:embed stock-data.json
var data []byte

// ItemData Структура данных в файле
type ItemData struct {
	Sku        int64 `json:"sku"`
	TotalCount int32 `json:"total_count"`
	Reserved   int32 `json:"reserved"`
}

// LoadData Загрузка товаров на склад
func LoadData(ctx context.Context, stock *model.Stock) error {
	var items []ItemData
	err := json.Unmarshal(data, &items)
	if err != nil {
		return err
	}
	stockItems := make([]model.StockItem, len(items))
	for i, itemData := range items {
		stockItems[i] = model.StockItem{
			SKU:        model.SKU(itemData.Sku),
			TotalCount: itemData.TotalCount,
			Reserved:   itemData.Reserved,
		}
	}
	return stock.Load(ctx, stockItems)
}
