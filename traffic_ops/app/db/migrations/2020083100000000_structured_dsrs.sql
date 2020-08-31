-- syntax:postgresql
/*
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
CREATE TABLE all_delivery_services (
	active boolean DEFAULT FALSE NOT NULL,
	anonymous_blocking_enabled boolean NOT NULL DEFAULT FALSE,
	backend_id bigserial PRIMARY KEY,
	cacheurl text,
	ccr_dns_ttl bigint,
	cdn_id bigint NOT NULL REFERENCES cdn(id) ON UPDATE RESTRICT ON DELETE RESTRICT,
	check_path text,
	consistent_hash_regex text,
	"original" boolean NOT NULL DEFAULT FALSE,
	deep_caching_type deep_caching_type NOT NULL DEFAULT 'NEVER',
	display_name text NOT NULL,
	dns_bypass_cname text,
	dns_bypass_ip text,
	dns_bypass_ip6 text,
	dns_bypass_ttl bigint,
	dscp bigint NOT NULL,
	ecs_enabled boolean NOT NULL DEFAULT FALSE,
	edge_header_rewrite text,
	first_header_rewrite text,
	fq_pacing_rate bigint DEFAULT 0,
	geo_limit smallint DEFAULT '0'::smallint,
	geo_limit_countries text,
	geo_provider smallint DEFAULT '0'::smallint,
	geolimit_redirect_url text,
	global_max_mbps bigint,
	global_max_tps bigint,
	http_bypass_fqdn text,
	id bigserial NOT NULL DEFAULT nextval('deliveryservice_id_seq'::regclass),
	info_url text,
	initial_dispersion bigint DEFAULT '1'::bigint,
	inner_header_rewrite text,
	ipv6_routing_enabled boolean DEFAULT FALSE,
	last_header_rewrite text,
	last_updated timestamp with time zone NOT NULL DEFAULT now(),
	logs_enabled boolean DEFAULT FALSE,
	long_desc text,
	long_desc_1 text,
	long_desc_2 text,
	max_dns_answers bigint DEFAULT '5'::bigint,
	max_origin_connections bigint NOT NULL DEFAULT 0 CHECK (max_origin_connections >= 0),
	mid_header_rewrite text,
	miss_lat numeric,
	miss_long numeric,
	multi_site_origin boolean DEFAULT FALSE,
	multi_site_origin_algorithm smallint,
	origin_shield text,
	profile bigint REFERENCES profile(id),
	protocol smallint DEFAULT '0'::smallint,
	qstring_ignore smallint,
	range_request_handling smallint DEFAULT '0'::smallint,
	range_slice_block_size integer DEFAULT NULL CHECK (range_slice_block_size >= 262144 AND range_slice_block_size <= 33554432),
	regex_remap text,
	regional_geo_blocking boolean DEFAULT FALSE NOT NULL,
	remap_text text,
	routing_name text NOT NULL DEFAULT 'cdn' CHECK (length(routing_name) > 0),
	service_category text REFERENCES service_category("name") ON UPDATE CASCADE,
	signing_algorithm deliveryservice_signature_type,
	ssl_key_version bigint DEFAULT '0'::bigint,
	tenant_id bigint NOT NULL REFERENCES tenant(id) MATCH FULL,
	tr_request_headers text,
	tr_response_headers text,
	"type" bigint NOT NULL REFERENCES "type"(id),
	xml_id text NOT NULL
);

INSERT INTO all_delivery_services (
	active,
	anonymous_blocking_enabled,
	cacheurl,
	ccr_dns_ttl,
	cdn_id,
	check_path,
	consistent_hash_regex,
	"original",
	deep_caching_type,
	display_name,
	dns_bypass_cname,
	dns_bypass_ip,
	dns_bypass_ip6,
	dns_bypass_ttl,
	dscp,
	ecs_enabled,
	edge_header_rewrite,
	first_header_rewrite,
	fq_pacing_rate,
	geo_limit,
	geo_limit_countries,
	geo_provider,
	geolimit_redirect_url,
	global_max_mbps,
	global_max_tps,
	http_bypass_fqdn,
	id,
	info_url,
	initial_dispersion,
	inner_header_rewrite,
	ipv6_routing_enabled,
	last_header_rewrite,
	last_updated,
	logs_enabled,
	long_desc,
	long_desc_1,
	long_desc_2,
	max_dns_answers,
	max_origin_connections,
	mid_header_rewrite,
	miss_lat,
	miss_long,
	multi_site_origin,
	multi_site_origin_algorithm,
	origin_shield,
	profile,
	protocol,
	qstring_ignore,
	range_request_handling,
	range_slice_block_size,
	regex_remap,
	regional_geo_blocking,
	remap_text,
	routing_name,
	service_category,
	signing_algorithm,
	ssl_key_version,
	tenant_id,
	tr_request_headers,
	tr_response_headers,
	"type",
	xml_id
) VALUES
SELECT
	active,
	anonymous_blocking_enabled,
	cacheurl,
	ccr_dns_ttl,
	cdn_id,
	check_path,
	consistent_hash_regex,
	TRUE,
	deep_caching_type,
	display_name,
	dns_bypass_cname,
	dns_bypass_ip,
	dns_bypass_ip6,
	dns_bypass_ttl,
	dscp,
	ecs_enabled,
	edge_header_rewrite,
	first_header_rewrite,
	fq_pacing_rate,
	geo_limit,
	geo_limit_countries,
	geo_provider,
	geolimit_redirect_url,
	global_max_mbps,
	global_max_tps,
	http_bypass_fqdn,
	id,
	info_url,
	initial_dispersion,
	inner_header_rewrite,
	ipv6_routing_enabled,
	last_header_rewrite,
	last_updated,
	logs_enabled,
	long_desc,
	long_desc_1,
	long_desc_2,
	max_dns_answers,
	max_origin_connections,
	mid_header_rewrite,
	miss_lat,
	miss_long,
	multi_site_origin,
	multi_site_origin_algorithm,
	origin_shield,
	profile,
	protocol,
	qstring_ignore,
	range_request_handling,
	range_slice_block_size,
	regex_remap,
	regional_geo_blocking,
	remap_text,
	routing_name,
	service_category,
	signing_algorithm,
	ssl_key_version,
	tenant_id,
	tr_request_headers,
	tr_response_headers,
	"type",
	xml_id
FROM deliveryservice;

DROP INDEX idx_k_deliveryservice_tenant_idx;
DROP INDEX idx_89502_fk_cdn1;
DROP INDEX idx_89502_fk_deliveryservice_profile1;
DROP INDEX idx_89502_fk_deliveryservice_type1;

ALTER TABLE ONLY deliveryservice DROP CONSTRAINT fk_cdn1;
ALTER TABLE ONLY deliveryservice DROP CONSTRAINT fk_deliveryservice_profile1;
ALTER TABLE ONLY deliveryservice DROP CONSTRAINT fk_deliveryservice_type1;
ALTER TABLE ONLY deliveryservice DROP CONSTRAINT fk_tenantid;
ALTER TABLE ONLY deliveryservice_regex DROP CONSTRAINT fk_ds_to_regex_regex1;
ALTER TABLE ONLY deliveryservice_server DROP CONSTRAINT fk_ds_to_cs_deliveryservice1;
ALTER TABLE ONLY federation_deliveryservice DROP CONSTRAINT fk_federation_to_ds1;
ALTER TABLE ONLY job DROP CONSTRAINT fk_job_deliveryservice1;
ALTER TABLE ONLY staticdnsentry DROP CONSTRAINT fk_staticdnsentry_ds;
ALTER TABLE ONLY steering_target DROP CONSTRAINT fk_steering_target_delivery_service;
ALTER TABLE ONLY steering_target DROP CONSTRAINT fk_steering_target_target;
ALTER TABLE ONLY deliveryservice_tmuser DROP CONSTRAINT fk_tm_user_ds;
ALTER TABLE ONLY origin DROP CONSTRAINT origin_deliveryservice_fkey;

DROP TABLE deliveryservice;
CREATE VIEW deliveryservice AS
SELECT
	active,
	anonymous_blocking_enabled,
	cacheurl,
	ccr_dns_ttl,
	cdn_id,
	check_path,
	consistent_hash_regex,
	deep_caching_type,
	display_name,
	dns_bypass_cname,
	dns_bypass_ip,
	dns_bypass_ip6,
	dns_bypass_ttl,
	dscp,
	ecs_enabled,
	edge_header_rewrite,
	first_header_rewrite,
	fq_pacing_rate,
	geo_limit,
	geo_limit_countries,
	geo_provider,
	geolimit_redirect_url,
	global_max_mbps,
	global_max_tps,
	http_bypass_fqdn,
	id,
	info_url,
	initial_dispersion,
	inner_header_rewrite,
	ipv6_routing_enabled,
	last_header_rewrite,
	last_updated,
	logs_enabled,
	long_desc,
	long_desc_1,
	long_desc_2,
	max_dns_answers,
	max_origin_connections,
	mid_header_rewrite,
	miss_lat,
	miss_long,
	multi_site_origin,
	multi_site_origin_algorithm,
	origin_shield,
	profile,
	protocol,
	qstring_ignore,
	range_request_handling,
	range_slice_block_size,
	regex_remap,
	regional_geo_blocking,
	remap_text,
	routing_name,
	service_category,
	signing_algorithm,
	ssl_key_version,
	tenant_id,
	tr_request_headers,
	tr_response_headers,
	"type",
	xml_id
FROM all_delivery_services
WHERE all_delivery_services."original" IS TRUE;

CREATE INDEX idx_deliveryservice_tenant ON all_delivery_services USING btree (tenant_id);
CREATE INDEX idx_deliveryservice_cdn ON all_delivery_services USING btree (cdn_id);
CREATE INDEX idx_deliveryservice_profile ON all_delivery_services USING btree (profile);
CREATE INDEX idx_deliveryservice_type ON all_delivery_services USING btree("type");

CREATE UNIQUE INDEX idx_deliveryservice_id ON all_delivery_services(id) WHERE "original" IS TRUE;
CREATE UNIQUE INDEX idx_deliveryservice_xml_id ON all_delivery_services(xml_id) WHERE "original" IS TRUE;

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON all_delivery_services FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

ALTER TABLE ONLY deliveryservice_regex ADD CONSTRAINT fk_deliveryservice_regex_to_deliveryservice FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY deliveryservice_server ADD CONSTRAINT fk_deliveryservice_server_to_deliveryservice FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY federation_deliveryservice ADD CONSTRAINT fk_federation_deliveryservice_to_deliveryservice FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY job ADD CONSTRAINT fk_job_to_deliveryservice FOREIGN KEY (job_deliveryservice) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY staticdnsentry ADD CONSTRAINT fk_staticdnsentry_to_deliveryservice FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY steering_target ADD CONSTRAINT fk_steering_target_to_deliveryservice FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY steering_target ADD CONSTRAINT fk_steering_target_target_to_deliveryservice FOREIGN KEY (target) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY deliveryservice_tmuser ADD CONSTRAINT fk_deliveryservice_tmuser_to_deliveryservice FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY origin ADD CONSTRAINT fk_origin_to_deliveryservice FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON DELETE CASCADE;

INSERT INTO all_delivery_services (
	active,
	anonymous_blocking_enabled,
	cacheurl,
	ccr_dns_ttl,
	cdn_id,
	check_path,
	consistent_hash_regex,
	"original",
	deep_caching_type,
	display_name,
	dns_bypass_cname,
	dns_bypass_ip,
	dns_bypass_ip6,
	dns_bypass_ttl,
	dscp,
	ecs_enabled,
	edge_header_rewrite,
	first_header_rewrite,
	fq_pacing_rate,
	geo_limit,
	geo_limit_countries,
	geo_provider,
	geolimit_redirect_url,
	global_max_mbps,
	global_max_tps,
	http_bypass_fqdn,
	id,
	info_url,
	initial_dispersion,
	inner_header_rewrite,
	ipv6_routing_enabled,
	last_header_rewrite,
	last_updated,
	logs_enabled,
	long_desc,
	long_desc_1,
	long_desc_2,
	max_dns_answers,
	max_origin_connections,
	mid_header_rewrite,
	miss_lat,
	miss_long,
	multi_site_origin,
	multi_site_origin_algorithm,
	origin_shield,
	profile,
	protocol,
	qstring_ignore,
	range_request_handling,
	range_slice_block_size,
	regex_remap,
	regional_geo_blocking,
	remap_text,
	routing_name,
	service_category,
	signing_algorithm,
	ssl_key_version,
	tenant_id,
	tr_request_headers,
	tr_response_headers,
	"type",
	xml_id
) VALUES
SELECT (
	active,
	anonymous_blocking_enabled,
	cacheurl,
	ccr_dns_ttl,
	cdn_id,
	check_path,
	consistent_hash_regex,
	FALSE,
	deep_caching_type,
	display_name,
	dns_bypass_cname,
	dns_bypass_ip,
	dns_bypass_ip6,
	dns_bypass_ttl,
	dscp,
	ecs_enabled,
	edge_header_rewrite,
	first_header_rewrite,
	fq_pacing_rate,
	geo_limit,
	geo_limit_countries,
	geo_provider,
	geolimit_redirect_url,
	global_max_mbps,
	global_max_tps,
	http_bypass_fqdn,
	id,
	info_url,
	initial_dispersion,
	inner_header_rewrite,
	ipv6_routing_enabled,
	last_header_rewrite,
	last_updated,
	logs_enabled,
	long_desc,
	long_desc_1,
	long_desc_2,
	max_dns_answers,
	max_origin_connections,
	mid_header_rewrite,
	miss_lat,
	miss_long,
	multi_site_origin,
	multi_site_origin_algorithm,
	origin_shield,
	profile,
	protocol,
	qstring_ignore,
	range_request_handling,
	range_slice_block_size,
	regex_remap,
	regional_geo_blocking,
	remap_text,
	routing_name,
	service_category,
	signing_algorithm,
	ssl_key_version,
	tenant_id,
	tr_request_headers,
	tr_response_headers,
	"type",
	xml_id
FROM (
	SELECT jsonb_to_record(deliveryservice) AS (
		active boolean,
		anonymous_blocking_enabled boolean,
		cacheurl text,
		ccr_dns_ttl bigint,
		cdn_id bigint,
		check_path text,
		consistent_hash_regex text,
		deep_caching_type deep_caching_type,
		display_name text,
		dns_bypass_cname text,
		dns_bypass_ip text,
		dns_bypass_ip6 text,
		dns_bypass_ttl bigint,
		dscp bigint,
		ecs_enabled boolean,
		edge_header_rewrite text,
		first_header_rewrite text,
		fq_pacing_rate bigint,
		geo_limit smallint,
		geo_limit_countries text,
		geo_provider smallint,
		geolimit_redirect_url text,
		global_max_mbps bigint,
		global_max_tps bigint,
		http_bypass_fqdn text,
		id bigint,
		info_url text,
		initial_dispersion bigint,
		inner_header_rewrite text,
		ipv6_routing_enabled boolean,
		last_header_rewrite text,
		last_updated timestamp with time zone,
		logs_enabled boolean,
		long_desc text,
		long_desc_1 text,
		long_desc_2 text,
		max_dns_answers bigint,
		max_origin_connections bigint,
		mid_header_rewrite text,
		miss_lat numeric,
		miss_long numeric,
		multi_site_origin boolean,
		multi_site_origin_algorithm smallint,
		origin_shield text,
		profile bigint,
		protocol smallint,
		qstring_ignore smallint,
		range_request_handling smallint,
		range_slice_block_size integer,
		regex_remap text,
		regional_geo_blocking boolean,
		remap_text text,
		routing_name text,
		service_category text,
		signing_algorithm deliveryservice_signature_type,
		ssl_key_version bigint,
		tenant_id bigint,
		tr_request_headers text,
		tr_response_headers text,
		"type" bigint,
		xml_id text
	)
	FROM deliveryservice_request
);

ALTER TABLE deliveryservice_request
DROP COLUMN deliveryservice;
ALTER TABLE deliveryservice_request
ADD COLUMN requested bigint NOT NULL REFERENCES all_delivery_services(backend_id);
ALTER TABLE deliveryservice_request
ADD COLUMN original bigint REFERENCES all_delivery_services(backend_id);

-- +goose Down

ALTER TABLE deliveryservice_request
ADD COLUMN deliveryservice jsonb;

WITH requested AS (
	SELECT *
	FROM all_delivery_services
	WHERE all_delivery_services."original" IS TRUE
)
UPDATE deliveryservice_request
SET deliveryservice = row_to_json(
	SELECT
		active,
		anonymous_blocking_enabled,
		cacheurl,
		ccr_dns_ttl,
		cdn_id,
		check_path,
		consistent_hash_regex,
		deep_caching_type,
		display_name,
		dns_bypass_cname,
		dns_bypass_ip,
		dns_bypass_ip6,
		dns_bypass_ttl,
		dscp,
		ecs_enabled,
		edge_header_rewrite,
		first_header_rewrite,
		fq_pacing_rate,
		geo_limit,
		geo_limit_countries,
		geo_provider,
		geolimit_redirect_url,
		global_max_mbps,
		global_max_tps,
		http_bypass_fqdn,
		id,
		info_url,
		initial_dispersion,
		inner_header_rewrite,
		ipv6_routing_enabled,
		last_header_rewrite,
		last_updated,
		logs_enabled,
		long_desc,
		long_desc_1,
		long_desc_2,
		max_dns_answers,
		max_origin_connections,
		mid_header_rewrite,
		miss_lat,
		miss_long,
		multi_site_origin,
		multi_site_origin_algorithm,
		origin_shield,
		profile,
		protocol,
		qstring_ignore,
		range_request_handling,
		range_slice_block_size,
		regex_remap,
		regional_geo_blocking,
		remap_text,
		routing_name,
		service_category,
		signing_algorithm,
		ssl_key_version,
		tenant_id,
		tr_request_headers,
		tr_response_headers,
		"type",
		xml_id
	FROM requested
	WHERE requested.backend_id = deliveryservice_request.requested
);

DELETE FROM all_delivery_services
WHERE original IS FALSE;

ALTER TABLE deliveryservice_request
DROP COLUMN requested;
ALTER TABLE deliveryservice_request
DROP COLUMN original;

ALTER TABLE deliveryservice_request
ALTER COLUMN deliveryservice
SET NOT NULL;

ALTER TABLE ONLY deliveryservice_regex DROP CONSTRAINT fk_deliveryservice_regex_to_deliveryservice;
ALTER TABLE ONLY deliveryservice_server DROP CONSTRAINT fk_deliveryservice_server_to_deliveryservice;
ALTER TABLE ONLY federation_deliveryservice DROP CONSTRAINT fk_federation_deliveryservice_to_deliveryservice;
ALTER TABLE ONLY job DROP CONSTRAINT fk_job_to_deliveryservice;
ALTER TABLE ONLY staticdnsentry DROP CONSTRAINT fk_staticdnsentry_to_deliveryservice;
ALTER TABLE ONLY steering_target DROP CONSTRAINT fk_steering_target_to_deliveryservice;
ALTER TABLE ONLY steering_target DROP CONSTRAINT fk_steering_target_target_to_deliveryservice;
ALTER TABLE ONLY deliveryservice_tmuser DROP CONSTRAINT fk_deliveryservice_tmuser_to_deliveryservice;
ALTER TABLE ONLY origin DROP CONSTRAINT fk_origin_to_deliveryservice;

DROP VIEW deliveryservice;
DROP INDEX idx_deliveryservice_id;
DROP INDEX idx_deliveryservice_xml_id;
ALTER TABLE all_delivery_services DROP COLUMN original;
ALTER TABLE all_delivery_services DROP CONSTRAINT all_delivery_services_pkey;
ALTER TABLE all_delivery_services ADD PRIMARY KEY(id);
ALTER TABLE all_delivery_services RENAME TO deliveryservice;
CREATE UNIQUE INDEX idx_deliveryservice_xml_id ON deliveryservice USING btree(xml_id);

ALTER TABLE ONLY deliveryservice_regex ADD CONSTRAINT fk_deliveryservice_regex_to_deliveryservice FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY deliveryservice_server ADD CONSTRAINT fk_deliveryservice_server_to_deliveryservice FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY federation_deliveryservice ADD CONSTRAINT fk_federation_deliveryservice_to_deliveryservice FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY job ADD CONSTRAINT fk_job_to_deliveryservice FOREIGN KEY (job_deliveryservice) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY staticdnsentry ADD CONSTRAINT fk_staticdnsentry_to_deliveryservice FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY steering_target ADD CONSTRAINT fk_steering_target_to_deliveryservice FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY steering_target ADD CONSTRAINT fk_steering_target_target_to_deliveryservice FOREIGN KEY (target) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY deliveryservice_tmuser ADD CONSTRAINT fk_deliveryservice_tmuser_to_deliveryservice FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY origin ADD CONSTRAINT fk_origin_to_deliveryservice FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON DELETE CASCADE;
