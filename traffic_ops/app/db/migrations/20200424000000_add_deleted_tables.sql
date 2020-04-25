-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

-- deleted_type
CREATE TABLE IF NOT EXISTS deleted_type (
    id bigint NOT NULL,
    name text,
    description text,
    use_in_table text,
    last_updated timestamp with time zone DEFAULT now(),
    deleted_time timestamp with time zone NOT NULL DEFAULT now()
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS deleted_type;