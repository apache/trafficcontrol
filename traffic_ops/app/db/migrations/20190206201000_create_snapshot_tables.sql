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
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE cachegroup_snapshot ( LIKE cachegroup );
ALTER TABLE  cachegroup_snapshot ADD COLUMN deleted boolean NOT NULL default false;
ALTER TABLE  cachegroup_snapshot ADD PRIMARY KEY (id, type, last_updated);
CREATE INDEX idx_k_cachegroup_snapshot_deleted_idx ON cachegroup_snapshot USING btree (deleted);
CREATE INDEX idx_k_cachegroup_snapshot_last_updated_idx ON cachegroup_snapshot USING btree (last_updated);

CREATE TABLE cachegroup_fallbacks_snapshot ( LIKE cachegroup_fallbacks );
ALTER TABLE  cachegroup_fallbacks_snapshot ADD COLUMN deleted boolean NOT NULL default false;
ALTER TABLE  cachegroup_fallbacks_snapshot ADD COLUMN last_updated timestamp with time zone NOT NULL default now();
ALTER TABLE  cachegroup_fallbacks_snapshot ADD PRIMARY KEY (primary_cg, backup_cg, last_updated);
CREATE INDEX idx_k_cachegroup_fallbacks_snapshot_deleted_idx ON cachegroup_fallbacks_snapshot USING btree (deleted);
CREATE INDEX idx_k_cachegroup_fallbacks_snapshot_last_updated_idx ON cachegroup_fallbacks_snapshot USING btree (last_updated);

CREATE TABLE cachegroup_localization_method_snapshot ( LIKE cachegroup_localization_method);
ALTER TABLE  cachegroup_localization_method_snapshot ADD COLUMN deleted boolean NOT NULL default false;
ALTER TABLE  cachegroup_localization_method_snapshot ADD COLUMN last_updated timestamp with time zone NOT NULL default now();
ALTER TABLE  cachegroup_localization_method_snapshot ADD PRIMARY KEY (cachegroup, method, last_updated);
CREATE INDEX idx_k_cachegroup_localization_method_snapshot_deleted_idx ON cachegroup_localization_method_snapshot USING btree (deleted);
CREATE INDEX idx_k_cachegroup_localization_method_snapshot_last_updated_idx ON cachegroup_localization_method_snapshot USING btree (last_updated);

CREATE TABLE cdn_snapshot ( LIKE cdn );
ALTER TABLE  cdn_snapshot ADD COLUMN deleted boolean NOT NULL default false;
ALTER TABLE  cdn_snapshot ADD PRIMARY KEY (name, last_updated);
CREATE INDEX idx_k_cdn_snapshot_deleted_idx ON cdn_snapshot USING btree (deleted);
CREATE INDEX idx_k_cdn_snapshot_last_updated_idx ON cdn_snapshot USING btree (last_updated);

CREATE TABLE coordinate_snapshot ( LIKE coordinate );
ALTER TABLE  coordinate_snapshot ADD COLUMN deleted boolean NOT NULL default false;
ALTER TABLE  coordinate_snapshot ADD PRIMARY KEY (id, last_updated);
CREATE INDEX idx_k_coordinate_snapshot_deleted_idx ON coordinate_snapshot USING btree (deleted);
CREATE INDEX idx_k_coordinate_snapshot_last_updated_idx ON coordinate_snapshot USING btree (last_updated);

CREATE TABLE deliveryservice_snapshot ( LIKE deliveryservice );
ALTER TABLE  deliveryservice_snapshot ADD COLUMN deleted boolean NOT NULL default false;
ALTER TABLE  deliveryservice_snapshot ADD PRIMARY KEY (xml_id, last_updated);
CREATE INDEX idx_k_deliveryservice_snapshot_deleted_idx ON deliveryservice_snapshot USING btree (deleted);
CREATE INDEX idx_k_deliveryservice_snapshot_last_updated_idx ON deliveryservice_snapshot USING btree (last_updated);

CREATE TABLE deliveryservice_regex_snapshot ( LIKE deliveryservice_regex );
ALTER TABLE  deliveryservice_regex_snapshot ADD COLUMN deleted boolean NOT NULL default false;
ALTER TABLE  deliveryservice_regex_snapshot ADD PRIMARY KEY (deliveryservice, regex, last_updated);
CREATE INDEX idx_k_deliveryservice_regex_snapshot_deleted_idx ON deliveryservice_regex_snapshot USING btree (deleted);
CREATE INDEX idx_k_deliveryservice_regex_snapshot_last_updated_idx ON deliveryservice_regex_snapshot USING btree (last_updated);

CREATE TABLE deliveryservice_server_snapshot ( LIKE deliveryservice_server );
ALTER TABLE  deliveryservice_server_snapshot ADD COLUMN deleted boolean NOT NULL default false;
ALTER TABLE  deliveryservice_server_snapshot ADD PRIMARY KEY (deliveryservice, server, last_updated);
CREATE INDEX idx_k_deliveryservice_server_snapshot_deleted_idx ON deliveryservice_server_snapshot USING btree (deleted);
CREATE INDEX idx_k_deliveryservice_server_snapshot_last_updated_idx ON deliveryservice_server_snapshot USING btree (last_updated);

