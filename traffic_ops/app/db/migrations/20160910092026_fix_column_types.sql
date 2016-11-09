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
	ALTER COLUMN dnssec_enabled TYPE boolean
		USING CASE WHEN dnssec_enabled = 0 THEN FALSE
			WHEN dnssec_enabled = 1 THEN TRUE
			ELSE NULL;

ALTER TABLE deliveryservice ALTER COLUMN miss_lat                     TYPE numeric;
ALTER TABLE deliveryservice ALTER COLUMN miss_long                    TYPE numeric;

ALTER TABLE deliveryservice
	ALTER COLUMN active TYPE boolean
		USING CASE WHEN active = 0 THEN FALSE
			WHEN active = 1 THEN TRUE
			ELSE NULL;

ALTER TABLE deliveryservice
	ALTER COLUMN signed TYPE boolean
		USING CASE WHEN signed = 0 THEN FALSE
			WHEN signed = 1 THEN TRUE
			ELSE NULL;

ALTER TABLE deliveryservice
	ALTER COLUMN ipv6_routing_enabled TYPE boolean
		USING CASE WHEN ipv6_routing_enabled = 0 THEN FALSE
			WHEN ipv6_routing_enabled = 1 THEN TRUE
			ELSE NULL;

ALTER TABLE deliveryservice
	ALTER COLUMN multi_site_origin TYPE boolean
		USING CASE WHEN multi_site_origin = 0 THEN FALSE
			WHEN multi_site_origin = 1 THEN TRUE
			ELSE NULL;

ALTER TABLE deliveryservice
	ALTER COLUMN regional_geo_blocking TYPE boolean
		USING CASE WHEN regional_geo_blocking = 0 THEN FALSE
			WHEN regional_geo_blocking = 1 THEN TRUE
			ELSE NULL;

ALTER TABLE deliveryservice
	ALTER COLUMN logs_enabled TYPE boolean
		USING CASE WHEN logs_enabled = 0 THEN FALSE
			WHEN logs_enabled = 1 THEN TRUE
			ELSE NULL;

ALTER TABLE deliveryservice
	ALTER COLUMN logs_enabled TYPE boolean
		USING CASE WHEN logs_enabled = 0 THEN FALSE
			WHEN logs_enabled = 1 THEN TRUE
			ELSE NULL;

ALTER TABLE parameter
	ALTER COLUMN secure TYPE boolean
		USING CASE WHEN secure = 0 THEN FALSE
			WHEN secure = 1 THEN TRUE
			ELSE NULL;

ALTER TABLE server
	ALTER COLUMN upd_pending TYPE boolean
		USING CASE WHEN upd_pending = 0 THEN FALSE
			WHEN upd_pending = 1 THEN TRUE
			ELSE NULL;

ALTER TABLE tm_user
	ALTER COLUMN new_user TYPE boolean
		USING CASE WHEN new_user = 0 THEN FALSE
			WHEN new_user = 1 THEN TRUE
			ELSE NULL;

ALTER TABLE to_extension
	ALTER COLUMN isactive TYPE boolean
		USING CASE WHEN isactive = 0 THEN FALSE
			WHEN isactive = 1 THEN TRUE
			ELSE NULL;
