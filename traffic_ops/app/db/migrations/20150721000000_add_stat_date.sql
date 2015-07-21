-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
alter table stats_summary add column stat_date date;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
