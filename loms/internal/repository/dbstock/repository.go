// Package dbstock Репозиторий складов, работающий с Postgres
package dbstock

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"route256.ozon.ru/project/loms/internal/model"
	"route256.ozon.ru/project/loms/internal/repository/dbstock/db"
)

// Repository Репозиторий
type Repository struct {
	master *pgx.Conn
	slave  *pgx.Conn
}

// GetCurrentStock Получение текущего склада
func (r *Repository) GetCurrentStock() (*model.Stock, error) {
	st := stock{repo: r}

	return &model.Stock{
		Reserve:        st.reserve,
		AvailableCount: st.availableCount,
		CancelReserve:  st.cancelReserve,
		WriteOff:       st.writeOff,
		Load:           st.load,
	}, nil
}

// New Конструктор
func New(master *pgx.Conn, slave *pgx.Conn) *Repository {
	return &Repository{
		master: master,
		slave:  slave,
	}
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
