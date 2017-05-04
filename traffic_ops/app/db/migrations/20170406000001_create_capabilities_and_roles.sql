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


-- capability
CREATE TABLE capability (
    name text primary key UNIQUE NOT NULL,
    description text,
    last_updated timestamp with time zone DEFAULT now()
);

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON capability FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

-- http_method_t (enum)
CREATE TYPE http_method_t as ENUM ('GET', 'POST', 'PUT', 'PATCH', 'DELETE');

-- api_capability

CREATE TABLE api_capability (
    id BIGSERIAL primary key NOT NULL,
    http_method http_method_t NOT NULL,
    route text NOT NULL,
    capability text NOT NULL,
    CONSTRAINT fk_capability FOREIGN KEY (capability) REFERENCES capability(name) ON DELETE RESTRICT,
    last_updated timestamp with time zone DEFAULT now()
);

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON api_capability FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

-- role_capability
CREATE TABLE role_capability (
    role_id bigint NOT NULL,
    CONSTRAINT fk_role_id FOREIGN KEY (role_id) REFERENCES role(id) ON DELETE CASCADE,  
    cap_name text NOT NULL,
    CONSTRAINT fk_cap_name FOREIGN KEY (cap_name) REFERENCES capability(name) ON DELETE RESTRICT,
    last_updated timestamp with time zone DEFAULT now()
);

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON role_capability FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

-- user_role
CREATE TABLE user_role (
    user_id bigint NOT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES tm_user(id) ON DELETE CASCADE,
    role_id bigint NOT NULL,
    CONSTRAINT fk_role_id FOREIGN KEY (role_id) REFERENCES role(id) ON DELETE RESTRICT,
    last_updated timestamp with time zone DEFAULT now()
);

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON user_role FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back


DROP TRIGGER on_update_current_timestamp ON user_role;

DROP TABLE user_role;

DROP TRIGGER on_update_current_timestamp ON role_capability;

DROP TABLE role_capability;

DROP TRIGGER on_update_current_timestamp ON api_capability;

DROP TABLE api_capability;

DROP TRIGGER on_update_current_timestamp ON capability;

DROP TABLE capability;



