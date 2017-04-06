-- -- -- -- -- -- -- /*

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
insert into cdn (name, dnssec_enabled, domain_name) values ('ALL', false, '-') ON CONFLICT (name) DO NOTHING;

-- job agents
insert into job_agent (name, description, active) values ('dummy', 'Description of Purge Agent', 1) ON CONFLICT (name) DO NOTHING;

-- job statuses
insert into job_status (name, description) values ('PENDING', 'Job is queued, but has not been picked up by any agents yet') ON CONFLICT (name) DO NOTHING;

-- parameters
insert into parameter (name, config_file, value) values ('ttl_max_hours', 'regex_revalidate.config', '672') ON CONFLICT (name, config_file, value) DO NOTHING;
insert into parameter (name, config_file, value) values ('ttl_min_hours', 'regex_revalidate.config', '48') ON CONFLICT (name, config_file, value) DO NOTHING;
insert into parameter (name, config_file, value) values ('maxRevalDurationDays', 'regex_revalidate.config', '90') ON CONFLICT (name, config_file, value) DO NOTHING;
insert into parameter (name, config_file, value) values ('CacheStats', 'traffic_stats.config', 'bandwidth') ON CONFLICT (name, config_file, value) DO NOTHING;
insert into parameter (name, config_file, value) values ('CacheStats', 'traffic_stats.config', 'maxKbps') ON CONFLICT (name, config_file, value) DO NOTHING;
insert into parameter (name, config_file, value) values ('CacheStats', 'traffic_stats.config', 'ats.proxy.process.http.current_client_connections') ON CONFLICT (name, config_file, value) DO NOTHING;
insert into parameter (name, config_file, value) values ('DsStats', 'traffic_stats.config', 'kbps') ON CONFLICT (name, config_file, value) DO NOTHING;
insert into parameter (name, config_file, value) values ('DsStats', 'traffic_stats.config', 'status_4xx') ON CONFLICT (name, config_file, value) DO NOTHING;
insert into parameter (name, config_file, value) values ('DsStats', 'traffic_stats.config', 'status_5xx') ON CONFLICT (name, config_file, value) DO NOTHING;
insert into parameter (name, config_file, value) values ('DsStats', 'traffic_stats.config', 'tps_2xx') ON CONFLICT (name, config_file, value) DO NOTHING;
insert into parameter (name, config_file, value) values ('DsStats', 'traffic_stats.config', 'tps_3xx') ON CONFLICT (name, config_file, value) DO NOTHING;
insert into parameter (name, config_file, value) values ('DsStats', 'traffic_stats.config', 'tps_4xx') ON CONFLICT (name, config_file, value) DO NOTHING;
insert into parameter (name, config_file, value) values ('DsStats', 'traffic_stats.config', 'tps_5xx') ON CONFLICT (name, config_file, value) DO NOTHING;
insert into parameter (name, config_file, value) values ('DsStats', 'traffic_stats.config', 'tps_total') ON CONFLICT (name, config_file, value) DO NOTHING;

-- profiles
insert into profile (name, description, type) values ('GLOBAL', 'Global Traffic Ops profile, DO NOT DELETE', 'UNK_PROFILE') ON CONFLICT (name) DO NOTHING;
insert into profile (name, description, type) values ('TRAFFIC_ANALYTICS', 'Traffic Analytics profile', 'UNK_PROFILE') ON CONFLICT (name) DO NOTHING;
insert into profile (name, description, type) values ('TRAFFIC_OPS', 'Traffic Ops profile', 'UNK_PROFILE') ON CONFLICT (name) DO NOTHING;
insert into profile (name, description, type) values ('TRAFFIC_OPS_DB', 'Traffic Ops DB profile', 'UNK_PROFILE') ON CONFLICT (name) DO NOTHING;
insert into profile (name, description, type) values ('TRAFFIC_PORTAL', 'Traffic Portal profile', 'TP_PROFILE') ON CONFLICT (name) DO NOTHING;
insert into profile (name, description, type) values ('TRAFFIC_STATS', 'Traffic Stats profile', 'TS_PROFILE') ON CONFLICT (name) DO NOTHING;
insert into profile (name, description, type) values ('INFLUXDB', 'InfluxDb profile', 'INFLUXDB_PROFILE') ON CONFLICT (name) DO NOTHING;
insert into profile (name, description, type) values ('RIAK_ALL', 'Riak profile for all CDNs', 'RIAK_PROFILE') ON CONFLICT (name) DO NOTHING;

-- profile_parameters
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'TRAFFIC_STATS'), (select id from parameter where name = 'CacheStats' and config_file = 'traffic_stats.config' and value = 'bandwidth') ) ON CONFLICT (profile, parameter) DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'TRAFFIC_STATS'), (select id from parameter where name = 'CacheStats' and config_file = 'traffic_stats.config' and value = 'maxKbps') ) ON CONFLICT (profile, parameter) DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'TRAFFIC_STATS'), (select id from parameter where name = 'CacheStats' and config_file = 'traffic_stats.config' and value = 'ats.proxy.process.http.current_client_connections') ) ON CONFLICT (profile, parameter) DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'TRAFFIC_STATS'), (select id from parameter where name = 'DsStats' and config_file = 'traffic_stats.config' and value = 'kbps') ) ON CONFLICT (profile, parameter) DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'TRAFFIC_STATS'), (select id from parameter where name = 'DsStats' and config_file = 'traffic_stats.config' and value = 'tps_2xx') ) ON CONFLICT (profile, parameter) DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'TRAFFIC_STATS'), (select id from parameter where name = 'DsStats' and config_file = 'traffic_stats.config' and value = 'status_4xx') ) ON CONFLICT (profile, parameter) DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'TRAFFIC_STATS'), (select id from parameter where name = 'DsStats' and config_file = 'traffic_stats.config' and value = 'status_5xx') ) ON CONFLICT (profile, parameter) DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'TRAFFIC_STATS'), (select id from parameter where name = 'DsStats' and config_file = 'traffic_stats.config' and value = 'tps_3xx') ) ON CONFLICT (profile, parameter) DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'TRAFFIC_STATS'), (select id from parameter where name = 'DsStats' and config_file = 'traffic_stats.config' and value = 'tps_4xx') ) ON CONFLICT (profile, parameter) DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'TRAFFIC_STATS'), (select id from parameter where name = 'DsStats' and config_file = 'traffic_stats.config' and value = 'tps_5xx') ) ON CONFLICT (profile, parameter) DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'TRAFFIC_STATS'), (select id from parameter where name = 'DsStats' and config_file = 'traffic_stats.config' and value = 'tps_total') ) ON CONFLICT (profile, parameter) DO NOTHING;

-- statuses
insert into status (name, description) values ('OFFLINE', 'Server is Offline. Not active in any configuration.') ON CONFLICT (name) DO NOTHING;
insert into status (name, description) values ('ONLINE', 'Server is online.') ON CONFLICT (name) DO NOTHING;
insert into status (name, description) values ('REPORTED', 'Server is online and reporeted in the health protocol.') ON CONFLICT (name) DO NOTHING;
insert into status (name, description) values ('ADMIN_DOWN', 'Sever is administrative down and does not receive traffic.') ON CONFLICT (name) DO NOTHING;
insert into status (name, description) values ('CCR_IGNORE', 'Server is ignored by traffic router.') ON CONFLICT (name) DO NOTHING;
insert into status (name, description) values ('PRE_PROD', 'Pre Production. Not active in any configuration.') ON CONFLICT (name) DO NOTHING;

-- roles
insert into role (name, description, priv_level) values ('admin', 'super-user', 30) ON CONFLICT (name) DO NOTHING;
insert into role (name, description, priv_level) values ('operations', 'Operations user', 20) ON CONFLICT (name) DO NOTHING;
insert into role (name, description, priv_level) values ('migrations', 'database migrations user - DO NOT REMOVE', 20) ON CONFLICT (name) DO NOTHING;
insert into role (name, description, priv_level) values ('federation', 'Role for Secondary CZF', 15) ON CONFLICT (name) DO NOTHING;
insert into role (name, description, priv_level) values ('steering', 'Role for Steering Delivery Services', 15) ON CONFLICT (name) DO NOTHING;
insert into role (name, description, priv_level) values ('read-only user', 'Read-Only user', 10) ON CONFLICT (name) DO NOTHING;
insert into role (name, description, priv_level) values ('portal', 'Portal User', 2) ON CONFLICT (name) DO NOTHING;
insert into role (name, description, priv_level) values ('disallowed', 'Block all access', 0) ON CONFLICT (name) DO NOTHING;
insert into role (name, description, priv_level) values ('root', 'Role for full capabilities - super-user ', 30) ON CONFLICT DO NOTHING;

-- tenants
insert into tenant (name, active, parent_id) values ('root', true, null) ON CONFLICT DO NOTHING;

-- capabilities
insert into capability (name, description) values ('all-read', 'Full read access') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('all-write', 'Full write access') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('asn-read', 'View ASN configuration') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('asn-write', 'Create, edit or delete ASN configuration') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('basic-read', 'Basic read operations. Every user should have this capability') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('basic-write', 'Basic write operations. Every user should have this capability') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('cache-config-files-read', 'View the generated cache configuration files') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('cache-group-read', 'View cache-group configuration') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('cache-group-write', 'Create, edit or delete cache-group configuration') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('cache-stats-read', 'View Cache statistics read access') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('cdn-config-snapshot-read', 'View config snapshot at CDN level') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('cdn-config-snapshot-write', 'Config snapshot write access at CDN level') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('cdn-health-read', 'View CDN health') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('cdn-read', 'View CDN configuration') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('cdn-write', 'Create, edit or delete CDN configuration') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('cdn-security-keys-read', 'View CDN DNSSEC keys') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('cdn-security-keys-write', 'Create, edit or delete CDN DNSSEC keys') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('cdn-stats-read', 'View CDN statistics') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('cdn-stats-write', 'Create, edit or delete CDN statistics') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('change-log-read', 'View change-log') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('change-log-write', 'Create change-log entries') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('division-read', 'View division configuration') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('division-write', 'Create, edit or delete division configuration') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('ds-cache-read', 'View delivery-service cache assignment') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('ds-cache-read', 'Create, edit or delete delivery-service cache assignment') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('ds-health-read', 'View delivery-service health') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('ds-read', 'View delivery-service configuration') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('ds-write', 'Create, edit or delete delivery-service configuration') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('ds-security-keys-read', 'View delivery-service security keys') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('ds-security-keys-write', 'Create, edit or delete delivery-service security keys') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('ds-stats-read', 'View delivery-service statistics') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('ds-steering-read', 'View delivery-service steering configuration') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('ds-steering-write', 'Create, edit or delete delivery-service steering configuration') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('federation-routing-read', 'View federation routing') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('federation-routing-write', 'Create, edit or delete federation routing') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('job-read', 'View jobs') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('job-write', 'Create, edit or delete jobs') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('params-read', 'View parameters') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('params-write', 'Create, edit or delete parameters') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('phys-location-read', 'View physical location configuration') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('phys-location-write', 'Create, edit or delete physical location configuration') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('profile-read', 'View profiles') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('profile-write', 'Create, edit or delete profiles') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('queue-updates-write', 'Queue updates to caches') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('region-read', 'View region configuration') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('region-write', 'Create, edit or delete region configuration') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('role-read', 'View role configuration') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('role-write', 'Create, edit or delete role configuration') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('security-keys-read', 'View security keys') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('security-keys-write', 'Create, edit or delete security keys') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('server-pull-updates-read', 'Read server update indication') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('server-pull-updates-write', 'Write server update indication') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('server-read', 'View server configuration') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('server-write', 'Create, edit or delete server configuration') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('static-dns-read', 'View static DNS configuration') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('static-dns-write', 'Create, edit or delete static DNS configuration') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('status-read', 'View the list of defined statuses') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('to-extension-read', 'View Traffic Ops extensions') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('to-extension-write', 'Create, edit or delete Traffic Ops extensions') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('type-read', 'View types configuration') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('type-write', 'Create, edit or delete type configuration') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('user-read', 'View user configuration') ON CONFLICT DO NOTHING;
insert into capability (name, description) values ('user-write', 'Create, edit or delete user configuration') ON CONFLICT DO NOTHING;

-- roles_capabilities
insert into role_capability (role_id, cap_name) values (10, 'all-read') ON CONFLICT DO NOTHING;
insert into role_capability (role_id, cap_name) values (10, 'all-write') ON CONFLICT DO NOTHING;

-- api_capabilities
insert into api_capability (http_method, route, capability) values ('GET', '/', 'all-read') ON CONFLICT DO NOTHING;
insert into api_capability (http_method, route, capability) values ('POST', '/', 'all-write') ON CONFLICT DO NOTHING;
insert into api_capability (http_method, route, capability) values ('PUT', '/', 'all-write') ON CONFLICT DO NOTHING;
insert into api_capability (http_method, route, capability) values ('PATCH', '/', 'all-write') ON CONFLICT DO NOTHING;
insert into api_capability (http_method, route, capability) values ('DELETE', '/', 'all-write') ON CONFLICT DO NOTHING;

insert into api_capability (http_method, route, capability) values ('GET', '/api/*/asns', 'asn-read') ON CONFLICT DO NOTHING; -- 4
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/asns/*', 'asn-read') ON CONFLICT DO NOTHING; -- 5
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/asns', 'asn-write') ON CONFLICT DO NOTHING; -- 6
insert into api_capability (http_method, route, capability) values ('PUT', '/api/*/asns/*', 'asn-write') ON CONFLICT DO NOTHING; -- 7
insert into api_capability (http_method, route, capability) values ('DELETE', '/api/*/asns/*', 'asn-write') ON CONFLICT DO NOTHING; -- 8
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/cache_stats', 'cache-stats-read') ON CONFLICT DO NOTHING; -- 11
insert into api_capability (http_method, route, capability) values ('GET', '/internal/api/*/daily_summary', 'cache-stats-read') ON CONFLICT DO NOTHING; -- 12
insert into api_capability (http_method, route, capability) values ('GET', '/internal/api/*/current_stats', 'cache-stats-read') ON CONFLICT DO NOTHING; -- 13
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/cachegroups', 'cache-group-read') ON CONFLICT DO NOTHING; -- 16
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/cachegroups/list', 'cache-group-read') ON CONFLICT DO NOTHING; -- 17
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/cachegroups/trimmed', 'cache-group-read') ON CONFLICT DO NOTHING; -- 18
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/cachegroups/*', 'cache-group-read') ON CONFLICT DO NOTHING; -- 19
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/cachegroups', 'cache-group-write') ON CONFLICT DO NOTHING; -- 20
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/cachegroups/create', 'cache-group-write') ON CONFLICT DO NOTHING; -- 21
insert into api_capability (http_method, route, capability) values ('PUT', '/api/*/cachegroups/*', 'cache-group-write') ON CONFLICT DO NOTHING; -- 22
insert into api_capability (http_method, route, capability) values ('PUT', '/api/*/cachegroups/*/update', 'cache-group-write') ON CONFLICT DO NOTHING; -- 23
insert into api_capability (http_method, route, capability) values ('DELETE', '/api/*/cachegroups/*', 'cache-group-write') ON CONFLICT DO NOTHING; -- 24
insert into api_capability (http_method, route, capability) values ('DELETE', '/api/*/cachegroups/*/delete', 'cache-group-write') ON CONFLICT DO NOTHING; -- 25
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/cachegroups/*/queue_update', 'queue-updates-write') ON CONFLICT DO NOTHING; -- 26
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/cachegroups/*/queue_update', 'queue-updates-write') ON CONFLICT DO NOTHING; -- 27
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/cachegroups/*/parameters', 'params-read') ON CONFLICT DO NOTHING; -- 28
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/cdns', 'cdn-read') ON CONFLICT DO NOTHING; -- 31
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/cdns/*', 'cdn-read') ON CONFLICT DO NOTHING; -- 32
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/cdns/name/*', 'cdn-read') ON CONFLICT DO NOTHING; -- 33
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/cdns', 'cdn-write') ON CONFLICT DO NOTHING; -- 34
insert into api_capability (http_method, route, capability) values ('PUT', '/api/*/cdns/*', 'cdn-write') ON CONFLICT DO NOTHING; -- 35
insert into api_capability (http_method, route, capability) values ('DELETE', '/api/*/cdns/*', 'cdn-write') ON CONFLICT DO NOTHING; -- 36
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/cdns/*/queue_update', 'queue-updates-write') ON CONFLICT DO NOTHING; -- 37
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/cdns/*/queue_update', 'queue-updates-write') ON CONFLICT DO NOTHING; -- 38
insert into api_capability (http_method, route, capability) values ('PUT', '/api/*/cdns/*/snapshot', 'cdn-config-snapshot-write') ON CONFLICT DO NOTHING; -- 40
insert into api_capability (http_method, route, capability) values ('PUT', '/api/*/snapshot/*', 'cdn-config-snapshot-write') ON CONFLICT DO NOTHING; -- 41
insert into api_capability (http_method, route, capability) values ('PUT', '/api/*/snapshot/*', 'cdn-config-snapshot-write') ON CONFLICT DO NOTHING; -- 42
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/cdns/configs', 'cdn-read') ON CONFLICT DO NOTHING; -- 44
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/cdns/*/configs/routing', 'cdn-read') ON CONFLICT DO NOTHING; -- 45
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/cdns/*/configs/monitoring', 'cdn-read') ON CONFLICT DO NOTHING; -- 46
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/cdns/domains', 'cdn-read') ON CONFLICT DO NOTHING; -- 47
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/cdns/health', 'cdn-health-read') ON CONFLICT DO NOTHING; -- 48
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/cdns/*/health', 'cdn-health-read') ON CONFLICT DO NOTHING; -- 49
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/cdns/capacity', 'cdn-health-read') ON CONFLICT DO NOTHING; -- 50
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/cdns/routing', 'cdn-read') ON CONFLICT DO NOTHING; -- 51
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/cdns/name/*/sslkeys', 'cdn-security-keys-read') ON CONFLICT DO NOTHING; -- 52
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/cdns/usage/overview', 'cdn-stats-read') ON CONFLICT DO NOTHING; -- 54
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/logs', 'change-log-read') ON CONFLICT DO NOTHING; -- 57
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/logs/*/days', 'change-log-read') ON CONFLICT DO NOTHING; -- 58
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/logs/newcount', 'change-log-read') ON CONFLICT DO NOTHING; -- 60
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/deliveryservices', 'ds-read') ON CONFLICT DO NOTHING; -- 69
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/deliveryservices/list', 'ds-read') ON CONFLICT DO NOTHING; -- 70
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/deliveryservices/*', 'ds-read') ON CONFLICT DO NOTHING; -- 71
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/deliveryservices/*/get', 'ds-read') ON CONFLICT DO NOTHING; -- 72
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/deliveryservices', 'ds-write') ON CONFLICT DO NOTHING; -- 73
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/*/deliveryservices/create', 'ds-write') ON CONFLICT DO NOTHING; -- 74
insert into api_capability (http_method, route, capability) values ('PUT', '/api/*/deliveryservices/*', 'ds-write') ON CONFLICT DO NOTHING; -- 75
insert into api_capability (http_method, route, capability) values ('PUT', '/api/*/*/deliveryservices/*/update', 'ds-write') ON CONFLICT DO NOTHING; -- 76
insert into api_capability (http_method, route, capability) values ('DELETE', '/api/*/deliveryservices/*', 'ds-write') ON CONFLICT DO NOTHING; -- 77
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/deliveryservices/*/health', 'ds-health-read') ON CONFLICT DO NOTHING; -- 78
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/deliveryservices/*/capacity', 'ds-health-read') ON CONFLICT DO NOTHING; -- 79
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/deliveryservices/*/routing', 'ds-read') ON CONFLICT DO NOTHING; -- 80
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/deliveryservices/*/state', 'ds-read') ON CONFLICT DO NOTHING; -- 81
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/deliveryservice_stats', 'ds-stats-read') ON CONFLICT DO NOTHING; -- 82
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/deliveryservices/request', 'ds-read') ON CONFLICT DO NOTHING; -- 83
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/deliveryservices/xmlId/*/sslkeys', 'ds-security-keys-read') ON CONFLICT DO NOTHING; -- 84
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/deliveryservices/hostname/#hostname/sslkeys', 'ds-security-keys-read') ON CONFLICT DO NOTHING; -- 85
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/deliveryservices/sslkeys/generate', 'ds-security-keys-write') ON CONFLICT DO NOTHING; -- 86
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/deliveryservices/sslkeys/add', 'ds-security-keys-write') ON CONFLICT DO NOTHING; -- 87
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/deliveryservices/xmlId/*/sslkeys/delete', 'ds-security-keys-write') ON CONFLICT DO NOTHING; -- 88
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/deliveryservices/xmlId/*/urlkeys', 'ds-security-keys-read') ON CONFLICT DO NOTHING; -- 89
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/deliveryservices/xmlId/*/urlkeys/generate', 'ds-security-keys-write') ON CONFLICT DO NOTHING; -- 90
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/deliveryservices/*/regexes', 'ds-read') ON CONFLICT DO NOTHING; -- 91
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/deliveryservice_matches', 'ds-read') ON CONFLICT DO NOTHING; -- 92
insert into api_capability (http_method, route, capability) values ('GET', '/internal/api/*/steering', 'ds-steering-read') ON CONFLICT DO NOTHING; -- 96
insert into api_capability (http_method, route, capability) values ('GET', '/internal/api/*/steering/*', 'ds-steering-read') ON CONFLICT DO NOTHING; -- 97
insert into api_capability (http_method, route, capability) values ('POST', '/internal/api/*/steering', 'ds-steering-write') ON CONFLICT DO NOTHING; -- 98
insert into api_capability (http_method, route, capability) values ('PUT', '/internal/api/*/steering/*', 'ds-steering-write') ON CONFLICT DO NOTHING; -- 99
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/deliveryserviceserver', 'ds-cache-read') ON CONFLICT DO NOTHING; -- 103
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/deliveryservices/*/servers', 'ds-cache-write') ON CONFLICT DO NOTHING; -- 106
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/deliveryservices_regexes', 'ds-read') ON CONFLICT DO NOTHING; -- 109
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/deliveryservices/*/regexes', 'ds-read') ON CONFLICT DO NOTHING; -- 110
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/deliveryservices/*/regexes/*', 'ds-read') ON CONFLICT DO NOTHING; -- 111
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/deliveryservices/*/regexes', 'ds-write') ON CONFLICT DO NOTHING; -- 112
insert into api_capability (http_method, route, capability) values ('PUT', '/api/*/deliveryservices/*/regexes/*', 'ds-write') ON CONFLICT DO NOTHING; -- 113
insert into api_capability (http_method, route, capability) values ('DELETE', '/api/*/deliveryservices/*/regexes/*', 'ds-write') ON CONFLICT DO NOTHING; -- 114
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/divisions', 'division-read') ON CONFLICT DO NOTHING; -- 120
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/divisions/*', 'division-read') ON CONFLICT DO NOTHING; -- 121
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/divisions', 'division-write') ON CONFLICT DO NOTHING; -- 122
insert into api_capability (http_method, route, capability) values ('PUT', '/api/*/divisions/*', 'division-write') ON CONFLICT DO NOTHING; -- 123
insert into api_capability (http_method, route, capability) values ('DELETE', '/api/*/divisions/*', 'division-write') ON CONFLICT DO NOTHING; -- 124
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/cdns/name/*/dnsseckeys', 'cdn-security-keys-read') ON CONFLICT DO NOTHING; -- 127
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/cdns/dnsseckeys/generate', 'cdn-security-keys-write') ON CONFLICT DO NOTHING; -- 128
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/cdns/name/*/dnsseckeys/delete', 'cdn-security-keys-write') ON CONFLICT DO NOTHING; -- 129
insert into api_capability (http_method, route, capability) values ('GET', '/internal/api/*/cdns/dnsseckeys/refresh', 'cdn-security-keys-read') ON CONFLICT DO NOTHING; -- 130
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/to_extensions', 'to-extension-read') ON CONFLICT DO NOTHING; -- 134
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/to_extensions', 'to-extension-write') ON CONFLICT DO NOTHING; -- 135
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/to_extensions/*/delete', 'to-extension-write') ON CONFLICT DO NOTHING; -- 136
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/federations', 'federation-routing-read') ON CONFLICT DO NOTHING; -- 139
insert into api_capability (http_method, route, capability) values ('GET', '/internal/api/*/federations', 'federation-routing-read') ON CONFLICT DO NOTHING; -- 140
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/federations', 'federation-routing-write') ON CONFLICT DO NOTHING; -- 141
insert into api_capability (http_method, route, capability) values ('PUT', '/api/*/federations', 'federation-routing-write') ON CONFLICT DO NOTHING; -- 142
insert into api_capability (http_method, route, capability) values ('DELETE', '/api/*/federations', 'federation-routing-write') ON CONFLICT DO NOTHING; -- 143
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/hwinfo', 'all-read') ON CONFLICT DO NOTHING; -- 148
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/parameters', 'params-read') ON CONFLICT DO NOTHING; -- 164
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/parameters/*', 'params-read') ON CONFLICT DO NOTHING; -- 168
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/parameters', 'params-write') ON CONFLICT DO NOTHING; -- 169
insert into api_capability (http_method, route, capability) values ('PUT', '/api/*/parameters/*', 'params-write') ON CONFLICT DO NOTHING; -- 170
insert into api_capability (http_method, route, capability) values ('DELETE', '/api/*/parameters/*', 'params-write') ON CONFLICT DO NOTHING; -- 171
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/parameters/*/validate', 'params-write') ON CONFLICT DO NOTHING; -- 172
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/profiles/*/parameters', 'params-read') ON CONFLICT DO NOTHING; -- 173
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/profiles/name/*/parameters', 'params-read') ON CONFLICT DO NOTHING; -- 174a
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/parameters/profile', 'params-read') ON CONFLICT DO NOTHING; -- 174b
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/profiles/name/*/parameters', 'params-write') ON CONFLICT DO NOTHING; -- 175
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/profiles/*/parameters', 'params-write') ON CONFLICT DO NOTHING; -- 176
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/profileparameters', 'params-read') ON CONFLICT DO NOTHING; -- 181
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/profileparameters', 'params-write') ON CONFLICT DO NOTHING; -- 182
insert into api_capability (http_method, route, capability) values ('DELETE', '/api/*/profileparameters/*/*', 'params-write') ON CONFLICT DO NOTHING; -- 183
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/cachegroupparameters', 'params-read') ON CONFLICT DO NOTHING; -- 186
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/phys_locations', 'phys-location-read') ON CONFLICT DO NOTHING; -- 191
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/phys_locations/trimmed', 'phys-location-read') ON CONFLICT DO NOTHING; -- 192
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/phys_locations/*', 'phys-location-read') ON CONFLICT DO NOTHING; -- 193
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/phys_locations', 'phys-location-write') ON CONFLICT DO NOTHING; -- 194
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/regions/*/phys_locations', 'phys-location-write') ON CONFLICT DO NOTHING; -- 195
insert into api_capability (http_method, route, capability) values ('PUT', '/api/*/phys_locations/*', 'phys-location-write') ON CONFLICT DO NOTHING; -- 196
insert into api_capability (http_method, route, capability) values ('DELETE', '/api/*/phys_locations/*', 'phys-location-write') ON CONFLICT DO NOTHING; -- 197
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/profiles', 'profile-read') ON CONFLICT DO NOTHING; -- 200
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/profiles/trimmed', 'profile-read') ON CONFLICT DO NOTHING; -- 201
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/profiles/*', 'profile-read') ON CONFLICT DO NOTHING; -- 202
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/profiles', 'profile-write') ON CONFLICT DO NOTHING; -- 203
insert into api_capability (http_method, route, capability) values ('PUT', '/api/*/profiles/*', 'profile-write') ON CONFLICT DO NOTHING; -- 204
insert into api_capability (http_method, route, capability) values ('DELETE', '/api/*/profiles/*', 'profile-write') ON CONFLICT DO NOTHING; -- 205
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/profiles/name/*/copy/*', 'profile-write') ON CONFLICT DO NOTHING; -- 206
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/regions', 'region-read') ON CONFLICT DO NOTHING; -- 213
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/regions/*', 'region-read') ON CONFLICT DO NOTHING; -- 214
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/regions', 'region-write') ON CONFLICT DO NOTHING; -- 215
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/divisions/*/regions', 'region-write') ON CONFLICT DO NOTHING; -- 216
insert into api_capability (http_method, route, capability) values ('PUT', '/api/*/regions/*', 'region-write') ON CONFLICT DO NOTHING; -- 217
insert into api_capability (http_method, route, capability) values ('DELETE', '/api/*/regions/*', 'region-write') ON CONFLICT DO NOTHING; -- 218
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/riak/ping', 'cdn-security-keys-write') ON CONFLICT DO NOTHING; -- 221
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/keys/ping', 'security-keys-write') ON CONFLICT DO NOTHING; -- 222
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/riak/bucket/#bucket/key/#key/values', 'security-keys-read') ON CONFLICT DO NOTHING; -- 223
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/riak/stats', 'security-keys-read') ON CONFLICT DO NOTHING; -- 224
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/roles', 'role-read') ON CONFLICT DO NOTHING; -- 227
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/servers', 'server-read') ON CONFLICT DO NOTHING; -- 230
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/servers?dsId=*', 'server-read') ON CONFLICT DO NOTHING; -- 231
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/servers?type=*', 'server-read') ON CONFLICT DO NOTHING; -- 232
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/servers?status=*', 'server-read') ON CONFLICT DO NOTHING; -- 233
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/servers?profileId=*', 'server-read') ON CONFLICT DO NOTHING; -- 234
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/servers/*', 'server-read') ON CONFLICT DO NOTHING; -- 235
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/servers', 'server-write') ON CONFLICT DO NOTHING; -- 237
insert into api_capability (http_method, route, capability) values ('PUT', '/api/*/servers/*', 'server-write') ON CONFLICT DO NOTHING; -- 238
insert into api_capability (http_method, route, capability) values ('DELETE', '/api/*/servers/*', 'server-write') ON CONFLICT DO NOTHING; -- 239
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/servers/details', 'server-read') ON CONFLICT DO NOTHING; -- 247
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/servers/totals', 'server-read') ON CONFLICT DO NOTHING; -- 249
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/servers/checks', 'server-read') ON CONFLICT DO NOTHING; -- 250a
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/servercheck/aadata', 'server-read') ON CONFLICT DO NOTHING; -- 250b
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/servercheck', 'server-write') ON CONFLICT DO NOTHING; -- 251
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/servers/*/queue_update', 'queue-updates-write') ON CONFLICT DO NOTHING; -- 252
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/servers/*/queue_update', 'queue-updates-write') ON CONFLICT DO NOTHING; -- 253
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/*/stats_summary', 'cdn-stats-read') ON CONFLICT DO NOTHING; -- 258
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/*/stats_summary?lastSummaryDate=true', 'cdn-stats-read') ON CONFLICT DO NOTHING; -- 259
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/*/stats_summary/create', 'cdn-stats-write') ON CONFLICT DO NOTHING; -- 260
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/statuses', 'status-read') ON CONFLICT DO NOTHING; -- 263
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/statuses/*', 'status-read') ON CONFLICT DO NOTHING; -- 264
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/staticdnsentries', 'static-dns-read') ON CONFLICT DO NOTHING; -- 270
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/system/info', 'basic-read') ON CONFLICT DO NOTHING; -- 275
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/types', 'type-read') ON CONFLICT DO NOTHING; -- 278
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/types/trimmed', 'type-read') ON CONFLICT DO NOTHING; -- 279
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/types/*', 'type-read') ON CONFLICT DO NOTHING; -- 280
insert into api_capability (http_method, route, capability) values ('POST', '/api/*/types', 'type-write') ON CONFLICT DO NOTHING; -- 281
insert into api_capability (http_method, route, capability) values ('PUT', '/api/*/types/*', 'type-write') ON CONFLICT DO NOTHING; -- 282
insert into api_capability (http_method, route, capability) values ('DELETE', '/api/*/types/*', 'type-write') ON CONFLICT DO NOTHING; -- 283
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/users', 'user-read') ON CONFLICT DO NOTHING; -- 289
insert into api_capability (http_method, route, capability) values ('GET', '/api/*/users/*', 'user-read') ON CONFLICT DO NOTHING; -- 290
insert into api_capability (http_method, route, capability) values ('PUT', '/api/*/users/*', 'user-write') ON CONFLICT DO NOTHING; -- 292


