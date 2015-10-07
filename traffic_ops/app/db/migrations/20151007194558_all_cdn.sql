
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

insert into cdn (name) values ('ALL');
update server set cdn_id = (select id from cdn where name='ALL') where cdn_id is NULL;
alter table server drop foreign key fk_cdn2;
alter table server modify cdn_id INT(11) NOT NULL;
alter table server add CONSTRAINT `fk_cdn2` FOREIGN KEY (`cdn_id`) REFERENCES `cdn` (`id`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

update server set cdn_id = NULL where cdn_id = (select id from cdn where name='ALL');
delete from cdn where name = 'ALL';
alter table server drop foreign key fk_cdn2;
alter table server modify cdn_id INT(11);
alter table server add CONSTRAINT `fk_cdn2` FOREIGN KEY (`cdn_id`) REFERENCES `cdn` (`id`);
