
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
insert into type (name, description, use_in_table) values ('CLIENT_STEERING', 'Client-Controlled Steering Delivery Service', 'deliveryservice');


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
delete from type where name = 'CLIENT_STEERING';
