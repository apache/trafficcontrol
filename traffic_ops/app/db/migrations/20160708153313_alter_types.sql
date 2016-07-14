-- +goose Up

-- deliveryservice
ALTER TABLE deliveryservice ALTER COLUMN signed TYPE INT2 USING (signed::integer);
ALTER TABLE deliveryservice ALTER COLUMN qstring_ignore TYPE INT2 USING (qstring_ignore::integer);
ALTER TABLE deliveryservice ALTER COLUMN miss_lat TYPE numeric;
ALTER TABLE deliveryservice ALTER COLUMN miss_long TYPE numeric;


ALTER TABLE cachegroup ALTER COLUMN latitude TYPE numeric;
ALTER TABLE cachegroup ALTER COLUMN longitude TYPE numeric;
