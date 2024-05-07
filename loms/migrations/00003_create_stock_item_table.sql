-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS stock_item
(
    id       BIGSERIAL PRIMARY KEY,
    sku      BIGINT NOT NULL,
    total_count    INT    NOT NULL DEFAULT 0,
    reserved INT    NOT NULL DEFAULT 0,
    CONSTRAINT count_positive CHECK ( total_count >= 0 ),
    CONSTRAINT reserved_positive CHECK ( reserved >= 0 AND reserved <= total_count )
);
CREATE UNIQUE INDEX IF NOT EXISTS u_idx_stock_item_sku ON stock_item (sku);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stock_item;
-- +goose StatementEnd
