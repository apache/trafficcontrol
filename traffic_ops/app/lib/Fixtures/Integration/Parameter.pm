package Fixtures::Integration::Parameter;
#
# Copyright 2015 Comcast Cable Communications Management, LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
use Moose;
extends 'DBIx::Class::EasyFixture';
use namespace::autoclean;
use Digest::SHA1 qw(sha1_hex);

my %definition_for = (
	'0-domain_name' => {
		new   => 'Parameter',
		using => {
			id          => 1,
			name        => 'domain_name',
			config_file => 'CRConfig.xml',
			value       => 'cdn1.kabletown.net',
		},
	},
	'1-GeolocationURL' => {
		new   => 'Parameter',
		using => {
			id          => 2,
			name        => 'GeolocationURL',
			config_file => 'CRConfig.xml',
			value       => 'http://aux.cdnlab.kabletown.net:8080/GeoLiteCity.dat.gz',
		},
	},
	'2-CacheHealthTimeout' => {
		new   => 'Parameter',
		using => {
			id          => 3,
			name        => 'CacheHealthTimeout',
			config_file => 'CRConfig.xml',
			value       => '70',
		},
	},
	'3-CoverageZoneMapURL' => {
		new   => 'Parameter',
		using => {
			id          => 4,
			name        => 'CoverageZoneMapURL',
			config_file => 'CRConfig.xml',
			value       => 'http://aux.cdnlab.kabletown.net/logs/production/reports/czf/current/kabletown_cdn_czf.xml',
		},
	},
	'4-CoverageZoneMapRefreshPeriodHours' => {
		new   => 'Parameter',
		using => {
			id          => 5,
			name        => 'CoverageZoneMapRefreshPeriodHours',
			config_file => 'CRConfig.xml',
			value       => '24',
		},
	},
	'5-Drive_Prefix' => {
		new   => 'Parameter',
		using => {
			id          => 11,
			name        => 'Drive_Prefix',
			config_file => 'storage.config',
			value       => '/dev/sd',
		},
	},
	'6-Drive_Letters' => {
		new   => 'Parameter',
		using => {
			id          => 12,
			name        => 'Drive_Letters',
			config_file => 'storage.config',
			value       => '0,1,2,3,4,5,6',
		},
	},
	'7-Drive_Letters' => {
		new   => 'Parameter',
		using => {
			id          => 13,
			name        => 'Drive_Letters',
			config_file => 'storage.config',
			value       => 'b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q,r,s,t,u,v,w,x,y',
		},
	},
	'8-CONFIG-proxy.config.proxy_name' => {
		new   => 'Parameter',
		using => {
			id          => 14,
			name        => 'CONFIG proxy.config.proxy_name',
			config_file => 'records.config',
			value       => 'STRING __HOSTNAME__',
		},
	},
	'9-CONFIG-proxy.config.config_dir' => {
		new   => 'Parameter',
		using => {
			id          => 15,
			name        => 'CONFIG proxy.config.config_dir',
			config_file => 'records.config',
			value       => 'STRING etc/trafficserver',
		},
	},
	'10-CONFIG-proxy.config.proxy_binary_opts' => {
		new   => 'Parameter',
		using => {
			id          => 16,
			name        => 'CONFIG proxy.config.proxy_binary_opts',
			config_file => 'records.config',
			value       => 'STRING -M',
		},
	},
	'11-CONFIG-proxy.config.env_prep' => {
		new   => 'Parameter',
		using => {
			id          => 17,
			name        => 'CONFIG proxy.config.env_prep',
			config_file => 'records.config',
			value       => 'STRING example_prep.sh',
		},
	},
	'12-CONFIG-proxy.config.temp_dir' => {
		new   => 'Parameter',
		using => {
			id          => 18,
			name        => 'CONFIG proxy.config.temp_dir',
			config_file => 'records.config',
			value       => 'STRING /tmp',
		},
	},
	'13-CONFIG-proxy.config.alarm_email' => {
		new   => 'Parameter',
		using => {
			id          => 19,
			name        => 'CONFIG proxy.config.alarm_email',
			config_file => 'records.config',
			value       => 'STRING ats',
		},
	},
	'14-CONFIG-proxy.config.syslog_facility' => {
		new   => 'Parameter',
		using => {
			id          => 20,
			name        => 'CONFIG proxy.config.syslog_facility',
			config_file => 'records.config',
			value       => 'STRING LOG_DAEMON',
		},
	},
	'15-CONFIG-proxy.config.output.logfile' => {
		new   => 'Parameter',
		using => {
			id          => 21,
			name        => 'CONFIG proxy.config.output.logfile',
			config_file => 'records.config',
			value       => 'STRING traffic.out',
		},
	},
	'16-CONFIG-proxy.config.snapshot_dir' => {
		new   => 'Parameter',
		using => {
			id          => 22,
			name        => 'CONFIG proxy.config.snapshot_dir',
			config_file => 'records.config',
			value       => 'STRING snapshots',
		},
	},
	'17-CONFIG-proxy.config.system.mmap_max' => {
		new   => 'Parameter',
		using => {
			id          => 23,
			name        => 'CONFIG proxy.config.system.mmap_max',
			config_file => 'records.config',
			value       => 'INT 2097152',
		},
	},
	'18-CONFIG-proxy.config.exec_thread.autoconfig' => {
		new   => 'Parameter',
		using => {
			id          => 24,
			name        => 'CONFIG proxy.config.exec_thread.autoconfig',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'19-CONFIG-proxy.config.exec_thread.autoconfig.scale' => {
		new   => 'Parameter',
		using => {
			id          => 25,
			name        => 'CONFIG proxy.config.exec_thread.autoconfig.scale',
			config_file => 'records.config',
			value       => 'FLOAT 1.5',
		},
	},
	'20-CONFIG-proxy.config.exec_thread.limit' => {
		new   => 'Parameter',
		using => {
			id          => 26,
			name        => 'CONFIG proxy.config.exec_thread.limit',
			config_file => 'records.config',
			value       => 'INT 2',
		},
	},
	'21-CONFIG-proxy.config.accept_threads' => {
		new   => 'Parameter',
		using => {
			id          => 27,
			name        => 'CONFIG proxy.config.accept_threads',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'22-CONFIG-proxy.config.admin.admin_user' => {
		new   => 'Parameter',
		using => {
			id          => 28,
			name        => 'CONFIG proxy.config.admin.admin_user',
			config_file => 'records.config',
			value       => 'STRING admin',
		},
	},
	'23-CONFIG-proxy.config.admin.number_config_bak' => {
		new   => 'Parameter',
		using => {
			id          => 29,
			name        => 'CONFIG proxy.config.admin.number_config_bak',
			config_file => 'records.config',
			value       => 'INT 3',
		},
	},
	'24-CONFIG-proxy.config.admin.user_id' => {
		new   => 'Parameter',
		using => {
			id          => 30,
			name        => 'CONFIG proxy.config.admin.user_id',
			config_file => 'records.config',
			value       => 'STRING ats',
		},
	},
	'25-CONFIG-proxy.config.admin.autoconf_port' => {
		new   => 'Parameter',
		using => {
			id          => 31,
			name        => 'CONFIG proxy.config.admin.autoconf_port',
			config_file => 'records.config',
			value       => 'INT 8083',
		},
	},
	'26-CONFIG-proxy.config.process_manager.mgmt_port' => {
		new   => 'Parameter',
		using => {
			id          => 32,
			name        => 'CONFIG proxy.config.process_manager.mgmt_port',
			config_file => 'records.config',
			value       => 'INT 8084',
		},
	},
	'27-CONFIG-proxy.config.alarm.bin' => {
		new   => 'Parameter',
		using => {
			id          => 33,
			name        => 'CONFIG proxy.config.alarm.bin',
			config_file => 'records.config',
			value       => 'STRING example_alarm_bin.sh',
		},
	},
	'28-CONFIG-proxy.config.alarm.abs_path' => {
		new   => 'Parameter',
		using => {
			id          => 34,
			name        => 'CONFIG proxy.config.alarm.abs_path',
			config_file => 'records.config',
			value       => 'STRING NULL',
		},
	},
	'29-CONFIG-proxy.config.http.server_ports' => {
		new   => 'Parameter',
		using => {
			id          => 35,
			name        => 'CONFIG proxy.config.http.server_ports',
			config_file => 'records.config',
			value       => 'STRING 80 80:ipv6',
		},
	},
	'30-CONFIG-proxy.config.http.connect_ports' => {
		new   => 'Parameter',
		using => {
			id          => 36,
			name        => 'CONFIG proxy.config.http.connect_ports',
			config_file => 'records.config',
			value       => 'STRING 443 563',
		},
	},
	'31-CONFIG-proxy.config.http.insert_request_via_str' => {
		new   => 'Parameter',
		using => {
			id          => 37,
			name        => 'CONFIG proxy.config.http.insert_request_via_str',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'32-CONFIG-proxy.config.http.insert_response_via_str' => {
		new   => 'Parameter',
		using => {
			id          => 38,
			name        => 'CONFIG proxy.config.http.insert_response_via_str',
			config_file => 'records.config',
			value       => 'INT 3',
		},
	},
	'33-CONFIG-proxy.config.http.response_server_enabled' => {
		new   => 'Parameter',
		using => {
			id          => 39,
			name        => 'CONFIG proxy.config.http.response_server_enabled',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'34-CONFIG-proxy.config.http.insert_age_in_response' => {
		new   => 'Parameter',
		using => {
			id          => 40,
			name        => 'CONFIG proxy.config.http.insert_age_in_response',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'35-CONFIG-proxy.config.http.enable_url_expandomatic' => {
		new   => 'Parameter',
		using => {
			id          => 41,
			name        => 'CONFIG proxy.config.http.enable_url_expandomatic',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'36-CONFIG-proxy.config.http.no_dns_just_forward_to_parent' => {
		new   => 'Parameter',
		using => {
			id          => 42,
			name        => 'CONFIG proxy.config.http.no_dns_just_forward_to_parent',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'37-CONFIG-proxy.config.http.uncacheable_requests_bypass_parent' => {
		new   => 'Parameter',
		using => {
			id          => 43,
			name        => 'CONFIG proxy.config.http.uncacheable_requests_bypass_parent',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'38-CONFIG-proxy.config.http.keep_alive_enabled_in' => {
		new   => 'Parameter',
		using => {
			id          => 44,
			name        => 'CONFIG proxy.config.http.keep_alive_enabled_in',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'39-CONFIG-proxy.config.http.keep_alive_enabled_out' => {
		new   => 'Parameter',
		using => {
			id          => 45,
			name        => 'CONFIG proxy.config.http.keep_alive_enabled_out',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'40-CONFIG-proxy.config.http.chunking_enabled' => {
		new   => 'Parameter',
		using => {
			id          => 46,
			name        => 'CONFIG proxy.config.http.chunking_enabled',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'41-CONFIG-proxy.config.http.send_http11_requests' => {
		new   => 'Parameter',
		using => {
			id          => 47,
			name        => 'CONFIG proxy.config.http.send_http11_requests',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'42-CONFIG-proxy.config.http.share_server_sessions' => {
		new   => 'Parameter',
		using => {
			id          => 48,
			name        => 'CONFIG proxy.config.http.share_server_sessions',
			config_file => 'records.config',
			value       => 'INT 2',
		},
	},
	'43-CONFIG-proxy.config.http.origin_server_pipeline' => {
		new   => 'Parameter',
		using => {
			id          => 49,
			name        => 'CONFIG proxy.config.http.origin_server_pipeline',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'44-CONFIG-proxy.config.http.user_agent_pipeline' => {
		new   => 'Parameter',
		using => {
			id          => 50,
			name        => 'CONFIG proxy.config.http.user_agent_pipeline',
			config_file => 'records.config',
			value       => 'INT 8',
		},
	},
	'45-CONFIG-proxy.config.http.referer_filter' => {
		new   => 'Parameter',
		using => {
			id          => 51,
			name        => 'CONFIG proxy.config.http.referer_filter',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'46-CONFIG-proxy.config.http.referer_format_redirect' => {
		new   => 'Parameter',
		using => {
			id          => 52,
			name        => 'CONFIG proxy.config.http.referer_format_redirect',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'47-CONFIG-proxy.config.http.referer_default_redirect' => {
		new   => 'Parameter',
		using => {
			id          => 53,
			name        => 'CONFIG proxy.config.http.referer_default_redirect',
			config_file => 'records.config',
			value       => 'STRING http://www.example.com/',
		},
	},
	'48-CONFIG-proxy.config.http.parent_proxy_routing_enable' => {
		new   => 'Parameter',
		using => {
			id          => 54,
			name        => 'CONFIG proxy.config.http.parent_proxy_routing_enable',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'49-CONFIG-proxy.config.http.parent_proxy.retry_time' => {
		new   => 'Parameter',
		using => {
			id          => 55,
			name        => 'CONFIG proxy.config.http.parent_proxy.retry_time',
			config_file => 'records.config',
			value       => 'INT 300',
		},
	},
	'50-CONFIG-proxy.config.http.parent_proxy.fail_threshold' => {
		new   => 'Parameter',
		using => {
			id          => 56,
			name        => 'CONFIG proxy.config.http.parent_proxy.fail_threshold',
			config_file => 'records.config',
			value       => 'INT 10',
		},
	},
	'51-CONFIG-proxy.config.http.parent_proxy.total_connect_attempts' => {
		new   => 'Parameter',
		using => {
			id          => 57,
			name        => 'CONFIG proxy.config.http.parent_proxy.total_connect_attempts',
			config_file => 'records.config',
			value       => 'INT 4',
		},
	},
	'52-CONFIG-proxy.config.http.parent_proxy.per_parent_connect_attempts' => {
		new   => 'Parameter',
		using => {
			id          => 58,
			name        => 'CONFIG proxy.config.http.parent_proxy.per_parent_connect_attempts',
			config_file => 'records.config',
			value       => 'INT 2',
		},
	},
	'53-CONFIG-proxy.config.http.parent_proxy.connect_attempts_timeout' => {
		new   => 'Parameter',
		using => {
			id          => 59,
			name        => 'CONFIG proxy.config.http.parent_proxy.connect_attempts_timeout',
			config_file => 'records.config',
			value       => 'INT 30',
		},
	},
	'54-CONFIG-proxy.config.http.forward.proxy_auth_to_parent' => {
		new   => 'Parameter',
		using => {
			id          => 60,
			name        => 'CONFIG proxy.config.http.forward.proxy_auth_to_parent',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'55-CONFIG-proxy.config.http.keep_alive_no_activity_timeout_in' => {
		new   => 'Parameter',
		using => {
			id          => 61,
			name        => 'CONFIG proxy.config.http.keep_alive_no_activity_timeout_in',
			config_file => 'records.config',
			value       => 'INT 115',
		},
	},
	'56-CONFIG-proxy.config.http.keep_alive_no_activity_timeout_out' => {
		new   => 'Parameter',
		using => {
			id          => 62,
			name        => 'CONFIG proxy.config.http.keep_alive_no_activity_timeout_out',
			config_file => 'records.config',
			value       => 'INT 120',
		},
	},
	'57-CONFIG-proxy.config.http.transaction_no_activity_timeout_in' => {
		new   => 'Parameter',
		using => {
			id          => 63,
			name        => 'CONFIG proxy.config.http.transaction_no_activity_timeout_in',
			config_file => 'records.config',
			value       => 'INT 30',
		},
	},
	'58-CONFIG-proxy.config.http.transaction_no_activity_timeout_out' => {
		new   => 'Parameter',
		using => {
			id          => 64,
			name        => 'CONFIG proxy.config.http.transaction_no_activity_timeout_out',
			config_file => 'records.config',
			value       => 'INT 30',
		},
	},
	'59-CONFIG-proxy.config.http.transaction_active_timeout_in' => {
		new   => 'Parameter',
		using => {
			id          => 65,
			name        => 'CONFIG proxy.config.http.transaction_active_timeout_in',
			config_file => 'records.config',
			value       => 'INT 900',
		},
	},
	'60-CONFIG-proxy.config.http.transaction_active_timeout_out' => {
		new   => 'Parameter',
		using => {
			id          => 66,
			name        => 'CONFIG proxy.config.http.transaction_active_timeout_out',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'61-CONFIG-proxy.config.http.accept_no_activity_timeout' => {
		new   => 'Parameter',
		using => {
			id          => 67,
			name        => 'CONFIG proxy.config.http.accept_no_activity_timeout',
			config_file => 'records.config',
			value       => 'INT 120',
		},
	},
	'62-CONFIG-proxy.config.http.background_fill_active_timeout' => {
		new   => 'Parameter',
		using => {
			id          => 68,
			name        => 'CONFIG proxy.config.http.background_fill_active_timeout',
			config_file => 'records.config',
			value       => 'INT 60',
		},
	},
	'63-CONFIG-proxy.config.http.background_fill_completed_threshold' => {
		new   => 'Parameter',
		using => {
			id          => 69,
			name        => 'CONFIG proxy.config.http.background_fill_completed_threshold',
			config_file => 'records.config',
			value       => 'FLOAT 0.5',
		},
	},
	'64-CONFIG-proxy.config.http.connect_attempts_max_retries' => {
		new   => 'Parameter',
		using => {
			id          => 70,
			name        => 'CONFIG proxy.config.http.connect_attempts_max_retries',
			config_file => 'records.config',
			value       => 'INT 6',
		},
	},
	'65-CONFIG-proxy.config.http.connect_attempts_max_retries_dead_server' => {
		new   => 'Parameter',
		using => {
			id          => 71,
			name        => 'CONFIG proxy.config.http.connect_attempts_max_retries_dead_server',
			config_file => 'records.config',
			value       => 'INT 3',
		},
	},
	'66-CONFIG-proxy.config.http.connect_attempts_rr_retries' => {
		new   => 'Parameter',
		using => {
			id          => 72,
			name        => 'CONFIG proxy.config.http.connect_attempts_rr_retries',
			config_file => 'records.config',
			value       => 'INT 3',
		},
	},
	'67-CONFIG-proxy.config.http.connect_attempts_timeout' => {
		new   => 'Parameter',
		using => {
			id          => 73,
			name        => 'CONFIG proxy.config.http.connect_attempts_timeout',
			config_file => 'records.config',
			value       => 'INT 30',
		},
	},
	'68-CONFIG-proxy.config.http.post_connect_attempts_timeout' => {
		new   => 'Parameter',
		using => {
			id          => 74,
			name        => 'CONFIG proxy.config.http.post_connect_attempts_timeout',
			config_file => 'records.config',
			value       => 'INT 1800',
		},
	},
	'69-CONFIG-proxy.config.http.down_server.cache_time' => {
		new   => 'Parameter',
		using => {
			id          => 75,
			name        => 'CONFIG proxy.config.http.down_server.cache_time',
			config_file => 'records.config',
			value       => 'INT 300',
		},
	},
	'70-CONFIG-proxy.config.http.down_server.abort_threshold' => {
		new   => 'Parameter',
		using => {
			id          => 76,
			name        => 'CONFIG proxy.config.http.down_server.abort_threshold',
			config_file => 'records.config',
			value       => 'INT 10',
		},
	},
	'71-CONFIG-proxy.config.http.congestion_control.enabled' => {
		new   => 'Parameter',
		using => {
			id          => 77,
			name        => 'CONFIG proxy.config.http.congestion_control.enabled',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'72-CONFIG-proxy.config.http.negative_caching_enabled' => {
		new   => 'Parameter',
		using => {
			id          => 78,
			name        => 'CONFIG proxy.config.http.negative_caching_enabled',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'73-CONFIG-proxy.config.http.negative_caching_lifetime' => {
		new   => 'Parameter',
		using => {
			id          => 79,
			name        => 'CONFIG proxy.config.http.negative_caching_lifetime',
			config_file => 'records.config',
			value       => 'INT 1800',
		},
	},
	'74-CONFIG-proxy.config.http.anonymize_remove_from' => {
		new   => 'Parameter',
		using => {
			id          => 80,
			name        => 'CONFIG proxy.config.http.anonymize_remove_from',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'75-CONFIG-proxy.config.http.anonymize_remove_referer' => {
		new   => 'Parameter',
		using => {
			id          => 81,
			name        => 'CONFIG proxy.config.http.anonymize_remove_referer',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'76-CONFIG-proxy.config.http.anonymize_remove_user_agent' => {
		new   => 'Parameter',
		using => {
			id          => 82,
			name        => 'CONFIG proxy.config.http.anonymize_remove_user_agent',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'77-CONFIG-proxy.config.http.anonymize_remove_cookie' => {
		new   => 'Parameter',
		using => {
			id          => 83,
			name        => 'CONFIG proxy.config.http.anonymize_remove_cookie',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'78-CONFIG-proxy.config.http.anonymize_remove_client_ip' => {
		new   => 'Parameter',
		using => {
			id          => 84,
			name        => 'CONFIG proxy.config.http.anonymize_remove_client_ip',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'79-CONFIG-proxy.config.http.anonymize_insert_client_ip' => {
		new   => 'Parameter',
		using => {
			id          => 85,
			name        => 'CONFIG proxy.config.http.anonymize_insert_client_ip',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'80-CONFIG-proxy.config.http.anonymize_other_header_list' => {
		new   => 'Parameter',
		using => {
			id          => 86,
			name        => 'CONFIG proxy.config.http.anonymize_other_header_list',
			config_file => 'records.config',
			value       => 'STRING NULL',
		},
	},
	'81-CONFIG-proxy.config.http.insert_squid_x_forwarded_for' => {
		new   => 'Parameter',
		using => {
			id          => 87,
			name        => 'CONFIG proxy.config.http.insert_squid_x_forwarded_for',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'82-CONFIG-proxy.config.http.push_method_enabled' => {
		new   => 'Parameter',
		using => {
			id          => 88,
			name        => 'CONFIG proxy.config.http.push_method_enabled',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'83-CONFIG-proxy.config.http.cache.http' => {
		new   => 'Parameter',
		using => {
			id          => 89,
			name        => 'CONFIG proxy.config.http.cache.http',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'84-CONFIG-proxy.config.http.cache.ignore_client_no_cache' => {
		new   => 'Parameter',
		using => {
			id          => 90,
			name        => 'CONFIG proxy.config.http.cache.ignore_client_no_cache',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'85-CONFIG-proxy.config.http.cache.ims_on_client_no_cache' => {
		new   => 'Parameter',
		using => {
			id          => 91,
			name        => 'CONFIG proxy.config.http.cache.ims_on_client_no_cache',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'86-CONFIG-proxy.config.http.cache.ignore_server_no_cache' => {
		new   => 'Parameter',
		using => {
			id          => 92,
			name        => 'CONFIG proxy.config.http.cache.ignore_server_no_cache',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'87-CONFIG-proxy.config.http.cache.ignore_client_cc_max_age' => {
		new   => 'Parameter',
		using => {
			id          => 93,
			name        => 'CONFIG proxy.config.http.cache.ignore_client_cc_max_age',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'88-CONFIG-proxy.config.http.normalize_ae_gzip' => {
		new   => 'Parameter',
		using => {
			id          => 94,
			name        => 'CONFIG proxy.config.http.normalize_ae_gzip',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'89-CONFIG-proxy.config.http.cache.cache_responses_to_cookies' => {
		new   => 'Parameter',
		using => {
			id          => 95,
			name        => 'CONFIG proxy.config.http.cache.cache_responses_to_cookies',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'90-CONFIG-proxy.config.http.cache.ignore_authentication' => {
		new   => 'Parameter',
		using => {
			id          => 96,
			name        => 'CONFIG proxy.config.http.cache.ignore_authentication',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'91-CONFIG-proxy.config.http.cache.cache_urls_that_look_dynamic' => {
		new   => 'Parameter',
		using => {
			id          => 97,
			name        => 'CONFIG proxy.config.http.cache.cache_urls_that_look_dynamic',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'92-CONFIG-proxy.config.http.cache.enable_default_vary_headers' => {
		new   => 'Parameter',
		using => {
			id          => 98,
			name        => 'CONFIG proxy.config.http.cache.enable_default_vary_headers',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'93-CONFIG-proxy.config.http.cache.when_to_revalidate' => {
		new   => 'Parameter',
		using => {
			id          => 99,
			name        => 'CONFIG proxy.config.http.cache.when_to_revalidate',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'94-CONFIG-proxy.config.http.cache.when_to_add_no_cache_to_msie_requests' => {
		new   => 'Parameter',
		using => {
			id          => 100,
			name        => 'CONFIG proxy.config.http.cache.when_to_add_no_cache_to_msie_requests',
			config_file => 'records.config',
			value       => 'INT -1',
		},
	},
	'95-CONFIG-proxy.config.http.cache.required_headers' => {
		new   => 'Parameter',
		using => {
			id          => 101,
			name        => 'CONFIG proxy.config.http.cache.required_headers',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'96-CONFIG-proxy.config.http.cache.max_stale_age' => {
		new   => 'Parameter',
		using => {
			id          => 102,
			name        => 'CONFIG proxy.config.http.cache.max_stale_age',
			config_file => 'records.config',
			value       => 'INT 604800',
		},
	},
	'97-CONFIG-proxy.config.http.cache.range.lookup' => {
		new   => 'Parameter',
		using => {
			id          => 103,
			name        => 'CONFIG proxy.config.http.cache.range.lookup',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'98-CONFIG-proxy.config.http.cache.heuristic_min_lifetime' => {
		new   => 'Parameter',
		using => {
			id          => 104,
			name        => 'CONFIG proxy.config.http.cache.heuristic_min_lifetime',
			config_file => 'records.config',
			value       => 'INT 3600',
		},
	},
	'99-CONFIG-proxy.config.http.cache.heuristic_max_lifetime' => {
		new   => 'Parameter',
		using => {
			id          => 105,
			name        => 'CONFIG proxy.config.http.cache.heuristic_max_lifetime',
			config_file => 'records.config',
			value       => 'INT 86400',
		},
	},
	'100-CONFIG-proxy.config.http.cache.heuristic_lm_factor' => {
		new   => 'Parameter',
		using => {
			id          => 106,
			name        => 'CONFIG proxy.config.http.cache.heuristic_lm_factor',
			config_file => 'records.config',
			value       => 'FLOAT 0.10',
		},
	},
	'101-CONFIG-proxy.config.http.cache.fuzz.time' => {
		new   => 'Parameter',
		using => {
			id          => 107,
			name        => 'CONFIG proxy.config.http.cache.fuzz.time',
			config_file => 'records.config',
			value       => 'INT 240',
		},
	},
	'102-CONFIG-proxy.config.http.cache.fuzz.probability' => {
		new   => 'Parameter',
		using => {
			id          => 108,
			name        => 'CONFIG proxy.config.http.cache.fuzz.probability',
			config_file => 'records.config',
			value       => 'FLOAT 0.005',
		},
	},
	'103-CONFIG-proxy.config.http.cache.vary_default_text' => {
		new   => 'Parameter',
		using => {
			id          => 109,
			name        => 'CONFIG proxy.config.http.cache.vary_default_text',
			config_file => 'records.config',
			value       => 'STRING NULL',
		},
	},
	'104-CONFIG-proxy.config.http.cache.vary_default_images' => {
		new   => 'Parameter',
		using => {
			id          => 110,
			name        => 'CONFIG proxy.config.http.cache.vary_default_images',
			config_file => 'records.config',
			value       => 'STRING NULL',
		},
	},
	'105-CONFIG-proxy.config.http.cache.vary_default_other' => {
		new   => 'Parameter',
		using => {
			id          => 111,
			name        => 'CONFIG proxy.config.http.cache.vary_default_other',
			config_file => 'records.config',
			value       => 'STRING NULL',
		},
	},
	'106-CONFIG-proxy.config.http.enable_http_stats' => {
		new   => 'Parameter',
		using => {
			id          => 112,
			name        => 'CONFIG proxy.config.http.enable_http_stats',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'107-CONFIG-proxy.config.body_factory.enable_customizations' => {
		new   => 'Parameter',
		using => {
			id          => 113,
			name        => 'CONFIG proxy.config.body_factory.enable_customizations',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'108-CONFIG-proxy.config.body_factory.enable_logging' => {
		new   => 'Parameter',
		using => {
			id          => 114,
			name        => 'CONFIG proxy.config.body_factory.enable_logging',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'109-CONFIG-proxy.config.body_factory.response_suppression_mode' => {
		new   => 'Parameter',
		using => {
			id          => 115,
			name        => 'CONFIG proxy.config.body_factory.response_suppression_mode',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'110-CONFIG-proxy.config.net.connections_throttle' => {
		new   => 'Parameter',
		using => {
			id          => 116,
			name        => 'CONFIG proxy.config.net.connections_throttle',
			config_file => 'records.config',
			value       => 'INT 500000',
		},
	},
	'111-CONFIG-proxy.config.net.defer_accept' => {
		new   => 'Parameter',
		using => {
			id          => 117,
			name        => 'CONFIG proxy.config.net.defer_accept',
			config_file => 'records.config',
			value       => 'INT 45',
		},
	},
	'112-LOCAL-proxy.local.cluster.type' => {
		new   => 'Parameter',
		using => {
			id          => 118,
			name        => 'LOCAL proxy.local.cluster.type',
			config_file => 'records.config',
			value       => 'INT 3',
		},
	},
	'113-CONFIG-proxy.config.cluster.cluster_port' => {
		new   => 'Parameter',
		using => {
			id          => 119,
			name        => 'CONFIG proxy.config.cluster.cluster_port',
			config_file => 'records.config',
			value       => 'INT 8086',
		},
	},
	'114-CONFIG-proxy.config.cluster.rsport' => {
		new   => 'Parameter',
		using => {
			id          => 120,
			name        => 'CONFIG proxy.config.cluster.rsport',
			config_file => 'records.config',
			value       => 'INT 8088',
		},
	},
	'115-CONFIG-proxy.config.cluster.mcport' => {
		new   => 'Parameter',
		using => {
			id          => 121,
			name        => 'CONFIG proxy.config.cluster.mcport',
			config_file => 'records.config',
			value       => 'INT 8089',
		},
	},
	'116-CONFIG-proxy.config.cluster.mc_group_addr' => {
		new   => 'Parameter',
		using => {
			id          => 122,
			name        => 'CONFIG proxy.config.cluster.mc_group_addr',
			config_file => 'records.config',
			value       => 'STRING 224.0.1.37',
		},
	},
	'117-CONFIG-proxy.config.cluster.mc_ttl' => {
		new   => 'Parameter',
		using => {
			id          => 123,
			name        => 'CONFIG proxy.config.cluster.mc_ttl',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'118-CONFIG-proxy.config.cluster.log_bogus_mc_msgs' => {
		new   => 'Parameter',
		using => {
			id          => 124,
			name        => 'CONFIG proxy.config.cluster.log_bogus_mc_msgs',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'119-CONFIG-proxy.config.cluster.ethernet_interface' => {
		new   => 'Parameter',
		using => {
			id          => 125,
			name        => 'CONFIG proxy.config.cluster.ethernet_interface',
			config_file => 'records.config',
			value       => 'STRING lo',
		},
	},
	'120-CONFIG-proxy.config.cache.permit.pinning' => {
		new   => 'Parameter',
		using => {
			id          => 126,
			name        => 'CONFIG proxy.config.cache.permit.pinning',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'121-CONFIG-proxy.config.cache.ram_cache.size' => {
		new   => 'Parameter',
		using => {
			id          => 127,
			name        => 'CONFIG proxy.config.cache.ram_cache.size',
			config_file => 'records.config',
			value       => 'INT 21474836480',
		},
	},
	'122-CONFIG-proxy.config.cache.ram_cache_cutoff' => {
		new   => 'Parameter',
		using => {
			id          => 128,
			name        => 'CONFIG proxy.config.cache.ram_cache_cutoff',
			config_file => 'records.config',
			value       => 'INT 4194304',
		},
	},
	'123-CONFIG-proxy.config.cache.ram_cache.algorithm' => {
		new   => 'Parameter',
		using => {
			id          => 129,
			name        => 'CONFIG proxy.config.cache.ram_cache.algorithm',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'124-CONFIG-proxy.config.cache.ram_cache.use_seen_filter' => {
		new   => 'Parameter',
		using => {
			id          => 130,
			name        => 'CONFIG proxy.config.cache.ram_cache.use_seen_filter',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'125-CONFIG-proxy.config.cache.ram_cache.compress' => {
		new   => 'Parameter',
		using => {
			id          => 131,
			name        => 'CONFIG proxy.config.cache.ram_cache.compress',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'126-CONFIG-proxy.config.cache.limits.http.max_alts' => {
		new   => 'Parameter',
		using => {
			id          => 132,
			name        => 'CONFIG proxy.config.cache.limits.http.max_alts',
			config_file => 'records.config',
			value       => 'INT 5',
		},
	},
	'127-CONFIG-proxy.config.cache.target_fragment_size' => {
		new   => 'Parameter',
		using => {
			id          => 133,
			name        => 'CONFIG proxy.config.cache.target_fragment_size',
			config_file => 'records.config',
			value       => 'INT 1048576',
		},
	},
	'128-CONFIG-proxy.config.cache.max_doc_size' => {
		new   => 'Parameter',
		using => {
			id          => 134,
			name        => 'CONFIG proxy.config.cache.max_doc_size',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'129-CONFIG-proxy.config.cache.enable_read_while_writer' => {
		new   => 'Parameter',
		using => {
			id          => 135,
			name        => 'CONFIG proxy.config.cache.enable_read_while_writer',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'130-CONFIG-proxy.config.cache.min_average_object_size' => {
		new   => 'Parameter',
		using => {
			id          => 136,
			name        => 'CONFIG proxy.config.cache.min_average_object_size',
			config_file => 'records.config',
			value       => 'INT 131072',
		},
	},
	'131-CONFIG-proxy.config.cache.threads_per_disk' => {
		new   => 'Parameter',
		using => {
			id          => 137,
			name        => 'CONFIG proxy.config.cache.threads_per_disk',
			config_file => 'records.config',
			value       => 'INT 8',
		},
	},
	'132-CONFIG-proxy.config.cache.mutex_retry_delay' => {
		new   => 'Parameter',
		using => {
			id          => 138,
			name        => 'CONFIG proxy.config.cache.mutex_retry_delay',
			config_file => 'records.config',
			value       => 'INT 2',
		},
	},
	'133-CONFIG-proxy.config.dns.search_default_domains' => {
		new   => 'Parameter',
		using => {
			id          => 139,
			name        => 'CONFIG proxy.config.dns.search_default_domains',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'134-CONFIG-proxy.config.dns.splitDNS.enabled' => {
		new   => 'Parameter',
		using => {
			id          => 140,
			name        => 'CONFIG proxy.config.dns.splitDNS.enabled',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'135-CONFIG-proxy.config.dns.max_dns_in_flight' => {
		new   => 'Parameter',
		using => {
			id          => 141,
			name        => 'CONFIG proxy.config.dns.max_dns_in_flight',
			config_file => 'records.config',
			value       => 'INT 2048',
		},
	},
	'136-CONFIG-proxy.config.dns.url_expansions' => {
		new   => 'Parameter',
		using => {
			id          => 142,
			name        => 'CONFIG proxy.config.dns.url_expansions',
			config_file => 'records.config',
			value       => 'STRING NULL',
		},
	},
	'137-CONFIG-proxy.config.dns.round_robin_nameservers' => {
		new   => 'Parameter',
		using => {
			id          => 143,
			name        => 'CONFIG proxy.config.dns.round_robin_nameservers',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'138-CONFIG-proxy.config.dns.nameservers' => {
		new   => 'Parameter',
		using => {
			id          => 144,
			name        => 'CONFIG proxy.config.dns.nameservers',
			config_file => 'records.config',
			value       => 'STRING NULL',
		},
	},
	'139-CONFIG-proxy.config.dns.resolv_conf' => {
		new   => 'Parameter',
		using => {
			id          => 145,
			name        => 'CONFIG proxy.config.dns.resolv_conf',
			config_file => 'records.config',
			value       => 'STRING /etc/resolv.conf',
		},
	},
	'140-CONFIG-proxy.config.dns.validate_query_name' => {
		new   => 'Parameter',
		using => {
			id          => 146,
			name        => 'CONFIG proxy.config.dns.validate_query_name',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'141-CONFIG-proxy.config.hostdb.size' => {
		new   => 'Parameter',
		using => {
			id          => 147,
			name        => 'CONFIG proxy.config.hostdb.size',
			config_file => 'records.config',
			value       => 'INT 120000',
		},
	},
	'142-CONFIG-proxy.config.hostdb.storage_size' => {
		new   => 'Parameter',
		using => {
			id          => 148,
			name        => 'CONFIG proxy.config.hostdb.storage_size',
			config_file => 'records.config',
			value       => 'INT 33554432',
		},
	},
	'143-CONFIG-proxy.config.hostdb.ttl_mode' => {
		new   => 'Parameter',
		using => {
			id          => 149,
			name        => 'CONFIG proxy.config.hostdb.ttl_mode',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'144-CONFIG-proxy.config.hostdb.timeout' => {
		new   => 'Parameter',
		using => {
			id          => 150,
			name        => 'CONFIG proxy.config.hostdb.timeout',
			config_file => 'records.config',
			value       => 'INT 1440',
		},
	},
	'145-CONFIG-proxy.config.hostdb.strict_round_robin' => {
		new   => 'Parameter',
		using => {
			id          => 151,
			name        => 'CONFIG proxy.config.hostdb.strict_round_robin',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'146-CONFIG-proxy.config.log.logging_enabled' => {
		new   => 'Parameter',
		using => {
			id          => 152,
			name        => 'CONFIG proxy.config.log.logging_enabled',
			config_file => 'records.config',
			value       => 'INT 3',
		},
	},
	'147-CONFIG-proxy.config.log.max_secs_per_buffer' => {
		new   => 'Parameter',
		using => {
			id          => 153,
			name        => 'CONFIG proxy.config.log.max_secs_per_buffer',
			config_file => 'records.config',
			value       => 'INT 5',
		},
	},
	'148-CONFIG-proxy.config.log.max_space_mb_for_logs' => {
		new   => 'Parameter',
		using => {
			id          => 154,
			name        => 'CONFIG proxy.config.log.max_space_mb_for_logs',
			config_file => 'records.config',
			value       => 'INT 25000',
		},
	},
	'149-CONFIG-proxy.config.log.max_space_mb_for_orphan_logs' => {
		new   => 'Parameter',
		using => {
			id          => 155,
			name        => 'CONFIG proxy.config.log.max_space_mb_for_orphan_logs',
			config_file => 'records.config',
			value       => 'INT 25',
		},
	},
	'150-CONFIG-proxy.config.log.max_space_mb_headroom' => {
		new   => 'Parameter',
		using => {
			id          => 156,
			name        => 'CONFIG proxy.config.log.max_space_mb_headroom',
			config_file => 'records.config',
			value       => 'INT 1000',
		},
	},
	'151-CONFIG-proxy.config.log.hostname' => {
		new   => 'Parameter',
		using => {
			id          => 157,
			name        => 'CONFIG proxy.config.log.hostname',
			config_file => 'records.config',
			value       => 'STRING localhost',
		},
	},
	'152-CONFIG-proxy.config.log.logfile_dir' => {
		new   => 'Parameter',
		using => {
			id          => 158,
			name        => 'CONFIG proxy.config.log.logfile_dir',
			config_file => 'records.config',
			value       => 'STRING var/log/trafficserver',
		},
	},
	'153-CONFIG-proxy.config.log.logfile_perm' => {
		new   => 'Parameter',
		using => {
			id          => 159,
			name        => 'CONFIG proxy.config.log.logfile_perm',
			config_file => 'records.config',
			value       => 'STRING rw-r--r--',
		},
	},
	'154-CONFIG-proxy.config.log.custom_logs_enabled' => {
		new   => 'Parameter',
		using => {
			id          => 160,
			name        => 'CONFIG proxy.config.log.custom_logs_enabled',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'155-CONFIG-proxy.config.log.squid_log_enabled' => {
		new   => 'Parameter',
		using => {
			id          => 161,
			name        => 'CONFIG proxy.config.log.squid_log_enabled',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'156-CONFIG-proxy.config.log.squid_log_is_ascii' => {
		new   => 'Parameter',
		using => {
			id          => 162,
			name        => 'CONFIG proxy.config.log.squid_log_is_ascii',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'157-CONFIG-proxy.config.log.squid_log_name' => {
		new   => 'Parameter',
		using => {
			id          => 163,
			name        => 'CONFIG proxy.config.log.squid_log_name',
			config_file => 'records.config',
			value       => 'STRING squid',
		},
	},
	'158-CONFIG-proxy.config.log.squid_log_header' => {
		new   => 'Parameter',
		using => {
			id          => 164,
			name        => 'CONFIG proxy.config.log.squid_log_header',
			config_file => 'records.config',
			value       => 'STRING NULL',
		},
	},
	'159-CONFIG-proxy.config.log.common_log_enabled' => {
		new   => 'Parameter',
		using => {
			id          => 165,
			name        => 'CONFIG proxy.config.log.common_log_enabled',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'160-CONFIG-proxy.config.log.common_log_is_ascii' => {
		new   => 'Parameter',
		using => {
			id          => 166,
			name        => 'CONFIG proxy.config.log.common_log_is_ascii',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'161-CONFIG-proxy.config.log.common_log_name' => {
		new   => 'Parameter',
		using => {
			id          => 167,
			name        => 'CONFIG proxy.config.log.common_log_name',
			config_file => 'records.config',
			value       => 'STRING common',
		},
	},
	'162-CONFIG-proxy.config.log.common_log_header' => {
		new   => 'Parameter',
		using => {
			id          => 168,
			name        => 'CONFIG proxy.config.log.common_log_header',
			config_file => 'records.config',
			value       => 'STRING NULL',
		},
	},
	'163-CONFIG-proxy.config.log.extended_log_enabled' => {
		new   => 'Parameter',
		using => {
			id          => 169,
			name        => 'CONFIG proxy.config.log.extended_log_enabled',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'164-CONFIG-proxy.config.log.extended_log_is_ascii' => {
		new   => 'Parameter',
		using => {
			id          => 170,
			name        => 'CONFIG proxy.config.log.extended_log_is_ascii',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'165-CONFIG-proxy.config.log.extended_log_name' => {
		new   => 'Parameter',
		using => {
			id          => 171,
			name        => 'CONFIG proxy.config.log.extended_log_name',
			config_file => 'records.config',
			value       => 'STRING extended',
		},
	},
	'166-CONFIG-proxy.config.log.extended_log_header' => {
		new   => 'Parameter',
		using => {
			id          => 172,
			name        => 'CONFIG proxy.config.log.extended_log_header',
			config_file => 'records.config',
			value       => 'STRING NULL',
		},
	},
	'167-CONFIG-proxy.config.log.extended2_log_enabled' => {
		new   => 'Parameter',
		using => {
			id          => 173,
			name        => 'CONFIG proxy.config.log.extended2_log_enabled',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'168-CONFIG-proxy.config.log.extended2_log_is_ascii' => {
		new   => 'Parameter',
		using => {
			id          => 174,
			name        => 'CONFIG proxy.config.log.extended2_log_is_ascii',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'169-CONFIG-proxy.config.log.extended2_log_name' => {
		new   => 'Parameter',
		using => {
			id          => 175,
			name        => 'CONFIG proxy.config.log.extended2_log_name',
			config_file => 'records.config',
			value       => 'STRING extended2',
		},
	},
	'170-CONFIG-proxy.config.log.extended2_log_header' => {
		new   => 'Parameter',
		using => {
			id          => 176,
			name        => 'CONFIG proxy.config.log.extended2_log_header',
			config_file => 'records.config',
			value       => 'STRING NULL',
		},
	},
	'171-CONFIG-proxy.config.log.separate_icp_logs' => {
		new   => 'Parameter',
		using => {
			id          => 177,
			name        => 'CONFIG proxy.config.log.separate_icp_logs',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'172-CONFIG-proxy.config.log.separate_host_logs' => {
		new   => 'Parameter',
		using => {
			id          => 178,
			name        => 'CONFIG proxy.config.log.separate_host_logs',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'173-LOCAL-proxy.local.log.collation_mode' => {
		new   => 'Parameter',
		using => {
			id          => 179,
			name        => 'LOCAL proxy.local.log.collation_mode',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'174-CONFIG-proxy.config.log.collation_host' => {
		new   => 'Parameter',
		using => {
			id          => 180,
			name        => 'CONFIG proxy.config.log.collation_host',
			config_file => 'records.config',
			value       => 'STRING NULL',
		},
	},
	'175-CONFIG-proxy.config.log.collation_port' => {
		new   => 'Parameter',
		using => {
			id          => 181,
			name        => 'CONFIG proxy.config.log.collation_port',
			config_file => 'records.config',
			value       => 'INT 8085',
		},
	},
	'176-CONFIG-proxy.config.log.collation_secret' => {
		new   => 'Parameter',
		using => {
			id          => 182,
			name        => 'CONFIG proxy.config.log.collation_secret',
			config_file => 'records.config',
			value       => 'STRING foobar',
		},
	},
	'177-CONFIG-proxy.config.log.collation_host_tagged' => {
		new   => 'Parameter',
		using => {
			id          => 183,
			name        => 'CONFIG proxy.config.log.collation_host_tagged',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'178-CONFIG-proxy.config.log.collation_retry_sec' => {
		new   => 'Parameter',
		using => {
			id          => 184,
			name        => 'CONFIG proxy.config.log.collation_retry_sec',
			config_file => 'records.config',
			value       => 'INT 5',
		},
	},
	'179-CONFIG-proxy.config.log.rolling_enabled' => {
		new   => 'Parameter',
		using => {
			id          => 185,
			name        => 'CONFIG proxy.config.log.rolling_enabled',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'180-CONFIG-proxy.config.log.rolling_interval_sec' => {
		new   => 'Parameter',
		using => {
			id          => 186,
			name        => 'CONFIG proxy.config.log.rolling_interval_sec',
			config_file => 'records.config',
			value       => 'INT 86400',
		},
	},
	'181-CONFIG-proxy.config.log.rolling_offset_hr' => {
		new   => 'Parameter',
		using => {
			id          => 187,
			name        => 'CONFIG proxy.config.log.rolling_offset_hr',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'182-CONFIG-proxy.config.log.rolling_size_mb' => {
		new   => 'Parameter',
		using => {
			id          => 188,
			name        => 'CONFIG proxy.config.log.rolling_size_mb',
			config_file => 'records.config',
			value       => 'INT 10',
		},
	},
	'183-CONFIG-proxy.config.log.auto_delete_rolled_files' => {
		new   => 'Parameter',
		using => {
			id          => 189,
			name        => 'CONFIG proxy.config.log.auto_delete_rolled_files',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'184-CONFIG-proxy.config.log.sampling_frequency' => {
		new   => 'Parameter',
		using => {
			id          => 190,
			name        => 'CONFIG proxy.config.log.sampling_frequency',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'185-CONFIG-proxy.config.reverse_proxy.enabled' => {
		new   => 'Parameter',
		using => {
			id          => 191,
			name        => 'CONFIG proxy.config.reverse_proxy.enabled',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'186-CONFIG-proxy.config.header.parse.no_host_url_redirect' => {
		new   => 'Parameter',
		using => {
			id          => 192,
			name        => 'CONFIG proxy.config.header.parse.no_host_url_redirect',
			config_file => 'records.config',
			value       => 'STRING NULL',
		},
	},
	'187-CONFIG-proxy.config.url_remap.default_to_server_pac' => {
		new   => 'Parameter',
		using => {
			id          => 193,
			name        => 'CONFIG proxy.config.url_remap.default_to_server_pac',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'188-CONFIG-proxy.config.url_remap.default_to_server_pac_port' => {
		new   => 'Parameter',
		using => {
			id          => 194,
			name        => 'CONFIG proxy.config.url_remap.default_to_server_pac_port',
			config_file => 'records.config',
			value       => 'INT -1',
		},
	},
	'189-CONFIG-proxy.config.url_remap.remap_required' => {
		new   => 'Parameter',
		using => {
			id          => 195,
			name        => 'CONFIG proxy.config.url_remap.remap_required',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'190-CONFIG-proxy.config.url_remap.pristine_host_hdr' => {
		new   => 'Parameter',
		using => {
			id          => 196,
			name        => 'CONFIG proxy.config.url_remap.pristine_host_hdr',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'191-CONFIG-proxy.config.ssl.number.threads' => {
		new   => 'Parameter',
		using => {
			id          => 197,
			name        => 'CONFIG proxy.config.ssl.number.threads',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'192-CONFIG-proxy.config.ssl.SSLv2' => {
		new   => 'Parameter',
		using => {
			id          => 198,
			name        => 'CONFIG proxy.config.ssl.SSLv2',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'193-CONFIG-proxy.config.ssl.SSLv3' => {
		new   => 'Parameter',
		using => {
			id          => 199,
			name        => 'CONFIG proxy.config.ssl.SSLv3',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'194-CONFIG-proxy.config.ssl.TLSv1' => {
		new   => 'Parameter',
		using => {
			id          => 200,
			name        => 'CONFIG proxy.config.ssl.TLSv1',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'195-CONFIG-proxy.config.ssl.server.cipher_suite' => {
		new   => 'Parameter',
		using => {
			id          => 201,
			name        => 'CONFIG proxy.config.ssl.server.cipher_suite',
			config_file => 'records.config',
			value       => 'STRING RC4-SHA:AES128-SHA:DES-CBC3-SHA:AES256-SHA:ALL:!aNULL:!EXP:!LOW:!MD5:!SSLV2:!NULL',
		},
	},
	'196-CONFIG-proxy.config.ssl.server.honor_cipher_order' => {
		new   => 'Parameter',
		using => {
			id          => 202,
			name        => 'CONFIG proxy.config.ssl.server.honor_cipher_order',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'197-CONFIG-proxy.config.ssl.compression' => {
		new   => 'Parameter',
		using => {
			id          => 203,
			name        => 'CONFIG proxy.config.ssl.compression',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'198-CONFIG-proxy.config.ssl.client.certification_level' => {
		new   => 'Parameter',
		using => {
			id          => 204,
			name        => 'CONFIG proxy.config.ssl.client.certification_level',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'199-CONFIG-proxy.config.ssl.server.cert_chain.filename' => {
		new   => 'Parameter',
		using => {
			id          => 205,
			name        => 'CONFIG proxy.config.ssl.server.cert_chain.filename',
			config_file => 'records.config',
			value       => 'STRING NULL',
		},
	},
	'200-CONFIG-proxy.config.ssl.server.cert.path' => {
		new   => 'Parameter',
		using => {
			id          => 206,
			name        => 'CONFIG proxy.config.ssl.server.cert.path',
			config_file => 'records.config',
			value       => 'STRING etc/trafficserver',
		},
	},
	'201-CONFIG-proxy.config.ssl.server.private_key.path' => {
		new   => 'Parameter',
		using => {
			id          => 207,
			name        => 'CONFIG proxy.config.ssl.server.private_key.path',
			config_file => 'records.config',
			value       => 'STRING etc/trafficserver',
		},
	},
	'202-CONFIG-proxy.config.ssl.CA.cert.filename' => {
		new   => 'Parameter',
		using => {
			id          => 208,
			name        => 'CONFIG proxy.config.ssl.CA.cert.filename',
			config_file => 'records.config',
			value       => 'STRING NULL',
		},
	},
	'203-CONFIG-proxy.config.ssl.CA.cert.path' => {
		new   => 'Parameter',
		using => {
			id          => 209,
			name        => 'CONFIG proxy.config.ssl.CA.cert.path',
			config_file => 'records.config',
			value       => 'STRING etc/trafficserver',
		},
	},
	'204-CONFIG-proxy.config.ssl.client.verify.server' => {
		new   => 'Parameter',
		using => {
			id          => 210,
			name        => 'CONFIG proxy.config.ssl.client.verify.server',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'205-CONFIG-proxy.config.ssl.client.cert.filename' => {
		new   => 'Parameter',
		using => {
			id          => 211,
			name        => 'CONFIG proxy.config.ssl.client.cert.filename',
			config_file => 'records.config',
			value       => 'STRING NULL',
		},
	},
	'206-CONFIG-proxy.config.ssl.client.cert.path' => {
		new   => 'Parameter',
		using => {
			id          => 212,
			name        => 'CONFIG proxy.config.ssl.client.cert.path',
			config_file => 'records.config',
			value       => 'STRING etc/trafficserver',
		},
	},
	'207-CONFIG-proxy.config.ssl.client.private_key.filename' => {
		new   => 'Parameter',
		using => {
			id          => 213,
			name        => 'CONFIG proxy.config.ssl.client.private_key.filename',
			config_file => 'records.config',
			value       => 'STRING NULL',
		},
	},
	'208-CONFIG-proxy.config.ssl.client.private_key.path' => {
		new   => 'Parameter',
		using => {
			id          => 214,
			name        => 'CONFIG proxy.config.ssl.client.private_key.path',
			config_file => 'records.config',
			value       => 'STRING etc/trafficserver',
		},
	},
	'209-CONFIG-proxy.config.ssl.client.CA.cert.filename' => {
		new   => 'Parameter',
		using => {
			id          => 215,
			name        => 'CONFIG proxy.config.ssl.client.CA.cert.filename',
			config_file => 'records.config',
			value       => 'STRING NULL',
		},
	},
	'210-CONFIG-proxy.config.ssl.client.CA.cert.path' => {
		new   => 'Parameter',
		using => {
			id          => 216,
			name        => 'CONFIG proxy.config.ssl.client.CA.cert.path',
			config_file => 'records.config',
			value       => 'STRING etc/trafficserver',
		},
	},
	'211-CONFIG-proxy.config.icp.enabled' => {
		new   => 'Parameter',
		using => {
			id          => 217,
			name        => 'CONFIG proxy.config.icp.enabled',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'212-CONFIG-proxy.config.icp.icp_interface' => {
		new   => 'Parameter',
		using => {
			id          => 218,
			name        => 'CONFIG proxy.config.icp.icp_interface',
			config_file => 'records.config',
			value       => 'STRING NULL',
		},
	},
	'213-CONFIG-proxy.config.icp.icp_port' => {
		new   => 'Parameter',
		using => {
			id          => 219,
			name        => 'CONFIG proxy.config.icp.icp_port',
			config_file => 'records.config',
			value       => 'INT 3130',
		},
	},
	'214-CONFIG-proxy.config.icp.multicast_enabled' => {
		new   => 'Parameter',
		using => {
			id          => 220,
			name        => 'CONFIG proxy.config.icp.multicast_enabled',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'215-CONFIG-proxy.config.icp.query_timeout' => {
		new   => 'Parameter',
		using => {
			id          => 221,
			name        => 'CONFIG proxy.config.icp.query_timeout',
			config_file => 'records.config',
			value       => 'INT 2',
		},
	},
	'216-CONFIG-proxy.config.update.enabled' => {
		new   => 'Parameter',
		using => {
			id          => 222,
			name        => 'CONFIG proxy.config.update.enabled',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'217-CONFIG-proxy.config.update.force' => {
		new   => 'Parameter',
		using => {
			id          => 223,
			name        => 'CONFIG proxy.config.update.force',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'218-CONFIG-proxy.config.update.retry_count' => {
		new   => 'Parameter',
		using => {
			id          => 224,
			name        => 'CONFIG proxy.config.update.retry_count',
			config_file => 'records.config',
			value       => 'INT 10',
		},
	},
	'219-CONFIG-proxy.config.update.retry_interval' => {
		new   => 'Parameter',
		using => {
			id          => 225,
			name        => 'CONFIG proxy.config.update.retry_interval',
			config_file => 'records.config',
			value       => 'INT 2',
		},
	},
	'220-CONFIG-proxy.config.update.concurrent_updates' => {
		new   => 'Parameter',
		using => {
			id          => 226,
			name        => 'CONFIG proxy.config.update.concurrent_updates',
			config_file => 'records.config',
			value       => 'INT 100',
		},
	},
	'221-CONFIG-proxy.config.net.sock_send_buffer_size_in' => {
		new   => 'Parameter',
		using => {
			id          => 227,
			name        => 'CONFIG proxy.config.net.sock_send_buffer_size_in',
			config_file => 'records.config',
			value       => 'INT 262144',
		},
	},
	'222-CONFIG-proxy.config.net.sock_recv_buffer_size_in' => {
		new   => 'Parameter',
		using => {
			id          => 228,
			name        => 'CONFIG proxy.config.net.sock_recv_buffer_size_in',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'223-CONFIG-proxy.config.net.sock_send_buffer_size_out' => {
		new   => 'Parameter',
		using => {
			id          => 229,
			name        => 'CONFIG proxy.config.net.sock_send_buffer_size_out',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'224-CONFIG-proxy.config.net.sock_recv_buffer_size_out' => {
		new   => 'Parameter',
		using => {
			id          => 230,
			name        => 'CONFIG proxy.config.net.sock_recv_buffer_size_out',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'225-CONFIG-proxy.config.core_limit' => {
		new   => 'Parameter',
		using => {
			id          => 231,
			name        => 'CONFIG proxy.config.core_limit',
			config_file => 'records.config',
			value       => 'INT -1',
		},
	},
	'226-CONFIG-proxy.config.diags.debug.enabled' => {
		new   => 'Parameter',
		using => {
			id          => 232,
			name        => 'CONFIG proxy.config.diags.debug.enabled',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'227-CONFIG-proxy.config.diags.debug.tags' => {
		new   => 'Parameter',
		using => {
			id          => 233,
			name        => 'CONFIG proxy.config.diags.debug.tags',
			config_file => 'records.config',
			value       => 'STRING http.*|dns.*',
		},
	},
	'228-CONFIG-proxy.config.dump_mem_info_frequency' => {
		new   => 'Parameter',
		using => {
			id          => 234,
			name        => 'CONFIG proxy.config.dump_mem_info_frequency',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'229-CONFIG-proxy.config.http.slow.log.threshold' => {
		new   => 'Parameter',
		using => {
			id          => 235,
			name        => 'CONFIG proxy.config.http.slow.log.threshold',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'230-CONFIG-proxy.config.task_threads' => {
		new   => 'Parameter',
		using => {
			id          => 236,
			name        => 'CONFIG proxy.config.task_threads',
			config_file => 'records.config',
			value       => 'INT 2',
		},
	},
	'231-location' => {
		new   => 'Parameter',
		using => {
			id          => 263,
			name        => 'location',
			config_file => 'cache.config',
			value       => '/opt/trafficserver/etc/trafficserver/',
		},
	},
	'232-location' => {
		new   => 'Parameter',
		using => {
			id          => 264,
			name        => 'location',
			config_file => 'hosting.config',
			value       => '/opt/trafficserver/etc/trafficserver/',
		},
	},
	'233-location' => {
		new   => 'Parameter',
		using => {
			id          => 265,
			name        => 'location',
			config_file => 'parent.config',
			value       => '/opt/trafficserver/etc/trafficserver/',
		},
	},
	'234-location' => {
		new   => 'Parameter',
		using => {
			id          => 266,
			name        => 'location',
			config_file => 'plugin.config',
			value       => '/opt/trafficserver/etc/trafficserver/',
		},
	},
	'235-location' => {
		new   => 'Parameter',
		using => {
			id          => 267,
			name        => 'location',
			config_file => 'records.config',
			value       => '/opt/trafficserver/etc/trafficserver/',
		},
	},
	'236-location' => {
		new   => 'Parameter',
		using => {
			id          => 268,
			name        => 'location',
			config_file => 'remap.config',
			value       => '/opt/trafficserver/etc/trafficserver/',
		},
	},
	'237-location' => {
		new   => 'Parameter',
		using => {
			id          => 269,
			name        => 'location',
			config_file => 'storage.config',
			value       => '/opt/trafficserver/etc/trafficserver/',
		},
	},
	'238-location' => {
		new   => 'Parameter',
		using => {
			id          => 270,
			name        => 'location',
			config_file => 'volume.config',
			value       => '/opt/trafficserver/etc/trafficserver/',
		},
	},
	'239-location' => {
		new   => 'Parameter',
		using => {
			id          => 273,
			name        => 'location',
			config_file => '50-ats.rules',
			value       => '/etc/udev/rules.d/',
		},
	},
	'240-location' => {
		new   => 'Parameter',
		using => {
			id          => 276,
			name        => 'location',
			config_file => 'CRConfig.xml',
			value       => 'XMPP CRConfig node',
		},
	},
	'241-location' => {
		new   => 'Parameter',
		using => {
			id          => 277,
			name        => 'location',
			config_file => 'dns.zone',
			value       => '/etc/kabletown/zones/<zonename>.info',
		},
	},
	'242-CONFIG-proxy.config.http.parent_proxy_routing_enable' => {
		new   => 'Parameter',
		using => {
			id          => 278,
			name        => 'CONFIG proxy.config.http.parent_proxy_routing_enable',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'243-CONFIG-proxy.config.url_remap.remap_required' => {
		new   => 'Parameter',
		using => {
			id          => 279,
			name        => 'CONFIG proxy.config.url_remap.remap_required',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'244-location' => {
		new   => 'Parameter',
		using => {
			id          => 291,
			name        => 'location',
			config_file => 'http-log4j.properties',
			value       => '/etc/kabletown',
		},
	},
	'245-location' => {
		new   => 'Parameter',
		using => {
			id          => 292,
			name        => 'location',
			config_file => 'dns-log4j.properties',
			value       => '/etc/kabletown',
		},
	},
	'246-location' => {
		new   => 'Parameter',
		using => {
			id          => 293,
			name        => 'location',
			config_file => 'geolocation.properties',
			value       => '/etc/kabletown',
		},
	},
	'247-domain_name' => {
		new   => 'Parameter',
		using => {
			id          => 295,
			name        => 'domain_name',
			config_file => 'CRConfig.xml',
			value       => 'cdn2.kabletown.net',
		},
	},
	'248-CONFIG-proxy.config.http.parent_proxy.file' => {
		new   => 'Parameter',
		using => {
			id          => 325,
			name        => 'CONFIG proxy.config.http.parent_proxy.file',
			config_file => 'records.config',
			value       => 'STRING parent.config',
		},
	},
	'249-CONFIG-proxy.config.url_remap.filename' => {
		new   => 'Parameter',
		using => {
			id          => 326,
			name        => 'CONFIG proxy.config.url_remap.filename',
			config_file => 'records.config',
			value       => 'STRING remap.config',
		},
	},
	'250-location' => {
		new   => 'Parameter',
		using => {
			id          => 327,
			name        => 'location',
			config_file => 'ip_allow.config',
			value       => '/opt/trafficserver/etc/trafficserver',
		},
	},
	'251-CONFIG-proxy.config.cluster.cluster_configuration-' => {
		new   => 'Parameter',
		using => {
			id          => 328,
			name        => 'CONFIG proxy.config.cluster.cluster_configuration ',
			config_file => 'records.config',
			value       => 'STRING cluster.config',
		},
	},
	'252-Drive_Prefix' => {
		new   => 'Parameter',
		using => {
			id          => 329,
			name        => 'Drive_Prefix',
			config_file => 'storage.config',
			value       => '/dev/ram',
		},
	},
	'253-ramdisk_size' => {
		new   => 'Parameter',
		using => {
			id          => 330,
			name        => 'ramdisk_size',
			config_file => 'grub.conf',
			value       => 'ramdisk_size=16777216',
		},
	},
	'254-cron_syncds' => {
		new   => 'Parameter',
		using => {
			id          => 331,
			name        => 'cron_syncds',
			config_file => 'crontab_root',
			value       => '*/15 * * * * /opt/ort/ipcdn_install_ort.pl syncds error &amp;gt; /tmp/ort/syncds.log 2&amp;gt;&amp;amp;1',
		},
	},
	'255-location' => {
		new   => 'Parameter',
		using => {
			id          => 332,
			name        => 'location',
			config_file => 'crontab_root',
			value       => '/var/spool/cron',
		},
	},
	'256-CONFIG-proxy.config.http.insert_age_in_response' => {
		new   => 'Parameter',
		using => {
			id          => 333,
			name        => 'CONFIG proxy.config.http.insert_age_in_response',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'257-monitor:///opt/tomcat/logs/access.log' => {
		new   => 'Parameter',
		using => {
			id          => 334,
			name        => 'monitor:///opt/tomcat/logs/access.log',
			config_file => 'inputs.conf',
			value       => 'index=index_odol_test;sourcetype=access_ccr',
		},
	},
	'258-location' => {
		new   => 'Parameter',
		using => {
			id          => 341,
			name        => 'location',
			config_file => 'CRConfig.xml',
			value       => 'XMPP CRConfigOTT node',
		},
	},
	'277-purge_allow_ip' => {
		new   => 'Parameter',
		using => {
			id          => 360,
			name        => 'purge_allow_ip',
			config_file => 'ip_allow.config',
			value       => '33.101.99.100',
		},
	},
	'278-astats_over_http.so' => {
		new   => 'Parameter',
		using => {
			id          => 361,
			name        => 'astats_over_http.so',
			config_file => 'plugin.config',
			value       => '_astats 33.101.99.100,172.39.19.39,172.39.19.49,172.39.19.49,172.39.29.49',
		},
	},
	'279-health.threshold.loadavg' => {
		new   => 'Parameter',
		using => {
			id          => 363,
			name        => 'health.threshold.loadavg',
			config_file => 'rascal.properties',
			value       => '25.0',
		},
	},
	'280-health.threshold.availableBandwidthInKbps' => {
		new   => 'Parameter',
		using => {
			id          => 364,
			name        => 'health.threshold.availableBandwidthInKbps',
			config_file => 'rascal.properties',
			value       => '>1750000',
		},
	},
	'281-history.count' => {
		new   => 'Parameter',
		using => {
			id          => 366,
			name        => 'history.count',
			config_file => 'rascal.properties',
			value       => '30',
		},
	},
	'282-cacheurl.so' => {
		new   => 'Parameter',
		using => {
			id          => 367,
			name        => 'cacheurl.so',
			config_file => 'plugin.config',
			value       => '',
		},
	},
	'283-location' => {
		new   => 'Parameter',
		using => {
			id          => 368,
			name        => 'location',
			config_file => 'cacheurl.config',
			value       => '/opt/trafficserver/etc/trafficserver/',
		},
	},
	'284-CONFIG-proxy.config.cache.control.filename' => {
		new   => 'Parameter',
		using => {
			id          => 369,
			name        => 'CONFIG proxy.config.cache.control.filename',
			config_file => 'records.config',
			value       => 'STRING cache.config',
		},
	},
	'285-LogFormat.Name' => {
		new   => 'Parameter',
		using => {
			id          => 370,
			name        => 'LogFormat.Name',
			config_file => 'logs_xml.config',
			value       => 'custom_ats_2',
		},
	},
	'286-LogFormat.Format' => {
		new   => 'Parameter',
		using => {
			id          => 371,
			name        => 'LogFormat.Format',
			config_file => 'logs_xml.config',
			value =>
				'%<chi> %<caun> [%<cqtq>] "%<cqtx>" %<pssc> %<pscl> %<sssc> %<sscl> %<cqbl> %<pqbl> %<cqhl> %<pshl> %<ttms> %<pqhl> %<sshl> %<phr> %<cfsc> %<pfsc> %<crc> "%<{User-Agent}cqh>"',
		},
	},
	'287-LogObject.Format' => {
		new   => 'Parameter',
		using => {
			id          => 372,
			name        => 'LogObject.Format',
			config_file => 'logs_xml.config',
			value       => 'custom_ats_2',
		},
	},
	'288-LogObject.Filename' => {
		new   => 'Parameter',
		using => {
			id          => 373,
			name        => 'LogObject.Filename',
			config_file => 'logs_xml.config',
			value       => 'custom_ats_2',
		},
	},
	'289-LogObject.RollingEnabled' => {
		new   => 'Parameter',
		using => {
			id          => 374,
			name        => 'LogObject.RollingEnabled',
			config_file => 'logs_xml.config',
			value       => '3',
		},
	},
	'290-LogObject.RollingIntervalSec' => {
		new   => 'Parameter',
		using => {
			id          => 375,
			name        => 'LogObject.RollingIntervalSec',
			config_file => 'logs_xml.config',
			value       => '86400',
		},
	},
	'291-LogObject.RollingOffsetHr' => {
		new   => 'Parameter',
		using => {
			id          => 376,
			name        => 'LogObject.RollingOffsetHr',
			config_file => 'logs_xml.config',
			value       => '11',
		},
	},
	'292-LogObject.RollingSizeMb' => {
		new   => 'Parameter',
		using => {
			id          => 377,
			name        => 'LogObject.RollingSizeMb',
			config_file => 'logs_xml.config',
			value       => '1024',
		},
	},
	'293-location' => {
		new   => 'Parameter',
		using => {
			id          => 378,
			name        => 'location',
			config_file => 'logs_xml.config',
			value       => '/opt/trafficserver/etc/trafficserver',
		},
	},
	'294-CDN_name' => {
		new   => 'Parameter',
		using => {
			id          => 379,
			name        => 'CDN_name',
			config_file => 'rascal-config.txt',
			value       => 'cdn_number_1',
		},
	},
	'295-CDN_name' => {
		new   => 'Parameter',
		using => {
			id          => 380,
			name        => 'CDN_name',
			config_file => 'rascal-config.txt',
			value       => 'cdn_number_2',
		},
	},
	'296-CONFIG-proxy.config.log.xml_config_file' => {
		new   => 'Parameter',
		using => {
			id          => 381,
			name        => 'CONFIG proxy.config.log.xml_config_file',
			config_file => 'records.config',
			value       => 'STRING logs_xml.config',
		},
	},
	'297-health.polling.interval' => {
		new   => 'Parameter',
		using => {
			id          => 382,
			name        => 'health.polling.interval',
			config_file => 'rascal-config.txt',
			value       => '8000',
		},
	},
	'298-tm.crConfig.polling.url' => {
		new   => 'Parameter',
		using => {
			id          => 383,
			name        => 'tm.crConfig.polling.url',
			config_file => 'rascal-config.txt',
			value       => 'https://${tmHostname}/CRConfig-Snapshots/${cdnName}/CRConfig.xml',
		},
	},
	'299-tm.dataServer.polling.url' => {
		new   => 'Parameter',
		using => {
			id          => 384,
			name        => 'tm.dataServer.polling.url',
			config_file => 'rascal-config.txt',
			value       => 'https://${tmHostname}/dataserver/orderby/id',
		},
	},
	'300-tm.healthParams.polling.url' => {
		new   => 'Parameter',
		using => {
			id          => 385,
			name        => 'tm.healthParams.polling.url',
			config_file => 'rascal-config.txt',
			value       => 'https://${tmHostname}/health/${cdnName}',
		},
	},
	'301-tm.polling.interval' => {
		new   => 'Parameter',
		using => {
			id          => 386,
			name        => 'tm.polling.interval',
			config_file => 'rascal-config.txt',
			value       => '60000',
		},
	},
	'302-location' => {
		new   => 'Parameter',
		using => {
			id          => 387,
			name        => 'location',
			config_file => 'rascal-config.txt',
			value       => '/opt/traffic_monitor/conf',
		},
	},
	'303-health.threshold.queryTime' => {
		new   => 'Parameter',
		using => {
			id          => 388,
			name        => 'health.threshold.queryTime',
			config_file => 'rascal.properties',
			value       => '1000',
		},
	},
	'304-health.polling.url' => {
		new   => 'Parameter',
		using => {
			id          => 389,
			name        => 'health.polling.url',
			config_file => 'rascal.properties',
			value       => 'http://${hostname}/_astats?application=&inf.name=${interface_name}',
		},
	},
	'305-health.threadPool' => {
		new   => 'Parameter',
		using => {
			id          => 390,
			name        => 'health.threadPool',
			config_file => 'rascal-config.txt',
			value       => '4',
		},
	},
	'306-health.event-count' => {
		new   => 'Parameter',
		using => {
			id          => 391,
			name        => 'health.event-count',
			config_file => 'rascal-config.txt',
			value       => '200',
		},
	},
	'307-hack.ttl' => {
		new   => 'Parameter',
		using => {
			id          => 392,
			name        => 'hack.ttl',
			config_file => 'rascal-config.txt',
			value       => '30',
		},
	},
	'308-RAM_Drive_Prefix' => {
		new   => 'Parameter',
		using => {
			id          => 393,
			name        => 'RAM_Drive_Prefix',
			config_file => 'storage.config',
			value       => '/dev/ram',
		},
	},
	'309-RAM_Drive_Letters' => {
		new   => 'Parameter',
		using => {
			id          => 394,
			name        => 'RAM_Drive_Letters',
			config_file => 'storage.config',
			value       => '0,1,2,3,4,5,6,7',
		},
	},
	'310-RAM_Volume' => {
		new   => 'Parameter',
		using => {
			id          => 395,
			name        => 'RAM_Volume',
			config_file => 'storage.config',
			value       => '2',
		},
	},
	'311-Disk_Volume' => {
		new   => 'Parameter',
		using => {
			id          => 396,
			name        => 'Disk_Volume',
			config_file => 'storage.config',
			value       => '1',
		},
	},
	'312-CONFIG-proxy.config.cache.hosting_filename' => {
		new   => 'Parameter',
		using => {
			id          => 397,
			name        => 'CONFIG proxy.config.cache.hosting_filename',
			config_file => 'records.config',
			value       => 'STRING hosting.config',
		},
	},
	'313-CoverageZoneJsonURL' => {
		new   => 'Parameter',
		using => {
			id          => 398,
			name        => 'CoverageZoneJsonURL',
			config_file => 'CRConfig.xml',
			value       => 'http://staging.cdnlab.kabletown.net/ipcdn/CZF/current/kabletown_ipcdn_czf-current.json',
		},
	},
	'314-health.connection.timeout' => {
		new   => 'Parameter',
		using => {
			id          => 399,
			name        => 'health.connection.timeout',
			config_file => 'rascal.properties',
			value       => '2000',
		},
	},
	'315-geolocation.polling.url' => {
		new   => 'Parameter',
		using => {
			id          => 400,
			name        => 'geolocation.polling.url',
			config_file => 'CRConfig.json',
			value       => 'https://tm.kabletown.net/MaxMind/GeoLiteCity.dat.gz',
		},
	},
	'316-geolocation.polling.interval' => {
		new   => 'Parameter',
		using => {
			id          => 401,
			name        => 'geolocation.polling.interval',
			config_file => 'CRConfig.json',
			value       => '86400000',
		},
	},
	'317-coveragezone.polling.interval' => {
		new   => 'Parameter',
		using => {
			id          => 402,
			name        => 'coveragezone.polling.interval',
			config_file => 'CRConfig.json',
			value       => '86400000',
		},
	},
	'318-coveragezone.polling.url' => {
		new   => 'Parameter',
		using => {
			id          => 403,
			name        => 'coveragezone.polling.url',
			config_file => 'CRConfig.json',
			value       => 'http://staging.cdnlab.kabletown.net/ipcdn/CZF/current/kabletown_ipcdn_czf-current.json',
		},
	},
	'319-domain_name' => {
		new   => 'Parameter',
		using => {
			id          => 404,
			name        => 'domain_name',
			config_file => 'CRConfig.json',
			value       => 'cdn1.kabletown.net',
		},
	},
	'320-domain_name' => {
		new   => 'Parameter',
		using => {
			id          => 405,
			name        => 'domain_name',
			config_file => 'CRConfig.json',
			value       => 'cdn2.kabletown.net',
		},
	},
	'321-location' => {
		new   => 'Parameter',
		using => {
			id          => 406,
			name        => 'location',
			config_file => '12M_facts',
			value       => '/opt/ort',
		},
	},
	'322-CONFIG-proxy.config.http.cache.ignore_accept_encoding_mismatch' => {
		new   => 'Parameter',
		using => {
			id          => 407,
			name        => 'CONFIG proxy.config.http.cache.ignore_accept_encoding_mismatch',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'323-health.timepad' => {
		new   => 'Parameter',
		using => {
			id          => 408,
			name        => 'health.timepad',
			config_file => 'rascal-config.txt',
			value       => '30',
		},
	},
	'342-tm.url' => {
		new   => 'Parameter',
		using => {
			id          => 502,
			name        => 'tm.url',
			config_file => 'global',
			value       => 'https://tm.kabletown.net/',
		},
	},
	'344-tm.toolname' => {
		new   => 'Parameter',
		using => {
			id          => 504,
			name        => 'tm.toolname',
			config_file => 'global',
			value       => 'Traffic Ops',
		},
	},
	'345-tm.infourl' => {
		new   => 'Parameter',
		using => {
			id          => 505,
			name        => 'tm.infourl',
			config_file => 'global',
			value       => 'http://staging.cdnlab.kabletown.net/tm/info',
		},
	},
	'346-tm.logourl' => {
		new   => 'Parameter',
		using => {
			id          => 506,
			name        => 'tm.logourl',
			config_file => 'global',
			value       => '/images/tc_logo.png',
		},
	},
	'347-tld.ttls.AAAA' => {
		new   => 'Parameter',
		using => {
			id          => 507,
			name        => 'tld.ttls.AAAA',
			config_file => 'CRConfig.json',
			value       => '3600',
		},
	},
	'348-tld.ttls.SOA' => {
		new   => 'Parameter',
		using => {
			id          => 508,
			name        => 'tld.ttls.SOA',
			config_file => 'CRConfig.json',
			value       => '86400',
		},
	},
	'349-tld.ttls.A' => {
		new   => 'Parameter',
		using => {
			id          => 509,
			name        => 'tld.ttls.A',
			config_file => 'CRConfig.json',
			value       => '3600',
		},
	},
	'350-tld.ttls.NS' => {
		new   => 'Parameter',
		using => {
			id          => 510,
			name        => 'tld.ttls.NS',
			config_file => 'CRConfig.json',
			value       => '3600',
		},
	},
	'351-tld.soa.expire' => {
		new   => 'Parameter',
		using => {
			id          => 511,
			name        => 'tld.soa.expire',
			config_file => 'CRConfig.json',
			value       => '604800',
		},
	},
	'352-tld.soa.minimum' => {
		new   => 'Parameter',
		using => {
			id          => 512,
			name        => 'tld.soa.minimum',
			config_file => 'CRConfig.json',
			value       => '86400',
		},
	},
	'353-tld.soa.admin' => {
		new   => 'Parameter',
		using => {
			id          => 513,
			name        => 'tld.soa.admin',
			config_file => 'CRConfig.json',
			value       => 'traffic_ops',
		},
	},
	'354-tld.soa.retry' => {
		new   => 'Parameter',
		using => {
			id          => 514,
			name        => 'tld.soa.retry',
			config_file => 'CRConfig.json',
			value       => '7200',
		},
	},
	'355-tld.soa.refresh' => {
		new   => 'Parameter',
		using => {
			id          => 515,
			name        => 'tld.soa.refresh',
			config_file => 'CRConfig.json',
			value       => '28800',
		},
	},
	'356-key0' => {
		new   => 'Parameter',
		using => {
			id          => 551,
			name        => 'key0',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'HOOJ3Ghq1x4gChp3iQkqVTcPlOj8UCi3',
		},
	},
	'357-key1' => {
		new   => 'Parameter',
		using => {
			id          => 552,
			name        => 'key1',
			config_file => 'url_sig_cdl-c2.config',
			value       => '_9LZYkRnfCS0rCBF7fTQzM9Scwlp2FhO',
		},
	},
	'358-key2' => {
		new   => 'Parameter',
		using => {
			id          => 553,
			name        => 'key2',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'AFpkxfc4oTiyFSqtY6_ohjt3V80aAIxS',
		},
	},
	'359-key3' => {
		new   => 'Parameter',
		using => {
			id          => 554,
			name        => 'key3',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'AL9kzs_SXaRZjPWH8G5e2m4ByTTzkzlc',
		},
	},
	'360-key4' => {
		new   => 'Parameter',
		using => {
			id          => 555,
			name        => 'key4',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'poP3n3szbD1U4vx1xQXV65BvkVgWzfN8',
		},
	},
	'361-key5' => {
		new   => 'Parameter',
		using => {
			id          => 556,
			name        => 'key5',
			config_file => 'url_sig_cdl-c2.config',
			value       => '1ir32ng4C4w137p5oq72kd2wqmIZUrya',
		},
	},
	'362-key6' => {
		new   => 'Parameter',
		using => {
			id          => 557,
			name        => 'key6',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'B1qLptn2T1b_iXeTCWDcVuYvANtH139f',
		},
	},
	'363-key7' => {
		new   => 'Parameter',
		using => {
			id          => 558,
			name        => 'key7',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'PiCV_5OODMzBbsNFMWsBxcQ8v1sK0TYE',
		},
	},
	'364-key8' => {
		new   => 'Parameter',
		using => {
			id          => 559,
			name        => 'key8',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'Ggpv6DqXDvt2s1CETPBpNKwaLk4fTM9l',
		},
	},
	'365-key9' => {
		new   => 'Parameter',
		using => {
			id          => 560,
			name        => 'key9',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'qPlVT_s6kL37aqb6hipDm4Bt55S72mI7',
		},
	},
	'366-key10' => {
		new   => 'Parameter',
		using => {
			id          => 561,
			name        => 'key10',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'BsI5A9EmWrobIS1FeuOs1z9fm2t2WSBe',
		},
	},
	'367-key11' => {
		new   => 'Parameter',
		using => {
			id          => 562,
			name        => 'key11',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'A54y66NCIj897GjS4yA9RrsSPtCUnQXP',
		},
	},
	'368-key12' => {
		new   => 'Parameter',
		using => {
			id          => 563,
			name        => 'key12',
			config_file => 'url_sig_cdl-c2.config',
			value       => '2jZH0NDPSJttIr4c2KP510f47EKqTQAu',
		},
	},
	'369-key13' => {
		new   => 'Parameter',
		using => {
			id          => 564,
			name        => 'key13',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'XduT2FBjBmmVID5JRB5LEf9oR5QDtBgC',
		},
	},
	'370-key14' => {
		new   => 'Parameter',
		using => {
			id          => 565,
			name        => 'key14',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'D9nH0SvK_0kP5w8QNd1UFJ28ulFkFKPn',
		},
	},
	'371-key15' => {
		new   => 'Parameter',
		using => {
			id          => 566,
			name        => 'key15',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'udKXWYNwbXXweaaLzaKDGl57OixnIIcm',
		},
	},
	'372-location' => {
		new   => 'Parameter',
		using => {
			id          => 567,
			name        => 'location',
			config_file => 'url_sig_cdl-c2.config',
			value       => '/opt/trafficserver/etc/trafficserver',
		},
	},
	'373-error_url' => {
		new   => 'Parameter',
		using => {
			id          => 568,
			name        => 'error_url',
			config_file => 'url_sig_cdl-c2.config',
			value       => '403',
		},
	},
	'374-tld.ttls.NS' => {
		new   => 'Parameter',
		using => {
			id          => 591,
			name        => 'tld.ttls.NS',
			config_file => 'CRConfig.json',
			value       => '3600',
		},
	},
	'375-geolocation6.polling.url' => {
		new   => 'Parameter',
		using => {
			id          => 592,
			name        => 'geolocation6.polling.url',
			config_file => 'CRConfig.json',
			value       => 'https://tm.kabletown.net/MaxMind/GeoLiteCityv6.dat.gz',
		},
	},
	'376-geolocation6.polling.interval' => {
		new   => 'Parameter',
		using => {
			id          => 593,
			name        => 'geolocation6.polling.interval',
			config_file => 'CRConfig.json',
			value       => '86400000',
		},
	},
	'377-trafficserver' => {
		new   => 'Parameter',
		using => {
			id          => 594,
			name        => 'trafficserver',
			config_file => 'chkconfig',
			value       => '0:off	1:off	2:on	3:on	4:on	5:on	6:off',
		},
	},
	'378-astats_over_http' => {
		new   => 'Parameter',
		using => {
			id          => 595,
			name        => 'astats_over_http',
			config_file => 'package',
			value       => '1.1-2.el6.x86_64',
		},
	},
	'379-cacheurl' => {
		new   => 'Parameter',
		using => {
			id          => 596,
			name        => 'cacheurl',
			config_file => 'package',
			value       => '1.0-1.el6.x86_64',
		},
	},
	'380-dscp_remap' => {
		new   => 'Parameter',
		using => {
			id          => 597,
			name        => 'dscp_remap',
			config_file => 'package',
			value       => '1.0-1.el6.x86_64',
		},
	},
	'381-regex_revalidate' => {
		new   => 'Parameter',
		using => {
			id          => 598,
			name        => 'regex_revalidate',
			config_file => 'package',
			value       => '1.0-1.el6.x86_64',
		},
	},
	'382-remap_stats' => {
		new   => 'Parameter',
		using => {
			id          => 599,
			name        => 'remap_stats',
			config_file => 'package',
			value       => '1.0-1.el6.x86_64',
		},
	},
	'383-url_sign' => {
		new   => 'Parameter',
		using => {
			id          => 600,
			name        => 'url_sign',
			config_file => 'package',
			value       => '1.0-1.el6.x86_64',
		},
	},
	'384-trafficserver' => {
		new   => 'Parameter',
		using => {
			id          => 601,
			name        => 'trafficserver',
			config_file => 'package',
			value       => '4.0.2-2.el6.x86_64',
		},
	},
	'385-astats_over_http' => {
		new   => 'Parameter',
		using => {
			id          => 602,
			name        => 'astats_over_http',
			config_file => 'package',
			value       => '3.2.0-4114.el6.x86_64',
		},
	},
	'386-cacheurl' => {
		new   => 'Parameter',
		using => {
			id          => 603,
			name        => 'cacheurl',
			config_file => 'package',
			value       => '3.2.0-5628.el6.x86_64',
		},
	},
	'387-dscp_remap' => {
		new   => 'Parameter',
		using => {
			id          => 604,
			name        => 'dscp_remap',
			config_file => 'package',
			value       => '3.2.0-4613.el6.x86_64',
		},
	},
	'388-regex_revalidate' => {
		new   => 'Parameter',
		using => {
			id          => 605,
			name        => 'regex_revalidate',
			config_file => 'package',
			value       => '3.2.0-5695.el6.x86_64',
		},
	},
	'389-remap_stats' => {
		new   => 'Parameter',
		using => {
			id          => 606,
			name        => 'remap_stats',
			config_file => 'package',
			value       => '3.2.0-2.el6.x86_64',
		},
	},
	'390-url_sign' => {
		new   => 'Parameter',
		using => {
			id          => 607,
			name        => 'url_sign',
			config_file => 'package',
			value       => '3.2.0-4130.el6.x86_64',
		},
	},
	'391-trafficserver' => {
		new   => 'Parameter',
		using => {
			id          => 608,
			name        => 'trafficserver',
			config_file => 'package',
			value       => '3.2.0-4812.el6.x86_64',
		},
	},
	'392-CONFIG-proxy.config.allocator.debug_filter' => {
		new   => 'Parameter',
		using => {
			id          => 609,
			name        => 'CONFIG proxy.config.allocator.debug_filter',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'393-CONFIG-proxy.config.allocator.enable_reclaim' => {
		new   => 'Parameter',
		using => {
			id          => 610,
			name        => 'CONFIG proxy.config.allocator.enable_reclaim',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'394-CONFIG-proxy.config.allocator.max_overage' => {
		new   => 'Parameter',
		using => {
			id          => 611,
			name        => 'CONFIG proxy.config.allocator.max_overage',
			config_file => 'records.config',
			value       => 'INT 3',
		},
	},
	'395-CONFIG-proxy.config.diags.show_location' => {
		new   => 'Parameter',
		using => {
			id          => 612,
			name        => 'CONFIG proxy.config.diags.show_location',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'396-CONFIG-proxy.config.http.cache.allow_empty_doc' => {
		new   => 'Parameter',
		using => {
			id          => 613,
			name        => 'CONFIG proxy.config.http.cache.allow_empty_doc',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'397-LOCAL-proxy.config.cache.interim.storage' => {
		new   => 'Parameter',
		using => {
			id          => 614,
			name        => 'LOCAL proxy.config.cache.interim.storage',
			config_file => 'records.config',
			value       => 'STRING NULL',
		},
	},
	'398-tld.ttls.SOA' => {
		new   => 'Parameter',
		using => {
			id          => 615,
			name        => 'tld.ttls.SOA',
			config_file => 'CRConfig.json',
			value       => '86400',
		},
	},
	'399-regex_revalidate.so' => {
		new   => 'Parameter',
		using => {
			id          => 616,
			name        => 'regex_revalidate.so',
			config_file => 'plugin.config',
			value       => '--config regex_revalidate.config',
		},
	},
	'400-location' => {
		new   => 'Parameter',
		using => {
			id          => 618,
			name        => 'location',
			config_file => 'regex_revalidate.config',
			value       => '/opt/trafficserver/etc/trafficserver',
		},
	},
	'419-remap_stats.so' => {
		new   => 'Parameter',
		using => {
			id          => 640,
			name        => 'remap_stats.so',
			config_file => 'plugin.config',
			value       => '',
		},
	},
	'420-remap_stats.so' => {
		new   => 'Parameter',
		using => {
			id          => 641,
			name        => 'remap_stats.so',
			config_file => 'plugin.config',
			value       => '',
		},
	},
	'421-remap_stats' => {
		new   => 'Parameter',
		using => {
			id          => 642,
			name        => 'remap_stats',
			config_file => 'package',
			value       => '3.2.0-3.el6.x86_64',
		},
	},
	'422-CONFIG-proxy.config.stack_dump_enabled' => {
		new   => 'Parameter',
		using => {
			id          => 643,
			name        => 'CONFIG proxy.config.stack_dump_enabled',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'423-CONFIG-proxy.config.stack_dump_enabled' => {
		new   => 'Parameter',
		using => {
			id          => 645,
			name        => 'CONFIG proxy.config.stack_dump_enabled',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'442-location' => {
		new   => 'Parameter',
		using => {
			id          => 666,
			name        => 'location',
			config_file => 'drop_qstring.config',
			value       => '/opt/trafficserver/etc/trafficserver',
		},
	},
	'443-Drive_Letters' => {
		new   => 'Parameter',
		using => {
			id          => 667,
			name        => 'Drive_Letters',
			config_file => 'storage.config',
			value       => 'b,c,d,e,f,g,h,i,j,k,l,m,n,o',
		},
	},
	'444-CONFIG-proxy.config.hostdb.ttl_mode' => {
		new   => 'Parameter',
		using => {
			id          => 668,
			name        => 'CONFIG proxy.config.hostdb.ttl_mode',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'445-CONFIG-proxy.config.dns.lookup_timeout' => {
		new   => 'Parameter',
		using => {
			id          => 669,
			name        => 'CONFIG proxy.config.dns.lookup_timeout',
			config_file => 'records.config',
			value       => 'INT 2',
		},
	},
	'446-CONFIG-proxy.config.hostdb.serve_stale_for' => {
		new   => 'Parameter',
		using => {
			id          => 670,
			name        => 'CONFIG proxy.config.hostdb.serve_stale_for',
			config_file => 'records.config',
			value       => 'INT 6',
		},
	},
	'447-trafficserver' => {
		new   => 'Parameter',
		using => {
			id          => 671,
			name        => 'trafficserver',
			config_file => 'package',
			value       => '4.2.1-6.el6.x86_64',
		},
	},
	'448-CONFIG-proxy.config.cache.enable_read_while_writer' => {
		new   => 'Parameter',
		using => {
			id          => 678,
			name        => 'CONFIG proxy.config.cache.enable_read_while_writer',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'449-CONFIG-proxy.config.http.background_fill_active_timeout' => {
		new   => 'Parameter',
		using => {
			id          => 679,
			name        => 'CONFIG proxy.config.http.background_fill_active_timeout',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'450-CONFIG-proxy.config.http.background_fill_completed_threshold' => {
		new   => 'Parameter',
		using => {
			id          => 680,
			name        => 'CONFIG proxy.config.http.background_fill_completed_threshold',
			config_file => 'records.config',
			value       => 'FLOAT 0.0',
		},
	},
	'451-CONFIG-proxy.config.log.extended2_log_enabled' => {
		new   => 'Parameter',
		using => {
			id          => 681,
			name        => 'CONFIG proxy.config.log.extended2_log_enabled',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'452-CONFIG-proxy.config.exec_thread.affinity' => {
		new   => 'Parameter',
		using => {
			id          => 682,
			name        => 'CONFIG proxy.config.exec_thread.affinity',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'453-CONFIG-proxy.config.exec_thread.autoconfig' => {
		new   => 'Parameter',
		using => {
			id          => 683,
			name        => 'CONFIG proxy.config.exec_thread.autoconfig',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'454-CONFIG-proxy.config.exec_thread.limit' => {
		new   => 'Parameter',
		using => {
			id          => 684,
			name        => 'CONFIG proxy.config.exec_thread.limit',
			config_file => 'records.config',
			value       => 'INT 32',
		},
	},
	'455-CONFIG-proxy.config.allocator.thread_freelist_size' => {
		new   => 'Parameter',
		using => {
			id          => 685,
			name        => 'CONFIG proxy.config.allocator.thread_freelist_size',
			config_file => 'records.config',
			value       => 'INT 1024',
		},
	},
	'456-CONFIG-proxy.config.cache.ram_cache.size' => {
		new   => 'Parameter',
		using => {
			id          => 687,
			name        => 'CONFIG proxy.config.cache.ram_cache.size',
			config_file => 'records.config',
			value       => 'INT 16106127360',
		},
	},
	'457-CONFIG-proxy.config.mlock_enabled' => {
		new   => 'Parameter',
		using => {
			id          => 688,
			name        => 'CONFIG proxy.config.mlock_enabled',
			config_file => 'records.config',
			value       => 'INT 2',
		},
	},
	'458-LogFormat.Format' => {
		new   => 'Parameter',
		using => {
			id          => 689,
			name        => 'LogFormat.Format',
			config_file => 'logs_xml.config',
			value =>
				'%<cqtq> chi=%<chi> phn=%<phn> shn=%<shn> url=%<cquuc> cqhm=%<cqhm> cqhv=%<cqhv> pssc=%<pssc> ttms=%<ttms> b=%<pscl> sssc=%<sssc> sscl=%<sscl> cfsc=%<cfsc> pfsc=%<pfsc> crc=%<crc> phr=%<phr> uas="%<{User-Agent}cqh>"',
		},
	},
	'459-LogFormat.Name' => {
		new   => 'Parameter',
		using => {
			id          => 690,
			name        => 'LogFormat.Name',
			config_file => 'logs_xml.config',
			value       => 'custom_ats_2',
		},
	},
	'460-LogObject.Format' => {
		new   => 'Parameter',
		using => {
			id          => 691,
			name        => 'LogObject.Format',
			config_file => 'logs_xml.config',
			value       => 'custom_ats_2',
		},
	},
	'461-LogObject.Filename' => {
		new   => 'Parameter',
		using => {
			id          => 692,
			name        => 'LogObject.Filename',
			config_file => 'logs_xml.config',
			value       => 'custom_ats_2',
		},
	},
	'462-url_sig' => {
		new   => 'Parameter',
		using => {
			id          => 693,
			name        => 'url_sig',
			config_file => 'package',
			value       => '1.0-3.el6.x86_64',
		},
	},
	'463-astats_over_http' => {
		new   => 'Parameter',
		using => {
			id          => 694,
			name        => 'astats_over_http',
			config_file => 'package',
			value       => '1.2-3.el6.x86_64',
		},
	},
	'464-cacheurl' => {
		new   => 'Parameter',
		using => {
			id          => 695,
			name        => 'cacheurl',
			config_file => 'package',
			value       => '1.0-3.el6.x86_64',
		},
	},
	'465-dscp_remap' => {
		new   => 'Parameter',
		using => {
			id          => 696,
			name        => 'dscp_remap',
			config_file => 'package',
			value       => '1.0-3.el6.x86_64',
		},
	},
	'466-remap_stats' => {
		new   => 'Parameter',
		using => {
			id          => 697,
			name        => 'remap_stats',
			config_file => 'package',
			value       => '1.0-4.el6.x86_64',
		},
	},
	'467-regex_revalidate' => {
		new   => 'Parameter',
		using => {
			id          => 698,
			name        => 'regex_revalidate',
			config_file => 'package',
			value       => '1.0-4.el6.x86_64',
		},
	},
	'468-CONFIG-proxy.config.cache.ram_cache.size' => {
		new   => 'Parameter',
		using => {
			id          => 699,
			name        => 'CONFIG proxy.config.cache.ram_cache.size',
			config_file => 'records.config',
			value       => 'INT 34359738368',
		},
	},
	'469-header_rewrite' => {
		new   => 'Parameter',
		using => {
			id          => 700,
			name        => 'header_rewrite',
			config_file => 'package',
			value       => '4.0.2-1.el6.x86_64',
		},
	},
	'470-api.port' => {
		new   => 'Parameter',
		using => {
			id          => 701,
			name        => 'api.port',
			config_file => 'server.xml',
			value       => '8080',
		},
	},
	'471-api.port' => {
		new   => 'Parameter',
		using => {
			id          => 702,
			name        => 'api.port',
			config_file => 'server.xml',
			value       => '8080',
		},
	},
	'472-astats_over_http.so' => {
		new   => 'Parameter',
		using => {
			id          => 703,
			name        => 'astats_over_http.so',
			config_file => 'plugin.config',
			value       => '',
		},
	},
	'473-allow_ip' => {
		new   => 'Parameter',
		using => {
			id          => 704,
			name        => 'allow_ip',
			config_file => 'astats.config',
			value       => '127.0.0.1,172.39.0.0/16,33.101.99.0/24',
		},
	},
	'474-allow_ip6' => {
		new   => 'Parameter',
		using => {
			id          => 705,
			name        => 'allow_ip6',
			config_file => 'astats.config',
			value       => '::1,2033:D011:3300::336/64,2033:D011:3300::335/64,2033:D021:3300::333/64,2033:D021:3300::334/64',
		},
	},
	'475-record_types' => {
		new   => 'Parameter',
		using => {
			id          => 706,
			name        => 'record_types',
			config_file => 'astats.config',
			value       => '144',
		},
	},
	'476-location' => {
		new   => 'Parameter',
		using => {
			id          => 707,
			name        => 'location',
			config_file => 'astats.config',
			value       => '/opt/trafficserver/etc/trafficserver',
		},
	},
	'477-path' => {
		new   => 'Parameter',
		using => {
			id          => 708,
			name        => 'path',
			config_file => 'astats.config',
			value       => '_astats',
		},
	},
	'478-CONFIG-proxy.config.cache.http.compatibility.4-2-0-fixup' => {
		new   => 'Parameter',
		using => {
			id          => 709,
			name        => 'CONFIG proxy.config.cache.http.compatibility.4-2-0-fixup',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'479-location' => {
		new   => 'Parameter',
		using => {
			id          => 710,
			name        => 'location',
			config_file => 'hdr_rw_cdl-c2.config',
			value       => '/opt/trafficserver/etc/trafficserver',
		},
	},
	'480-location' => {
		new   => 'Parameter',
		using => {
			id          => 711,
			name        => 'location',
			config_file => 'hdr_rw_movies-c1.config',
			value       => '/opt/trafficserver/etc/trafficserver',
		},
	},
	'481-location' => {
		new   => 'Parameter',
		using => {
			id          => 715,
			name        => 'location',
			config_file => 'hdr_rw_images-c1.config',
			value       => '/opt/trafficserver/etc/trafficserver',
		},
	},
	'482-algorithm' => {
		new   => 'Parameter',
		using => {
			id          => 716,
			name        => 'algorithm',
			config_file => 'parent.config',
			value       => 'consistent_hash',
		},
	},
	'483-DsStats' => {
		new   => 'Parameter',
		using => {
			id          => 717,
			name        => 'DsStats',
			config_file => 'redis.config',
			value       => 'kbps',
		},
	},
	'484-DsStats' => {
		new   => 'Parameter',
		using => {
			id          => 718,
			name        => 'DsStats',
			config_file => 'redis.config',
			value       => 'tps_2xx',
		},
	},
	'485-DsStats' => {
		new   => 'Parameter',
		using => {
			id          => 719,
			name        => 'DsStats',
			config_file => 'redis.config',
			value       => 'status_3xx',
		},
	},
	'486-DsStats' => {
		new   => 'Parameter',
		using => {
			id          => 720,
			name        => 'DsStats',
			config_file => 'redis.config',
			value       => 'status_4xx',
		},
	},
	'487-DsStats' => {
		new   => 'Parameter',
		using => {
			id          => 721,
			name        => 'DsStats',
			config_file => 'redis.config',
			value       => 'status_5xx',
		},
	},
	'488-CacheStats' => {
		new   => 'Parameter',
		using => {
			id          => 722,
			name        => 'CacheStats',
			config_file => 'redis.config',
			value       => 'bandwidth',
		},
	},
	'489-CacheStats' => {
		new   => 'Parameter',
		using => {
			id          => 723,
			name        => 'CacheStats',
			config_file => 'redis.config',
			value       => 'maxKbps',
		},
	},
	'490-CacheStats' => {
		new   => 'Parameter',
		using => {
			id          => 724,
			name        => 'CacheStats',
			config_file => 'redis.config',
			value       => 'ats.proxy.process.http.current_client_connections',
		},
	},
	'491-maxRevalDurationDays' => {
		new   => 'Parameter',
		using => {
			id          => 725,
			name        => 'maxRevalDurationDays',
			config_file => 'regex_revalidate.config',
			value       => '90',
		},
	},
	'492-tm.instance_name' => {
		new   => 'Parameter',
		using => {
			id          => 726,
			name        => 'tm.instance_name',
			config_file => 'global',
			value       => 'Kabletown CDN',
		},
	},
	'493-CONFIG-proxy.config.cache.ram_cache_cutoff' => {
		new   => 'Parameter',
		using => {
			id          => 727,
			name        => 'CONFIG proxy.config.cache.ram_cache_cutoff',
			config_file => 'records.config',
			value       => '268435456',
		},
	},
	'494-CONFIG-proxy.config.cache.ram_cache_cutoff' => {
		new   => 'Parameter',
		using => {
			id          => 728,
			name        => 'CONFIG proxy.config.cache.ram_cache_cutoff',
			config_file => 'records.config',
			value       => 'INT 268435456',
		},
	},
	'496-health.threshold.availableBandwidthInKbps' => {
		new   => 'Parameter',
		using => {
			id          => 731,
			name        => 'health.threshold.availableBandwidthInKbps',
			config_file => 'rascal.properties',
			value       => '1062500',
		},
	},
	'497-health.threshold.availableBandwidthInKbps' => {
		new   => 'Parameter',
		using => {
			id          => 732,
			name        => 'health.threshold.availableBandwidthInKbps',
			config_file => 'rascal.properties',
			value       => '1062500',
		},
	},
	'498-health.threshold.availableBandwidthInKbps' => {
		new   => 'Parameter',
		using => {
			id          => 733,
			name        => 'health.threshold.availableBandwidthInKbps',
			config_file => 'rascal.properties',
			value       => '1062500',
		},
	},
	'499-health.threshold.availableBandwidthInKbps' => {
		new   => 'Parameter',
		using => {
			id          => 734,
			name        => 'health.threshold.availableBandwidthInKbps',
			config_file => 'rascal.properties',
			value       => '>11500000',
		},
	},
	'500-DsStats' => {
		new   => 'Parameter',
		using => {
			id          => 735,
			name        => 'DsStats',
			config_file => 'redis.config',
			value       => 'tps_3xx',
		},
	},
	'501-DsStats' => {
		new   => 'Parameter',
		using => {
			id          => 736,
			name        => 'DsStats',
			config_file => 'redis.config',
			value       => 'tps_4xx',
		},
	},
	'502-DsStats' => {
		new   => 'Parameter',
		using => {
			id          => 737,
			name        => 'DsStats',
			config_file => 'redis.config',
			value       => 'tps_5xx',
		},
	},
	'503-DsStats' => {
		new   => 'Parameter',
		using => {
			id          => 738,
			name        => 'DsStats',
			config_file => 'redis.config',
			value       => 'tps_total',
		},
	},

	'558-ttl_max_hours' => {
		new   => 'Parameter',
		using => {
			id          => 793,
			name        => 'ttl_max_hours',
			config_file => 'regex_revalidate.config',
			value       => '672',
		},
	},
	'559-ttl_min_hours' => {
		new   => 'Parameter',
		using => {
			id          => 794,
			name        => 'ttl_min_hours',
			config_file => 'regex_revalidate.config',
			value       => '48',
		},
	},
	'560-snapshot_dir' => {
		new   => 'Parameter',
		using => {
			id          => 795,
			name        => 'snapshot_dir',
			config_file => 'regex_revalidate.config',
			value       => 'public/Trafficserver-Snapshots/',
		},
	},
	'561-CONFIG-proxy.config.ssl.server.cipher_suite' => {
		new   => 'Parameter',
		using => {
			id          => 796,
			name        => 'CONFIG proxy.config.ssl.server.cipher_suite',
			config_file => 'records.config',
			value =>
				'STRING ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES128-SHA256:ECDHE-RSA-AES256-SHA384:AES128-GCM-SHA256:AES256-GCM-SHA384:ECDHE-RSA-RC4-SHA:ECDHE-RSA-AES128-SHA:ECDHE-RSA-AES256-SHA:RC4-SHA:RC4-MD5:AES128-SHA:AES256-SHA:DES-CBC3-SHA!SRP:!DSS:!PSK:!aNULL:!eNULL:!SSLv2',
		},
	},
	'562-CONFIG-proxy.config.ssl.server.honor_cipher_order' => {
		new   => 'Parameter',
		using => {
			id          => 797,
			name        => 'CONFIG proxy.config.ssl.server.honor_cipher_order',
			config_file => 'records.config',
			value       => 'INT 1',
		},
	},
	'563-CONFIG-proxy.config.ssl.server.cert.path' => {
		new   => 'Parameter',
		using => {
			id          => 798,
			name        => 'CONFIG proxy.config.ssl.server.cert.path',
			config_file => 'records.config',
			value       => 'STRING etc/trafficserver/ssl',
		},
	},
	'564-CONFIG-proxy.config.http.server_ports' => {
		new   => 'Parameter',
		using => {
			id          => 799,
			name        => 'CONFIG proxy.config.http.server_ports',
			config_file => 'records.config',
			value       => 'STRING 80 80:ipv6 443:ssl 443:ipv6:ssl',
		},
	},
	'565-CONFIG-proxy.config.ssl.server.private_key.path' => {
		new   => 'Parameter',
		using => {
			id          => 800,
			name        => 'CONFIG proxy.config.ssl.server.private_key.path',
			config_file => 'records.config',
			value       => 'STRING etc/trafficserver/ssl',
		},
	},
	'566-CONFIG-proxy.config.ssl.client.CA.cert.path' => {
		new   => 'Parameter',
		using => {
			id          => 801,
			name        => 'CONFIG proxy.config.ssl.client.CA.cert.path',
			config_file => 'records.config',
			value       => 'STRING etc/trafficserver/ssl',
		},
	},
	'567-CONFIG-proxy.config.ssl.client.private_key.path' => {
		new   => 'Parameter',
		using => {
			id          => 802,
			name        => 'CONFIG proxy.config.ssl.client.private_key.path',
			config_file => 'records.config',
			value       => 'STRING etc/trafficserver/ssl',
		},
	},
	'568-CONFIG-proxy.config.ssl.client.cert.path' => {
		new   => 'Parameter',
		using => {
			id          => 803,
			name        => 'CONFIG proxy.config.ssl.client.cert.path',
			config_file => 'records.config',
			value       => 'STRING etc/trafficserver/ssl',
		},
	},
	'569-CONFIG-proxy.config.ssl.CA.cert.path' => {
		new   => 'Parameter',
		using => {
			id          => 804,
			name        => 'CONFIG proxy.config.ssl.CA.cert.path',
			config_file => 'records.config',
			value       => 'STRING etc/trafficserver/ssl',
		},
	},
	'570-CONFIG-proxy.config.ssl.SSLv3' => {
		new   => 'Parameter',
		using => {
			id          => 805,
			name        => 'CONFIG proxy.config.ssl.SSLv3',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'571-CONFIG-proxy.config.ssl.server.multicert.filename' => {
		new   => 'Parameter',
		using => {
			id          => 808,
			name        => 'CONFIG proxy.config.ssl.server.multicert.filename',
			config_file => 'records.config',
			value       => 'STRING ssl_multicert.config',
		},
	},
	'572-location' => {
		new   => 'Parameter',
		using => {
			id          => 816,
			name        => 'location',
			config_file => 'hdr_rw_games-c1.config',
			value       => '/opt/trafficserver/etc/trafficserver',
		},
	},
	'817-proxy' => {
		new   => 'Parameter',
		using => {
			id          => 817,
			name        => 'traffic_mon_fwd_proxy',
			config_file => 'global',
			value       => 'http://proxy.kabletown.net:81',
		},
	},
	'818-proxy' => {
		new   => 'Parameter',
		using => {
			id          => 818,
			name        => 'traffic_rtr_fwd_proxy',
			config_file => 'global',
			value       => 'http://proxy.kabletown.net:81',
		},
	},
	'819-weight' => {
		new   => 'Parameter',
		using => {
			id          => 819,
			name        => 'weight',
			config_file => 'parent.config',
			value       => '1.0',
		},
	},
	'820-location' => {
		new   => 'Parameter',
		using => {
			id          => 820,
			name        => 'location',
			config_file => 'hdr_rw_mid_movies-c1.config',
			value       => '/opt/trafficserver/etc/trafficserver',
		},
	},
);

sub name {
	return "Parameter";
}

sub get_definition {
	my ( $self, $name ) = @_;
	return $definition_for{$name};
}

sub all_fixture_names {
	return keys %definition_for;
}

__PACKAGE__->meta->make_immutable;

1;
