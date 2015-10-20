
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

alter table cdn modify name varchar(127);
alter table cdn add constraint cdn_cdn_UNIQUE unique (name);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

alter table cdn drop index cdn_cdn_UNIQUE;
alter table cdn modify name varchar(1024);
