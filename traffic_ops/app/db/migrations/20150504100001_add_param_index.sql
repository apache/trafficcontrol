
-- +goose Up
create index parameter_name_value_idx on parameter (name(512),value(512));

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

