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
ALTER TABLE interface ADD COLUMN router_host_name text NOT NULL DEFAULT '';
ALTER TABLE interface ADD COLUMN router_port_name text NOT NULL DEFAULT '';

UPDATE interface
SET router_host_name = COALESCE(server.router_host_name, ''),
router_port_name = COALESCE(server.router_port_name, '')
FROM server
WHERE server = server.id;

ALTER TABLE server DROP COLUMN router_host_name;
ALTER TABLE server DROP COLUMN router_port_name;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE server ADD COLUMN router_host_name text DEFAULT '';
ALTER TABLE server ADD COLUMN router_port_name text DEFAULT '';

UPDATE server s SET (router_host_name, router_port_name) =
(SELECT interface.router_host_name, interface.router_port_name FROM interface
JOIN ip_address ip ON ip.interface = name
JOIN server on ip.server = server.id
WHERE ip.service_address = true
AND s.id = ip.server AND s.id = interface.server LIMIT 1);

ALTER TABLE interface DROP COLUMN router_host_name;
ALTER TABLE interface DROP COLUMN router_port_name;
