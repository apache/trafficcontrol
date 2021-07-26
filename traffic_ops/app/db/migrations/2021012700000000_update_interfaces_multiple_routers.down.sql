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
