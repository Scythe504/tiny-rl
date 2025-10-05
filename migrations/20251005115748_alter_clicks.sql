-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE clicks 
ADD COLUMN country VARCHAR(100), 
ADD COLUMN country_iso_code VARCHAR(3);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE clicks 
DROP COLUMN country,
DROP COLUMN country_iso_code;
-- +goose StatementEnd
