
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
alter table server add column guid varchar(45) DEFAULT NULL AFTER router_port_name;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table server drop column guid;