CREATE TABLE parameter_snapshot ( LIKE parameter );
ALTER TABLE  parameter_snapshot ADD COLUMN deleted boolean NOT NULL default false;
ALTER TABLE  parameter_snapshot ADD PRIMARY KEY (name, config_file, value, last_updated);
CREATE INDEX idx_k_parameter_snapshot_deleted_idx ON parameter_snapshot USING btree (deleted);
CREATE INDEX idx_k_parameter_snapshot_last_updated_idx ON parameter_snapshot USING btree (last_updated);

CREATE TABLE profile_snapshot ( LIKE profile );
ALTER TABLE  profile_snapshot ADD COLUMN deleted boolean NOT NULL default false;
ALTER TABLE  profile_snapshot ADD PRIMARY KEY (name, last_updated);
CREATE INDEX idx_k_profile_snapshot_deleted_idx ON profile_snapshot USING btree (deleted);
CREATE INDEX idx_k_profile_snapshot_last_updated_idx ON profile_snapshot USING btree (last_updated);

CREATE TABLE profile_parameter_snapshot ( LIKE profile_parameter );
ALTER TABLE  profile_parameter_snapshot ADD COLUMN deleted boolean NOT NULL default false;
ALTER TABLE  profile_parameter_snapshot ADD PRIMARY KEY (profile, parameter, last_updated);
CREATE INDEX idx_k_profile_parameter_snapshot_deleted_idx ON profile_parameter_snapshot USING btree (deleted);
CREATE INDEX idx_k_profile_parameter_snapshot_last_updated_idx ON profile_parameter_snapshot USING btree (last_updated);

CREATE TABLE regex_snapshot ( LIKE regex );
ALTER TABLE  regex_snapshot ADD COLUMN deleted boolean NOT NULL default false;
ALTER TABLE  regex_snapshot ADD PRIMARY KEY (id, last_updated);
CREATE INDEX idx_k_regex_snapshot_deleted_idx ON regex_snapshot USING btree (deleted);
CREATE INDEX idx_k_regex_snapshot_last_updated_idx ON regex_snapshot USING btree (last_updated);

CREATE TABLE server_snapshot ( LIKE server );
ALTER TABLE  server_snapshot ADD COLUMN deleted boolean NOT NULL default false;
ALTER TABLE  server_snapshot ADD PRIMARY KEY (ip_address, profile, last_updated);
CREATE INDEX idx_k_server_snapshot_deleted_idx ON server_snapshot USING btree (deleted);
CREATE INDEX idx_k_server_snapshot_last_updated_idx ON server_snapshot USING btree (last_updated);

CREATE TABLE staticdnsentry_snapshot ( LIKE staticdnsentry );
ALTER TABLE  staticdnsentry_snapshot ADD COLUMN deleted boolean NOT NULL default false;
ALTER TABLE  staticdnsentry_snapshot ADD PRIMARY KEY (host, address, deliveryservice, cachegroup, last_updated);
CREATE INDEX idx_k_staticdnsentry_snapshot_deleted_idx ON staticdnsentry_snapshot USING btree (deleted);
CREATE INDEX idx_k_staticdnsentry_snapshot_last_updated_idx ON staticdnsentry_snapshot USING btree (last_updated);

CREATE TABLE status_snapshot ( LIKE status );
ALTER TABLE  status_snapshot ADD COLUMN deleted boolean NOT NULL default false;
ALTER TABLE  status_snapshot ADD PRIMARY KEY (name, last_updated);
CREATE INDEX idx_k_status_snapshot_deleted_idx ON status_snapshot USING btree (deleted);
CREATE INDEX idx_k_status_snapshot_last_updated_idx ON status_snapshot USING btree (last_updated);

CREATE TABLE type_snapshot ( LIKE type );
ALTER TABLE  type_snapshot ADD COLUMN deleted boolean NOT NULL default false;
ALTER TABLE  type_snapshot ADD PRIMARY KEY (name, last_updated);
CREATE INDEX idx_k_type_snapshot_deleted_idx ON type_snapshot USING btree (deleted);
CREATE INDEX idx_k_type_snapshot_last_updated_idx ON type_snapshot USING btree (last_updated);

CREATE INDEX idx_k_cdn_snapshot_name_idx ON cdn_snapshot USING btree (name);

CREATE INDEX idx_k_cachegroup_snapshot_id_idx ON cachegroup_snapshot USING btree (id);
CREATE INDEX idx_k_cachegroup_snapshot_coordinate_idx ON cachegroup_snapshot USING btree (coordinate);

