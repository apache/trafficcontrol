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

-- add built-in roles
INSERT INTO role (name, description, priv_level) VALUES ('operations', 'Has all reads and most write capabilities', 20) ON CONFLICT (name) DO NOTHING;
INSERT INTO role (name, description, priv_level) VALUES ('read-only', 'Has access to all read capabilities', 10) ON CONFLICT (name) DO NOTHING;
INSERT INTO role (name, description, priv_level) values ('disallowed', 'Block all access', 0) ON CONFLICT (name) DO NOTHING;
INSERT INTO role (name, description, priv_level) VALUES ('portal','Portal User', 2) ON CONFLICT DO NOTHING;
INSERT INTO role (name, description, priv_level) VALUES ('steering','Steering User', 15) ON CONFLICT DO NOTHING;
INSERT INTO role (name, description, priv_level) VALUES ('federation','Role for Secondary CZF', 15) ON CONFLICT DO NOTHING;

-- add permissions to roles
INSERT INTO role_capability SELECT id, perm FROM public.role CROSS JOIN (VALUES
('ASN:CREATE'),
('ASN:DELETE'),
('ASN:UPDATE'),
('CACHE-GROUP:CREATE'),
('CACHE-GROUP:DELETE'),
('CACHE-GROUP:UPDATE'),
('CDN-LOCK:CREATE'),
('CDN-LOCK:DELETE'),
('CDN-SNAPSHOT:CREATE'),
('CDN:CREATE'),
('CDN:DELETE'),
('CDN:UPDATE'),
('COORDINATE:CREATE'),
('COORDINATE:UPDATE'),
('COORDINATE:DELETE'),
('DELIVERY-SERVICE-SAFE:UPDATE'),
('DELIVERY-SERVICE:CREATE'),
('DELIVERY-SERVICE:DELETE'),
('DIVISION:CREATE'),
('DIVISION:DELETE'),
('DIVISION:UPDATE'),
('DNS-SEC:UPDATE'),
('ISO:GENERATE'),
('ORIGIN:CREATE'),
('ORIGIN:DELETE'),
('ORIGIN:UPDATE'),
('PARAMETER:CREATE'),
('PARAMETER:DELETE'),
('PARAMETER:UPDATE'),
('PHYSICAL-LOCATION:CREATE'),
('PHYSICAL-LOCATION:DELETE'),
('PHYSICAL-LOCATION:UPDATE'),
('PROFILE:CREATE'),
('PROFILE:DELETE'),
('PROFILE:UPDATE'),
('REGION:CREATE'),
('REGION:DELETE'),
('REGION:UPDATE'),
('SERVER-CAPABILITY:CREATE'),
('SERVER-CAPABILITY:DELETE'),
('SERVER-CAPABILITY:UPDATE'),
('SERVER:CREATE'),
('SERVER:DELETE'),
('SERVER:QUEUE'),
('SERVER:UPDATE'),
('SERVICE-CATEGORY:CREATE'),
('SERVICE-CATEGORY:DELETE'),
('SERVICE-CATEGORY:UPDATE'),
('STATIC-DN:CREATE'),
('STATIC-DN:DELETE'),
('STATIC-DN:UPDATE'),
('STATUS:CREATE'),
('STATUS:DELETE'),
('STATUS:UPDATE'),
('TENANT:CREATE'),
('TENANT:DELETE'),
('TENANT:UPDATE'),
('TOPOLOGY:CREATE'),
('TOPOLOGY:DELETE'),
('TOPOLOGY:UPDATE'),
('TYPE:CREATE'),
('TYPE:DELETE'),
('TYPE:UPDATE'),
('USER:CREATE'),
('USER:UPDATE'),
('SERVER-CHECK:CREATE'),
('SERVER-CHECK:DELETE')) AS perms(perm)
WHERE priv_level >= 20;

INSERT INTO role_capability SELECT id, perm FROM public.role CROSS JOIN (VALUES
('FEDERATION:CREATE'),
('FEDERATION:UPDATE'),
('FEDERATION:DELETE'),
('FEDERATION-RESOLVER:CREATE'),
('FEDERATION-RESOLVER:DELETE'),
('DELIVERY-SERVICE:UPDATE'),
('JOB:CREATE'),
('JOB:UPDATE'),
('JOB:DELETE'),
('DS-REQUEST:UPDATE'),
('DS-REQUEST:CREATE'),
('DS-REQUEST:DELETE'),
('STEERING:CREATE'),
('STEERING:UPDATE'),
('STEERING:DELETE')) AS perms(perm)
WHERE priv_level >= 15;

INSERT INTO role_capability SELECT id, perm FROM public.role CROSS JOIN (VALUES
('ASN:READ'),
('ASYNC-STATUS:READ'),
('CACHE-GROUP:READ'),
('CAPABILITY:READ'),
('CDN-SNAPSHOT:READ'),
('CDN:READ'),
('COORDINATE:READ'),
('DELIVERY-SERVICE:READ'),
('DIVISION:READ'),
('DS-REQUEST:READ'),
('DS-SECURITY-KEY:READ'),
('FEDERATION:READ'),
('FEDERATION-RESOLVER:READ'),
('ISO:READ'),
('JOB:READ'),
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
('SERVER:READ'),
('SERVICE-CATEGORY:READ'),
('STATIC-DN:READ'),
('STATUS:READ'),
('SERVER-CHECK:READ'),
('STEERING:READ'),
('STAT:READ'),
('TENANT:READ'),
('TOPOLOGY:READ'),
('TRAFFIC-VAULT:READ'),
('TYPE:READ'),
('USER:READ'),
('STAT:CREATE')) AS perms(perm)
WHERE priv_level >= 10;

INSERT INTO role_capability (role_id, cap_name) SELECT * FROM (SELECT (SELECT role FROM tm_user WHERE username='extension'), 'SERVER-CHECK:CREATE') AS i(role_id, cap_name) WHERE EXISTS (SELECT 1 FROM tm_user WHERE username='extension');
INSERT INTO role_capability (role_id, cap_name) SELECT * FROM (SELECT (SELECT role FROM tm_user WHERE username='extension'), 'SERVER-CHECK:DELETE') AS i(role_id, cap_name) WHERE EXISTS (SELECT 1 FROM tm_user WHERE username='extension');
INSERT INTO role_capability (role_id, cap_name) SELECT * FROM (SELECT (SELECT role FROM tm_user WHERE username='extension'), 'SERVER-CHECK:READ') AS i(role_id, cap_name) WHERE EXISTS (SELECT 1 FROM tm_user WHERE username='extension');
INSERT INTO role_capability (role_id, cap_name) SELECT * FROM (SELECT (SELECT role FROM tm_user WHERE username='extension'), 'SERVER:READ') AS i(role_id, cap_name) WHERE EXISTS (SELECT 1 FROM tm_user WHERE username='extension');