-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE clicks ADD COLUMN browser VARCHAR;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE clicks DROP COLUMN browser;
-- +goose StatementEnd
