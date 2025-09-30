-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE link_map (
  id text PRIMARY KEY,
  url text,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE link_map;
-- +goose StatementEnd
