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

-- acme
INSERT INTO capability (name, description) VALUES ('ACME:READ', 'Ability to view acme keys and data') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('ACME:CREATE', 'Ability to create acme keys and data') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('ACME:UPDATE', 'Ability to modify acme keys and data') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('ACME:DELETE', 'Ability to delete acme keys and data') ON CONFLICT (name) DO NOTHING;

-- delivery services
INSERT INTO capability (name, description) VALUES ('DELIVERY-SERVICE:READ', 'Ability to view delivery services') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('DELIVERY-SERVICE:CREATE', 'Ability to create delivery services') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('DELIVERY-SERVICE:UPDATE', 'Ability to modify delivery services') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('DELIVERY-SERVICE:DELETE', 'Ability to delete delivery services') ON CONFLICT (name) DO NOTHING;

-- delivery service safe
INSERT INTO capability (name, description) VALUES ('DELIVERY-SERVICE-SAFE:UPDATE', 'Ability to modify delivery services safely') ON CONFLICT (name) DO NOTHING;

-- delivery service security keys
INSERT INTO capability (name, description) VALUES ('DS-SECURITY-KEY:READ', 'Ability to view delivery service security keys') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('DS-SECURITY-KEY:CREATE', 'Ability to create delivery service security keys') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('DS-SECURITY-KEY:UPDATE', 'Ability to modify delivery service security keys') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('DS-SECURITY-KEY:DELETE', 'Ability to delete delivery service security keys') ON CONFLICT (name) DO NOTHING;

-- asns
INSERT INTO capability (name, description) VALUES ('ASN:READ', 'Ability to view asns') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('ASN:CREATE', 'Ability to create asns') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('ASN:UPDATE', 'Ability to modify asns') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('ASN:DELETE', 'Ability to delete asns') ON CONFLICT (name) DO NOTHING;

-- cachegroups
INSERT INTO capability (name, description) VALUES ('CACHE-GROUP:READ', 'Ability to view cache groups') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('CACHE-GROUP:CREATE', 'Ability to create cache groups') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('CACHE-GROUP:UPDATE', 'Ability to modify cache groups') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('CACHE-GROUP:DELETE', 'Ability to delete cache groups') ON CONFLICT (name) DO NOTHING;

-- async_status
INSERT INTO capability (name, description) VALUES ('ASYNC-STATUS:READ', 'Ability to view async status') ON CONFLICT (name) DO NOTHING;

-- cdns
INSERT INTO capability (name, description) VALUES ('CDN:READ', 'Ability to view cdns') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('CDN:CREATE', 'Ability to create cdns') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('CDN:UPDATE', 'Ability to modify cdns') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('CDN:DELETE', 'Ability to delete cdns') ON CONFLICT (name) DO NOTHING;

-- types
INSERT INTO capability (name, description) VALUES ('TYPE:READ', 'Ability to view types') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('TYPE:CREATE', 'Ability to create types') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('TYPE:UPDATE', 'Ability to modify types') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('TYPE:DELETE', 'Ability to delete types') ON CONFLICT (name) DO NOTHING;

-- servers
INSERT INTO capability (name, description) VALUES ('SERVER:READ', 'Ability to view servers') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('SERVER:CREATE', 'Ability to create servers') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('SERVER:UPDATE', 'Ability to modify servers') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('SERVER:DELETE', 'Ability to delete servers') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('SERVER:QUEUE', 'Ability to queue updates on servers') ON CONFLICT (name) DO NOTHING;

-- profiles
INSERT INTO capability (name, description) VALUES ('PROFILE:READ', 'Ability to view profiles') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('PROFILE:CREATE', 'Ability to create profiles') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('PROFILE:UPDATE', 'Ability to modify profiles') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('PROFILE:DELETE', 'Ability to delete profiles') ON CONFLICT (name) DO NOTHING;

-- capabilities
INSERT INTO capability (name, description) VALUES ('CAPABILITY:READ', 'Ability to view capabilities') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('CAPABILITY:CREATE', 'Ability to create capabilities') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('CAPABILITY:UPDATE', 'Ability to modify capabilities') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('CAPABILITY:DELETE', 'Ability to delete capabilities') ON CONFLICT (name) DO NOTHING;

