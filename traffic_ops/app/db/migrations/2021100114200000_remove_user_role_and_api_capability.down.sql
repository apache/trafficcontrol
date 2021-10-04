/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with this
 * work for additional information regarding copyright ownership.  The ASF
 * licenses this file to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  See the
 * License for the specific language governing permissions and limitations under
 * the License.
 */

CREATE TABLE IF NOT EXISTS user_role (
user_id bigint NOT NULL,
role_id bigint NOT NULL,
last_updated timestamp with time zone NOT NULL DEFAULT now()
);

ALTER TABLE user_role OWNER TO traffic_ops;

CREATE TABLE IF NOT EXISTS api_capability (
id bigserial PRIMARY KEY,
http_method http_method_t NOT NULL,
route text NOT NULL,
capability text NOT NULL,
last_updated timestamp with time zone NOT NULL DEFAULT now(),
UNIQUE (http_method, route, capability)
);

ALTER TABLE api_capability OWNER TO traffic_ops;


CREATE OR REPLACE FUNCTION create_constraint_if_not_exists (c_name text, t_name text, constraint_string text)
RETURNS void AS
$$
BEGIN
    IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = c_name AND table_name = t_name) then execute constraint_string;
END IF;
END;
$$ LANGUAGE PLPGSQL;

SELECT create_constraint_if_not_exists('fk_user_id', 'user_role', 'ALTER TABLE ONLY user_role ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES tm_user (id) ON DELETE CASCADE;');
SELECT create_constraint_if_not_exists('fk_role_id', 'user_role', 'ALTER TABLE ONLY user_role ADD CONSTRAINT fk_role_id FOREIGN KEY (role_id) REFERENCES role (id) ON DELETE RESTRICT;');
SELECT create_constraint_if_not_exists('fk_capability', 'api_capability', 'ALTER TABLE ONLY api_capability ADD CONSTRAINT fk_capability FOREIGN KEY (capability) REFERENCES capability (name) ON DELETE RESTRICT;');