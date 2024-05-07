-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_item
(
    id       BIGSERIAL PRIMARY KEY,
    order_id BIGINT REFERENCES "order" (id) NOT NULL,
    sku      BIGINT                         NOT NULL,
    count    INT                            NOT NULL,
    CONSTRAINT count_non_zero CHECK ( order_item.count > 0 )
);
CREATE UNIQUE INDEX IF NOT EXISTS u_idx_order_item_order_id_sku ON order_item (order_id, sku);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_item;
-- +goose StatementEnd