-- cdn-locks
INSERT INTO capability (name, description) VALUES ('CDN-LOCK:READ', 'Ability to view locks') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('CDN-LOCK:CREATE', 'Ability to create locks') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('CDN-LOCK:DELETE-OTHERS', 'Ability to delete locks of other users') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('CDN-LOCK:DELETE', 'Ability to delete locks') ON CONFLICT (name) DO NOTHING;

-- dns-sec
INSERT INTO capability (name, description) VALUES ('DNS-SEC:READ', 'Ability to view dns-sec keys') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('DNS-SEC:CREATE', 'Ability to create dns-sec keys') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('DNS-SEC:UPDATE', 'Ability to modify dns-sec keys') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('DNS-SEC:DELETE', 'Ability to delete dns-sec keys') ON CONFLICT (name) DO NOTHING;

-- parameters
INSERT INTO capability (name, description) VALUES ('PARAMETER:READ', 'Ability to view parameters') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('PARAMETER:CREATE', 'Ability to create parameters') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('PARAMETER:UPDATE', 'Ability to modify parameters') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('PARAMETER:DELETE', 'Ability to delete parameters') ON CONFLICT (name) DO NOTHING;

-- monitor-config
INSERT INTO capability (name, description) VALUES ('MONITOR-CONFIG:READ', 'Ability to view monitor configs') ON CONFLICT (name) DO NOTHING;

-- federations
INSERT INTO capability (name, description) VALUES ('FEDERATION:READ', 'Ability to view federations') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('FEDERATION:CREATE', 'Ability to create federations') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('FEDERATION:UPDATE', 'Ability to modify federations') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('FEDERATION:DELETE', 'Ability to delete federations') ON CONFLICT (name) DO NOTHING;

-- snapshots
INSERT INTO capability (name, description) VALUES ('CDN-SNAPSHOT:READ', 'Ability to view CDN snapshots') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('CDN-SNAPSHOT:CREATE', 'Ability to create CDN snapshots') ON CONFLICT (name) DO NOTHING;

-- coordinates
INSERT INTO capability (name, description) VALUES ('COORDINATE:READ', 'Ability to view coordinates') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('COORDINATE:CREATE', 'Ability to create coordinates') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('COORDINATE:UPDATE', 'Ability to modify coordinates') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('COORDINATE:DELETE', 'Ability to delete coordinates') ON CONFLICT (name) DO NOTHING;

-- regions
INSERT INTO capability (name, description) VALUES ('REGION:READ', 'Ability to view regions') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('REGION:CREATE', 'Ability to create regions') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('REGION:UPDATE', 'Ability to modify regions') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('REGION:DELETE', 'Ability to delete regions') ON CONFLICT (name) DO NOTHING;

-- divisions
INSERT INTO capability (name, description) VALUES ('DIVISION:READ', 'Ability to view divisions') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('DIVISION:CREATE', 'Ability to create divisions') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('DIVISION:UPDATE', 'Ability to modify divisions') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('DIVISION:DELETE', 'Ability to delete divisions') ON CONFLICT (name) DO NOTHING;

-- physical-locations
INSERT INTO capability (name, description) VALUES ('PHYSICAL-LOCATION:READ', 'Ability to view physical locations') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('PHYSICAL-LOCATION:CREATE', 'Ability to create physical locations') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('PHYSICAL-LOCATION:UPDATE', 'Ability to modify physical locations') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('PHYSICAL-LOCATION:DELETE', 'Ability to delete physical locations') ON CONFLICT (name) DO NOTHING;

-- dbdump
INSERT INTO capability (name, description) VALUES ('DBDUMP:READ', 'Ability to view the db dump') ON CONFLICT (name) DO NOTHING;

