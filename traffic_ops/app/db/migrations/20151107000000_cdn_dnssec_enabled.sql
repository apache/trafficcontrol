
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
 alter table cdn add column `dnssec_enabled` tinyint(4) default 0;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

