
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
alter table server add column st_chg_reason varchar(256) DEFAULT NULL AFTER status;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table server drop column st_chg_reason;
