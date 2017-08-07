-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE tenant ADD COLUMN email text;
UPDATE tenant SET email = 'defaulttenant@comcast.com' WHERE name = 'root';
ALTER TABLE tenant ALTER COLUMN email SET NOT NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE tenant DROP COLUMN email;