-- types
-- delivery service types
insert into type (name, description, use_in_table) values ('HTTP', 'HTTP Content Routing', 'deliveryservice') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('HTTP_NO_CACHE', 'HTTP Content Routing, no caching', 'deliveryservice') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('HTTP_LIVE', 'HTTP Content routing cache in RAM', 'deliveryservice') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('HTTP_LIVE_NATNL', 'HTTP Content routing, RAM cache, National', 'deliveryservice') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('DNS', 'DNS Content Routing', 'deliveryservice') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('DNS_LIVE', 'DNS Content routing, RAM cache, Local', 'deliveryservice') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('DNS_LIVE_NATNL', 'DNS Content routing, RAM cache, National', 'deliveryservice') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('ANY_MAP', 'No Content Routing - arbitrary remap at the edge, no Traffic Router config', 'deliveryservice') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('STEERING', 'Steering Delivery Service', 'deliveryservice') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('CLIENT_STEERING', 'Client-Controlled Steering Delivery Service', 'deliveryservice') ON CONFLICT (name) DO NOTHING;

-- server types
insert into type (name, description, use_in_table) values ('EDGE', 'Edge Cache', 'server') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('MID', 'Mid Tier Cache', 'server') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('ORG', 'Origin', 'server') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('CCR', 'Traffic Router', 'server') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('RASCAL', 'Traffic Monitor', 'server') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('RIAK', 'Riak keystore', 'server') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('INFLUXDB', 'influxDb server', 'server') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('TRAFFIC_ANALYTICS', 'traffic analytics server', 'server') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('TRAFFIC_OPS', 'traffic ops server', 'server') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('TRAFFIC_OPS_DB', 'traffic ops DB server', 'server') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('TRAFFIC_PORTAL', 'traffic portal server', 'server') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('TRAFFIC_STATS', 'traffic stats server', 'server') ON CONFLICT (name) DO NOTHING;

