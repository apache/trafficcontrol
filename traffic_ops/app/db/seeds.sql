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


-- THIS FILE INCLUDES STATIC DATA REQUIRED OF TRAFFIC OPS

-- cdns
INSERT INTO public.cdn ("name", dnssec_enabled, domain_name) VALUES ('ALL', FALSE, '-') ON CONFLICT ("name") DO NOTHING;

-- parameters
-- Moved into postinstall global parameters
INSERT INTO public.profile ("name", "description", "type", cdn) VALUES ('GLOBAL', 'Global Traffic Ops profile, DO NOT DELETE', 'UNK_PROFILE', (SELECT id FROM cdn WHERE "name"='ALL')) ON CONFLICT ("name") DO NOTHING;

---------------------------------

-- global parameters (settings)
---------------------------------
DO
$do$
BEGIN
	IF NOT EXISTS (SELECT id FROM public.parameter WHERE "name" = 'tm.instance_name' AND config_file = 'global') THEN
		INSERT INTO public.parameter ("name", config_file, "value") VALUES ('tm.instance_name', 'global', 'Traffic Ops CDN');
		INSERT INTO public.profile_parameter ("profile", parameter) VALUES ( (SELECT id FROM public.profile WHERE "name" = 'GLOBAL'), (SELECT id FROM public.parameter WHERE "name" = 'tm.instance_name' AND config_file = 'global' AND "value" = 'Traffic Ops CDN') ) ON CONFLICT ("profile", parameter) DO NOTHING;
	END IF;
	IF NOT EXISTS (SELECT id FROM public.parameter WHERE "name" = 'tm.toolname' AND config_file = 'global') THEN
		INSERT INTO public.parameter ("name", config_file, "value") VALUES ('tm.toolname', 'global', 'Traffic Ops');
		INSERT INTO public.profile_parameter ("profile", parameter) VALUES ( (SELECT id FROM public.profile WHERE "name" = 'GLOBAL'), (SELECT id FROM public.parameter WHERE "name" = 'tm.toolname' AND config_file = 'global' AND "value" = 'Traffic Ops') ) ON CONFLICT ("profile", parameter) DO NOTHING;
	END IF;
	IF NOT EXISTS (SELECT id FROM public.parameter WHERE "name" = 'maxRevalDurationDays' AND config_file = 'regex_revalidate.config') THEN
		INSERT INTO public.parameter ("name", config_file, "value") VALUES ('maxRevalDurationDays', 'regex_revalidate.config', '90');
		INSERT INTO public.profile_parameter ("profile", parameter) VALUES ( (SELECT id FROM public.profile WHERE "name" = 'GLOBAL'), (SELECT id FROM public.parameter WHERE "name" = 'maxRevalDurationDays' AND config_file = 'regex_revalidate.config' AND "value" = '90') ) ON CONFLICT ("profile", parameter) DO NOTHING;
	END IF;
END
$do$;

-- parameters
---------------------------------
INSERT INTO public.parameter ("name", config_file, "value") VALUES ('mso.parent_retry', 'parent.config', 'simple_retry') ON CONFLICT DO NOTHING;
INSERT INTO public.parameter ("name", config_file, "value") VALUES ('mso.parent_retry', 'parent.config', 'unavailable_server_retry') ON CONFLICT DO NOTHING;
INSERT INTO public.parameter ("name", config_file, "value") VALUES ('mso.parent_retry', 'parent.config', 'both') ON CONFLICT DO NOTHING;

-- profiles
---------------------------------
INSERT INTO public.profile ("name", "description", "type", cdn) VALUES ('TRAFFIC_ANALYTICS', 'Traffic Analytics profile', 'UNK_PROFILE', (SELECT id FROM cdn WHERE "name"='ALL')) ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.profile ("name", "description", "type", cdn) VALUES ('TRAFFIC_OPS', 'Traffic Ops profile', 'UNK_PROFILE', (SELECT id FROM cdn WHERE "name"='ALL')) ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.profile ("name", "description", "type", cdn) VALUES ('TRAFFIC_OPS_DB', 'Traffic Ops DB profile', 'UNK_PROFILE', (SELECT id FROM cdn WHERE "name"='ALL')) ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.profile ("name", "description", "type", cdn) VALUES ('TRAFFIC_PORTAL', 'Traffic Portal profile', 'TP_PROFILE', (SELECT id FROM cdn WHERE "name"='ALL')) ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.profile ("name", "description", "type", cdn) VALUES ('TRAFFIC_STATS', 'Traffic Stats profile', 'TS_PROFILE', (SELECT id FROM cdn WHERE "name"='ALL')) ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.profile ("name", "description", "type", cdn) VALUES ('INFLUXDB', 'InfluxDb profile', 'INFLUXDB_PROFILE', (SELECT id FROM cdn WHERE "name"='ALL')) ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.profile ("name", "description", "type", cdn) VALUES ('RIAK_ALL', 'Riak profile for all CDNs', 'RIAK_PROFILE', (SELECT id FROM cdn WHERE "name"='ALL')) ON CONFLICT ("name") DO NOTHING;

