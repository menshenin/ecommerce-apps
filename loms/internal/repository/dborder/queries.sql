-- name: GerByID :one
SELECT *
FROM "order"
WHERE id = $1
LIMIT 1;

-- name: GetOrderItems :many
SELECT *
FROM order_item
WHERE order_id = $1;

-- name: UpdateStatus :exec
UPDATE "order"
SET status = $2
WHERE id = $1;

-- name: CreateOrder :one
INSERT INTO "order" (user_id)
VALUES ($1)
RETURNING *;

-- name: CreateOrderItems :batchexec
INSERT INTO "order_item" (order_id, sku, count)
VALUES ($1, $2, $3);