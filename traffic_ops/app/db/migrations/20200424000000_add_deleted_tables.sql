
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

-- last_deleted
CREATE TABLE IF NOT EXISTS last_deleted (
  tab_name text NOT NULL,
  last_updated timestamp with time zone NOT NULL DEFAULT now()
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS last_deleted;