-- statuses
INSERT INTO public.status ("name", "description") VALUES ('ONLINE', 'Server is online.') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.status ("name", "description") VALUES ('OFFLINE', 'Server is Offline. Not active in any configuration.') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.status ("name", "description") VALUES ('REPORTED', 'Server is online and reported in the health protocol.') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.status ("name", "description") VALUES ('ADMIN_DOWN', 'Sever is administrative down and does not receive traffic.') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.status ("name", "description") VALUES ('PRE_PROD', 'Pre Production. Not active in any configuration.') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.status ("name", "description") VALUES ('CCR_IGNORE', 'Server is ignored by traffic router.') ON CONFLICT ("name") DO NOTHING;

-- tenants
INSERT INTO public.tenant ("name", active, parent_id) VALUES ('root', true, NULL) ON CONFLICT DO NOTHING;
INSERT INTO public.tenant ("name", active, parent_id) VALUES ('unassigned', true, (SELECT id FROM public.tenant WHERE "name"='root')) ON CONFLICT DO NOTHING;

-- roles
-- out of the box, only 4 roles are defined. Other roles can be created by the admin as needed.
INSERT INTO public.role ("name", "description", priv_level) VALUES ('admin', 'Has access to everything - cannot be modified or deleted', 30) ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.role ("name", "description", priv_level) VALUES ('operations', 'Has all reads and most write capabilities', 20) ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.role ("name", "description", priv_level) VALUES ('read-only', 'Has access to all read capabilities', 10) ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.role ("name", "description", priv_level) VALUES ('disallowed', 'Block all access', 0) ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.role ("name", "description", priv_level) VALUES ('portal','Portal User', 2) ON CONFLICT DO NOTHING;
INSERT INTO public.role ("name", "description", priv_level) VALUES ('steering','Steering User', 15) ON CONFLICT DO NOTHING;
INSERT INTO public.role ("name", "description", priv_level) VALUES ('federation','Role for Secondary CZF', 15) ON CONFLICT DO NOTHING;

-- roles_capabilities
-- out of the box, the admin role has ALL capabilities
INSERT INTO public.role_capability (role_id, cap_name)
SELECT id, 'ALL'
FROM public.role
WHERE "name" = 'admin'
ON CONFLICT DO NOTHING;

INSERT INTO public.role_capability
SELECT id, 'DELIVERY-SERVICE-SAFE:UPDATE'
FROM public.role
WHERE "name" in ('operations', 'read-only', 'portal', 'federation', 'steering')
ON CONFLICT DO NOTHING;

-- Using role 'read-only'
INSERT INTO public.role_capability (role_id, cap_name)
SELECT id, perm
FROM public.role
CROSS JOIN ( VALUES
	('ACME:READ'),
	('ASN:READ'),
	('ASYNC-STATUS:READ'),
	('CACHE-GROUP:READ'),
	('CAPABILITY:READ'),
	('CDN-SNAPSHOT:READ'),
	('CDN:READ'),
	('CDNI-ADMIN:READ'),
	('CDNI-CAPACITY:READ'),
	('COORDINATE:READ'),
	('DELIVERY-SERVICE:READ'),
	('DELIVERY-SERVICE-SAFE:UPDATE'),
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
	('SSL-KEY-EXPIRATION:READ'),
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
	('STAT:CREATE')
) AS perms(perm)
WHERE "name" IN ('operations', 'portal', 'read-only', 'federation', 'steering')
ON CONFLICT DO NOTHING;

