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

UPDATE api_capability SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE api_capability ALTER COLUMN last_updated SET NOT NULL;

UPDATE asn SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE asn ALTER COLUMN last_updated SET NOT NULL;

UPDATE cachegroup SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE cachegroup ALTER COLUMN last_updated SET NOT NULL;

UPDATE cachegroup_parameter SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE cachegroup_parameter ALTER COLUMN last_updated SET NOT NULL;

UPDATE capability SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE capability ALTER COLUMN last_updated SET NOT NULL;

UPDATE deliveryservice SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE deliveryservice ALTER COLUMN last_updated SET NOT NULL;

ALTER TABLE deliveryservice_regex ADD COLUMN last_updated timestamp with time zone NOT NULL DEFAULT now();
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON deliveryservice_regex FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

UPDATE deliveryservice_server SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE deliveryservice_server ALTER COLUMN last_updated SET NOT NULL;

UPDATE deliveryservice_tmuser SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE deliveryservice_tmuser ALTER COLUMN last_updated SET NOT NULL;

UPDATE division SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE division ALTER COLUMN last_updated SET NOT NULL;

UPDATE federation SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE federation ALTER COLUMN last_updated SET NOT NULL;

UPDATE federation_deliveryservice SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE federation_deliveryservice ALTER COLUMN last_updated SET NOT NULL;

UPDATE federation_federation_resolver SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE federation_federation_resolver ALTER COLUMN last_updated SET NOT NULL;

UPDATE federation_resolver SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE federation_resolver ALTER COLUMN last_updated SET NOT NULL;

UPDATE federation_tmuser SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE federation_tmuser ALTER COLUMN last_updated SET NOT NULL;

UPDATE hwinfo SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE hwinfo ALTER COLUMN last_updated SET NOT NULL;

UPDATE job SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE job ALTER COLUMN last_updated SET NOT NULL;

UPDATE job_agent SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE job_agent ALTER COLUMN last_updated SET NOT NULL;

UPDATE job_result SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE job_result ALTER COLUMN last_updated SET NOT NULL;

UPDATE job_status SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE job_status ALTER COLUMN last_updated SET NOT NULL;

UPDATE log SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE log ALTER COLUMN last_updated SET NOT NULL;

UPDATE parameter SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE parameter ALTER COLUMN last_updated SET NOT NULL;

UPDATE phys_location SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE phys_location ALTER COLUMN last_updated SET NOT NULL;

UPDATE profile SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE profile ALTER COLUMN last_updated SET NOT NULL;

UPDATE profile_parameter SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE profile_parameter ALTER COLUMN last_updated SET NOT NULL;

UPDATE regex SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE regex ALTER COLUMN last_updated SET NOT NULL;

UPDATE region SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE region ALTER COLUMN last_updated SET NOT NULL;

ALTER TABLE role ADD COLUMN last_updated timestamp with time zone NOT NULL DEFAULT now();
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON role FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

UPDATE role_capability SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE role_capability ALTER COLUMN last_updated SET NOT NULL;

UPDATE server SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE server ALTER COLUMN last_updated SET NOT NULL;

UPDATE servercheck SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE servercheck ALTER COLUMN last_updated SET NOT NULL;

UPDATE snapshot SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE snapshot ALTER COLUMN last_updated SET NOT NULL;

UPDATE staticdnsentry SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE staticdnsentry ALTER COLUMN last_updated SET NOT NULL;

UPDATE status SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE status ALTER COLUMN last_updated SET NOT NULL;

UPDATE steering_target SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE steering_target ALTER COLUMN last_updated SET NOT NULL;

UPDATE tenant SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE tenant ALTER COLUMN last_updated SET NOT NULL;

UPDATE tm_user SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE tm_user ALTER COLUMN last_updated SET NOT NULL;

UPDATE type SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE type ALTER COLUMN last_updated SET NOT NULL;

UPDATE user_role SET last_updated = now() WHERE last_updated IS NULL;
ALTER TABLE user_role ALTER COLUMN last_updated SET NOT NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TRIGGER on_update_current_timestamp ON deliveryservice_regex;
DROP TRIGGER on_update_current_timestamp ON role;
ALTER TABLE api_capability ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE asn ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE cachegroup ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE cachegroup_parameter ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE capability ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE deliveryservice ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE deliveryservice_server ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE deliveryservice_tmuser ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE division ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE federation ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE federation_deliveryservice ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE federation_federation_resolver ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE federation_resolver ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE federation_tmuser ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE hwinfo ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE job ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE job_agent ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE job_result ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE job_status ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE log ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE parameter ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE phys_location ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE profile ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE profile_parameter ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE regex ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE region ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE role_capability ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE server ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE servercheck ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE snapshot ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE staticdnsentry ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE status ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE steering_target ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE tenant ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE tm_user ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE type ALTER COLUMN last_updated DROP NOT NULL;
ALTER TABLE user_role ALTER COLUMN last_updated DROP NOT NULL;