-- cachegroup types
insert into type (name, description, use_in_table) values ('EDGE_LOC', 'Edge Logical Location', 'cachegroup') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('MID_LOC', 'Mid Logical Location', 'cachegroup') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('ORG_LOC', 'Origin Logical Site', 'cachegroup') ON CONFLICT (name) DO NOTHING;

-- to_extension types
insert into type (name, description, use_in_table) values ('CHECK_EXTENSION_BOOL', 'Extension for checkmark in Server Check', 'to_extension') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('CHECK_EXTENSION_NUM', 'Extension for int value in Server Check', 'to_extension') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('CHECK_EXTENSION_OPEN_SLOT', 'Open slot for check in Server Status', 'to_extension') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('CONFIG_EXTENSION', 'Extension for additional configuration file', 'to_extension') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('STATISTIC_EXTENSION', 'Extension source for 12M graphs', 'to_extension') ON CONFLICT (name) DO NOTHING;

-- regex types
insert into type (name, description, use_in_table) values ('HOST_REGEXP', 'Host header regular expression', 'regex') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('HEADER_REGEXP', 'HTTP header regular expression', 'regex') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('PATH_REGEXP', 'URL path regular expression', 'regex') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('STEERING_REGEXP', 'Steering target filter regular expression', 'regex') ON CONFLICT (name) DO NOTHING;