-- Traditionally the 'portal'/'federations'/'steering' Role(s)
INSERT INTO public.role_capability
SELECT id, perm
FROM public.role
CROSS JOIN ( VALUES
	('DELIVERY-SERVICE:UPDATE'),
	('JOB:CREATE'),
	('JOB:UPDATE'),
	('JOB:DELETE'),
	('DS-REQUEST:UPDATE'),
	('DS-REQUEST:CREATE'),
	('DS-REQUEST:DELETE'),
	('STEERING:CREATE'),
	('STEERING:UPDATE'),
	('STEERING:DELETE')
) AS perms(perm)
WHERE "name" IN ('operations', 'portal', 'federation', 'steering')
ON CONFLICT DO NOTHING;

-- Federation and Steering Role Permissions (also given to operators).
INSERT INTO public.role_capability
SELECT id, perm
FROM public.role
CROSS JOIN ( VALUES
	('FEDERATION:CREATE'),
	('FEDERATION:UPDATE'),
	('FEDERATION:DELETE'),
	('FEDERATION-RESOLVER:CREATE'),
	('FEDERATION-RESOLVER:DELETE')
) AS perms(perm)
WHERE "name" IN ('operations', 'federation', 'steering')
ON CONFLICT DO NOTHING;

-- Using role 'operations'
INSERT INTO public.role_capability
SELECT id, perm
FROM public.role
CROSS JOIN ( VALUES
	('ACME:CREATE'),
	('ACME:DELETE'),
	('ACME:READ'),
	('ACME:UPDATE'),
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
	('DNS-SEC:READ'),
	('DNS-SEC:UPDATE'),
	('DNS-SEC:DELETE'),
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
	('SECURE-SERVER:READ'),
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
	('SERVER-CHECK:DELETE')
) AS perms(perm)
WHERE "name" = 'operations'
ON CONFLICT DO NOTHING;

-- types

-- delivery service types
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('HTTP', 'HTTP Content Routing', 'deliveryservice') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('HTTP_NO_CACHE', 'HTTP Content Routing, no caching', 'deliveryservice') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('HTTP_LIVE', 'HTTP Content routing cache in RAM', 'deliveryservice') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('HTTP_LIVE_NATNL', 'HTTP Content routing, RAM cache, National', 'deliveryservice') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('DNS', 'DNS Content Routing', 'deliveryservice') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('DNS_LIVE', 'DNS Content routing, RAM cache, Local', 'deliveryservice') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('DNS_LIVE_NATNL', 'DNS Content routing, RAM cache, National', 'deliveryservice') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('ANY_MAP', 'No Content Routing - arbitrary remap at the edge, no Traffic Router config', 'deliveryservice') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('STEERING', 'Steering Delivery Service', 'deliveryservice') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('CLIENT_STEERING', 'Client-Controlled Steering Delivery Service', 'deliveryservice') ON CONFLICT ("name") DO NOTHING;

-- server types
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('EDGE', 'Edge Cache', 'server') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('MID', 'Mid Tier Cache', 'server') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('ORG', 'Origin', 'server') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('CCR', 'Traffic Router', 'server') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('RASCAL', 'Traffic Monitor', 'server') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('RIAK', 'Riak keystore', 'server') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('INFLUXDB', 'influxDb server', 'server') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('TRAFFIC_ANALYTICS', 'traffic analytics server', 'server') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('TRAFFIC_OPS', 'traffic ops server', 'server') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('TRAFFIC_OPS_DB', 'traffic ops DB server', 'server') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('TRAFFIC_PORTAL', 'traffic portal server', 'server') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('TRAFFIC_STATS', 'traffic stats server', 'server') ON CONFLICT ("name") DO NOTHING;

-- cachegroup types
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('EDGE_LOC', 'Edge Logical Location', 'cachegroup') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('MID_LOC', 'Mid Logical Location', 'cachegroup') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('ORG_LOC', 'Origin Logical Site', 'cachegroup') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('TR_LOC', 'Traffic Router Logical Location', 'cachegroup') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('TC_LOC', 'Traffic Control Component Location', 'cachegroup') ON CONFLICT ("name") DO NOTHING;

-- to_extension types
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('CHECK_EXTENSION_BOOL', 'Extension for checkmark in Server Check', 'to_extension') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('CHECK_EXTENSION_NUM', 'Extension for int value in Server Check', 'to_extension') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('CHECK_EXTENSION_OPEN_SLOT', 'Open slot for check in Server Status', 'to_extension') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('CONFIG_EXTENSION', 'Extension for additional configuration file', 'to_extension') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('STATISTIC_EXTENSION', 'Extension source for 12M graphs', 'to_extension') ON CONFLICT ("name") DO NOTHING;

