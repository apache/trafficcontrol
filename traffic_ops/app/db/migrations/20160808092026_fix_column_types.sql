
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE cachegroup ALTER COLUMN latitude  TYPE numeric;
ALTER TABLE cachegroup ALTER COLUMN longitude TYPE numeric;

ALTER TABLE deliveryservice ALTER COLUMN signed TYPE smallint USING (signed::int::smallint);
ALTER TABLE deliveryservice ALTER COLUMN qstring_ignore TYPE smallint USING (qstring_ignore::int::smallint);

ALTER TABLE deliveryservice
  ALTER COLUMN geo_limit DROP DEFAULT,
  ALTER COLUMN geo_limit TYPE smallint USING (geo_limit::int::smallint),
  ALTER COLUMN geo_limit SET DEFAULT '0';

ALTER TABLE deliveryservice ALTER COLUMN miss_lat                     TYPE numeric;
ALTER TABLE deliveryservice ALTER COLUMN miss_long                    TYPE numeric;
ALTER TABLE deliveryservice ALTER COLUMN multi_site_origin            TYPE smallint USING (multi_site_origin::int::smallint);
ALTER TABLE deliveryservice ALTER COLUMN regional_geo_blocking        TYPE smallint USING (regional_geo_blocking::int::smallint);
ALTER TABLE deliveryservice ALTER COLUMN logs_enabled                 TYPE smallint USING (logs_enabled::int::smallint);
ALTER TABLE deliveryservice ALTER COLUMN multi_site_origin_algorithm  TYPE smallint USING (multi_site_origin_algorithm::int::smallint);

ALTER TABLE parameter
  ALTER COLUMN secure DROP DEFAULT,
  ALTER COLUMN secure TYPE smallint USING (secure::int::smallint),
  ALTER COLUMN secure SET DEFAULT '0';

ALTER TABLE server
  ALTER COLUMN upd_pending DROP DEFAULT,
  ALTER COLUMN upd_pending TYPE smallint USING (upd_pending::int::smallint),
  ALTER COLUMN upd_pending SET DEFAULT '0';

ALTER TABLE tm_user
  ALTER COLUMN new_user DROP DEFAULT,
  ALTER COLUMN new_user TYPE smallint USING (new_user::int::smallint),
  ALTER COLUMN new_user SET DEFAULT '1';

ALTER TABLE to_extension  ALTER COLUMN isactive TYPE smallint USING (isactive::int::smallint);



-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE cachegroup ALTER COLUMN latitude  TYPE double precision;
ALTER TABLE cachegroup ALTER COLUMN longitude TYPE double precision;

ALTER TABLE deliveryservice ALTER COLUMN signed         TYPE boolean USING (signed::int::boolean);
ALTER TABLE deliveryservice ALTER COLUMN qstring_ignore TYPE boolean USING (qstring_ignore::int::boolean);

ALTER TABLE deliveryservice
  ALTER COLUMN geo_limit DROP DEFAULT,
  ALTER COLUMN geo_limit TYPE boolean USING (geo_limit::int::boolean),
  ALTER COLUMN geo_limit SET DEFAULT FALSE;

ALTER TABLE deliveryservice ALTER COLUMN miss_lat                     TYPE double precision;
ALTER TABLE deliveryservice ALTER COLUMN miss_long                    TYPE double precision;
ALTER TABLE deliveryservice ALTER COLUMN multi_site_origin            TYPE boolean USING (multi_site_origin::int::boolean);
ALTER TABLE deliveryservice ALTER COLUMN regional_geo_blocking        TYPE boolean USING (regional_geo_blocking::int::boolean);
ALTER TABLE deliveryservice ALTER COLUMN multi_site_origin_algorithm  TYPE boolean USING (multi_site_origin_algorithm::int::boolean);

ALTER TABLE parameter
  ALTER COLUMN secure DROP DEFAULT,
  ALTER COLUMN secure TYPE boolean USING (secure::int::boolean),
  ALTER COLUMN secure SET DEFAULT FALSE;

ALTER TABLE server
  ALTER COLUMN upd_pending DROP DEFAULT,
  ALTER COLUMN upd_pending TYPE boolean USING (upd_pending::int::boolean),
  ALTER COLUMN upd_pending SET DEFAULT FALSE;

ALTER TABLE tm_user
  ALTER COLUMN new_user DROP DEFAULT,
  ALTER COLUMN new_user TYPE boolean USING (new_user::int::boolean),
  ALTER COLUMN new_user SET DEFAULT TRUE;

ALTER TABLE to_extension  ALTER COLUMN isactive TYPE boolean USING (isactive::int::boolean);
