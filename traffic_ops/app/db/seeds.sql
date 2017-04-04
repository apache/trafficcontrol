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

-- roles
insert into role (id, name, description, priv_level) values (1, 'disallowed', 'Block all access',0) ON CONFLICT DO NOTHING;
insert into role (id, name, description, priv_level) values (2, 'read-only user', 'Block all access', 10) ON CONFLICT DO NOTHING;
insert into role (id, name, description, priv_level) values (3, 'operations', 'Block all access', 20) ON CONFLICT DO NOTHING;
insert into role (id, name, description, priv_level) values (4, 'admin', 'super-user', 30) ON CONFLICT DO NOTHING;
insert into role (id, name, description, priv_level) values (5, 'portal', 'Portal User', 2) ON CONFLICT DO NOTHING;
insert into role (id, name, description, priv_level) values (6, 'migrations', 'database migrations user - DO NOT REMOVE', 20) ON CONFLICT DO NOTHING;
insert into role (id, name, description, priv_level) values (7, 'federation', 'Role for Secondary CZF', 15) ON CONFLICT DO NOTHING;
insert into role (id, name, description, priv_level) values (8, 'steering', 'Role for Steering Delivery Services', 15) ON CONFLICT DO NOTHING;

-- types
-- delivery service types
insert into type (name, description, use_in_table) values ('HTTP', 'HTTP Content Routing', 'deliveryservice') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('HTTP_NO_CACHE', 'HTTP Content Routing, no caching', 'deliveryservice') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('HTTP_LIVE', 'HTTP Content routing cache in RAM', 'deliveryservice') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('HTTP_LIVE_NATNL', 'HTTP Content routing, RAM cache, National', 'deliveryservice') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('DNS', 'DNS Content Routing', 'deliveryservice') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('DNS_LIVE', 'DNS Content routing, RAM cache, Local', 'deliveryservice') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('DNS_LIVE_NATNL', 'DNS Content routing, RAM cache, National', 'deliveryservice') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('ANY_MAP', 'No Content Routing - arbitrary remap at the edge, no Traffic Router config', 'deliveryservice') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('STEERING', 'Steering Delivery Service', 'deliveryservice') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('CLIENT_STEERING', 'Client-Controlled Steering Delivery Service', 'deliveryservice') ON CONFLICT DO NOTHING;

-- server types
insert into type (name, description, use_in_table) values ('EDGE', 'Edge Cache', 'server') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('MID', 'Mid Tier Cache', 'server') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('ORG', 'Origin', 'server') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('CCR', 'Traffic Router', 'server') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('RASCAL', 'Traffic Monitor', 'server') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('RIAK', 'Riak keystore', 'server') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('INFLUXDB', 'influxDb server', 'server') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('TRAFFIC_ANALYTICS', 'traffic_analytics server', 'server') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('TRAFFIC_PORTAL', 'traffic_portal server', 'server') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('TRAFFIC_STATS', 'traffic_stats server', 'server') ON CONFLICT DO NOTHING;

-- cachegroup types
insert into type (name, description, use_in_table) values ('EDGE_LOC', 'Edge Logical Location', 'cachegroup') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('MID_LOC', 'Mid Logical Location', 'cachegroup') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('ORG_LOC', 'Origin Logical Site', 'cachegroup') ON CONFLICT DO NOTHING;

-- to_extension types
insert into type (name, description, use_in_table) values ('CHECK_EXTENSION_BOOL', 'Extension for checkmark in Server Check', 'to_extension') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('CHECK_EXTENSION_NUM', 'Extension for int value in Server Check', 'to_extension') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('CHECK_EXTENSION_OPEN_SLOT', 'Open slot for check in Server Status', 'to_extension') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('CONFIG_EXTENSION', 'Extension for additional configuration file', 'to_extension') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('STATISTIC_EXTENSION', 'Extension source for 12M graphs', 'to_extension') ON CONFLICT DO NOTHING;

-- regex types
insert into type (name, description, use_in_table) values ('HOST_REGEXP', 'Host header regular expression', 'regex') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('HEADER_REGEXP', 'HTTP header regular expression', 'regex') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('PATH_REGEXP', 'URL path regular expression', 'regex') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('STEERING_REGEXP', 'Steering target filter regular expression', 'regex') ON CONFLICT DO NOTHING;

-- federation types
insert into type (name, description, use_in_table) values ('RESOLVE4', 'federation type resolve4', 'federation') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('RESOLVE6', 'federation type resolve6', 'federation') ON CONFLICT DO NOTHING;

-- static dns entry types
insert into type (name, description, use_in_table) values ('A_RECORD', 'Static DNS A entry', 'staticdnsentry') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('AAAA_RECORD', 'Static DNS AAAA entry', 'staticdnsentry') ON CONFLICT DO NOTHING;
insert into type (name, description, use_in_table) values ('CNAME_RECORD', 'Static DNS CNAME entry', 'staticdnsentry') ON CONFLICT DO NOTHING;

-- statuses
insert into status (id, name, description) values (1, 'OFFLINE', 'Server is Offline. Not active in any configuration.') ON CONFLICT DO NOTHING;
insert into status (id, name, description) values (2, 'ONLINE', 'Server is online.') ON CONFLICT DO NOTHING;
insert into status (id, name, description) values (3, 'REPORTED', 'Server is online and reporeted in the health protocol.') ON CONFLICT DO NOTHING;
insert into status (id, name, description) values (4, 'ADMIN_DOWN', 'Sever is administrative down and does not receive traffic.') ON CONFLICT DO NOTHING;
insert into status (id, name, description) values (5, 'CCR_IGNORE', 'Server is ignored by traffic router.') ON CONFLICT DO NOTHING;
insert into status (id, name, description) values (6, 'PRE_PROD', 'Pre Production. Not active in any configuration.') ON CONFLICT DO NOTHING;

-- job agents
insert into job_agent (name, description, active) values ('dummy', 'Description of Purge Agent', '1') ON CONFLICT DO NOTHING;

-- job statuses
insert into job_status (name, description) values ('PENDING', 'Job is queued, but has not been picked up by any agents yet') ON CONFLICT DO NOTHING;
insert into job_status (name, description) values ('IN_PROGRESS', 'Job is being processed by agents') ON CONFLICT DO NOTHING;
insert into job_status (name, description) values ('COMPLETED', 'Job has finished') ON CONFLICT DO NOTHING;
insert into job_status (name, description) values ('CANCELLED', 'Job was cancelled') ON CONFLICT DO NOTHING;
insert into job_status (name, description) values ('PURGE', 'Initial Purge state') ON CONFLICT DO NOTHING;


-- profiles
insert into profile (name, description, type) values ('GLOBAL', 'Global Traffic Ops profile, DO NOT DELETE', 'UNK_PROFILE') ON CONFLICT DO NOTHING;
insert into profile (name, description, type) values ('RIAK_ALL', 'Riak profile for all CDNs', 'RIAK_PROFILE') ON CONFLICT DO NOTHING;
insert into profile (name, description, type) values ('TRAFFIC_STATS', 'Traffic_Stats profile', 'TS_PROFILE') ON CONFLICT DO NOTHING;
insert into profile (name, description, type) values ('TRAFFIC_PORTAL', 'Traffic_Portal profile', 'TP_PROFILE') ON CONFLICT DO NOTHING;
insert into profile (name, description, type) values ('INFLUXDB', 'InfluxDb profile', 'INFLUXDB_PROFILE') ON CONFLICT DO NOTHING;
insert into profile (name, description, type) values ('CCR_CDN', 'Kabletown Content Router for cdn1', 'TR_PROFILE') ON CONFLICT DO NOTHING;
insert into profile (name, description, type) values ('RASCAL_CDN', 'Traffic Monitor profile for cdn1', 'TM_PROFILE') ON CONFLICT DO NOTHING;
insert into profile (name, description, type) values ('HOSTMON', 'Hostmon pipe host', 'UNK_PROFILE') ON CONFLICT DO NOTHING;
insert into profile (name, description, type) values ('EDGE1_CDN_520', 'Dell R720xd, Edge, CDN, ATS v5.2.0', 'ATS_PROFILE') ON CONFLICT DO NOTHING;
insert into profile (name, description, type) values ('MID1_CDN_421', 'Dell R720xd, Mid, CDN, ATS v4.2.1', 'ATS_PROFILE') ON CONFLICT DO NOTHING;

