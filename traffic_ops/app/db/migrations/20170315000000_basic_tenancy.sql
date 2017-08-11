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
    active boolean NOT NULL DEFAULT false,
    parent_id bigint DEFAULT 1 CHECK (id != parent_id),
    CONSTRAINT fk_parentid FOREIGN KEY (parent_id) REFERENCES tenant(id),    
    last_updated timestamp with time zone DEFAULT now()
); 
CREATE INDEX idx_k_tenant_parent_tenant_idx ON tenant USING btree (parent_id);

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON tenant FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TRIGGER on_update_current_timestamp ON tenant;

DROP TABLE tenant;
