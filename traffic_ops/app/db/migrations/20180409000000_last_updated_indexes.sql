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

CREATE INDEX idx_k_cachegroup_last_updated_idx ON cachegroup USING btree (last_updated);
CREATE INDEX idx_k_cdn_last_updated_idx ON cdn USING btree (last_updated);
CREATE INDEX idx_k_deliveryservice_last_updated_idx ON deliveryservice USING btree (last_updated);
CREATE INDEX idx_k_deliveryservice_regex_last_updated_idx ON deliveryservice_regex USING btree (last_updated);
CREATE INDEX idx_k_parameter_last_updated_idx ON parameter USING btree (last_updated);
CREATE INDEX idx_k_profile_last_updated_idx ON profile USING btree (last_updated);
CREATE INDEX idx_k_profile_parameter_last_updated_idx ON profile_parameter USING btree (last_updated);
CREATE INDEX idx_k_regex_last_updated_idx ON regex USING btree (last_updated);
CREATE INDEX idx_k_server_last_updated_idx ON server USING btree (last_updated);
CREATE INDEX idx_k_staticdnsentry_last_updated_idx ON staticdnsentry USING btree (last_updated);
CREATE INDEX idx_k_type_last_updated_idx ON type USING btree (last_updated);
CREATE INDEX idx_k_status_last_updated_idx ON status USING btree (last_updated);
CREATE INDEX idx_k_snapshot_last_updated_idx ON snapshot USING btree (last_updated);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP INDEX idx_k_cachegroup_last_updated_idx;
DROP INDEX idx_k_cdn_last_updated_idx;
DROP INDEX idx_k_deliveryservice_last_updated_idx;
DROP INDEX idx_k_deliveryservice_regex_last_updated_idx;
DROP INDEX idx_k_parameter_last_updated_idx;
DROP INDEX idx_k_profile_last_updated_idx;
DROP INDEX idx_k_profile_parameter_last_updated_idx;
DROP INDEX idx_k_regex_last_updated_idx;
DROP INDEX idx_k_server_last_updated_idx;
DROP INDEX idx_k_staticdnsentry_last_updated_idx;
DROP INDEX idx_k_type_last_updated_idx;
DROP INDEX idx_k_status_last_updated_idx;
DROP INDEX idx_k_snapshot_last_updated_idx;
