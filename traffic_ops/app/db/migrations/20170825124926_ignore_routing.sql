
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE profile
    ADD routing_disabled BOOLEAN NOT NULL DEFAULT false;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE profile
DROP COLUMN routing_disabled;