CREATE INDEX idx_k_cachegroup_fallbacks_snapshot_backup_cg_idx ON cachegroup_fallbacks_snapshot USING btree (backup_cg);
CREATE INDEX idx_k_cachegroup_fallbacks_snapshot_set_order_idx ON cachegroup_fallbacks_snapshot USING btree (set_order);

CREATE INDEX idx_k_coordinate_snapshot_id_idx ON coordinate_snapshot USING btree (id);

CREATE INDEX idx_k_deliveryservice_snapshot_id_idx ON deliveryservice_snapshot USING btree (id);
CREATE INDEX idx_k_deliveryservice_snapshot_cdn_id_idx ON deliveryservice_snapshot USING btree (cdn_id);
CREATE INDEX idx_k_deliveryservice_snapshot_active_idx ON deliveryservice_snapshot USING btree (active);
CREATE INDEX idx_k_deliveryservice_snapshot_type_idx ON deliveryservice_snapshot USING btree (type);
CREATE INDEX idx_k_deliveryservice_snapshot_profile_idx ON deliveryservice_snapshot USING btree (profile);

CREATE INDEX idx_k_deliveryservice_regex_snapshot_regex_idx ON deliveryservice_regex_snapshot USING btree (regex);
CREATE INDEX idx_k_deliveryservice_regex_snapshot_deliveryservice_idx ON deliveryservice_regex_snapshot USING btree (deliveryservice);
CREATE INDEX idx_k_deliveryservice_regex_snapshot_set_number_idx ON deliveryservice_regex_snapshot USING btree (set_number);

CREATE INDEX idx_k_deliveryservice_server_snapshot_server_idx ON deliveryservice_server_snapshot USING btree (server);
CREATE INDEX idx_k_deliveryservice_server_snapshot_deliveryservice_idx ON deliveryservice_server_snapshot USING btree (deliveryservice);

CREATE INDEX idx_k_parameter_snapshot_id_idx ON parameter_snapshot USING btree (id);
CREATE INDEX idx_k_parameter_snapshot_config_file_idx ON parameter_snapshot USING btree (config_file);
CREATE INDEX idx_k_parameter_snapshot_name_idx ON parameter_snapshot USING btree (name);

CREATE INDEX idx_k_profile_snapshot_id_idx ON profile_snapshot USING btree (id);
CREATE INDEX idx_k_profile_snapshot_routing_disabled_idx ON profile_snapshot USING btree (routing_disabled);

CREATE INDEX idx_k_profile_parameter_snapshot_profile_idx ON profile_parameter_snapshot USING btree (profile);
CREATE INDEX idx_k_profile_parameter_snapshot_parameter_idx ON profile_parameter_snapshot USING btree (parameter);

CREATE INDEX idx_k_regex_snapshot_id_idx ON regex_snapshot USING btree (id);
CREATE INDEX idx_k_regex_snapshot_type_idx ON regex_snapshot USING btree (type);

CREATE INDEX idx_k_server_snapshot_id_idx ON server_snapshot USING btree (id);
CREATE INDEX idx_k_server_snapshot_cdn_id_idx ON server_snapshot USING btree (cdn_id);
CREATE INDEX idx_k_server_snapshot_cachegroup_idx ON server_snapshot USING btree (cachegroup);
CREATE INDEX idx_k_server_snapshot_type_idx ON server_snapshot USING btree (type);
CREATE INDEX idx_k_server_snapshot_status_idx ON server_snapshot USING btree (status);
CREATE INDEX idx_k_server_snapshot_profile_idx ON server_snapshot USING btree (profile);

CREATE INDEX idx_k_status_snapshot_id_idx ON status_snapshot USING btree (id);
CREATE INDEX idx_k_status_snapshot_name_idx ON status_snapshot USING btree (name);

CREATE INDEX idx_k_staticdnsentry_snapshot_deliveryservice_idx ON staticdnsentry_snapshot USING btree (deliveryservice);
CREATE INDEX idx_k_staticdnsentry_snapshot_type_idx ON staticdnsentry_snapshot USING btree (type);

CREATE INDEX idx_k_type_snapshot_id_idx ON type_snapshot USING btree (id);
CREATE INDEX idx_k_type_snapshot_name_idx ON type_snapshot USING btree (name);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE cachegroup_snapshot;
DROP TABLE cachegroup_fallbacks_snapshot;
DROP TABLE cachegroup_localization_method_snapshot;
DROP TABLE cdn_snapshot;
DROP TABLE coordinate_snapshot;
DROP TABLE deliveryservice_snapshot;
DROP TABLE deliveryservice_regex_snapshot;
DROP TABLE deliveryservice_server_snapshot;
DROP TABLE parameter_snapshot;
DROP TABLE profile_snapshot;
DROP TABLE profile_parameter_snapshot;
DROP TABLE regex_snapshot;
DROP TABLE server_snapshot;
DROP TABLE staticdnsentry_snapshot;
DROP TABLE status_snapshot;
DROP TABLE type_snapshot;