-- parameters
insert into parameter (name, config_file, value) values ('ttl_max_hours', 'regex_revalidate.config', '672') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('ttl_min_hours', 'regex_revalidate.config', '48') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('maxRevalDurationDays', 'regex_revalidate.config', '90') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('ramdisk_size', 'grub.conf', 'ramdisk_size=16777216') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('monitor:///opt/tomcat/logs/access.log', 'inputs.conf', 'index=index_test;sourcetype=access_ccr') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('purge_allow_ip', 'ip_allow.config', '10.10.10.10') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('health.threshold.loadavg', 'rascal.properties', '25.0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('health.threshold.availableBandwidthInKbps', 'rascal.properties', '>1750000') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('health.threshold.queryTime', 'rascal.properties', '1000') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('health.polling.url', 'rascal.properties', 'http://${hostname}/_astats?application=&inf.name=${interface_name}') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('health.threadPool', 'rascal-config.txt', '4') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('health.event-count', 'rascal-config.txt', '200') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('health.polling.interval', 'rascal-config.txt', '8000') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('health.timepad', 'rascal-config.txt', '30') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('history.count', 'rascal.properties', '30') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('cacheurl.so', 'plugin.config', '') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('LogObject.RollingEnabled', 'logs_xml.config', '3') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('LogObject.RollingIntervalSec', 'logs_xml.config', '86400') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('LogObject.RollingOffsetHr', 'logs_xml.config', '11') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('LogObject.RollingSizeMb', 'logs_xml.config', '1024') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('tm.dataServer.polling.url', 'rascal-config.txt', 'https://${tmHostname}/dataserver/orderby/id') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('tm.healthParams.polling.url', 'rascal-config.txt', 'https://${tmHostname}/health/${cdnName}') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('tm.polling.interval', 'rascal-config.txt', '60000') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('tm.url', 'global', 'https://tm.kabletown.net/') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('tm.toolname', 'global', 'Traffic Ops') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('tm.infourl', 'global', 'http://docs.cdnl.kabletown.net/traffic_control/html/') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('tm.logourl', 'global', '/images/tc_logo.png') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('tm.instance_name', 'global', 'kabletown CDN') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('tm.traffic_mon_fwd_proxy', 'global', 'http://tm.kabletown.net:81') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('hack.ttl', 'rascal-config.txt', '30') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('health.connection.timeout', 'rascal.properties', '2000') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('geolocation.polling.url', 'CRConfig.json', 'https://tm.kabletown.net/MaxMind/auto/GeoIP2-City.mmdb.gz') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('geolocation.polling.interval', 'CRConfig.json', '86400000') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('coveragezone.polling.interval', 'CRConfig.json', '86400000') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('coveragezone.polling.url', 'CRConfig.json', 'http://cdn-tools.cdnl.kabletown.net/cdn/CZF/current/kabletown_cdn_czf-current.json') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('domain_name', 'CRConfig.json', 'cdn.kabletown.net') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('tld.ttls.AAAA', 'CRConfig.json', '3600') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('tld.ttls.SOA', 'CRConfig.json', '86400') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('tld.ttls.A', 'CRConfig.json', '3600') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('tld.ttls.NS', 'CRConfig.json', '3600') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('tld.soa.expire', 'CRConfig.json', '604800') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('tld.soa.minimum', 'CRConfig.json', '86400') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('tld.soa.admin', 'CRConfig.json', 'twelve_monkeys') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('tld.soa.retry', 'CRConfig.json', '7200') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('tld.soa.refresh', 'CRConfig.json', '28800') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('trafficserver', 'chkconfig', '0:off\t1:off\t2:on\t3:on\t4:on\t5:on\t6:off') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('regex_revalidate.so', 'plugin.config', '--config regex_revalidate.config') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('remap_stats.so', 'plugin.config', '') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('LogFormat.Format', 'logs_xml.config', '%<cqtq> chi=%<chi> phn=%<phn> shn=%<shn> url=%<cquuc> cqhm=%<cqhm> cqhv=%<cqhv> pssc=%<pssc> ttms=%<ttms> b=%<pscl> sssc=%<sssc> sscl=%<sscl> cfsc=%<cfsc> pfsc=%<pfsc> crc=%<crc> phr=%<phr> uas="%<{User-Agent}cqh>"') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('LogFormat.Name', 'logs_xml.config', 'custom_ats_2') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('LogObject.Format', 'logs_xml.config', 'custom_ats_2') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('LogObject.Filename', 'logs_xml.config', 'custom_ats_2') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('RAM_Drive_Prefix', 'storage.config', '/dev/ram') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('RAM_Drive_Letters', 'storage.config', '0,1,2,3,4,5,6,7') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('RAM_Volume', 'storage.config', '2') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('Disk_Volume', 'storage.config', '1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('Drive_Prefix', 'storage.config', '/dev/sd') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('Drive_Letters', 'storage.config', 'b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q,r,s,t,u,v,w,x,y') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('api.port', 'server.xml', '3333') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('astats_over_http.so', 'plugin.config', '') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('allow_ip', 'astats.config', '127.0.0.1,10.10.10.10/16') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('allow_ip6', 'astats.config', '::1,d009:5dd:f0d8:18d::2/64') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('record_types', 'astats.config', '144') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('path', 'astats.config', '_astats') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('algorithm', 'parent.config', 'consistent_hash') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('record_types', 'astats.config', '122') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('cron_ort_syncds', 'traffic_ops_ort_syncds.cron', '*/15 * * * * root /opt/ort/traffic_ops_ort.pl syncds warn https://tm.kabletown.net admin:password > /tmp/ort/syncds.log 2>&1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('api.cache-control.max-age', 'CRConfig.json', '10') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('qstring', 'parent.config', 'ignore') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('geolocation.polling.url', 'CRConfig.json', 'http://cdn-tools.cdnl.kabletown.net/cdn/MaxMind/GeoLiteCity.dat.gz') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('geolocation6.polling.url', 'CRConfig.json', 'http://cdn-tools.cdnl.kabletown.net/cdn/MaxMind/GeoLiteCityv6.dat.gz') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('kickstart.files.location', 'mkisofs', '__KICKSTART-FILE__') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CacheStats', 'traffic_stats.config', 'bandwidth') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CacheStats', 'traffic_stats.config', 'maxKbps') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CacheStats', 'traffic_stats.config', 'ats.proxy.process.http.current_client_connections') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('DsStats', 'traffic_stats.config', 'kbps') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('DsStats', 'traffic_stats.config', 'status_4xx') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('DsStats', 'traffic_stats.config', 'status_5xx') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('DsStats', 'traffic_stats.config', 'tps_2xx') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('DsStats', 'traffic_stats.config', 'tps_3xx') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('DsStats', 'traffic_stats.config', 'tps_4xx') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('DsStats', 'traffic_stats.config', 'tps_5xx') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('DsStats', 'traffic_stats.config', 'tps_total') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('tm.crConfig.polling.url', 'rascal-config.txt', 'https://${tmHostname}/CRConfig-Snapshots/${cdnName}/CRConfig.xml') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'cache.config', '/opt/trafficserver/etc/trafficserver/') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'hosting.config', '/opt/trafficserver/etc/trafficserver/') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'parent.config', '/opt/trafficserver/etc/trafficserver/') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'plugin.config', '/opt/trafficserver/etc/trafficserver/') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'records.config', '/opt/trafficserver/etc/trafficserver/') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'remap.config', '/opt/trafficserver/etc/trafficserver/') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'storage.config', '/opt/trafficserver/etc/trafficserver/') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'volume.config', '/opt/trafficserver/etc/trafficserver/') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', '50-ats.rules', '/etc/udev/rules.d/') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'dns.zone', '/etc/kabletown/zones/<zonename>.info') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'http-log4j.properties', '/etc/kabletown') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'dns-log4j.properties', '/etc/kabletown') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'geolocation.properties', '/etc/kabletown') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'ip_allow.config', '/opt/trafficserver/etc/trafficserver') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'cacheurl.config', '/opt/trafficserver/etc/trafficserver/') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'logs_xml.config', '/opt/trafficserver/etc/trafficserver') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'rascal-config.txt', '/opt/traffic_monitor/conf') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', '12M_facts', '/opt/ort') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'regex_revalidate.config', '/opt/trafficserver/etc/trafficserver') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'drop_qstring.config', '/opt/trafficserver/etc/trafficserver') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'astats.config', '/opt/trafficserver/etc/trafficserver') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'traffic_ops_ort_syncds.cron', '/etc/cron.d') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'set_dscp_0.config', '/opt/trafficserver/etc/trafficserver/dscp') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'set_dscp_10.config', '/opt/trafficserver/etc/trafficserver/dscp') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'set_dscp_12.config', '/opt/trafficserver/etc/trafficserver/dscp') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'set_dscp_14.config', '/opt/trafficserver/etc/trafficserver/dscp') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'set_dscp_18.config', '/opt/trafficserver/etc/trafficserver/dscp') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'set_dscp_20.config', '/opt/trafficserver/etc/trafficserver/dscp') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'set_dscp_22.config', '/opt/trafficserver/etc/trafficserver/dscp') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'set_dscp_26.config', '/opt/trafficserver/etc/trafficserver/dscp') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'set_dscp_28.config', '/opt/trafficserver/etc/trafficserver/dscp') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'set_dscp_30.config', '/opt/trafficserver/etc/trafficserver/dscp') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'set_dscp_34.config', '/opt/trafficserver/etc/trafficserver/dscp') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'set_dscp_36.config', '/opt/trafficserver/etc/trafficserver/dscp') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'set_dscp_38.config', '/opt/trafficserver/etc/trafficserver/dscp') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'set_dscp_8.config', '/opt/trafficserver/etc/trafficserver/dscp') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'set_dscp_16.config', '/opt/trafficserver/etc/trafficserver/dscp') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'set_dscp_24.config', '/opt/trafficserver/etc/trafficserver/dscp') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'set_dscp_32.config', '/opt/trafficserver/etc/trafficserver/dscp') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'set_dscp_40.config', '/opt/trafficserver/etc/trafficserver/dscp') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'set_dscp_48.config', '/opt/trafficserver/etc/trafficserver/dscp') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'set_dscp_56.config', '/opt/trafficserver/etc/trafficserver/dscp') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'bg_fetch.config', '/opt/trafficserver/etc/trafficserver') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('location', 'ssl_multicert.config', '/opt/trafficserver/etc/trafficserver') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.proxy_name', 'records.config', 'STRING __HOSTNAME__') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.config_dir', 'records.config', 'STRING etc/trafficserver') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.proxy_binary_opts', 'records.config', 'STRING -M') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.env_prep', 'records.config', 'STRING example_prep.sh') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.temp_dir', 'records.config', 'STRING /tmp') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.alarm_email', 'records.config', 'STRING ats') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.syslog_facility', 'records.config', 'STRING LOG_DAEMON') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.output.logfile', 'records.config', 'STRING traffic.out') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.snapshot_dir', 'records.config', 'STRING snapshots') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.system.mmap_max', 'records.config', 'INT 2097152') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.exec_thread.autoconfig.scale', 'records.config', 'FLOAT 1.5') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.accept_threads', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.admin.admin_user', 'records.config', 'STRING admin') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.admin.number_config_bak', 'records.config', 'INT 3') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.admin.user_id', 'records.config', 'STRING ats') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.admin.autoconf_port', 'records.config', 'INT 8083') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.process_manager.mgmt_port', 'records.config', 'INT 8084') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.alarm.bin', 'records.config', 'STRING example_alarm_bin.sh') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.alarm.abs_path', 'records.config', 'STRING NULL') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.connect_ports', 'records.config', 'STRING 443 563') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.insert_request_via_str', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.insert_response_via_str', 'records.config', 'INT 3') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.response_server_enabled', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.insert_age_in_response', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.enable_url_expandomatic', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.no_dns_just_forward_to_parent', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.uncacheable_requests_bypass_parent', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.keep_alive_enabled_in', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.keep_alive_enabled_out', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.chunking_enabled', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.send_http11_requests', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.share_server_sessions', 'records.config', 'INT 2') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.origin_server_pipeline', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.user_agent_pipeline', 'records.config', 'INT 8') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.referer_filter', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.referer_format_redirect', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.referer_default_redirect', 'records.config', 'STRING http://www.example.com/') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.parent_proxy_routing_enable', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.parent_proxy.retry_time', 'records.config', 'INT 300') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.parent_proxy.fail_threshold', 'records.config', 'INT 10') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.parent_proxy.total_connect_attempts', 'records.config', 'INT 4') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.parent_proxy.per_parent_connect_attempts', 'records.config', 'INT 2') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.parent_proxy.connect_attempts_timeout', 'records.config', 'INT 30') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.forward.proxy_auth_to_parent', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.keep_alive_no_activity_timeout_in', 'records.config', 'INT 115') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.keep_alive_no_activity_timeout_out', 'records.config', 'INT 120') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.transaction_no_activity_timeout_in', 'records.config', 'INT 30') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.transaction_no_activity_timeout_out', 'records.config', 'INT 30') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.transaction_active_timeout_in', 'records.config', 'INT 900') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.transaction_active_timeout_out', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.accept_no_activity_timeout', 'records.config', 'INT 120') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.connect_attempts_max_retries', 'records.config', 'INT 6') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.connect_attempts_max_retries_dead_server', 'records.config', 'INT 3') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.connect_attempts_rr_retries', 'records.config', 'INT 3') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.connect_attempts_timeout', 'records.config', 'INT 30') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.post_connect_attempts_timeout', 'records.config', 'INT 1800') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.down_server.cache_time', 'records.config', 'INT 300') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.down_server.abort_threshold', 'records.config', 'INT 10') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.congestion_control.enabled', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.anonymize_remove_from', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.anonymize_remove_referer', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.anonymize_remove_user_agent', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.anonymize_remove_cookie', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.anonymize_remove_client_ip', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.anonymize_insert_client_ip', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.anonymize_other_header_list', 'records.config', 'STRING NULL') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.insert_squid_x_forwarded_for', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.push_method_enabled', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.cache.http', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.cache.ignore_client_no_cache', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.cache.ims_on_client_no_cache', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.cache.ignore_server_no_cache', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.cache.ignore_client_cc_max_age', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.normalize_ae_gzip', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.cache.cache_responses_to_cookies', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.cache.ignore_authentication', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.cache.cache_urls_that_look_dynamic', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.cache.enable_default_vary_headers', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.cache.when_to_revalidate', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.cache.when_to_add_no_cache_to_msie_requests', 'records.config', 'INT -1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.cache.required_headers', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.cache.max_stale_age', 'records.config', 'INT 604800') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.cache.range.lookup', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.cache.heuristic_min_lifetime', 'records.config', 'INT 3600') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.cache.heuristic_max_lifetime', 'records.config', 'INT 86400') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.cache.heuristic_lm_factor', 'records.config', 'FLOAT 0.10') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.cache.fuzz.time', 'records.config', 'INT 240') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.cache.fuzz.probability', 'records.config', 'FLOAT 0.005') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.cache.vary_default_text', 'records.config', 'STRING NULL') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.cache.vary_default_images', 'records.config', 'STRING NULL') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.cache.vary_default_other', 'records.config', 'STRING NULL') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.enable_http_stats', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.body_factory.enable_logging', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.body_factory.response_suppression_mode', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.net.connections_throttle', 'records.config', 'INT 500000') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.net.defer_accept', 'records.config', 'INT 45') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.cluster.cluster_port', 'records.config', 'INT 8086') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.cluster.rsport', 'records.config', 'INT 8088') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.cluster.mcport', 'records.config', 'INT 8089') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.cluster.mc_group_addr', 'records.config', 'STRING 224.0.1.37') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.cluster.mc_ttl', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.cluster.log_bogus_mc_msgs', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.cluster.ethernet_interface', 'records.config', 'STRING lo') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.cache.permit.pinning', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.cache.ram_cache.algorithm', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.cache.ram_cache.use_seen_filter', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.cache.ram_cache.compress', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.cache.limits.http.max_alts', 'records.config', 'INT 5') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.cache.target_fragment_size', 'records.config', 'INT 1048576') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.cache.max_doc_size', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.cache.min_average_object_size', 'records.config', 'INT 131072') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.cache.threads_per_disk', 'records.config', 'INT 8') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.cache.mutex_retry_delay', 'records.config', 'INT 2') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.dns.search_default_domains', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.dns.splitDNS.enabled', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.dns.max_dns_in_flight', 'records.config', 'INT 2048') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.dns.url_expansions', 'records.config', 'STRING NULL') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.dns.round_robin_nameservers', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.dns.nameservers', 'records.config', 'STRING NULL') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.dns.resolv_conf', 'records.config', 'STRING /etc/resolv.conf') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.dns.validate_query_name', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.hostdb.size', 'records.config', 'INT 120000') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.hostdb.storage_size', 'records.config', 'INT 33554432') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.hostdb.timeout', 'records.config', 'INT 1440') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.hostdb.strict_round_robin', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.logging_enabled', 'records.config', 'INT 3') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.max_secs_per_buffer', 'records.config', 'INT 5') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.max_space_mb_for_logs', 'records.config', 'INT 25000') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.max_space_mb_for_orphan_logs', 'records.config', 'INT 25') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.max_space_mb_headroom', 'records.config', 'INT 1000') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.hostname', 'records.config', 'STRING localhost') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.logfile_dir', 'records.config', 'STRING var/log/trafficserver') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.logfile_perm', 'records.config', 'STRING rw-r--r--') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.custom_logs_enabled', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.squid_log_enabled', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.squid_log_is_ascii', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.squid_log_name', 'records.config', 'STRING squid') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.squid_log_header', 'records.config', 'STRING NULL') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.common_log_enabled', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.common_log_is_ascii', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.common_log_name', 'records.config', 'STRING common') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.common_log_header', 'records.config', 'STRING NULL') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.extended_log_enabled', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.extended_log_is_ascii', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.extended_log_name', 'records.config', 'STRING extended') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.extended_log_header', 'records.config', 'STRING NULL') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.extended2_log_is_ascii', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.extended2_log_name', 'records.config', 'STRING extended2') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.extended2_log_header', 'records.config', 'STRING NULL') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.separate_icp_logs', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.separate_host_logs', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.collation_host', 'records.config', 'STRING NULL') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.collation_port', 'records.config', 'INT 8085') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.collation_secret', 'records.config', 'STRING foobar') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.collation_host_tagged', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.collation_retry_sec', 'records.config', 'INT 5') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.rolling_enabled', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.rolling_interval_sec', 'records.config', 'INT 86400') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.rolling_offset_hr', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.rolling_size_mb', 'records.config', 'INT 10') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.auto_delete_rolled_files', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.sampling_frequency', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.reverse_proxy.enabled', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.header.parse.no_host_url_redirect', 'records.config', 'STRING NULL') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.url_remap.default_to_server_pac', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.url_remap.default_to_server_pac_port', 'records.config', 'INT -1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.url_remap.remap_required', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.url_remap.pristine_host_hdr', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.ssl.number.threads', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.ssl.SSLv2', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.ssl.TLSv1', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.ssl.compression', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.ssl.client.certification_level', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.ssl.server.cert_chain.filename', 'records.config', 'STRING NULL') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.ssl.CA.cert.filename', 'records.config', 'STRING NULL') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.ssl.client.verify.server', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.ssl.client.cert.filename', 'records.config', 'STRING NULL') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.ssl.client.private_key.filename', 'records.config', 'STRING NULL') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.ssl.client.CA.cert.filename', 'records.config', 'STRING NULL') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.icp.enabled', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.icp.icp_interface', 'records.config', 'STRING NULL') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.icp.icp_port', 'records.config', 'INT 3130') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.icp.multicast_enabled', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.icp.query_timeout', 'records.config', 'INT 2') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.update.enabled', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.update.force', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.update.retry_count', 'records.config', 'INT 10') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.update.retry_interval', 'records.config', 'INT 2') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.update.concurrent_updates', 'records.config', 'INT 100') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.net.sock_send_buffer_size_in', 'records.config', 'INT 262144') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.net.sock_recv_buffer_size_in', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.net.sock_send_buffer_size_out', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.net.sock_recv_buffer_size_out', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.core_limit', 'records.config', 'INT -1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.diags.debug.enabled', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.diags.debug.tags', 'records.config', 'STRING http.*|dns.*') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.dump_mem_info_frequency', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.slow.log.threshold', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.task_threads', 'records.config', 'INT 2') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.parent_proxy.file', 'records.config', 'STRING parent.config') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.url_remap.filename', 'records.config', 'STRING remap.config') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.cluster.cluster_configuration ', 'records.config', 'STRING cluster.config') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.cache.control.filename', 'records.config', 'STRING cache.config') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.cache.hosting_filename', 'records.config', 'STRING hosting.config') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.allocator.debug_filter', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.allocator.enable_reclaim', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.allocator.max_overage', 'records.config', 'INT 3') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.diags.show_location', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.cache.allow_empty_doc', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.stack_dump_enabled', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.hostdb.ttl_mode', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.hostdb.serve_stale_for', 'records.config', 'INT 6') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.cache.enable_read_while_writer', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.background_fill_active_timeout', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.background_fill_completed_threshold', 'records.config', 'FLOAT 0.0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.extended2_log_enabled', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.exec_thread.affinity', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.exec_thread.autoconfig', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.exec_thread.limit', 'records.config', 'INT 32') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.allocator.thread_freelist_size', 'records.config', 'INT 1024') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.mlock_enabled', 'records.config', 'INT 2') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.cache.ram_cache.size', 'records.config', 'INT 34359738368') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.log.xml_config_file', 'records.config', 'STRING logs_xml.config') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.cache.ram_cache_cutoff', 'records.config', 'INT 268435456') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.ssl.server.cipher_suite', 'records.config', 'STRING ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES128-SHA256:ECDHE-RSA-AES256-SHA384:AES128-GCM-SHA256:AES256-GCM-SHA384:ECDHE-RSA-RC4-SHA:ECDHE-RSA-AES128-SHA:ECDHE-RSA-AES256-SHA:RC4-SHA:RC4-MD5:AES128-SHA:AES256-SHA:DES-CBC3-SHA!SRP:!DSS:!PSK:!aNULL:!eNULL:!SSLv2') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.ssl.server.honor_cipher_order', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.ssl.server.cert.path', 'records.config', 'STRING etc/trafficserver/ssl') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.ssl.server.private_key.path', 'records.config', 'STRING etc/trafficserver/ssl') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.ssl.client.CA.cert.path', 'records.config', 'STRING etc/trafficserver/ssl') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.ssl.client.private_key.path', 'records.config', 'STRING etc/trafficserver/ssl') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.ssl.client.cert.path', 'records.config', 'STRING etc/trafficserver/ssl') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.ssl.CA.cert.path', 'records.config', 'STRING etc/trafficserver/ssl') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.ssl.SSLv3', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.ssl.server.multicert.filename', 'records.config', 'STRING ssl_multicert.config') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.allocator.hugepages', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.negative_caching_enabled', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.body_factory.enable_customizations', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.negative_caching_lifetime', 'records.config', 'INT 1') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.cache.http.compatibility.4-2-0-fixup', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.body_factory.template_sets_dir', 'records.config', 'STRING etc/trafficserver/body_factory') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.server_ports', 'records.config', 'STRING 80 80:ipv6 443:ssl 443:ipv6:ssl') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.http.parent_proxy_routing_enable', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.url_remap.remap_required', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.dns.lookup_timeout', 'records.config', 'INT 2') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('CONFIG proxy.config.cache.ip_allow.filename', 'records.config', 'STRING ip_allow.config') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('LOCAL proxy.config.cache.interim.storage', 'records.config', 'STRING NULL') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('LOCAL proxy.local.log.collation_mode', 'records.config', 'INT 0') ON CONFLICT DO NOTHING;
insert into parameter (name, config_file, value) values ('LOCAL proxy.local.cluster.type', 'records.config', 'INT 3') ON CONFLICT DO NOTHING;

-- profile parameters
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'GLOBAL'), (select id from parameter where name = 'tm.url' and config_file = 'global' and value = 'https://tm.kabletown.net/') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'GLOBAL'), (select id from parameter where name = 'tm.toolname' and config_file = 'global' and value = 'Traffic Ops') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'GLOBAL'), (select id from parameter where name = 'tm.infourl' and config_file = 'global' and value = 'http://docs.cdnl.kabletown.net/traffic_control/html/') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'GLOBAL'), (select id from parameter where name = 'tm.logourl' and config_file = 'global' and value = '/images/tc_logo.png') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'GLOBAL'), (select id from parameter where name = 'tm.instance_name' and config_file = 'global' and value = 'kabletown CDN') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'GLOBAL'), (select id from parameter where name = 'tm.traffic_mon_fwd_proxy' and config_file = 'global' and value = 'http://tm.kabletown.net:81') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'GLOBAL'), (select id from parameter where name = 'geolocation.polling.url' and config_file = 'CRConfig.json' and value = 'http://cdn-tools.cdnl.kabletown.net/cdn/MaxMind/GeoLiteCity.dat.gz') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'GLOBAL'), (select id from parameter where name = 'geolocation6.polling.url' and config_file = 'CRConfig.json' and value = 'http://cdn-tools.cdnl.kabletown.net/cdn/MaxMind/GeoLiteCityv6.dat.gz') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'Drive_Prefix' and config_file = 'storage.config' and value = '/dev/sd') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'Drive_Letters' and config_file = 'storage.config' and value = 'b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q,r,s,t,u,v,w,x,y') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.proxy_name' and config_file = 'records.config' and value = 'STRING __HOSTNAME__') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.config_dir' and config_file = 'records.config' and value = 'STRING etc/trafficserver') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.proxy_binary_opts' and config_file = 'records.config' and value = 'STRING -M') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.env_prep' and config_file = 'records.config' and value = 'STRING example_prep.sh') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.temp_dir' and config_file = 'records.config' and value = 'STRING /tmp') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.alarm_email' and config_file = 'records.config' and value = 'STRING ats') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.syslog_facility' and config_file = 'records.config' and value = 'STRING LOG_DAEMON') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.output.logfile' and config_file = 'records.config' and value = 'STRING traffic.out') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.snapshot_dir' and config_file = 'records.config' and value = 'STRING snapshots') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.system.mmap_max' and config_file = 'records.config' and value = 'INT 2097152') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.exec_thread.autoconfig.scale' and config_file = 'records.config' and value = 'FLOAT 1.5') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.accept_threads' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.admin.admin_user' and config_file = 'records.config' and value = 'STRING admin') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.admin.number_config_bak' and config_file = 'records.config' and value = 'INT 3') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.admin.user_id' and config_file = 'records.config' and value = 'STRING ats') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.admin.autoconf_port' and config_file = 'records.config' and value = 'INT 8083') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.process_manager.mgmt_port' and config_file = 'records.config' and value = 'INT 8084') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.alarm.bin' and config_file = 'records.config' and value = 'STRING example_alarm_bin.sh') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.alarm.abs_path' and config_file = 'records.config' and value = 'STRING NULL') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.connect_ports' and config_file = 'records.config' and value = 'STRING 443 563') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.insert_request_via_str' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.insert_response_via_str' and config_file = 'records.config' and value = 'INT 3') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.response_server_enabled' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.insert_age_in_response' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.enable_url_expandomatic' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.no_dns_just_forward_to_parent' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.uncacheable_requests_bypass_parent' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.keep_alive_enabled_in' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.keep_alive_enabled_out' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.chunking_enabled' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.send_http11_requests' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.share_server_sessions' and config_file = 'records.config' and value = 'INT 2') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.origin_server_pipeline' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.user_agent_pipeline' and config_file = 'records.config' and value = 'INT 8') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.referer_filter' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.referer_format_redirect' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.referer_default_redirect' and config_file = 'records.config' and value = 'STRING http://www.example.com/') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.parent_proxy_routing_enable' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.parent_proxy.retry_time' and config_file = 'records.config' and value = 'INT 300') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.parent_proxy.fail_threshold' and config_file = 'records.config' and value = 'INT 10') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.parent_proxy.total_connect_attempts' and config_file = 'records.config' and value = 'INT 4') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.parent_proxy.per_parent_connect_attempts' and config_file = 'records.config' and value = 'INT 2') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.parent_proxy.connect_attempts_timeout' and config_file = 'records.config' and value = 'INT 30') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.forward.proxy_auth_to_parent' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.keep_alive_no_activity_timeout_in' and config_file = 'records.config' and value = 'INT 115') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.keep_alive_no_activity_timeout_out' and config_file = 'records.config' and value = 'INT 120') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.transaction_no_activity_timeout_in' and config_file = 'records.config' and value = 'INT 30') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.transaction_no_activity_timeout_out' and config_file = 'records.config' and value = 'INT 30') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.transaction_active_timeout_in' and config_file = 'records.config' and value = 'INT 900') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.transaction_active_timeout_out' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.accept_no_activity_timeout' and config_file = 'records.config' and value = 'INT 120') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.connect_attempts_max_retries' and config_file = 'records.config' and value = 'INT 6') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.connect_attempts_max_retries_dead_server' and config_file = 'records.config' and value = 'INT 3') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.connect_attempts_rr_retries' and config_file = 'records.config' and value = 'INT 3') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.connect_attempts_timeout' and config_file = 'records.config' and value = 'INT 30') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.post_connect_attempts_timeout' and config_file = 'records.config' and value = 'INT 1800') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.down_server.cache_time' and config_file = 'records.config' and value = 'INT 300') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.down_server.abort_threshold' and config_file = 'records.config' and value = 'INT 10') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.congestion_control.enabled' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.anonymize_remove_from' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.anonymize_remove_referer' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.anonymize_remove_user_agent' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.anonymize_remove_cookie' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.anonymize_remove_client_ip' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.anonymize_insert_client_ip' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.anonymize_other_header_list' and config_file = 'records.config' and value = 'STRING NULL') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.insert_squid_x_forwarded_for' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.push_method_enabled' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.http' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.ignore_client_no_cache' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.ims_on_client_no_cache' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.ignore_server_no_cache' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.ignore_client_cc_max_age' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.normalize_ae_gzip' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.cache_responses_to_cookies' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.ignore_authentication' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.cache_urls_that_look_dynamic' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.enable_default_vary_headers' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.when_to_revalidate' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.when_to_add_no_cache_to_msie_requests' and config_file = 'records.config' and value = 'INT -1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.required_headers' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.max_stale_age' and config_file = 'records.config' and value = 'INT 604800') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.range.lookup' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.heuristic_min_lifetime' and config_file = 'records.config' and value = 'INT 3600') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.heuristic_max_lifetime' and config_file = 'records.config' and value = 'INT 86400') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.heuristic_lm_factor' and config_file = 'records.config' and value = 'FLOAT 0.10') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.fuzz.time' and config_file = 'records.config' and value = 'INT 240') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.fuzz.probability' and config_file = 'records.config' and value = 'FLOAT 0.005') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.vary_default_text' and config_file = 'records.config' and value = 'STRING NULL') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.vary_default_images' and config_file = 'records.config' and value = 'STRING NULL') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.vary_default_other' and config_file = 'records.config' and value = 'STRING NULL') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.enable_http_stats' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.body_factory.enable_logging' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.body_factory.response_suppression_mode' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.net.connections_throttle' and config_file = 'records.config' and value = 'INT 500000') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.net.defer_accept' and config_file = 'records.config' and value = 'INT 45') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'LOCAL proxy.local.cluster.type' and config_file = 'records.config' and value = 'INT 3') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.cluster.cluster_port' and config_file = 'records.config' and value = 'INT 8086') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.cluster.rsport' and config_file = 'records.config' and value = 'INT 8088') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.cluster.mcport' and config_file = 'records.config' and value = 'INT 8089') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.cluster.mc_group_addr' and config_file = 'records.config' and value = 'STRING 224.0.1.37') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.cluster.mc_ttl' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.cluster.log_bogus_mc_msgs' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.cluster.ethernet_interface' and config_file = 'records.config' and value = 'STRING lo') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.cache.permit.pinning' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.cache.ram_cache.algorithm' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.cache.ram_cache.use_seen_filter' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.cache.ram_cache.compress' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.cache.limits.http.max_alts' and config_file = 'records.config' and value = 'INT 5') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.cache.target_fragment_size' and config_file = 'records.config' and value = 'INT 1048576') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.cache.max_doc_size' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.cache.min_average_object_size' and config_file = 'records.config' and value = 'INT 131072') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.cache.threads_per_disk' and config_file = 'records.config' and value = 'INT 8') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.cache.mutex_retry_delay' and config_file = 'records.config' and value = 'INT 2') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.dns.search_default_domains' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.dns.splitDNS.enabled' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.dns.max_dns_in_flight' and config_file = 'records.config' and value = 'INT 2048') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.dns.url_expansions' and config_file = 'records.config' and value = 'STRING NULL') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.dns.round_robin_nameservers' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.dns.nameservers' and config_file = 'records.config' and value = 'STRING NULL') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.dns.resolv_conf' and config_file = 'records.config' and value = 'STRING /etc/resolv.conf') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.dns.validate_query_name' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.hostdb.size' and config_file = 'records.config' and value = 'INT 120000') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.hostdb.storage_size' and config_file = 'records.config' and value = 'INT 33554432') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.hostdb.timeout' and config_file = 'records.config' and value = 'INT 1440') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.hostdb.strict_round_robin' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.logging_enabled' and config_file = 'records.config' and value = 'INT 3') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.max_secs_per_buffer' and config_file = 'records.config' and value = 'INT 5') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.max_space_mb_for_logs' and config_file = 'records.config' and value = 'INT 25000') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.max_space_mb_for_orphan_logs' and config_file = 'records.config' and value = 'INT 25') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.max_space_mb_headroom' and config_file = 'records.config' and value = 'INT 1000') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.hostname' and config_file = 'records.config' and value = 'STRING localhost') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.logfile_dir' and config_file = 'records.config' and value = 'STRING var/log/trafficserver') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.logfile_perm' and config_file = 'records.config' and value = 'STRING rw-r--r--') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.custom_logs_enabled' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.squid_log_enabled' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.squid_log_is_ascii' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.squid_log_name' and config_file = 'records.config' and value = 'STRING squid') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.squid_log_header' and config_file = 'records.config' and value = 'STRING NULL') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.common_log_enabled' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.common_log_is_ascii' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.common_log_name' and config_file = 'records.config' and value = 'STRING common') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.common_log_header' and config_file = 'records.config' and value = 'STRING NULL') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.extended_log_enabled' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.extended_log_is_ascii' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.extended_log_name' and config_file = 'records.config' and value = 'STRING extended') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.extended_log_header' and config_file = 'records.config' and value = 'STRING NULL') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.extended2_log_is_ascii' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.extended2_log_name' and config_file = 'records.config' and value = 'STRING extended2') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.extended2_log_header' and config_file = 'records.config' and value = 'STRING NULL') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.separate_icp_logs' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.separate_host_logs' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'LOCAL proxy.local.log.collation_mode' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.collation_host' and config_file = 'records.config' and value = 'STRING NULL') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.collation_port' and config_file = 'records.config' and value = 'INT 8085') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.collation_secret' and config_file = 'records.config' and value = 'STRING foobar') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.collation_host_tagged' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.collation_retry_sec' and config_file = 'records.config' and value = 'INT 5') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.rolling_enabled' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.rolling_interval_sec' and config_file = 'records.config' and value = 'INT 86400') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.rolling_offset_hr' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.rolling_size_mb' and config_file = 'records.config' and value = 'INT 10') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.auto_delete_rolled_files' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.sampling_frequency' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.reverse_proxy.enabled' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.header.parse.no_host_url_redirect' and config_file = 'records.config' and value = 'STRING NULL') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.url_remap.default_to_server_pac' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.url_remap.default_to_server_pac_port' and config_file = 'records.config' and value = 'INT -1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.url_remap.remap_required' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.url_remap.pristine_host_hdr' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.ssl.number.threads' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.ssl.SSLv2' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.ssl.TLSv1' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.ssl.compression' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.ssl.client.certification_level' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.ssl.server.cert_chain.filename' and config_file = 'records.config' and value = 'STRING NULL') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.ssl.CA.cert.filename' and config_file = 'records.config' and value = 'STRING NULL') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.ssl.client.verify.server' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.ssl.client.cert.filename' and config_file = 'records.config' and value = 'STRING NULL') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.ssl.client.private_key.filename' and config_file = 'records.config' and value = 'STRING NULL') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.ssl.client.CA.cert.filename' and config_file = 'records.config' and value = 'STRING NULL') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.icp.enabled' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.icp.icp_interface' and config_file = 'records.config' and value = 'STRING NULL') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.icp.icp_port' and config_file = 'records.config' and value = 'INT 3130') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.icp.multicast_enabled' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.icp.query_timeout' and config_file = 'records.config' and value = 'INT 2') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.update.enabled' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.update.force' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.update.retry_count' and config_file = 'records.config' and value = 'INT 10') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.update.retry_interval' and config_file = 'records.config' and value = 'INT 2') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.update.concurrent_updates' and config_file = 'records.config' and value = 'INT 100') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.net.sock_send_buffer_size_in' and config_file = 'records.config' and value = 'INT 262144') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.net.sock_recv_buffer_size_in' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.net.sock_send_buffer_size_out' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.net.sock_recv_buffer_size_out' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.core_limit' and config_file = 'records.config' and value = 'INT -1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.diags.debug.enabled' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.diags.debug.tags' and config_file = 'records.config' and value = 'STRING http.*|dns.*') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.dump_mem_info_frequency' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.slow.log.threshold' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.task_threads' and config_file = 'records.config' and value = 'INT 2') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'cache.config' and value = '/opt/trafficserver/etc/trafficserver/') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'hosting.config' and value = '/opt/trafficserver/etc/trafficserver/') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'parent.config' and value = '/opt/trafficserver/etc/trafficserver/') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'plugin.config' and value = '/opt/trafficserver/etc/trafficserver/') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'records.config' and value = '/opt/trafficserver/etc/trafficserver/') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'remap.config' and value = '/opt/trafficserver/etc/trafficserver/') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'storage.config' and value = '/opt/trafficserver/etc/trafficserver/') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'volume.config' and value = '/opt/trafficserver/etc/trafficserver/') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = '50-ats.rules' and value = '/etc/udev/rules.d/') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.parent_proxy.file' and config_file = 'records.config' and value = 'STRING parent.config') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.url_remap.filename' and config_file = 'records.config' and value = 'STRING remap.config') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'ip_allow.config' and value = '/opt/trafficserver/etc/trafficserver') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.cluster.cluster_configuration ' and config_file = 'records.config' and value = 'STRING cluster.config') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'ramdisk_size' and config_file = 'grub.conf' and value = 'ramdisk_size=16777216') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'purge_allow_ip' and config_file = 'ip_allow.config' and value = '10.10.10.10') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'health.threshold.loadavg' and config_file = 'rascal.properties' and value = '25.0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'health.threshold.availableBandwidthInKbps' and config_file = 'rascal.properties' and value = '>1750000') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'history.count' and config_file = 'rascal.properties' and value = '30') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'cacheurl.so' and config_file = 'plugin.config' and value = '') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'cacheurl.config' and value = '/opt/trafficserver/etc/trafficserver/') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.cache.control.filename' and config_file = 'records.config' and value = 'STRING cache.config') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'LogObject.RollingEnabled' and config_file = 'logs_xml.config' and value = '3') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'LogObject.RollingIntervalSec' and config_file = 'logs_xml.config' and value = '86400') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'LogObject.RollingOffsetHr' and config_file = 'logs_xml.config' and value = '11') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'LogObject.RollingSizeMb' and config_file = 'logs_xml.config' and value = '1024') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'logs_xml.config' and value = '/opt/trafficserver/etc/trafficserver') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.xml_config_file' and config_file = 'records.config' and value = 'STRING logs_xml.config') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'health.threshold.queryTime' and config_file = 'rascal.properties' and value = '1000') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'health.polling.url' and config_file = 'rascal.properties' and value = 'http://${hostname}/_astats?application=&inf.name=${interface_name}') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'RAM_Drive_Prefix' and config_file = 'storage.config' and value = '/dev/ram') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'RAM_Drive_Letters' and config_file = 'storage.config' and value = '0,1,2,3,4,5,6,7') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'RAM_Volume' and config_file = 'storage.config' and value = '2') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'domain_name' and config_file = 'CRConfig.json' and value = 'cdn.kabletown.net') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = '12M_facts' and value = '/opt/ort') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'Disk_Volume' and config_file = 'storage.config' and value = '1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.cache.hosting_filename' and config_file = 'records.config' and value = 'STRING hosting.config') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'health.connection.timeout' and config_file = 'rascal.properties' and value = '2000') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'trafficserver' and config_file = 'chkconfig' and value = '0:off\t1:off\t2:on\t3:on\t4:on\t5:on\t6:off') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.allocator.debug_filter' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.allocator.enable_reclaim' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.allocator.max_overage' and config_file = 'records.config' and value = 'INT 3') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.diags.show_location' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.allow_empty_doc' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'LOCAL proxy.config.cache.interim.storage' and config_file = 'records.config' and value = 'STRING NULL') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'regex_revalidate.so' and config_file = 'plugin.config' and value = '--config regex_revalidate.config') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'regex_revalidate.config' and value = '/opt/trafficserver/etc/trafficserver') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'remap_stats.so' and config_file = 'plugin.config' and value = '') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.stack_dump_enabled' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'drop_qstring.config' and value = '/opt/trafficserver/etc/trafficserver') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.hostdb.ttl_mode' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.hostdb.serve_stale_for' and config_file = 'records.config' and value = 'INT 6') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.cache.enable_read_while_writer' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.background_fill_active_timeout' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.background_fill_completed_threshold' and config_file = 'records.config' and value = 'FLOAT 0.0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.log.extended2_log_enabled' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.exec_thread.affinity' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.exec_thread.autoconfig' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.exec_thread.limit' and config_file = 'records.config' and value = 'INT 32') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.allocator.thread_freelist_size' and config_file = 'records.config' and value = 'INT 1024') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.mlock_enabled' and config_file = 'records.config' and value = 'INT 2') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'LogFormat.Format' and config_file = 'logs_xml.config' and value = '%<cqtq> chi=%<chi> phn=%<phn> shn=%<shn> url=%<cquuc> cqhm=%<cqhm> cqhv=%<cqhv> pssc=%<pssc> ttms=%<ttms> b=%<pscl> sssc=%<sssc> sscl=%<sscl> cfsc=%<cfsc> pfsc=%<pfsc> crc=%<crc> phr=%<phr> uas="%<{User-Agent}cqh>"') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'LogFormat.Name' and config_file = 'logs_xml.config' and value = 'custom_ats_2') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'LogObject.Format' and config_file = 'logs_xml.config' and value = 'custom_ats_2') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'LogObject.Filename' and config_file = 'logs_xml.config' and value = 'custom_ats_2') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.cache.ram_cache.size' and config_file = 'records.config' and value = 'INT 34359738368') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'astats_over_http.so' and config_file = 'plugin.config' and value = '') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'allow_ip' and config_file = 'astats.config' and value = '127.0.0.1,10.10.10.10/16') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'allow_ip6' and config_file = 'astats.config' and value = '::1,d009:5dd:f0d8:18d::2/64') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'record_types' and config_file = 'astats.config' and value = '144') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'astats.config' and value = '/opt/trafficserver/etc/trafficserver') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'path' and config_file = 'astats.config' and value = '_astats') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.cache.http.compatibility.4-2-0-fixup' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'algorithm' and config_file = 'parent.config' and value = 'consistent_hash') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'maxRevalDurationDays' and config_file = 'regex_revalidate.config' and value = '90') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.cache.ram_cache_cutoff' and config_file = 'records.config' and value = 'INT 268435456') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.ssl.server.cipher_suite' and config_file = 'records.config' and value = 'STRING ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES128-SHA256:ECDHE-RSA-AES256-SHA384:AES128-GCM-SHA256:AES256-GCM-SHA384:ECDHE-RSA-RC4-SHA:ECDHE-RSA-AES128-SHA:ECDHE-RSA-AES256-SHA:RC4-SHA:RC4-MD5:AES128-SHA:AES256-SHA:DES-CBC3-SHA!SRP:!DSS:!PSK:!aNULL:!eNULL:!SSLv2') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.ssl.server.honor_cipher_order' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.ssl.server.cert.path' and config_file = 'records.config' and value = 'STRING etc/trafficserver/ssl') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.ssl.server.private_key.path' and config_file = 'records.config' and value = 'STRING etc/trafficserver/ssl') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.ssl.client.CA.cert.path' and config_file = 'records.config' and value = 'STRING etc/trafficserver/ssl') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.ssl.client.private_key.path' and config_file = 'records.config' and value = 'STRING etc/trafficserver/ssl') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.ssl.client.cert.path' and config_file = 'records.config' and value = 'STRING etc/trafficserver/ssl') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.ssl.CA.cert.path' and config_file = 'records.config' and value = 'STRING etc/trafficserver/ssl') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.ssl.SSLv3' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.ssl.server.multicert.filename' and config_file = 'records.config' and value = 'STRING ssl_multicert.config') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.allocator.hugepages' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.negative_caching_enabled' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'record_types' and config_file = 'astats.config' and value = '122') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.body_factory.enable_customizations' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.negative_caching_lifetime' and config_file = 'records.config' and value = 'INT 1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'cron_ort_syncds' and config_file = 'traffic_ops_ort_syncds.cron' and value = '*/15 * * * * root /opt/ort/traffic_ops_ort.pl syncds warn https://tm.kabletown.net admin:password > /tmp/ort/syncds.log 2>&1') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'traffic_ops_ort_syncds.cron' and value = '/etc/cron.d') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'set_dscp_0.config' and value = '/opt/trafficserver/etc/trafficserver/dscp') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'set_dscp_10.config' and value = '/opt/trafficserver/etc/trafficserver/dscp') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'set_dscp_12.config' and value = '/opt/trafficserver/etc/trafficserver/dscp') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'set_dscp_14.config' and value = '/opt/trafficserver/etc/trafficserver/dscp') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'set_dscp_18.config' and value = '/opt/trafficserver/etc/trafficserver/dscp') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'set_dscp_20.config' and value = '/opt/trafficserver/etc/trafficserver/dscp') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'set_dscp_22.config' and value = '/opt/trafficserver/etc/trafficserver/dscp') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'set_dscp_26.config' and value = '/opt/trafficserver/etc/trafficserver/dscp') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'set_dscp_28.config' and value = '/opt/trafficserver/etc/trafficserver/dscp') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'set_dscp_30.config' and value = '/opt/trafficserver/etc/trafficserver/dscp') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'set_dscp_34.config' and value = '/opt/trafficserver/etc/trafficserver/dscp') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'set_dscp_36.config' and value = '/opt/trafficserver/etc/trafficserver/dscp') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'set_dscp_38.config' and value = '/opt/trafficserver/etc/trafficserver/dscp') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'set_dscp_8.config' and value = '/opt/trafficserver/etc/trafficserver/dscp') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'set_dscp_16.config' and value = '/opt/trafficserver/etc/trafficserver/dscp') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'set_dscp_24.config' and value = '/opt/trafficserver/etc/trafficserver/dscp') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'set_dscp_32.config' and value = '/opt/trafficserver/etc/trafficserver/dscp') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'set_dscp_40.config' and value = '/opt/trafficserver/etc/trafficserver/dscp') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'set_dscp_48.config' and value = '/opt/trafficserver/etc/trafficserver/dscp') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'set_dscp_56.config' and value = '/opt/trafficserver/etc/trafficserver/dscp') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.http.server_ports' and config_file = 'records.config' and value = 'STRING 80 80:ipv6 443:ssl 443:ipv6:ssl') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'ssl_multicert.config' and value = '/opt/trafficserver/etc/trafficserver') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'bg_fetch.config' and value = '/opt/trafficserver/etc/trafficserver') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'qstring' and config_file = 'parent.config' and value = 'ignore') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'CONFIG proxy.config.body_factory.template_sets_dir' and config_file = 'records.config' and value = 'STRING etc/trafficserver/body_factory') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'ssl_multicert.config' and value = '/opt/trafficserver/etc/trafficserver') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'bg_fetch.config' and value = '/opt/trafficserver/etc/trafficserver') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'qstring' and config_file = 'parent.config' and value = 'ignore') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'ssl_multicert.config' and value = '/opt/trafficserver/etc/trafficserver') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'location' and config_file = 'bg_fetch.config' and value = '/opt/trafficserver/etc/trafficserver') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'EDGE1_CDN_520'), (select id from parameter where name = 'qstring' and config_file = 'parent.config' and value = 'ignore') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'MID1_CDN_421'), (select id from parameter where name = 'CONFIG proxy.config.http.parent_proxy_routing_enable' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'MID1_CDN_421'), (select id from parameter where name = 'CONFIG proxy.config.url_remap.remap_required' and config_file = 'records.config' and value = 'INT 0') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'MID1_CDN_421'), (select id from parameter where name = 'CONFIG proxy.config.dns.lookup_timeout' and config_file = 'records.config' and value = 'INT 2') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'MID1_CDN_421'), (select id from parameter where name = 'CONFIG proxy.config.cache.ip_allow.filename' and config_file = 'records.config' and value = 'STRING ip_allow.config') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'HOSTMON'), (select id from parameter where name = 'location' and config_file = 'dns.zone' and value = '/etc/kabletown/zones/<zonename>.info') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'HOSTMON'), (select id from parameter where name = 'location' and config_file = 'http-log4j.properties' and value = '/etc/kabletown') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'HOSTMON'), (select id from parameter where name = 'location' and config_file = 'dns-log4j.properties' and value = '/etc/kabletown') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'HOSTMON'), (select id from parameter where name = 'geolocation.polling.interval' and config_file = 'CRConfig.json' and value = '86400000') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'HOSTMON'), (select id from parameter where name = 'location' and config_file = 'geolocation.properties' and value = '/etc/kabletown') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'HOSTMON'), (select id from parameter where name = 'coveragezone.polling.url' and config_file = 'CRConfig.json' and value = 'http://cdn-tools.cdnl.kabletown.net/cdn/CZF/current/kabletown_cdn_czf-current.json') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'HOSTMON'), (select id from parameter where name = 'geolocation.polling.url' and config_file = 'CRConfig.json' and value = 'https://tm.kabletown.net/MaxMind/auto/GeoIP2-City.mmdb.gz') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'HOSTMON'), (select id from parameter where name = 'monitor:///opt/tomcat/logs/access.log' and config_file = 'inputs.conf' and value = 'index=index_test;sourcetype=access_ccr') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'HOSTMON'), (select id from parameter where name = 'coveragezone.polling.interval' and config_file = 'CRConfig.json' and value = '86400000') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'HOSTMON'), (select id from parameter where name = 'tld.soa.admin' and config_file = 'CRConfig.json' and value = 'twelve_monkeys') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'HOSTMON'), (select id from parameter where name = 'tld.soa.minimum' and config_file = 'CRConfig.json' and value = '86400') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'HOSTMON'), (select id from parameter where name = 'tld.ttls.AAAA' and config_file = 'CRConfig.json' and value = '3600') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'HOSTMON'), (select id from parameter where name = 'tld.soa.retry' and config_file = 'CRConfig.json' and value = '7200') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'HOSTMON'), (select id from parameter where name = 'tld.soa.expire' and config_file = 'CRConfig.json' and value = '604800') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'HOSTMON'), (select id from parameter where name = 'tld.ttls.A' and config_file = 'CRConfig.json' and value = '3600') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'HOSTMON'), (select id from parameter where name = 'tld.soa.refresh' and config_file = 'CRConfig.json' and value = '28800') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'HOSTMON'), (select id from parameter where name = 'tld.ttls.NS' and config_file = 'CRConfig.json' and value = '3600') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'HOSTMON'), (select id from parameter where name = 'tld.ttls.SOA' and config_file = 'CRConfig.json' and value = '86400') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'HOSTMON'), (select id from parameter where name = 'api.port' and config_file = 'server.xml' and value = '3333') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'CCR_CDN'), (select id from parameter where name = 'location' and config_file = 'dns.zone' and value = '/etc/kabletown/zones/<zonename>.info') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'CCR_CDN'), (select id from parameter where name = 'location' and config_file = 'http-log4j.properties' and value = '/etc/kabletown') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'CCR_CDN'), (select id from parameter where name = 'geolocation.polling.url' and config_file = 'CRConfig.json' and value = 'https://tm.kabletown.net/MaxMind/auto/GeoIP2-City.mmdb.gz') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'CCR_CDN'), (select id from parameter where name = 'location' and config_file = 'dns-log4j.properties' and value = '/etc/kabletown') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'CCR_CDN'), (select id from parameter where name = 'location' and config_file = 'geolocation.properties' and value = '/etc/kabletown') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'CCR_CDN'), (select id from parameter where name = 'monitor:///opt/tomcat/logs/access.log' and config_file = 'inputs.conf' and value = 'index=index_test;sourcetype=access_ccr') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'CCR_CDN'), (select id from parameter where name = 'geolocation.polling.interval' and config_file = 'CRConfig.json' and value = '86400000') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'CCR_CDN'), (select id from parameter where name = 'coveragezone.polling.interval' and config_file = 'CRConfig.json' and value = '86400000') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'CCR_CDN'), (select id from parameter where name = 'coveragezone.polling.url' and config_file = 'CRConfig.json' and value = 'http://cdn-tools.cdnl.kabletown.net/cdn/CZF/current/kabletown_cdn_czf-current.json') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'CCR_CDN'), (select id from parameter where name = 'domain_name' and config_file = 'CRConfig.json' and value = 'cdn.kabletown.net') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'CCR_CDN'), (select id from parameter where name = 'tld.ttls.AAAA' and config_file = 'CRConfig.json' and value = '3600') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'CCR_CDN'), (select id from parameter where name = 'tld.ttls.SOA' and config_file = 'CRConfig.json' and value = '86400') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'CCR_CDN'), (select id from parameter where name = 'tld.ttls.A' and config_file = 'CRConfig.json' and value = '3600') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'CCR_CDN'), (select id from parameter where name = 'tld.ttls.NS' and config_file = 'CRConfig.json' and value = '3600') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'CCR_CDN'), (select id from parameter where name = 'tld.soa.expire' and config_file = 'CRConfig.json' and value = '604800') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'CCR_CDN'), (select id from parameter where name = 'tld.soa.minimum' and config_file = 'CRConfig.json' and value = '86400') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'CCR_CDN'), (select id from parameter where name = 'tld.soa.admin' and config_file = 'CRConfig.json' and value = 'twelve_monkeys') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'CCR_CDN'), (select id from parameter where name = 'tld.soa.retry' and config_file = 'CRConfig.json' and value = '7200') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'CCR_CDN'), (select id from parameter where name = 'tld.soa.refresh' and config_file = 'CRConfig.json' and value = '28800') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'CCR_CDN'), (select id from parameter where name = 'api.port' and config_file = 'server.xml' and value = '3333') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'CCR_CDN'), (select id from parameter where name = 'api.cache-control.max-age' and config_file = 'CRConfig.json' and value = '10') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'RASCAL_CDN'), (select id from parameter where name = 'tm.crConfig.polling.url' and config_file = 'rascal-config.txt' and value = 'https://${tmHostname}/CRConfig-Snapshots/${cdnName}/CRConfig.xml') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'RASCAL_CDN'), (select id from parameter where name = 'health.polling.interval' and config_file = 'rascal-config.txt' and value = '8000') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'RASCAL_CDN'), (select id from parameter where name = 'tm.dataServer.polling.url' and config_file = 'rascal-config.txt' and value = 'https://${tmHostname}/dataserver/orderby/id') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'RASCAL_CDN'), (select id from parameter where name = 'tm.healthParams.polling.url' and config_file = 'rascal-config.txt' and value = 'https://${tmHostname}/health/${cdnName}') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'RASCAL_CDN'), (select id from parameter where name = 'tm.polling.interval' and config_file = 'rascal-config.txt' and value = '60000') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'RASCAL_CDN'), (select id from parameter where name = 'location' and config_file = 'rascal-config.txt' and value = '/opt/traffic_monitor/conf') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'RASCAL_CDN'), (select id from parameter where name = 'health.threadPool' and config_file = 'rascal-config.txt' and value = '4') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'RASCAL_CDN'), (select id from parameter where name = 'health.event-count' and config_file = 'rascal-config.txt' and value = '200') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'RASCAL_CDN'), (select id from parameter where name = 'hack.ttl' and config_file = 'rascal-config.txt' and value = '30') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'RASCAL_CDN'), (select id from parameter where name = 'health.timepad' and config_file = 'rascal-config.txt' and value = '30') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'TRAFFIC_STATS'), (select id from parameter where name = 'CacheStats' and config_file = 'traffic_stats.config' and value = 'bandwidth') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'TRAFFIC_STATS'), (select id from parameter where name = 'CacheStats' and config_file = 'traffic_stats.config' and value = 'maxKbps') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'TRAFFIC_STATS'), (select id from parameter where name = 'CacheStats' and config_file = 'traffic_stats.config' and value = 'ats.proxy.process.http.current_client_connections') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'TRAFFIC_STATS'), (select id from parameter where name = 'DsStats' and config_file = 'traffic_stats.config' and value = 'kbps') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'TRAFFIC_STATS'), (select id from parameter where name = 'DsStats' and config_file = 'traffic_stats.config' and value = 'tps_2xx') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'TRAFFIC_STATS'), (select id from parameter where name = 'DsStats' and config_file = 'traffic_stats.config' and value = 'status_4xx') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'TRAFFIC_STATS'), (select id from parameter where name = 'DsStats' and config_file = 'traffic_stats.config' and value = 'status_5xx') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'TRAFFIC_STATS'), (select id from parameter where name = 'DsStats' and config_file = 'traffic_stats.config' and value = 'tps_3xx') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'TRAFFIC_STATS'), (select id from parameter where name = 'DsStats' and config_file = 'traffic_stats.config' and value = 'tps_4xx') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'TRAFFIC_STATS'), (select id from parameter where name = 'DsStats' and config_file = 'traffic_stats.config' and value = 'tps_5xx') ) ON CONFLICT DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'TRAFFIC_STATS'), (select id from parameter where name = 'DsStats' and config_file = 'traffic_stats.config' and value = 'tps_total') ) ON CONFLICT DO NOTHING;

-- servers
update server set https_port = 443 where https_port is null;

-- users
insert into tm_user (username, role,full_name) values ('portal',(select id from role where name='portal'), 'Portal User') ON CONFLICT DO NOTHING;
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
