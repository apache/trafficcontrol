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

-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE topology (
    name text PRIMARY KEY,
    description text NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL
);
DROP TRIGGER IF EXISTS on_update_current_timestamp ON topology;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON topology FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

CREATE TABLE topology_cachegroup (
    id BIGSERIAL PRIMARY KEY,
    topology text NOT NULL,
    cachegroup text NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT topology_cachegroup_cachegroup_fkey FOREIGN KEY (cachegroup) REFERENCES cachegroup(name) ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT topology_cachegroup_topology_fkey FOREIGN KEY (topology) REFERENCES topology(name) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT unique_topology_cachegroup UNIQUE (topology, cachegroup)
);
CREATE INDEX topology_cachegroup_cachegroup_fkey ON topology_cachegroup USING btree (cachegroup);
CREATE INDEX topology_cachegroup_topology_fkey ON topology_cachegroup USING btree (topology);
DROP TRIGGER IF EXISTS on_update_current_timestamp ON topology_cachegroup;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON topology_cachegroup FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

CREATE TABLE topology_cachegroup_parents (
    child bigint NOT NULL,
    parent bigint NOT NULL,
    rank integer NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT topology_cachegroup_parents_rank_check CHECK (rank = 1 OR rank = 2),
    CONSTRAINT topology_cachegroup_parents_child_fkey FOREIGN KEY (child) REFERENCES topology_cachegroup(id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT topology_cachegroup_parents_parent_fkey FOREIGN KEY (parent) REFERENCES topology_cachegroup(id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT unique_child_rank UNIQUE (child, rank),
    CONSTRAINT unique_child_parent UNIQUE (child, parent)
);
CREATE INDEX topology_cachegroup_parents_child_fkey ON topology_cachegroup_parents USING btree (child);
CREATE INDEX topology_cachegroup_parents_parents_fkey ON topology_cachegroup_parents USING btree (parent);
DROP TRIGGER IF EXISTS on_update_current_timestamp ON topology_cachegroup_parents;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON topology_cachegroup_parents FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

ALTER TABLE deliveryservice
    ADD COLUMN topology text,
    ADD CONSTRAINT deliveryservice_topology_fkey FOREIGN KEY (topology) REFERENCES topology (name) ON UPDATE CASCADE ON DELETE RESTRICT;
CREATE INDEX deliveryservice_topology_fkey ON deliveryservice USING btree (topology);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE deliveryservice DROP COLUMN topology;
DROP TABLE topology_cachegroup_parents;
DROP TABLE topology_cachegroup;
DROP TABLE topology;
