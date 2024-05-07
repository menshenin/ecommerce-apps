-- name: GetStockItems :many
SELECT *
FROM stock_item
WHERE sku = ANY (@skus::bigint[]);

-- name: GetStockItemBySKU :one
SELECT *
FROM stock_item
WHERE sku = $1;

-- name: UpdateStockItem :batchexec
UPDATE "stock_item"
SET reserved    = $2,
    total_count = $3
WHERE sku = $1;

-- name: Load :batchexec
INSERT INTO "stock_item" (sku, total_count, reserved)
VALUES ($1, $2, $3)
ON CONFLICT (sku) DO UPDATE SET total_count = EXCLUDED.total_count,
                                reserved    = EXCLUDED.reserved;