-- federation types
insert into type (name, description, use_in_table) values ('RESOLVE4', 'federation type resolve4', 'federation') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('RESOLVE6', 'federation type resolve6', 'federation') ON CONFLICT (name) DO NOTHING;

-- static dns entry types
insert into type (name, description, use_in_table) values ('A_RECORD', 'Static DNS A entry', 'staticdnsentry') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('AAAA_RECORD', 'Static DNS AAAA entry', 'staticdnsentry') ON CONFLICT (name) DO NOTHING;
insert into type (name, description, use_in_table) values ('CNAME_RECORD', 'Static DNS CNAME entry', 'staticdnsentry') ON CONFLICT (name) DO NOTHING;

-- users
insert into tm_user (username, role, full_name, token) values ('extension', 3, 'Extension User, DO NOT DELETE', '91504CE6-8E4A-46B2-9F9F-FE7C15228498') ON CONFLICT DO NOTHING;

-- to extensions
-- some of the old ones do not get a new place, and there will be 'gaps' in the column usage.... New to_extension add will have to take care of that.
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (1, 'ILO_PING', 'ILO', 'aa', '1.0.0', '-', 'ToPingCheck.pl', '1', '{ check_name: "ILO", "base_url": "https://localhost", "select": "ilo_ip_address", "cron": "9 * * * *" }',
        (select id from type where name='CHECK_EXTENSION_BOOL') ) ON CONFLICT DO NOTHING;
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (2, '10G_PING', '10G', 'ab', '1.0.0', '-', 'ToPingCheck.pl', '1', '{ check_name: "10G", "base_url": "https://localhost", "select": "ip_address", "cron": "18 * * * *" }',
        (select id from type where name='CHECK_EXTENSION_BOOL') ) ON CONFLICT DO NOTHING;
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (3, 'FQDN_PING', 'FQDN', 'ac', '1.0.0', '-', 'ToPingCheck.pl', '1', '{ check_name: "FQDN", "base_url": "https://localhost", "select": "host_name", "cron": "27 * * * *" }',
        (select id from type where name='CHECK_EXTENSION_BOOL') ) ON CONFLICT DO NOTHING;
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (4, 'CHECK_DSCP', 'DSCP', 'ad', '1.0.0', '-', 'ToDSCPCheck.pl', '1', '{ "check_name": "DSCP", "base_url": "https://localhost", "cron": "36 * * * *" }',
        (select id from type where name='CHECK_EXTENSION_BOOL') ) ON CONFLICT DO NOTHING;
