-- +goose Up
-- +goose StatementBegin
INSERT INTO "stock_item" ("sku", "total_count", "reserved")
VALUES (773297411, 150, 10),
       (1002, 200, 20),
       (1003, 250, 30),
       (1004, 300, 40),
       (1005, 350, 50);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE "stock_item";
-- +goose StatementEnd
