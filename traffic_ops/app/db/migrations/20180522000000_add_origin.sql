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

CREATE TYPE origin_protocol AS ENUM ('http', 'https');

CREATE TABLE origin (
    id bigserial PRIMARY KEY NOT NULL,
    name text UNIQUE NOT NULL,
    fqdn text NOT NULL,
    protocol origin_protocol NOT NULL DEFAULT 'http',
    is_primary boolean NOT NULL DEFAULT FALSE,
    port bigint,
    ip_address text,
    ip6_address text,
    deliveryservice bigint NOT NULL REFERENCES deliveryservice(id) ON DELETE CASCADE,
    coordinate bigint REFERENCES coordinate(id) ON DELETE RESTRICT,
    profile bigint REFERENCES profile(id) ON DELETE RESTRICT,
    cachegroup bigint REFERENCES cachegroup(id) ON DELETE RESTRICT,
    tenant bigint REFERENCES tenant(id) ON DELETE RESTRICT,
    last_updated timestamp WITH time zone NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX origin_is_primary_deliveryservice_constraint ON origin (is_primary, deliveryservice) WHERE is_primary;

CREATE INDEX origin_deliveryservice_fkey ON origin USING btree (deliveryservice);
CREATE INDEX origin_coordinate_fkey ON origin USING btree (coordinate);
CREATE INDEX origin_profile_fkey ON origin USING btree (profile);
CREATE INDEX origin_cachegroup_fkey ON origin USING btree (cachegroup);
CREATE INDEX origin_tenant_fkey ON origin USING btree (tenant);

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON origin FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS origin;

DROP TYPE IF EXISTS origin_protocol;
