
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
alter table server add offline_reason varchar(256) AFTER status;


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table server drop column offline_reason;
