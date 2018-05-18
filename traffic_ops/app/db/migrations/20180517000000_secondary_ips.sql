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

-- IPs
CREATE TABLE interface (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    server bigint NOT NULL,
    interface_name text NOT NULL,
    interface_mtu bigint DEFAULT '9000'::bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now(),
    CONSTRAINT fk_interface_server_id FOREIGN KEY (server) REFERENCES server(id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT unique_interface_server_interface_name UNIQUE (server, interface_name)
);

CREATE TABLE ip (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    server bigint NOT NULL,
    type bigint NOT NULL,
    interface bigint,
    ipv4 inet,
    ipv4_gateway inet,
    ipv6 inet,
    ipv6_gateway inet,
    last_updated timestamp with time zone DEFAULT now(),
    CONSTRAINT fk_ip_server_id FOREIGN KEY (server) REFERENCES server(id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_ip_type_id FOREIGN KEY (type) REFERENCES type(id),
    CONSTRAINT fk_ip_intf_id FOREIGN KEY (interface) REFERENCES interface(id) ON UPDATE CASCADE ON DELETE CASCADE
); 

ALTER TABLE deliveryservice_server 
    ADD COLUMN ip bigint,
    ADD CONSTRAINT fk_dss_ip_id FOREIGN KEY (ip) REFERENCES ip(id) ON UPDATE CASCADE ON DELETE CASCADE;

INSERT INTO type (name, description, use_in_table) VALUES ('IP_MANAGEMENT', 'Management IP', 'ip');
INSERT INTO type (name, description, use_in_table) VALUES ('IP_ILO', 'ILO IP', 'ip');
INSERT INTO type (name, description, use_in_table) VALUES ('IP_PRIMARY', 'Primary IP', 'ip');
INSERT INTO type (name, description, use_in_table) VALUES ('IP_SECONDARY', 'Secondary streaming IP', 'ip');

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DELETE FROM type WHERE use_in_table='ip';

ALTER TABLE deliveryservice_server 
    DROP CONSTRAINT fk_dss_ip_id,
    DROP COLUMN ip;

DROP TABLE ip;

DROP TABLE interface;
