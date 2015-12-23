
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

alter table tm_user add constraint tmuser_email_UNIQUE unique (email);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

alter table cdn drop index tmuser_email_UNIQUE;
