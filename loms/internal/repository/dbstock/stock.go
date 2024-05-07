package dbstock

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"route256.ozon.ru/project/loms/internal/model"
	"route256.ozon.ru/project/loms/internal/pkg/helpers"
	"route256.ozon.ru/project/loms/internal/repository/dbstock/db"
)

type stock struct {
	repo *Repository
}

func (s stock) updateStockItems(
	ctx context.Context,
	orderItems []model.OrderItem,
	updateFunc func(model.OrderItem, db.StockItem) (db.StockItem, error),
) (err error) {
	return s.repo.inTx(ctx, func(queries *db.Queries) error {
		stockItems, err := queries.GetStockItems(ctx, helpers.Map(orderItems, func(o model.OrderItem) int64 {
			return int64(o.SKU)
		}))
		if err != nil {
			return err
		}
		if len(stockItems) != len(orderItems) {
			return model.ErrSkuNotFound
		}
		for i, stockItem := range stockItems {
			for _, item := range orderItems {
				if item.SKU == model.SKU(stockItem.Sku) {
					stockItems[i], err = updateFunc(item, stockItem)
					if err != nil {
						return err
					}
					break
				}
			}
		}
		res := queries.UpdateStockItem(ctx, helpers.Map(stockItems, func(si db.StockItem) db.UpdateStockItemParams {
			return db.UpdateStockItemParams{
				Sku:        si.Sku,
				Reserved:   si.Reserved,
				TotalCount: si.TotalCount,
			}
		}))
		return res.Close()
	})
}

func (s stock) reserve(ctx context.Context, items []model.OrderItem) (err error) {
	return s.updateStockItems(ctx, items, func(orderItem model.OrderItem, stockItem db.StockItem) (si db.StockItem, err error) {
		stockItem.Reserved += int32(orderItem.Count)
		if stockItem.Reserved > stockItem.TotalCount {
			return si, model.ErrReserveMoreThenTotalCount
		}
		return stockItem, nil
	})
}

func (s stock) availableCount(ctx context.Context, sku model.SKU) (int32, error) {
	stockItem, err := db.New(s.repo.slave).GetStockItemBySKU(ctx, int64(sku))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, model.ErrSkuNotFound
		}
		return 0, err
	}
	return stockItem.TotalCount - stockItem.Reserved, nil
}

func (s stock) cancelReserve(ctx context.Context, items []model.OrderItem) (err error) {
	return s.updateStockItems(ctx, items, func(orderItem model.OrderItem, stockItem db.StockItem) (si db.StockItem, err error) {
		stockItem.Reserved -= int32(orderItem.Count)
		if stockItem.Reserved < 0 {
			return si, model.ErrCancelMoreThenReserved
		}

		return stockItem, nil
	})
}

func (s stock) writeOff(ctx context.Context, items []model.OrderItem) (err error) {
	return s.updateStockItems(ctx, items, func(orderItem model.OrderItem, stockItem db.StockItem) (si db.StockItem, err error) {
		stockItem.Reserved -= int32(orderItem.Count)
		stockItem.TotalCount -= int32(orderItem.Count)
		if stockItem.Reserved < 0 {
			return si, model.ErrWriteOffMoreThenReserved
		}
		return stockItem, nil
	})
}

func (s stock) load(ctx context.Context, items []model.StockItem) (err error) {
	res := db.New(s.repo.master).Load(ctx, helpers.Map(items, func(t model.StockItem) db.LoadParams {
		return db.LoadParams{
			Sku:        int64(t.SKU),
			TotalCount: t.TotalCount,
			Reserved:   t.Reserved,
		}
	}))

	return res.Close()
}
