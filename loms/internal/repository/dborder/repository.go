// Package dborder Репозиторий заказов, работающий с Postgres
package dborder

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"route256.ozon.ru/project/loms/internal/model"
	"route256.ozon.ru/project/loms/internal/repository/dborder/db"
)

// Repository Репозиторий
type Repository struct {
	master *pgx.Conn
	slave  *pgx.Conn
}

// GetByID Получение заказа по ID
func (r *Repository) GetByID(ctx context.Context, id model.OrderID) (*model.Order, error) {
	order, err := db.New(r.slave).GerByID(ctx, int64(id))
	if err != nil {
		return nil, err
	}

	return r.convertToDomainOrder(order), nil
}

// Create Создание заказа
func (r *Repository) Create(ctx context.Context, userID model.UserID, items []model.OrderItem) (order *model.Order, err error) {
	err = r.inTx(ctx, func(queries *db.Queries) error {
		dbOrder, err := queries.CreateOrder(ctx, int64(userID))
		oiParams := make([]db.CreateOrderItemsParams, len(items))
		for i, item := range items {
			oiParams[i] = db.CreateOrderItemsParams{
				OrderID: dbOrder.ID,
				Sku:     int64(item.SKU),
				Count:   int32(item.Count),
			}
		}
		results := queries.CreateOrderItems(ctx, oiParams)
		results.Exec(func(_ int, resErr error) {
			if resErr != nil {
				err = resErr
			}
		})
		if err != nil {
			return err
		}
		order = r.convertToDomainOrder(dbOrder)
		order.GetItems = func(_ context.Context) ([]model.OrderItem, error) {
			return items, nil
		}
		return nil
	})

	return
}

func (r *Repository) inTx(ctx context.Context, f func(queries *db.Queries) error) (err error) {
	tx, err := r.master.Begin(ctx)
	if err != nil {
		return
	}
	defer func() {
		rollBackErr := tx.Rollback(ctx)
		if !errors.Is(rollBackErr, pgx.ErrTxClosed) {
			if err != nil {
				err = fmt.Errorf("tx rollback: %w prev: %w", tx.Rollback(ctx), err)
			} else {
				err = rollBackErr
			}
		}
	}()
	err = f(db.New(tx))
	if err != nil {
		return
	}

	return tx.Commit(ctx)
}

// New Конструктор
func New(master *pgx.Conn, slave *pgx.Conn) *Repository {
	return &Repository{
		master: master,
		slave:  slave,
	}
}
