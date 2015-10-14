
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

update deliveryservice set cdn_id = (select id from cdn where name='ALL') where cdn_id is NULL;
alter table deliveryservice drop foreign key fk_cdn1;
alter table deliveryservice modify cdn_id INT(11) NOT NULL;
alter table deliveryservice add CONSTRAINT `fk_cdn1` FOREIGN KEY (`cdn_id`) REFERENCES `cdn` (`id`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

update deliveryservice set cdn_id = NULL where cdn_id = (select id from cdn where name='ALL');\
alter table deliveryservice drop foreign key fk_cdn1;
alter table deliveryservice modify cdn_id INT(11);
alter table deliveryservice add CONSTRAINT `fk_cdn1` FOREIGN KEY (`cdn_id`) REFERENCES `cdn` (`id`);
