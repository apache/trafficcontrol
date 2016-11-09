/*
	Copyright 2015 Comcast Cable Communications Management, LLC

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	    http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE cachegroup ALTER COLUMN latitude  TYPE numeric;
ALTER TABLE cachegroup ALTER COLUMN longitude TYPE numeric;

ALTER TABLE cdn
  ALTER COLUMN dnssec_enabled DROP DEFAULT,
	ALTER COLUMN dnssec_enabled TYPE boolean
		USING CASE WHEN dnssec_enabled = 0 THEN FALSE
			WHEN dnssec_enabled = 1 THEN TRUE
			ELSE NULL
			END,
  ALTER COLUMN dnssec_enabled SET DEFAULT FALSE;

ALTER TABLE deliveryservice ALTER COLUMN miss_lat                     TYPE numeric;
ALTER TABLE deliveryservice ALTER COLUMN miss_long                    TYPE numeric;

ALTER TABLE deliveryservice
  ALTER COLUMN active DROP DEFAULT,
	ALTER COLUMN active TYPE boolean
		USING CASE WHEN active = 0 THEN FALSE
			WHEN active = 1 THEN TRUE
			ELSE NULL
			END,
  ALTER COLUMN active SET DEFAULT FALSE;

ALTER TABLE deliveryservice
  ALTER COLUMN signed DROP DEFAULT,
	ALTER COLUMN signed TYPE boolean
		USING CASE WHEN signed = 0 THEN FALSE
			WHEN signed = 1 THEN TRUE
			ELSE NULL
			END,
  ALTER COLUMN signed SET DEFAULT FALSE;

ALTER TABLE deliveryservice
  ALTER COLUMN ipv6_routing_enabled DROP DEFAULT,
	ALTER COLUMN ipv6_routing_enabled TYPE boolean
		USING CASE WHEN ipv6_routing_enabled = 0 THEN FALSE
			WHEN ipv6_routing_enabled = 1 THEN TRUE
			ELSE NULL
			END,
  ALTER COLUMN ipv6_routing_enabled SET DEFAULT FALSE;

ALTER TABLE deliveryservice
  ALTER COLUMN multi_site_origin DROP DEFAULT,
	ALTER COLUMN multi_site_origin TYPE boolean
		USING CASE WHEN multi_site_origin = 0 THEN FALSE
			WHEN multi_site_origin = 1 THEN TRUE
			ELSE NULL
			END,
  ALTER COLUMN multi_site_origin SET DEFAULT FALSE;

ALTER TABLE deliveryservice
  ALTER COLUMN regional_geo_blocking DROP DEFAULT,
	ALTER COLUMN regional_geo_blocking TYPE boolean
		USING CASE WHEN regional_geo_blocking = 0 THEN FALSE
			WHEN regional_geo_blocking = 1 THEN TRUE
			ELSE NULL
			END,
  ALTER COLUMN regional_geo_blocking SET DEFAULT FALSE;

ALTER TABLE deliveryservice
  ALTER COLUMN logs_enabled DROP DEFAULT,
	ALTER COLUMN logs_enabled TYPE boolean
		USING CASE WHEN logs_enabled = 0 THEN FALSE
			WHEN logs_enabled = 1 THEN TRUE
			ELSE NULL
			END,
  ALTER COLUMN logs_enabled SET DEFAULT FALSE;

ALTER TABLE goose_db_version
  ALTER COLUMN is_applied DROP DEFAULT,
	ALTER COLUMN is_applied TYPE boolean
		USING CASE WHEN is_applied = 0 THEN FALSE
			WHEN is_applied = 1 THEN TRUE
			ELSE FALSE
			END,
  ALTER COLUMN is_applied SET DEFAULT FALSE;

ALTER TABLE parameter
  ALTER COLUMN secure DROP DEFAULT,
	ALTER COLUMN secure TYPE boolean
		USING CASE WHEN secure = 0 THEN FALSE
			WHEN secure = 1 THEN TRUE
			ELSE NULL
			END,
  ALTER COLUMN secure SET DEFAULT FALSE;

ALTER TABLE server
  ALTER COLUMN upd_pending DROP DEFAULT,
	ALTER COLUMN upd_pending TYPE boolean
		USING CASE WHEN upd_pending = 0 THEN FALSE
			WHEN upd_pending = 1 THEN TRUE
			ELSE NULL
			END,
  ALTER COLUMN upd_pending SET DEFAULT FALSE;

ALTER TABLE tm_user
  ALTER COLUMN new_user DROP DEFAULT,
	ALTER COLUMN new_user TYPE boolean
		USING CASE WHEN new_user = 0 THEN FALSE
			WHEN new_user = 1 THEN TRUE
			ELSE NULL
			END,
  ALTER COLUMN new_user SET DEFAULT FALSE;

ALTER TABLE to_extension
  ALTER COLUMN isactive DROP DEFAULT,
	ALTER COLUMN isactive TYPE boolean
		USING CASE WHEN isactive = 0 THEN FALSE
			WHEN isactive = 1 THEN TRUE
			ELSE NULL
			END,
    ALTER COLUMN isactive SET DEFAULT FALSE;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE cachegroup ALTER COLUMN latitude  TYPE double precision;
ALTER TABLE cachegroup ALTER COLUMN longitude TYPE double precision;

ALTER TABLE cdn
  ALTER COLUMN dnssec_enabled DROP DEFAULT,
	ALTER COLUMN dnssec_enabled TYPE SMALLINT
   USING CASE WHEN dnssec_enabled THEN 1 ELSE 0 END,
  ALTER COLUMN dnssec_enabled SET DEFAULT 0;

ALTER TABLE deliveryservice ALTER COLUMN miss_lat  TYPE double precision;
ALTER TABLE deliveryservice ALTER COLUMN miss_long TYPE double precision;

ALTER TABLE deliveryservice
  ALTER COLUMN active DROP DEFAULT,
	ALTER COLUMN active TYPE SMALLINT
   USING CASE WHEN active THEN 1 ELSE 0 END,
  ALTER COLUMN active SET DEFAULT 0;

ALTER TABLE deliveryservice
  ALTER COLUMN signed DROP DEFAULT,
	ALTER COLUMN signed TYPE SMALLINT
   USING CASE WHEN signed THEN 1 ELSE 0 END,
  ALTER COLUMN signed SET DEFAULT 0;

ALTER TABLE deliveryservice
  ALTER COLUMN ipv6_routing_enabled DROP DEFAULT,
	ALTER COLUMN ipv6_routing_enabled TYPE SMALLINT
   USING CASE WHEN ipv6_routing_enabled THEN 1 ELSE 0 END,
  ALTER COLUMN ipv6_routing_enabled SET DEFAULT 0;

ALTER TABLE deliveryservice
  ALTER COLUMN multi_site_origin DROP DEFAULT,
	ALTER COLUMN multi_site_origin TYPE SMALLINT
   USING CASE WHEN multi_site_origin THEN 1 ELSE 0 END,
  ALTER COLUMN multi_site_origin SET DEFAULT 0;

ALTER TABLE deliveryservice
  ALTER COLUMN regional_geo_blocking DROP DEFAULT,
	ALTER COLUMN regional_geo_blocking TYPE SMALLINT
   USING CASE WHEN regional_geo_blocking THEN 1 ELSE 0 END,
  ALTER COLUMN regional_geo_blocking SET DEFAULT 0;

ALTER TABLE deliveryservice
  ALTER COLUMN logs_enabled DROP DEFAULT,
	ALTER COLUMN logs_enabled TYPE SMALLINT
   USING CASE WHEN logs_enabled THEN 1 ELSE 0 END,
  ALTER COLUMN logs_enabled SET DEFAULT 0;

ALTER TABLE goose_db_version
  ALTER COLUMN is_applied DROP DEFAULT,
	ALTER COLUMN is_applied TYPE SMALLINT
   USING CASE WHEN is_applied THEN 1 ELSE 0 END,
  ALTER COLUMN is_applied SET DEFAULT 0;

ALTER TABLE parameter
  ALTER COLUMN secure DROP DEFAULT,
	ALTER COLUMN secure TYPE SMALLINT
   USING CASE WHEN secure THEN 1 ELSE 0 END,
  ALTER COLUMN secure SET DEFAULT 0;

ALTER TABLE server
  ALTER COLUMN upd_pending DROP DEFAULT,
	ALTER COLUMN upd_pending TYPE SMALLINT
   USING CASE WHEN upd_pending THEN 1 ELSE 0 END,
  ALTER COLUMN upd_pending SET DEFAULT 0;

ALTER TABLE tm_user
  ALTER COLUMN new_user DROP DEFAULT,
	ALTER COLUMN new_user TYPE SMALLINT
   USING CASE WHEN new_user THEN 1 ELSE 0 END,
  ALTER COLUMN new_user SET DEFAULT 0;

ALTER TABLE to_extension
  ALTER COLUMN isactive DROP DEFAULT,
	ALTER COLUMN isactive TYPE SMALLINT
   USING CASE WHEN isactive THEN 1 ELSE 0 END,
  ALTER COLUMN isactive SET DEFAULT 0;
