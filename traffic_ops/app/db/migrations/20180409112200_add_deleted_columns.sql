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

ALTER TABLE cachegroup ADD COLUMN deleted boolean NOT NULL DEFAULT false;
CREATE INDEX idx_k_cachegroup_deleted_idx ON cachegroup USING btree (deleted);
ALTER TABLE cdn ADD COLUMN deleted boolean NOT NULL DEFAULT false;
CREATE INDEX idx_k_cdn_deleted_idx ON cdn USING btree (deleted);
ALTER TABLE deliveryservice ADD COLUMN deleted boolean NOT NULL DEFAULT false;
CREATE INDEX idx_k_deliveryservice_deleted_idx ON deliveryservice USING btree (deleted);
ALTER TABLE deliveryservice_regex ADD COLUMN deleted boolean NOT NULL DEFAULT false;
CREATE INDEX idx_k_deliveryservice_regex_deleted_idx ON deliveryservice_regex USING btree (deleted);
ALTER TABLE parameter ADD COLUMN deleted boolean NOT NULL DEFAULT false;
CREATE INDEX idx_k_parameter_deleted_idx ON parameter USING btree (deleted);
ALTER TABLE profile ADD COLUMN deleted boolean NOT NULL DEFAULT false;
CREATE INDEX idx_k_profile_deleted_idx ON profile USING btree (deleted);
ALTER TABLE profile_parameter ADD COLUMN deleted boolean NOT NULL DEFAULT false;
CREATE INDEX idx_k_profile_parameter_deleted_idx ON profile_parameter USING btree (deleted);
ALTER TABLE regex ADD COLUMN deleted boolean NOT NULL DEFAULT false;
CREATE INDEX idx_k_regex_deleted_idx ON regex USING btree (deleted);
ALTER TABLE server ADD COLUMN deleted boolean NOT NULL DEFAULT false;
CREATE INDEX idx_k_server_deleted_idx ON server USING btree (deleted);
ALTER TABLE staticdnsentry ADD COLUMN deleted boolean NOT NULL DEFAULT false;
CREATE INDEX idx_k_staticdnsentry_deleted_idx ON staticdnsentry USING btree (deleted);
ALTER TABLE type ADD COLUMN deleted boolean NOT NULL DEFAULT false;
CREATE INDEX idx_k_type_deleted_idx ON type USING btree (deleted);
ALTER TABLE status ADD COLUMN deleted boolean NOT NULL DEFAULT false;
CREATE INDEX idx_k_status_deleted_idx ON status USING btree (deleted);
ALTER TABLE snapshot ADD COLUMN deleted boolean NOT NULL DEFAULT false;
CREATE INDEX idx_k_snapshot_deleted_idx ON snapshot USING btree (deleted);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP INDEX idx_k_cachegroup_deleted_idx;
ALTER TABLE cachegroup DROP COLUMN deleted;
DROP INDEX idx_k_cdn_deleted_idx;
ALTER TABLE cdn DROP COLUMN deleted;
DROP INDEX idx_k_deliveryservice_deleted_idx;
ALTER TABLE deliveryservice DROP COLUMN deleted;
DROP INDEX idx_k_deliveryservice_regex_deleted_idx;
ALTER TABLE deliveryservice_regex DROP COLUMN deleted;
DROP INDEX idx_k_parameter_deleted_idx;
ALTER TABLE parameter DROP COLUMN deleted;
DROP INDEX idx_k_profile_deleted_idx;
ALTER TABLE profile DROP COLUMN deleted;
DROP INDEX idx_k_profile_parameter_deleted_idx;
ALTER TABLE profile_parameter DROP COLUMN deleted;
DROP INDEX idx_k_regex_deleted_idx;
ALTER TABLE regex DROP COLUMN deleted;
DROP INDEX idx_k_server_deleted_idx;
ALTER TABLE server DROP COLUMN deleted;
DROP INDEX idx_k_staticdnsentry_deleted_idx;
ALTER TABLE staticdnsentry DROP COLUMN deleted;
DROP INDEX idx_k_type_deleted_idx;
ALTER TABLE type DROP COLUMN deleted;
DROP INDEX idx_k_status_deleted_idx;
ALTER TABLE status DROP COLUMN deleted;
DROP INDEX idx_k_snapshot_deleted_idx;
ALTER TABLE snapshot DROP COLUMN deleted;
