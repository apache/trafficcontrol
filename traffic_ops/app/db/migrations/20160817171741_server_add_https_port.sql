
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
alter table server add https_port int(10) unsigned default NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table server drop column https_port;