
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
insert into role (name, description, priv_level) select * from (select 'deploy', 'Deployment role', 15) as tmp where not exists (select name from role where name = 'deploy') limit 1;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
delete from role where name = 'deploy' and description = 'Deployment role' and priv_level = 15;
