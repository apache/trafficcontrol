ALTER TABLE cdn
  ALTER COLUMN dnssec_enabled DROP DEFAULT,
	ALTER COLUMN dnssec_enabled TYPE boolean
		USING CASE WHEN dnssec_enabled = 1 THEN TRUE
			ELSE FALSE
			END,
  ALTER COLUMN dnssec_enabled SET DEFAULT FALSE;

ALTER TABLE deliveryservice
  ALTER COLUMN active DROP DEFAULT,
	ALTER COLUMN active TYPE boolean
		USING CASE WHEN active = 1 THEN TRUE
			ELSE FALSE
			END,
  ALTER COLUMN active SET DEFAULT FALSE;

ALTER TABLE deliveryservice
  ALTER COLUMN signed DROP DEFAULT,
	ALTER COLUMN signed TYPE boolean
		USING CASE WHEN signed = 1 THEN TRUE
			ELSE FALSE
			END,
  ALTER COLUMN signed SET DEFAULT FALSE;

ALTER TABLE deliveryservice
  ALTER COLUMN ipv6_routing_enabled DROP DEFAULT,
	ALTER COLUMN ipv6_routing_enabled TYPE boolean
		USING CASE WHEN ipv6_routing_enabled = 1 THEN TRUE
			ELSE FALSE
			END,
  ALTER COLUMN ipv6_routing_enabled SET DEFAULT FALSE;

ALTER TABLE deliveryservice
  ALTER COLUMN multi_site_origin DROP DEFAULT,
	ALTER COLUMN multi_site_origin TYPE boolean
		USING CASE WHEN multi_site_origin = 1 THEN TRUE
			ELSE FALSE
			END,
  ALTER COLUMN multi_site_origin SET DEFAULT FALSE;

ALTER TABLE deliveryservice
  ALTER COLUMN regional_geo_blocking DROP DEFAULT,
	ALTER COLUMN regional_geo_blocking TYPE boolean
		USING CASE WHEN regional_geo_blocking = 1 THEN TRUE
			ELSE FALSE
			END,
  ALTER COLUMN regional_geo_blocking SET DEFAULT FALSE;

ALTER TABLE deliveryservice
  ALTER COLUMN logs_enabled DROP DEFAULT,
	ALTER COLUMN logs_enabled TYPE boolean
		USING CASE WHEN logs_enabled = 1 THEN TRUE
			ELSE FALSE
			END,
  ALTER COLUMN logs_enabled SET DEFAULT FALSE;

  ALTER TABLE goose_db_version
  ALTER COLUMN is_applied DROP DEFAULT,
	ALTER COLUMN is_applied TYPE boolean
		USING CASE WHEN is_applied = 1 THEN TRUE
			ELSE FALSE
			END,
  ALTER COLUMN is_applied SET DEFAULT FALSE;

ALTER TABLE parameter
  ALTER COLUMN secure DROP DEFAULT,
	ALTER COLUMN secure TYPE boolean
		USING CASE WHEN secure = 1 THEN TRUE
			ELSE FALSE
			END,
  ALTER COLUMN secure SET DEFAULT FALSE;

ALTER TABLE server
  ALTER COLUMN upd_pending DROP DEFAULT,
	ALTER COLUMN upd_pending TYPE boolean
		USING CASE WHEN upd_pending = 1 THEN TRUE
			ELSE FALSE
			END,
  ALTER COLUMN upd_pending SET DEFAULT FALSE;

ALTER TABLE tm_user
  ALTER COLUMN new_user DROP DEFAULT,
	ALTER COLUMN new_user TYPE boolean
		USING CASE WHEN new_user = 1 THEN TRUE
			ELSE FALSE
			END,
  ALTER COLUMN new_user SET DEFAULT FALSE;

ALTER TABLE to_extension
  ALTER COLUMN isactive DROP DEFAULT,
	ALTER COLUMN isactive TYPE boolean
		USING CASE WHEN isactive = 1 THEN TRUE
			ELSE FALSE
			END,
    ALTER COLUMN isactive SET DEFAULT FALSE;