-- open EF
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (5, 'OPEN', '', 'ae', '1.0.0', '-', '', '0', '', (select id from type where name='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (6, 'OPEN', '', 'af', '1.0.0', '-', '', '0', '', (select id from type where name='CHECK_EXTENSION_OPEN_SLOT'));
--
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (7, 'IPV6_PING', '10G6', 'ag', '1.0.0', '-', 'ToPingCheck.pl', '1', '{ "select": "ip6_address", "cron": "0 * * * *" }',
        (select id from type where name='CHECK_EXTENSION_BOOL') ) ON CONFLICT DO NOTHING;
-- upd_pending H -> open
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (8, 'OPEN', '', 'ah', '1.0.0', '', '', '0', '', (select id from type where name='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
-- open IJ
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (9, 'OPEN', '', 'ai', '1.0.0', '', '', '0', '', (select id from type where name='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (10, 'OPEN', '', 'aj', '1.0.0', '', '', '0', '', (select id from type where name='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
--
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (11, 'CHECK_MTU', 'MTU', 'ak', '1.0.0', '-', 'ToMtuCheck.pl', '1', '{ "check_name": "MTU", "base_url": "https://localhost", "cron": "45 * * * *" }',
        (select id from type where name='CHECK_EXTENSION_BOOL') ) ON CONFLICT DO NOTHING;
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (12, 'CHECK_TRAFFIC_ROUTER_STATUS', 'RTR', 'al', '1.0.0', '-', 'ToRTRCheck.pl', '1', '{  "check_name": "RTR", "base_url": "https://localhost", "cron": "10 * * * *" }',
        (select id from type where name='CHECK_EXTENSION_BOOL') ) ON CONFLICT DO NOTHING;
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (13, 'OPEN', '', 'am', '1.0.0', '-', '', '0', '', (select id from type where name='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (14, 'CACHE_HIT_RATIO_LAST_15', 'CHR', 'an', '1.0.0', '-', 'ToCHRCheck.pl', '1', '{ check_name: "CHR", "base_url": "https://localhost", cron": "0,15,30,45 * * * *" }',
        (select id from type where name='CHECK_EXTENSION_NUM') ) ON CONFLICT DO NOTHING;
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (15, 'DISK_UTILIZATION', 'CDU', 'ao', '1.0.0', '-', 'ToCDUCheck.pl', '1', '{ check_name: "CDU", "base_url": "https://localhost", cron": "20 * * * *" }',
        (select id from type where name='CHECK_EXTENSION_NUM') ) ON CONFLICT DO NOTHING;
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (16, 'ORT_ERROR_COUNT', 'ORT', 'ap', '1.0.0', '-', 'ToORTCheck.pl', '1', '{ check_name: "ORT", "base_url": "https://localhost", "cron": "40 * * * *" }',
        (select id from type where name='CHECK_EXTENSION_NUM') ) ON CONFLICT DO NOTHING;
-- rest open
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (17, 'OPEN', '', 'aq', '1.0.0', '-', '', '0', '', (select id from type where name='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (18, 'OPEN', '', 'ar', '1.0.0', '-', '', '0', '', (select id from type where name='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (19, 'OPEN', '', 'bf', '1.0.0', '-', '', '0', '', (select id from type where name='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (20, 'OPEN', '', 'at', '1.0.0', '-', '', '0', '', (select id from type where name='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (21, 'OPEN', '', 'au', '1.0.0', '-', '', '0', '', (select id from type where name='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (22, 'OPEN', '', 'av', '1.0.0', '-', '', '0', '', (select id from type where name='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (23, 'OPEN', '', 'aw', '1.0.0', '-', '', '0', '', (select id from type where name='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (24, 'OPEN', '', 'ax', '1.0.0', '-', '', '0', '', (select id from type where name='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (25, 'OPEN', '', 'ay', '1.0.0', '-', '', '0', '', (select id from type where name='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (26, 'OPEN', '', 'az', '1.0.0', '-', '', '0', '', (select id from type where name='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (27, 'OPEN', '', 'ba', '1.0.0', '-', '', '0', '', (select id from type where name='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (28, 'OPEN', '', 'bb', '1.0.0', '-', '', '0', '', (select id from type where name='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (29, 'OPEN', '', 'bc', '1.0.0', '-', '', '0', '', (select id from type where name='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (30, 'OPEN', '', 'bd', '1.0.0', '-', '', '0', '', (select id from type where name='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
insert into to_extension (id, name, servercheck_short_name, servercheck_column_name, version, info_url, script_file, isactive, additional_config_json, type)
values (31, 'OPEN', '', 'be', '1.0.0', '-', '', '0', '', (select id from type where name='CHECK_EXTENSION_OPEN_SLOT')) ON CONFLICT DO NOTHING;