-- delivery service requests
INSERT INTO capability (name, description) VALUES ('DS-REQUEST:READ', 'Ability to view delivery service requests') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('DS-REQUEST:CREATE', 'Ability to create delivery service requests') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('DS-REQUEST:UPDATE', 'Ability to modify delivery service requests') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('DS-REQUEST:DELETE', 'Ability to delete delivery service requests') ON CONFLICT (name) DO NOTHING;

-- users
INSERT INTO capability (name, description) VALUES ('USER:READ', 'Ability to view users') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('USER:CREATE', 'Ability to create users') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('USER:UPDATE', 'Ability to modify users') ON CONFLICT (name) DO NOTHING;

-- stats
INSERT INTO capability (name, description) VALUES ('STAT:READ', 'Ability to view stats') ON CONFLICT (name) DO NOTHING;

-- federation-resolvers
INSERT INTO capability (name, description) VALUES ('FEDERATION-RESOLVER:READ', 'Ability to view federation resolvers') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('FEDERATION-RESOLVER:CREATE', 'Ability to create federation resolvers') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('FEDERATION-RESOLVER:DELETE', 'Ability to delete federation resolvers') ON CONFLICT (name) DO NOTHING;

-- isos
INSERT INTO capability (name, description) VALUES ('ISO:READ', 'Ability to view isos') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('ISO:GENERATE', 'Ability to create isos') ON CONFLICT (name) DO NOTHING;

-- jobs
INSERT INTO capability (name, description) VALUES ('JOB:READ', 'Ability to view jobs') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('JOB:CREATE', 'Ability to create jobs') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('JOB:UPDATE', 'Ability to modify jobs') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('JOB:DELETE', 'Ability to delete jobs') ON CONFLICT (name) DO NOTHING;

-- logs
INSERT INTO capability (name, description) VALUES ('LOG:READ', 'Ability to view logs') ON CONFLICT (name) DO NOTHING;

-- origins
INSERT INTO capability (name, description) VALUES ('ORIGIN:READ', 'Ability to view origins') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('ORIGIN:CREATE', 'Ability to create origins') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('ORIGIN:UPDATE', 'Ability to modify origins') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('ORIGIN:DELETE', 'Ability to delete origins') ON CONFLICT (name) DO NOTHING;

-- plugins
INSERT INTO capability (name, description) VALUES ('PLUGIN:READ', 'Ability to view plugins') ON CONFLICT (name) DO NOTHING;

-- roles
INSERT INTO capability (name, description) VALUES ('ROLE:READ', 'Ability to view roles') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('ROLE:CREATE', 'Ability to create roles') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('ROLE:UPDATE', 'Ability to modify roles') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('ROLE:DELETE', 'Ability to delete roles') ON CONFLICT (name) DO NOTHING;

-- server-capabilities
INSERT INTO capability (name, description) VALUES ('SERVER-CAPABILITY:READ', 'Ability to view server capabilities') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('SERVER-CAPABILITY:CREATE', 'Ability to create server capabilities') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('SERVER-CAPABILITY:UPDATE', 'Ability to modify server capabilities') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('SERVER-CAPABILITY:DELETE', 'Ability to delete server capabilities') ON CONFLICT (name) DO NOTHING;

-- server-checks
INSERT INTO capability (name, description) VALUES ('SERVER-CHECK:READ', 'Ability to view server checks') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('SERVER-CHECK:CREATE', 'Ability to create server checks') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('SERVER-CHECK:DELETE', 'Ability to delete server checks') ON CONFLICT (name) DO NOTHING;

-- service-categories
INSERT INTO capability (name, description) VALUES ('SERVICE-CATEGORY:READ', 'Ability to view service categories') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('SERVICE-CATEGORY:CREATE', 'Ability to create service categories') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('SERVICE-CATEGORY:UPDATE', 'Ability to modify service categories') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('SERVICE-CATEGORY:DELETE', 'Ability to delete service categories') ON CONFLICT (name) DO NOTHING;

