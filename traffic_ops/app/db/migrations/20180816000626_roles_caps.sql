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

-- Using role 'read-only'

INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'auth' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;

INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'api-endpoints-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'asns-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'cache-config-files-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'cache-groups-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'capabilities-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'cdns-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'cdn-security-keys-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'change-logs-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'coordinates-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'delivery-services-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'delivery-service-requests-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'delivery-service-servers-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'divisions-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'to-extensions-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'federations-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'hwinfo-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'jobs-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'origins-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'params-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'phys-locations-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'profiles-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'regions-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'roles-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'servers-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'stats-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'statuses-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'static-dns-entries-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'steering-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'steering-targets-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'system-info-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'tenants-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'types-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'read-only'), 'users-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;

-- Using role 'federation'

INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'federation'), 'auth' WHERE EXISTS (SELECT id FROM role WHERE name = 'federation') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'federation'), 'federations-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'federation') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'federation'), 'federations-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'federation') ON CONFLICT DO NOTHING;

-- Using role 'operations'

INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'auth' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;

-- Includes 'readonly' endpoints:
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'api-endpoints-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'asns-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'cache-config-files-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'cache-groups-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'capabilities-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'cdns-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'cdn-security-keys-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'change-logs-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'coordinates-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'delivery-services-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'delivery-service-requests-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'delivery-service-servers-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'divisions-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'to-extensions-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'federations-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'hwinfo-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'jobs-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'origins-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'params-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'phys-locations-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'profiles-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'regions-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'roles-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'servers-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'stats-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'statuses-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'static-dns-entries-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'steering-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'steering-targets-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'system-info-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'tenants-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'types-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'users-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;

-- Explicitly require 'operations'
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'api-endpoints-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'asns-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'cache-groups-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'capabilities-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'cdns-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'cdns-snapshot' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'delivery-services-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'delivery-service-security-keys-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'delivery-service-servers-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'divisions-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'iso-generate' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
-- For 'params-write': related endpoints will check manually for 'parameters-secure-write', which currently doesn't exist yet. IMO params-write should also be renamed to parameters-write
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'params-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING; 
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'phys-locations-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'profiles-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'regions-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'roles-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'servers-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'stats-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'statuses-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'tenants-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'types-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'users-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;

-- Inherited endpoint from the 'privilege hierarchy' (really, just federations)
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'federations-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;

-- Outstanding capabilities that had to be thought about
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'coordinates-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'delivery-service-requests-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'to-extensions-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'jobs-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'steering-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'steering-targets-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'operations'), 'users-register' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;

-- Using role 'portal'

INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'auth' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;

-- Includes 'readonly' endpoints
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'api-endpoints-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'asns-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'cache-config-files-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'cache-groups-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'capabilities-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'cdns-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'cdn-security-keys-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'change-logs-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'coordinates-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'delivery-services-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'delivery-service-requests-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'delivery-service-servers-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'divisions-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'to-extensions-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'federations-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'hwinfo-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'jobs-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'origins-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'params-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'phys-locations-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'profiles-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'regions-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'roles-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'servers-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'stats-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'statuses-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'static-dns-entries-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'steering-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'steering-targets-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'system-info-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'tenants-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'types-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'users-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;

-- Explicitly requires 'portal'
-- none that we decided to keep

-- Outstanding capabilities that had to be thought about
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'delivery-service-requests-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'jobs-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'users-register' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;

-- Using role 'steering'

INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'steering'), 'auth' WHERE EXISTS (SELECT id FROM role WHERE name = 'steering') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'steering'), 'steering-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'steering') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'steering'), 'steering-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'steering') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'steering'), 'steering-targets-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'steering') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'steering'), 'steering-targets-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'steering') ON CONFLICT DO NOTHING;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DELETE FROM role_capability WHERE role_id = (SELECT id FROM role WHERE name = 'read-only');
DELETE FROM role_capability WHERE role_id = (SELECT id FROM role WHERE name = 'operations');
DELETE FROM role_capability WHERE role_id = (SELECT id FROM role WHERE name = 'federation');
DELETE FROM role_capability WHERE role_id = (SELECT id FROM role WHERE name = 'portal');
DELETE FROM role_capability WHERE role_id = (SELECT id FROM role WHERE name = 'steering');