-- regex types
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('HOST_REGEXP', 'Host header regular expression', 'regex') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('HEADER_REGEXP', 'HTTP header regular expression', 'regex') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('PATH_REGEXP', 'URL path regular expression', 'regex') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('STEERING_REGEXP', 'Steering target filter regular expression', 'regex') ON CONFLICT ("name") DO NOTHING;

-- federation types
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('RESOLVE4', 'federation type resolve4', 'federation') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('RESOLVE6', 'federation type resolve6', 'federation') ON CONFLICT ("name") DO NOTHING;

-- static dns entry types
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('A_RECORD', 'Static DNS A entry', 'staticdnsentry') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('AAAA_RECORD', 'Static DNS AAAA entry', 'staticdnsentry') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('CNAME_RECORD', 'Static DNS CNAME entry', 'staticdnsentry') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('TXT_RECORD', 'Static DNS TXT entry', 'staticdnsentry') ON CONFLICT ("name") DO NOTHING;

--steering_target types
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('STEERING_WEIGHT', 'Weighted steering target', 'steering_target') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('STEERING_ORDER', 'Ordered steering target', 'steering_target') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('STEERING_GEO_ORDER', 'Geo-ordered steering target', 'steering_target') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.type ("name", "description", use_in_table) VALUES ('STEERING_GEO_WEIGHT', 'Geo-weighted steering target', 'steering_target') ON CONFLICT ("name") DO NOTHING;

-- users
INSERT INTO public.tm_user (username, "role", full_name, token, tenant_id) VALUES ('extension',
    (SELECT id FROM public.role WHERE "name" = 'operations'), 'Extension User, DO NOT DELETE', '91504CE6-8E4A-46B2-9F9F-FE7C15228498',
    (SELECT id FROM public.tenant WHERE "name" = 'root')) ON CONFLICT DO NOTHING;