-- static-dns-entries
INSERT INTO capability (name, description) VALUES ('STATIC-DN:READ', 'Ability to view static dns entries') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('STATIC-DN:CREATE', 'Ability to create static dns entries') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('STATIC-DN:UPDATE', 'Ability to modify static dns entries') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('STATIC-DN:DELETE', 'Ability to delete static dns entries') ON CONFLICT (name) DO NOTHING;

-- statuses
INSERT INTO capability (name, description) VALUES ('STATUS:READ', 'Ability to view statuses') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('STATUS:CREATE', 'Ability to create statuses') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('STATUS:UPDATE', 'Ability to modify statuses') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('STATUS:DELETE', 'Ability to delete statuses') ON CONFLICT (name) DO NOTHING;

-- steering-targets
INSERT INTO capability (name, description) VALUES ('STEERING:READ', 'Ability to view steering targets') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('STEERING:CREATE', 'Ability to create steering targets') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('STEERING:UPDATE', 'Ability to modify steering targets') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('STEERING:DELETE', 'Ability to delete steering targets') ON CONFLICT (name) DO NOTHING;

-- tenants
INSERT INTO capability (name, description) VALUES ('TENANT:READ', 'Ability to view tenants') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('TENANT:CREATE', 'Ability to create tenants') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('TENANT:UPDATE', 'Ability to modify tenants') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('TENANT:DELETE', 'Ability to delete tenants') ON CONFLICT (name) DO NOTHING;

-- topologies
INSERT INTO capability (name, description) VALUES ('TOPOLOGY:READ', 'Ability to view topologies') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('TOPOLOGY:CREATE', 'Ability to create topologies') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('TOPOLOGY:UPDATE', 'Ability to modify topologies') ON CONFLICT (name) DO NOTHING;
INSERT INTO capability (name, description) VALUES ('TOPOLOGY:DELETE', 'Ability to delete topologies') ON CONFLICT (name) DO NOTHING;

-- traffic-vault
INSERT INTO capability (name, description) VALUES ('TRAFFIC-VAULT:READ', 'Ability to ping traffic vault') ON CONFLICT (name) DO NOTHING;


-- add permissions to roles
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'ACME:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'ASN:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'ASN:DELETE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'ASN:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'ASN:UPDATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'ASYNC-STATUS:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'CACHE-GROUP:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'CACHE-GROUP:DELETE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'CACHE-GROUP:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'CACHE-GROUP:UPDATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'CDN-LOCK:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'CDN-LOCK:DELETE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'CDN-SNAPSHOT:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'CDN-SNAPSHOT:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'CDN:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'CDN:DELETE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'CDN:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'CDN:UPDATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'COORDINATE:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'COORDINATE:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'COORDINATE:UPDATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'DELIVERY-SERVICE-SAFE:UPDATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'DELIVERY-SERVICE:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'DELIVERY-SERVICE:DELETE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'DELIVERY-SERVICE:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'DELIVERY-SERVICE:UPDATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'DIVISION:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'DIVISION:DELETE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'DIVISION:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'DIVISION:UPDATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'DNS-SEC:UPDATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'DS-REQUEST:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'DS-REQUEST:UPDATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'DS-SECURITY-KEY:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'DS-SECURITY-KEY:DELETE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'DS-SECURITY-KEY:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'DS-SECURITY-KEY:UPDATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'ISO:GENERATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'ISO:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'ORIGIN:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'ORIGIN:DELETE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'ORIGIN:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'ORIGIN:UPDATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'PARAMETER:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'PARAMETER:DELETE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'PARAMETER:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'PARAMETER:UPDATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'PHYSICAL-LOCATION:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'PHYSICAL-LOCATION:DELETE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'PHYSICAL-LOCATION:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'PHYSICAL-LOCATION:UPDATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'PROFILE:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'PROFILE:DELETE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'PROFILE:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'PROFILE:UPDATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'REGION:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'REGION:DELETE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'REGION:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'REGION:UPDATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'SERVER-CAPABILITY:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'SERVER-CAPABILITY:DELETE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'SERVER-CAPABILITY:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'SERVER-CAPABILITY:UPDATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'SERVER:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'SERVER:DELETE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'SERVER:QUEUE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'SERVER:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'SERVER:UPDATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'SERVICE-CATEGORY:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'SERVICE-CATEGORY:DELETE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'SERVICE-CATEGORY:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'SERVICE-CATEGORY:UPDATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'STATIC-DN:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'STATIC-DN:DELETE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'STATIC-DN:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'STATIC-DN:UPDATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'STATUS:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'STATUS:DELETE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'STATUS:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'STATUS:UPDATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'TENANT:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'TENANT:DELETE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'TENANT:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'TENANT:UPDATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'TOPOLOGY:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'TOPOLOGY:DELETE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'TOPOLOGY:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'TOPOLOGY:UPDATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'TYPE:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'TYPE:DELETE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'TYPE:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'TYPE:UPDATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'USER:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'USER:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'USER:UPDATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'SERVER-CHECK:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'SERVER-CHECK:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='operations'), 'SERVER-CHECK:DELETE' WHERE EXISTS (SELECT id FROM role WHERE name = 'operations') ON CONFLICT DO NOTHING;

INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='steering'), 'DELIVERY-SERVICE:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'steering') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='steering'), 'CDN:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'steering') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='steering'), 'TYPE:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'steering') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='steering'), 'STEERING:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'steering') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='steering'), 'STEERING:DELETE' WHERE EXISTS (SELECT id FROM role WHERE name = 'steering') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='steering'), 'STEERING:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'steering') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='steering'), 'STEERING:UPDATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'steering') ON CONFLICT DO NOTHING;

INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'ACME:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'ASN:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'ASYNC-STATUS:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'CACHE-GROUP:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'CDN-SNAPSHOT:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'CDN:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'COORDINATE:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'DELIVERY-SERVICE:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'DIVISION:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'DS-REQUEST:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'DS-SECURITY-KEY:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'ISO:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'ORIGIN:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'PARAMETER:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'PHYSICAL-LOCATION:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'PROFILE:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'REGION:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'SERVER-CAPABILITY:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'SERVER:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'SERVICE-CATEGORY:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'STATIC-DN:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'STATUS:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'TENANT:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'TOPOLOGY:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'TYPE:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'USER:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'SERVER-CHECK:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'SERVER-CHECK:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name='read-only'), 'STEERING:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'read-only') ON CONFLICT DO NOTHING;

-- For role 'federation'
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'federation'), 'FEDERATION:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'federation') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'federation'), 'FEDERATION:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'federation') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'federation'), 'FEDERATION:UPDATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'federation') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'federation'), 'FEDERATION:DELETE' WHERE EXISTS (SELECT id FROM role WHERE name = 'federation') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'federation'), 'FEDERATION-RESOLVER:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'federation') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'federation'), 'FEDERATION-RESOLVER:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'federation') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'federation'), 'FEDERATION-RESOLVER:DELETE' WHERE EXISTS (SELECT id FROM role WHERE name = 'federation') ON CONFLICT DO NOTHING;

-- For role 'portal'
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'ASN:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'CACHE-GROUP:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'CAPABILITY:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'CDN:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'DS-SECURITY-KEY:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'LOG:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'COORDINATE:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'DELIVERY-SERVICE:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'DIVISION:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'SERVER-CHECK:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'FEDERATION:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'JOB:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'JOB:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'JOB:UPDATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'JOB:DELETE' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'ORIGIN:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'PARAMETER:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'PHYSICAL-LOCATION:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'PROFILE:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'REGION:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'ROLE:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'SERVER:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'STAT:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'STATUS:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'STATIC-DN:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'STEERING:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'system-info-read' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'TENANT-READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'TYPE-READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'USER-READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'USER:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'DS-REQUEST:READ' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'DS-REQUEST:CREATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'DS-REQUEST:UPDATE' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) SELECT (SELECT id FROM role WHERE name = 'portal'), 'DS-REQUEST:DELETE' WHERE EXISTS (SELECT id FROM role WHERE name = 'portal') ON CONFLICT DO NOTHING;