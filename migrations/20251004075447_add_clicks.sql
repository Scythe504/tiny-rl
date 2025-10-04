-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE link_map
RENAME COLUMN id TO short_code;

CREATE TABLE clicks
(
  id SERIAL PRIMARY KEY,
  short_code text REFERENCES link_map(short_code),
  ip_addr VARCHAR,
  user_agent text,
  referrer text,
  clicked_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS clicks;
ALTER TABLE link_map
RENAME COLUMN short_code TO id;
-- +goose StatementEnd
