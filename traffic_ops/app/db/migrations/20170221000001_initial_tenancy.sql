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


-- tenant
CREATE TABLE tenant (
    id BIGSERIAL primary key NOT NULL,
    name text UNIQUE NOT NULL,
    parent_id bigint DEFAULT 1 CHECK (id != parent_id),
    CONSTRAINT fk_parentid FOREIGN KEY (parent_id) REFERENCES tenant(id),    
    last_updated timestamp with time zone DEFAULT now()
);
CREATE INDEX idx_k_tenant_parent_tenant_idx ON tenant USING btree (parent_id);

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON tenant FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

ALTER TABLE tm_user
    ADD tenant_id BIGINT,
    ADD CONSTRAINT fk_tenantid FOREIGN KEY (tenant_id) REFERENCES tenant (id) MATCH FULL,
    ALTER COLUMN tenant_id SET DEFAULT NULL;
CREATE INDEX idx_k_tm_user_tenant_idx ON tm_user USING btree (tenant_id);

ALTER TABLE cdn
    ADD tenant_id BIGINT,
    ADD CONSTRAINT fk_tenantid FOREIGN KEY (tenant_id) REFERENCES tenant (id) MATCH FULL,
    ALTER COLUMN tenant_id SET DEFAULT NULL;
CREATE INDEX idx_k_cdn_tenant_idx ON cdn USING btree (tenant_id);

ALTER TABLE deliveryservice
    ADD tenant_id BIGINT,
    ADD CONSTRAINT fk_tenantid FOREIGN KEY (tenant_id) REFERENCES tenant (id) MATCH FULL,
    ALTER COLUMN tenant_id SET DEFAULT NULL;
CREATE INDEX idx_k_deliveryservice_tenant_idx ON deliveryservice USING btree (tenant_id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE deliveryservice
DROP COLUMN tenant_id;

ALTER TABLE cdn
DROP COLUMN tenant_id;

ALTER TABLE tm_user
DROP COLUMN tenant_id;

DROP TRIGGER on_update_current_timestamp ON tenant;

DROP TABLE tenant;
