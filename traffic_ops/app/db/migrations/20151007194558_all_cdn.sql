
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

insert into cdn (name) values ('ALL');
update server set cdn_id = (select id from cdn where name='ALL') where cdn_id is NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

update server set cdn_id = NULL where cdn_id = (select id from cdn where name='ALL');
delete from cdn where name = 'ALL';
