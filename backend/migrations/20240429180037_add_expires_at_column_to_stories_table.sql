-- +goose Up
-- +goose StatementBegin
ALTER TABLE stories ADD COLUMN expires_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE stories DROP COLUMN expires_at;
-- +goose StatementEnd
