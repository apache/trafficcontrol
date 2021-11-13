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
INSERT INTO role_capability (role_id, cap_name)
SELECT id, perm
FROM public.role
CROSS JOIN ( VALUES
	('DELIVERY-SERVICE:READ'),
	('DELIVERY-SERVICE:UPDATE'),
	('FEDERATION:CREATE'),
	('FEDERATION:UPDATE'),
	('FEDERATION:DELETE'),
	('FEDERATION-RESOLVER:CREATE'),
	('FEDERATION-RESOLVER:DELETE')
) AS perms(perm)
WHERE "name" = 'federation'
ON CONFLICT DO NOTHING;

-- For role 'portal'
INSERT INTO role_capability (role_id, cap_name)
SELECT id, perm
FROM public.role
CROSS JOIN ( VALUES
	('ASN:READ'),
	('ASYNC-STATUS:READ'),
	('CACHE-GROUP:READ'),
	('CAPABILITY:READ'),
	('CDN-SNAPSHOT:READ'),
	('CDN:READ'),
	('COORDINATE:READ'),
	('DELIVERY-SERVICE:READ'),
	('DELIVERY-SERVICE:UPDATE'),
	('DIVISION:READ'),
	('DS-REQUEST:CREATE'),
	('DS-REQUEST:DELETE'),
	('DS-REQUEST:READ'),
	('DS-REQUEST:UPDATE'),
	('DS-SECURITY-KEY:READ'),
	('FEDERATION-RESOLVER:READ'),
	('FEDERATION:READ'),
	('ISO:READ'),
	('JOB:CREATE'),
	('JOB:DELETE'),
	('JOB:READ'),
	('JOB:UPDATE'),
	('LOG:READ'),
	('MONITOR-CONFIG:READ'),
	('ORIGIN:READ'),
	('PARAMETER:READ'),
	('PHYSICAL-LOCATION:READ'),
	('PLUGIN-READ'),
	('PROFILE:READ'),
	('REGION:READ'),
	('ROLE:READ'),
	('SERVER-CAPABILITY:READ'),
	('SERVER-CHECK:READ'),
	('SERVER:READ'),
	('SERVICE-CATEGORY:READ'),
	('STAT:CREATE'),
	('STAT:READ'),
	('STATIC-DN:READ'),
	('STATUS:READ'),
	('STEERING:CREATE'),
	('STEERING:DELETE'),
	('STEERING:READ'),
	('STEERING:UPDATE'),
	('TENANT:READ'),
	('TOPOLOGY:READ'),
	('TRAFFIC-VAULT:READ'),
	('TYPE:READ'),
	('USER:READ')
) AS perms(perm)
WHERE "name" = 'portal'
ON CONFLICT DO NOTHING;

-- For role 'steering'
INSERT INTO role_capability (role_id, cap_name)
SELECT id, perm
FROM public.role
CROSS JOIN ( VALUES
	('DELIVERY-SERVICE:READ'),
	('DELIVERY-SERVICE:UPDATE'),
	('STEERING:CREATE'),
	('STEERING:DELETE'),
	('STEERING:UPDATE')
) AS perms(perm)
WHERE "name" = 'steering'
ON CONFLICT DO NOTHING;
