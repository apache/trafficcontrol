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

-- cacheassignmentgroup
CREATE TABLE cacheassignmentgroup (
    id bigserial primary key NOT NULL,
    name text UNIQUE NOT NULL,
    description text NOT NULL,
    cdn_id bigint NOT NULL,
    last_updated timestamp WITH time zone NOT NULL DEFAULT now()
);

ALTER TABLE cacheassignmentgroup
    ADD CONSTRAINT fk_cdnid FOREIGN KEY (cdn_id) REFERENCES cdn(id) ON DELETE CASCADE;

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON cacheassignmentgroup FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

-- cacheassignmentgroup_server
CREATE TABLE cacheassignmentgroup_server (
    cacheassignmentgroup bigint NOT NULL,
    server bigint NOT NULL,
    last_updated timestamp WITH time zone NOT NULL DEFAULT now()    
);

ALTER TABLE cacheassignmentgroup_server
    ADD CONSTRAINT idx_cag_srv_primary PRIMARY KEY (cacheassignmentgroup, server),
    ADD CONSTRAINT fk_cacheassignmentgroupid FOREIGN KEY (cacheassignmentgroup) REFERENCES cacheassignmentgroup(id) ON UPDATE CASCADE ON DELETE CASCADE,
    ADD CONSTRAINT fk_serverid FOREIGN KEY (server) REFERENCES server(id) ON UPDATE CASCADE ON DELETE CASCADE;

-- deliveryservice_cacheassignmentgroup
CREATE TABLE deliveryservice_cacheassignmentgroup (
    deliveryservice bigint NOT NULL,
    cacheassignmentgroup bigint NOT NULL, 
    last_updated timestamp WITH time zone NOT NULL DEFAULT now()   
);

ALTER TABLE deliveryservice_cacheassignmentgroup
    ADD CONSTRAINT idx_ds_cag_primary PRIMARY KEY (deliveryservice, cacheassignmentgroup),
    ADD CONSTRAINT fk_deliveryserviceid FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE,
    ADD CONSTRAINT fk_cacheassignmentgroupid FOREIGN KEY (cacheassignmentgroup) REFERENCES cacheassignmentgroup(id) ON UPDATE CASCADE ON DELETE CASCADE;

CREATE VIEW deliveryservice_assignedservers AS
  SELECT
    server.id AS server,
    ds_cag.deliveryservice AS deliveryservice,
    GREATEST(cags.last_updated, ds_cag.last_updated) AS last_updated
  FROM server
  INNER JOIN cacheassignmentgroup_server AS cags
    ON server.id=cags.server
  INNER JOIN deliveryservice_cacheassignmentgroup AS ds_cag
    ON ds_cag.cacheassignmentgroup = cags.cacheassignmentgroup
  UNION
    SELECT server, deliveryservice, last_updated  FROM deliveryservice_server;
    
-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE cacheassignmentgroup_server;
DROP TABLE deliveryservice_cacheassignmentgroup;
DROP TRIGGER IF EXISTS on_update_current_timestamp ON cacheassignmentgroup;
DROP TABLE cacheassignmentgroup;

