// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: batch.go

package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

var (
	ErrBatchAlreadyClosed = errors.New("batch already closed")
)

const createOrderItems = `-- name: CreateOrderItems :batchexec
INSERT INTO "order_item" (order_id, sku, count)
VALUES ($1, $2, $3)
`

type CreateOrderItemsBatchResults struct {
	br     pgx.BatchResults
	tot    int
	closed bool
}

type CreateOrderItemsParams struct {
	OrderID int64
	Sku     int64
	Count   int32
}

func (q *Queries) CreateOrderItems(ctx context.Context, arg []CreateOrderItemsParams) *CreateOrderItemsBatchResults {
	batch := &pgx.Batch{}
	for _, a := range arg {
		vals := []interface{}{
			a.OrderID,
			a.Sku,
			a.Count,
		}
		batch.Queue(createOrderItems, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &CreateOrderItemsBatchResults{br, len(arg), false}
}

func (b *CreateOrderItemsBatchResults) Exec(f func(int, error)) {
	defer b.br.Close()
	for t := 0; t < b.tot; t++ {
		if b.closed {
			if f != nil {
				f(t, ErrBatchAlreadyClosed)
			}
			continue
		}
		_, err := b.br.Exec()
		if f != nil {
			f(t, err)
		}
	}
}

func (b *CreateOrderItemsBatchResults) Close() error {
	b.closed = true
	return b.br.Close()
}
