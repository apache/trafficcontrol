
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE servercheck CHANGE `as` `bf` int(11) DEFAULT NULL;


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE servercheck CHANGE `bf` `as` int(11) DEFAULT NULL;
