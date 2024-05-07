-- +goose Up
-- +goose StatementBegin
DROP TYPE IF EXISTS order_status;
CREATE TYPE order_status AS ENUM ('new', 'failed', 'cancelled', 'payed', 'awaiting_payment');
CREATE TABLE IF NOT EXISTS "order"
(
    id      BIGSERIAL PRIMARY KEY,
    user_id BIGINT       NOT NULL,
    status  order_status NOT NULL DEFAULT 'new'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "order";
DROP TYPE IF EXISTS order_status;
-- +goose StatementEnd
