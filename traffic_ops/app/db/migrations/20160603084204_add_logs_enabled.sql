
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
alter table deliveryservice add logs_enabled tinyint(1);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table deliveryservice drop column logs_enabled;

