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


-- THIS FILE INCLUDES POST-MIGRATION DATA FIXES REQUIRED OF TRAFFIC OPS

-- Mapping roles and capabilities: (used to be in a migration)

-- For role 'federation'
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'federation'), 'auth' WHERE EXISTS (SELECT id FROM role WHERE name = 'federation') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'federation'), 'federations-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'federation') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'federation'), 'federations-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'federation') ON CONFLICT DO NOTHING;

-- For role 'portal'
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'auth' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'api-endpoints-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'asns-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'cache-config-files-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'cache-groups-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'capabilities-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'cdns-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'cdn-security-keys-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'change-logs-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'consistenthash-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
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
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'parameters-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
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
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'delivery-service-requests-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'jobs-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'users-register' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;

-- For role 'steering'
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'steering'), 'auth' WHERE EXISTS (SELECT id FROM role WHERE name = 'steering') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'steering'), 'steering-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'steering') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'steering'), 'steering-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'steering') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'steering'), 'steering-targets-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'steering') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'steering'), 'steering-targets-write' WHERE EXISTS (SELECT id FROM role WHERE name = 'steering') ON CONFLICT DO NOTHING;