-- to extensions
-- some of the old ones do not get a new place, and there will be 'gaps' in the column usage.... New to_extension add will have to take care of that.
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (1, 'ILO_PING', 'ILO', 'aa', '1.0.0', '-', 'ToPingCheck.pl', '1', '{ check_name: "ILO", "base_url": "https://localhost", "select": "ilo_ip_address", "cron": "9 * * * *" }',
	(SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_BOOL') ) ON CONFLICT DO NOTHING;
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (2, '10G_PING', '10G', 'ab', '1.0.0', '-', 'ToPingCheck.pl', '1', '{ check_name: "10G", "base_url": "https://localhost", "select": "ip_address", "cron": "18 * * * *" }',
	(SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_BOOL') ) ON CONFLICT DO NOTHING;
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (3, 'FQDN_PING', 'FQDN', 'ac', '1.0.0', '-', 'ToPingCheck.pl', '1', '{ check_name: "FQDN", "base_url": "https://localhost", "select": "host_name", "cron": "27 * * * *" }',
	(SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_BOOL') ) ON CONFLICT DO NOTHING;
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (4, 'CHECK_DSCP', 'DSCP', 'ad', '1.0.0', '-', 'ToDSCPCheck.pl', '1', '{ "check_name": "DSCP", "base_url": "https://localhost", "cron": "36 * * * *" }',
	(SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_BOOL') ) ON CONFLICT DO NOTHING;
-- open EF
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (5, 'OPEN', '', 'ae', '1.0.0', '-', '', '0', '', (SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (6, 'OPEN', '', 'af', '1.0.0', '-', '', '0', '', (SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
--
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (7, 'IPV6_PING', '10G6', 'ag', '1.0.0', '-', 'ToPingCheck.pl', '1', '{ "select": "ip6_address", "cron": "0 * * * *" }',
	(SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_BOOL') ) ON CONFLICT DO NOTHING;
-- upd_pending H -> open
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (8, 'OPEN', '', 'ah', '1.0.0', '', '', '0', '', (SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
-- open IJ
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (9, 'OPEN', '', 'ai', '1.0.0', '', '', '0', '', (SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (10, 'OPEN', '', 'aj', '1.0.0', '', '', '0', '', (SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
--
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (11, 'CHECK_MTU', 'MTU', 'ak', '1.0.0', '-', 'ToMtuCheck.pl', '1', '{ "check_name": "MTU", "base_url": "https://localhost", "cron": "45 * * * *" }',
	(SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_BOOL') ) ON CONFLICT DO NOTHING;
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (12, 'CHECK_TRAFFIC_ROUTER_STATUS', 'RTR', 'al', '1.0.0', '-', 'ToRTRCheck.pl', '1', '{  "check_name": "RTR", "base_url": "https://localhost", "cron": "10 * * * *" }',
	(SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_BOOL') ) ON CONFLICT DO NOTHING;
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (13, 'OPEN', '', 'am', '1.0.0', '-', '', '0', '', (SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (14, 'CACHE_HIT_RATIO_LAST_15', 'CHR', 'an', '1.0.0', '-', 'ToCHRCheck.pl', '1', '{ check_name: "CHR", "base_url": "https://localhost", cron": "0,15,30,45 * * * *" }',
	(SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_NUM') ) ON CONFLICT DO NOTHING;
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (15, 'DISK_UTILIZATION', 'CDU', 'ao', '1.0.0', '-', 'ToCDUCheck.pl', '1', '{ check_name: "CDU", "base_url": "https://localhost", cron": "20 * * * *" }',
	(SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_NUM') ) ON CONFLICT DO NOTHING;
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (16, 'ORT_ERROR_COUNT', 'ORT', 'ap', '1.0.0', '-', 'ToORTCheck.pl', '1', '{ check_name: "ORT", "base_url": "https://localhost", "cron": "40 * * * *" }',
	(SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_NUM') ) ON CONFLICT DO NOTHING;
-- rest open
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (17, 'OPEN', '', 'aq', '1.0.0', '-', '', '0', '', (SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (18, 'OPEN', '', 'ar', '1.0.0', '-', '', '0', '', (SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (19, 'OPEN', '', 'bf', '1.0.0', '-', '', '0', '', (SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (20, 'OPEN', '', 'at', '1.0.0', '-', '', '0', '', (SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (21, 'OPEN', '', 'au', '1.0.0', '-', '', '0', '', (SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (22, 'OPEN', '', 'av', '1.0.0', '-', '', '0', '', (SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (23, 'OPEN', '', 'aw', '1.0.0', '-', '', '0', '', (SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (24, 'OPEN', '', 'ax', '1.0.0', '-', '', '0', '', (SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (25, 'OPEN', '', 'ay', '1.0.0', '-', '', '0', '', (SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (26, 'OPEN', '', 'az', '1.0.0', '-', '', '0', '', (SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (27, 'OPEN', '', 'ba', '1.0.0', '-', '', '0', '', (SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (28, 'OPEN', '', 'bb', '1.0.0', '-', '', '0', '', (SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (29, 'OPEN', '', 'bc', '1.0.0', '-', '', '0', '', (SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (30, 'OPEN', '', 'bd', '1.0.0', '-', '', '0', '', (SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
INSERT INTO public.to_extension (id, "name", servercheck_short_name, servercheck_column_name, "version", info_url, script_file, isactive, additional_config_json, "type")
VALUES (31, 'OPEN', '', 'be', '1.0.0', '-', '', '0', '', (SELECT id FROM public.type WHERE "name"='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;

INSERT INTO public.last_deleted (table_name) VALUES ('api_capability') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('asn') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('cachegroup') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('cachegroup_fallbacks') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('cachegroup_localization_method') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('cachegroup_parameter') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('capability') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('cdn') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('coordinate') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('deliveryservice') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('deliveryservice_regex') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('deliveryservice_request') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('deliveryservice_request_comment') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('deliveryservice_server') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('division') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('federation') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('federation_deliveryservice') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('federation_federation_resolver') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('federation_resolver') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('federation_tmuser') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('hwinfo') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('job') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('log') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('origin') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('parameter') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('phys_location') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('profile') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('profile_parameter') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('regex') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('region') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('role') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('role_capability') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('server') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('servercheck') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('snapshot') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('staticdnsentry') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('stats_summary') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('status') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('steering_target') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('tenant') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('tm_user') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('topology') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('topology_cachegroup') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('topology_cachegroup_parents') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('to_extension') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('type') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('user_role') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('server_capability') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('server_server_capability') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('service_category') ON CONFLICT (table_name) DO NOTHING;
INSERT INTO public.last_deleted (table_name) VALUES ('deliveryservices_required_capability') ON CONFLICT (table_name) DO NOTHING;
