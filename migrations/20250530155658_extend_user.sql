-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS name       VARCHAR,
    ADD COLUMN IF NOT EXISTS avatar     VARCHAR; 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN IF EXISTS name,
    DROP COLUMN IF EXISTS avatar;
-- +goose StatementEnd
