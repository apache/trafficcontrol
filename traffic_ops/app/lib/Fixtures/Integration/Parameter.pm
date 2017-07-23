package Fixtures::Integration::Parameter;

# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
# 
#   http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.


# Note - removing the domain_name parameter wreaks all kinds of havoc because of ordering / id problems, so I renamed
# it to NODNAME - JvD

use Moose;
extends 'DBIx::Class::EasyFixture';
use namespace::autoclean;

	my %definition_for = (
	## id => 1
	'0' => {
		new => 'Parameter',
		using => {
			name => 'algorithm',
			value => 'consistent_hash',
			config_file => 'parent.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 2
	'1' => {
		new => 'Parameter',
		using => {
			name => 'allow_ip',
			last_updated => '2015-12-10 15:43:46',
			value => '127.0.0.1,172.39.0.0/16,33.101.99.0/24',
			config_file => 'astats.config',
		},
	},
	## id => 3
	'2' => {
		new => 'Parameter',
		using => {
			name => 'allow_ip6',
			value => '::1,2033:D011:3300::336/64,2033:D011:3300::335/64,2033:D021:3300::333/64,2033:D021:3300::334/64',
			config_file => 'astats.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 4
	'3' => {
		new => 'Parameter',
		using => {
			name => 'api.port',
			last_updated => '2015-12-10 15:43:46',
			value => '8080',
			config_file => 'server.xml',
		},
	},
	## id => 5
	'4' => {
		new => 'Parameter',
		using => {
			name => 'api.port',
			config_file => 'server.xml',
			last_updated => '2015-12-10 15:43:48',
			value => '8081',
		},
	},
	## id => 6
	'5' => {
		new => 'Parameter',
		using => {
			name => 'astats_over_http',
			value => '1.2-3.el6.x86_64',
			config_file => 'package',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 7
	'6' => {
		new => 'Parameter',
		using => {
			name => 'astats_over_http',
			last_updated => '2015-12-10 15:43:46',
			value => '3.2.0-4114.el6.x86_64',
			config_file => 'package',
		},
	},
	## id => 8
	'7' => {
		new => 'Parameter',
		using => {
			name => 'astats_over_http',
			value => '1.1-2.el6.x86_64',
			config_file => 'package',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 9
	'8' => {
		new => 'Parameter',
		using => {
			name => 'astats_over_http.so',
			config_file => 'plugin.config',
			last_updated => '2015-12-10 15:43:48',
			value => '',
		},
	},
	## id => 10
	'9' => {
		new => 'Parameter',
		using => {
			name => 'astats_over_http.so',
			config_file => 'plugin.config',
			last_updated => '2015-12-10 15:43:46',
			value => '_astats 33.101.99.100,172.39.19.39,172.39.19.49,172.39.19.49,172.39.29.49',
		},
	},
	## id => 11
	'10' => {
		new => 'Parameter',
		using => {
			name => 'CacheHealthTimeout',
			last_updated => '2015-12-10 15:43:46',
			value => '70',
			config_file => 'CRConfig.xml',
		},
	},
	## id => 12
	'11' => {
		new => 'Parameter',
		using => {
			name => 'CacheStats',
			config_file => 'traffic_stats.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'bandwidth',
		},
	},
	## id => 13
	'12' => {
		new => 'Parameter',
		using => {
			name => 'CacheStats',
			value => 'maxKbps',
			config_file => 'traffic_stats.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 14
	'13' => {
		new => 'Parameter',
		using => {
			name => 'CacheStats',
			last_updated => '2015-12-10 15:43:46',
			value => 'ats.proxy.process.http.current_client_connections',
			config_file => 'traffic_stats.config',
		},
	},
	## id => 15
	'14' => {
		new => 'Parameter',
		using => {
			name => 'cacheurl',
			config_file => 'package',
			last_updated => '2015-12-10 15:43:46',
			value => '1.0-3.el6.x86_64',
		},
	},
	## id => 16
	'15' => {
		new => 'Parameter',
		using => {
			name => 'cacheurl',
			config_file => 'package',
			last_updated => '2015-12-10 15:43:47',
			value => '3.2.0-5628.el6.x86_64',
		},
	},
	## id => 17
	'16' => {
		new => 'Parameter',
		using => {
			name => 'cacheurl',
			config_file => 'package',
			last_updated => '2015-12-10 15:43:48',
			value => '1.0-1.el6.x86_64',
		},
	},
	## id => 18
	'17' => {
		new => 'Parameter',
		using => {
			name => 'cacheurl.so',
			value => '',
			config_file => 'plugin.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 19
	'18' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.accept_threads',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 1',
			config_file => 'records.config',
		},
	},
	## id => 20
	'19' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.admin.admin_user',
			value => 'STRING admin',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 21
	'20' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.admin.autoconf_port',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 8083',
			config_file => 'records.config',
		},
	},
	## id => 22
	'21' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.admin.number_config_bak',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 3',
		},
	},
	## id => 23
	'22' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.admin.user_id',
			value => 'STRING ats',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 24
	'23' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.alarm.abs_path',
			last_updated => '2015-12-10 15:43:46',
			value => 'STRING NULL',
			config_file => 'records.config',
		},
	},
	## id => 25
	'24' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.alarm.bin',
			value => 'STRING example_alarm_bin.sh',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 26
	'25' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.alarm_email',
			value => 'STRING ats',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 27
	'26' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.allocator.debug_filter',
			value => 'INT 0',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 28
	'27' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.allocator.enable_reclaim',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 0',
		},
	},
	## id => 29
	'28' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.allocator.max_overage',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 3',
		},
	},
	## id => 30
	'29' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.allocator.thread_freelist_size',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 1024',
		},
	},
	## id => 31
	'30' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.body_factory.enable_customizations',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 0',
			config_file => 'records.config',
		},
	},
	## id => 32
	'31' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.body_factory.enable_logging',
			value => 'INT 0',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 33
	'32' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.body_factory.response_suppression_mode',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 0',
			config_file => 'records.config',
		},
	},
	## id => 34
	'33' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.cache.enable_read_while_writer',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 0',
		},
	},
	## id => 35
	'34' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.cache.enable_read_while_writer',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 1',
			config_file => 'records.config',
		},
	},
	## id => 36
	'35' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.cache.hosting_filename',
			last_updated => '2015-12-10 15:43:47',
			value => 'STRING hosting.config',
			config_file => 'records.config',
		},
	},
	## id => 37
	'36' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.cache.http.compatibility.4-2-0-fixup',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 0',
		},
	},
	## id => 38
	'37' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.cache.limits.http.max_alts',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 5',
		},
	},
	## id => 39
	'38' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.cache.max_doc_size',
			value => 'INT 0',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 40
	'39' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.cache.min_average_object_size',
			value => 'INT 131072',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 41
	'40' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.cache.mutex_retry_delay',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 2',
			config_file => 'records.config',
		},
	},
	## id => 42
	'41' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.cache.permit.pinning',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 0',
		},
	},
	## id => 43
	'42' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.cache.ram_cache.algorithm',
			value => 'INT 0',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 44
	'43' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.cache.ram_cache.compress',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 0',
		},
	},
	## id => 45
	'44' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.cache.ram_cache.size',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 21474836480',
		},
	},
	## id => 46
	'45' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.cache.ram_cache.size',
			value => 'INT 16106127360',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 47
	'46' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.cache.ram_cache.size',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 34359738368',
		},
	},
	## id => 48
	'47' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.cache.ram_cache.use_seen_filter',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 0',
		},
	},

	## id => 49
	'48' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.cache.ram_cache_cutoff',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 4194304',
		},
	},
	## id => 50
	'49' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.cache.ram_cache_cutoff',
			value => '268435456',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 51
	'50' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.cache.ram_cache_cutoff',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 268435456',
		},
	},

	## id => 52
	'51' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.cache.target_fragment_size',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 1048576',
		},
	},
	## id => 53
	'52' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.cache.threads_per_disk',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 8',
			config_file => 'records.config',
		},
	},
	## id => 54
	'53' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.cluster.cluster_configuration ',
			last_updated => '2015-12-10 15:43:46',
			value => 'STRING cluster.config',
			config_file => 'records.config',
		},
	},
	## id => 55
	'54' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.cluster.cluster_port',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 8086',
		},
	},
	## id => 56
	'55' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.cluster.ethernet_interface',
			value => 'STRING lo',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 57
	'56' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.cluster.log_bogus_mc_msgs',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 1',
			config_file => 'records.config',
		},
	},
	## id => 58
	'57' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.cluster.mc_group_addr',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'STRING 224.0.1.37',
		},
	},
	## id => 59
	'58' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.cluster.mc_ttl',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 1',
			config_file => 'records.config',
		},
	},
	## id => 60
	'59' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.cluster.mcport',
			value => 'INT 8089',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 61
	'60' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.cluster.rsport',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 8088',
		},
	},
	## id => 62
	'61' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.config_dir',
			last_updated => '2015-12-10 15:43:47',
			value => 'STRING etc/trafficserver',
			config_file => 'records.config',
		},
	},
	## id => 63
	'62' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.core_limit',
			value => 'INT -1',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 64
	'63' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.diags.debug.enabled',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 0',
		},
	},
	## id => 65
	'64' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.diags.debug.tags',
			last_updated => '2015-12-10 15:43:47',
			value => 'STRING http.*|dns.*',
			config_file => 'records.config',
		},
	},
	## id => 66
	'65' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.diags.show_location',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 0',
		},
	},
	## id => 67
	'66' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.dns.lookup_timeout',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 2',
			config_file => 'records.config',
		},
	},
	## id => 68
	'67' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.dns.max_dns_in_flight',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 2048',
			config_file => 'records.config',
		},
	},
	## id => 69
	'68' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.dns.nameservers',
			value => 'STRING NULL',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 70
	'69' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.dns.resolv_conf',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'STRING /etc/resolv.conf',
		},
	},
	## id => 71
	'70' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.dns.round_robin_nameservers',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 0',
		},
	},
	## id => 72
	'71' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.dns.search_default_domains',
			value => 'INT 0',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 73
	'72' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.dns.splitDNS.enabled',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 0',
			config_file => 'records.config',
		},
	},
	## id => 74
	'73' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.dns.url_expansions',
			value => 'STRING NULL',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 75
	'74' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.dns.validate_query_name',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 0',
		},
	},
	## id => 76
	'75' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.dump_mem_info_frequency',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 0',
			config_file => 'records.config',
		},
	},
	## id => 77
	'76' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.env_prep',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'STRING example_prep.sh',
		},
	},
	## id => 78
	'77' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.exec_thread.affinity',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 1',
			config_file => 'records.config',
		},
	},
	## id => 79
	'78' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.exec_thread.autoconfig',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 0',
			config_file => 'records.config',
		},
	},
	## id => 80
	'79' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.exec_thread.autoconfig',
			value => 'INT 1',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 81
	'80' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.exec_thread.autoconfig.scale',
			last_updated => '2015-12-10 15:43:47',
			value => 'FLOAT 1.5',
			config_file => 'records.config',
		},
	},
	## id => 82
	'81' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.exec_thread.limit',
			value => 'INT 32',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 83
	'82' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.exec_thread.limit',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 2',
		},
	},
	## id => 84
	'83' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.header.parse.no_host_url_redirect',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'STRING NULL',
		},
	},
	## id => 85
	'84' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.hostdb.serve_stale_for',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 6',
			config_file => 'records.config',
		},
	},
	## id => 86
	'85' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.hostdb.size',
			value => 'INT 120000',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 87
	'86' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.hostdb.storage_size',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 33554432',
			config_file => 'records.config',
		},
	},
	## id => 88
	'87' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.hostdb.strict_round_robin',
			value => 'INT 0',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 89
	'88' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.hostdb.timeout',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 1440',
			config_file => 'records.config',
		},
	},
	## id => 90
	'89' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.hostdb.ttl_mode',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 1',
		},
	},
	## id => 91
	'90' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.hostdb.ttl_mode',
			value => 'INT 0',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 92
	'91' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.accept_no_activity_timeout',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 120',
		},
	},
	## id => 93
	'92' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.anonymize_insert_client_ip',
			value => 'INT 1',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 94
	'93' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.anonymize_other_header_list',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'STRING NULL',
		},
	},
	## id => 95
	'94' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.anonymize_remove_client_ip',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 0',
			config_file => 'records.config',
		},
	},
	## id => 96
	'95' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.anonymize_remove_cookie',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 0',
		},
	},
	## id => 97
	'96' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.anonymize_remove_from',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 0',
		},
	},
	## id => 98
	'97' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.anonymize_remove_referer',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 0',
		},
	},
	## id => 99
	'98' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.anonymize_remove_user_agent',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 0',
		},
	},
	## id => 100
	'99' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.background_fill_active_timeout',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 0',
		},
	},
	## id => 101
	'100' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.background_fill_active_timeout',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 60',
		},
	},
	## id => 102
	'101' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.background_fill_completed_threshold',
			last_updated => '2015-12-10 15:43:46',
			value => 'FLOAT 0.0',
			config_file => 'records.config',
		},
	},
	## id => 103
	'102' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.background_fill_completed_threshold',
			last_updated => '2015-12-10 15:43:47',
			value => 'FLOAT 0.5',
			config_file => 'records.config',
		},
	},
	## id => 104
	'103' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.cache.allow_empty_doc',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 0',
			config_file => 'records.config',
		},
	},
	## id => 105
	'104' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.cache.cache_responses_to_cookies',
	 		config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 1',
		},
	},
	## id => 106
	'105' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.cache.cache_urls_that_look_dynamic',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 1',
		},
	},
	## id => 107
	'106' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.cache.enable_default_vary_headers',
			value => 'INT 0',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 108
	'107' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.cache.fuzz.probability',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'FLOAT 0.005',
		},
	},
	## id => 109
	'108' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.cache.fuzz.time',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 240',
			config_file => 'records.config',
		},
	},
	## id => 110
	'109' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.cache.heuristic_lm_factor',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'FLOAT 0.10',
		},
	},
	## id => 111
	'110' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.cache.heuristic_max_lifetime',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 86400',
		},
	},
	## id => 112
	'111' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.cache.heuristic_min_lifetime',
			value => 'INT 3600',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 113
	'112' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.cache.http',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 1',
			config_file => 'records.config',
		},
	},
	## id => 114
	'113' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.cache.ignore_accept_encoding_mismatch',
			value => 'INT 1',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 115
	'114' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.cache.ignore_authentication',
			value => 'INT 0',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 116
	'115' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.cache.ignore_client_cc_max_age',
			value => 'INT 0',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 117
	'116' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.cache.ignore_client_no_cache',
			value => 'INT 1',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 118
	'117' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.cache.ignore_server_no_cache',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 0',
		},
	},
	## id => 119
	'118' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.cache.ims_on_client_no_cache',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 1',
		},
	},
	## id => 120
	'119' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.cache.max_stale_age',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 604800',
		},
	},
	## id => 121
	'120' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.cache.range.lookup',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 1',
		},
	},
	## id => 122
	'121' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.cache.required_headers',
			value => 'INT 0',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 123
	'122' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.cache.vary_default_images',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'STRING NULL',
		},
	},
	## id => 124
	'123' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.cache.vary_default_other',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'STRING NULL',
		},
	},
	## id => 125
	'124' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.cache.vary_default_text',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'STRING NULL',
		},
	},
	## id => 126
	'125' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.cache.when_to_add_no_cache_to_msie_requests',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT -1',
		},
	},
	## id => 127
	'126' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.cache.when_to_revalidate',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 0',
		},
	},
	## id => 128
	'127' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.chunking_enabled',
			value => 'INT 1',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 129
	'128' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.congestion_control.enabled',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 0',
		},
	},
	## id => 130
	'129' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.connect_attempts_max_retries',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 6',
			config_file => 'records.config',
		},
	},
	## id => 131
	'130' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.connect_attempts_max_retries_dead_server',
			value => 'INT 3',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 132
	'131' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.connect_attempts_rr_retries',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 3',
			config_file => 'records.config',
		},
	},
	## id => 133
	'132' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.connect_attempts_timeout',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 30',
		},
	},
	## id => 134
	'133' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.connect_ports',
			value => 'STRING 443 563',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 135
	'134' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.down_server.abort_threshold',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 10',
		},
	},
	## id => 136
	'135' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.down_server.cache_time',
			value => 'INT 300',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 137
	'136' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.enable_http_stats',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 1',
		},
	},
	## id => 138
	'137' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.enable_url_expandomatic',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 0',
		},
	},
	## id => 139
	'138' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.forward.proxy_auth_to_parent',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 0',
		},
	},
	## id => 140
	'139' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.insert_age_in_response',
			value => 'INT 0',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 141
	'140' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.insert_age_in_response',
			value => 'INT 1',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 142
	'141' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.insert_request_via_str',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 1',
		},
	},
	## id => 143
	'142' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.insert_response_via_str',
			value => 'INT 3',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 144
	'143' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.insert_squid_x_forwarded_for',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 1',
			config_file => 'records.config',
		},
	},
	## id => 145
	'144' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.keep_alive_enabled_in',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 1',
		},
	},
	## id => 146
	'145' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.keep_alive_enabled_out',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 1',
			config_file => 'records.config',
		},
	},
	## id => 147
	'146' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.keep_alive_no_activity_timeout_in',
			value => 'INT 115',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 148
	'147' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.keep_alive_no_activity_timeout_out',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 120',
		},
	},
	## id => 149
		'148' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.negative_caching_enabled',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 0',
		},
	},
	## id => 150
		'149' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.negative_caching_lifetime',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 1800',
			config_file => 'records.config',
		},
	},
	## id => 151
		'150' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.no_dns_just_forward_to_parent',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 0',
			config_file => 'records.config',
		},
	},
	## id => 152
		'151' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.normalize_ae_gzip',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 0',
		},
	},
	## id => 153
	'152' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.origin_server_pipeline',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 1',
		},
	},
	## id => 154
	'153' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.parent_proxy.connect_attempts_timeout',
			value => 'INT 30',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
		## id => 155
	'154' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.parent_proxy.fail_threshold',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 10',
		},
	},
	## id => 156
	'155' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.parent_proxy.file',
			last_updated => '2015-12-10 15:43:46',
			value => 'STRING parent.config',
			config_file => 'records.config',
		},
	},
	## id => 157
	'156' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.parent_proxy.per_parent_connect_attempts',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 2',
		},
	},
	## id => 158
	'157' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.parent_proxy.retry_time',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 300',
		},
	},
	## id => 159
	'158' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.parent_proxy.total_connect_attempts',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 4',
		},
	},
	## id => 160
	'159' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.parent_proxy_routing_enable',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 0',
		},
	},
	## id => 161
	'160' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.parent_proxy_routing_enable',
			value => 'INT 1',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 162
	'161' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.post_connect_attempts_timeout',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 1800',
			config_file => 'records.config',
		},
	},
	## id => 163
	'162' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.push_method_enabled',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 0',
		},
	},
	## id => 164
	'163' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.referer_default_redirect',
			last_updated => '2015-12-10 15:43:46',
			value => 'STRING http://www.example.com/',
			config_file => 'records.config',
		},
	},
	## id => 165
	'164' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.referer_filter',
			value => 'INT 0',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 166
	'165' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.referer_format_redirect',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 0',
			config_file => 'records.config',
		},
	},
	## id => 167
	'166' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.response_server_enabled',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 1',
		},
	},
	## id => 168
	'167' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.send_http11_requests',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 1',
			config_file => 'records.config',
		},
	},
	## id => 169
	'168' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.server_ports',
			last_updated => '2015-12-10 15:43:48',
			value => 'STRING 80 80:ipv6',
			config_file => 'records.config',
		},
	},
	## id => 170
	'169' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.server_ports',
			value => 'STRING 80 80:ipv6 443:ssl 443:ipv6:ssl',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 171
	'170' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.share_server_sessions',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 2',
		},
	},
	## id => 172
	'171' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.slow.log.threshold',
			value => 'INT 0',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 173
	'172' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.transaction_active_timeout_in',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 900',
			config_file => 'records.config',
		},
	},
	## id => 174
	'173' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.transaction_active_timeout_out',
			value => 'INT 0',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 175
	'174' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.transaction_no_activity_timeout_in',
			value => 'INT 30',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 176
	'175' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.transaction_no_activity_timeout_out',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 30',
			config_file => 'records.config',
		},
	},
	## id => 177
	'176' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.uncacheable_requests_bypass_parent',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 1',
		},
	},
	## id => 178
	'177' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.http.user_agent_pipeline',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 8',
			config_file => 'records.config',
		},
	},
	## id => 179
	'178' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.icp.enabled',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 0',
		},
	},
	## id => 180
	'179' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.icp.icp_interface',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'STRING NULL',
		},
	},
	## id => 181
	'180' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.icp.icp_port',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 3130',
			config_file => 'records.config',
		},
	},
	## id => 182
	'181' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.icp.multicast_enabled',
			value => 'INT 0',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 183
	'182' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.icp.query_timeout',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 2',
		},
	},
	## id => 184
	'183' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.auto_delete_rolled_files',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 1',
		},
	},
	## id => 185
	'184' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.collation_host',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'STRING NULL',
		},
	},
	## id => 186
	'185' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.collation_host_tagged',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 0',
		},
	},
	## id => 187
	'186' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.collation_port',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 8085',
		},
	},
	## id => 188
	'187' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.collation_retry_sec',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 5',
			config_file => 'records.config',
		},
	},
	## id => 189
	'188' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.collation_secret',
			last_updated => '2015-12-10 15:43:48',
			value => 'STRING foobar',
			config_file => 'records.config',
		},
	},
	## id => 190
	'189' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.common_log_enabled',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 0',
		},
	},
	## id => 191
	'190' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.common_log_header',
			value => 'STRING NULL',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 192
	'191' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.common_log_is_ascii',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 1',
		},
	},
	## id => 193
	'192' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.common_log_name',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'STRING common',
		},
	},
	## id => 194
	'193' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.custom_logs_enabled',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 1',
		},
	},
	## id => 195
	'194' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.extended2_log_enabled',
			value => 'INT 0',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 196
	'195' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.extended2_log_enabled',
			value => 'INT 1',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 197
	'196' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.extended2_log_header',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'STRING NULL',
		},
	},
	## id => 198
	'197' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.extended2_log_is_ascii',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 1',
			config_file => 'records.config',
		},
	},
	## id => 199
	'198' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.extended2_log_name',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'STRING extended2',
		},
	},
	## id => 200
	'199' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.extended_log_enabled',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 0',
		},
	},
	## id => 201
	'200' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.extended_log_header',
			last_updated => '2015-12-10 15:43:47',
			value => 'STRING NULL',
			config_file => 'records.config',
		},
	},
	## id => 202
	'201' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.extended_log_is_ascii',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 0',
			config_file => 'records.config',
		},
	},
	## id => 203
	'202' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.extended_log_name',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'STRING extended',
		},
	},
	## id => 204
	'203' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.hostname',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'STRING localhost',
		},
	},
	## id => 205
	'204' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.logfile_dir',
			last_updated => '2015-12-10 15:43:48',
			value => 'STRING var/log/trafficserver',
			config_file => 'records.config',
		},
	},
	## id => 206
	'205' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.logfile_perm',
			value => 'STRING rw-r--r--',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 207
	'206' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.logging_enabled',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 3',
			config_file => 'records.config',
		},
	},
	## id => 208
	'207' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.max_secs_per_buffer',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 5',
		},
	},
	## id => 209
	'208' => {
			new => 'Parameter',
			using => {
				name => 'CONFIG proxy.config.log.max_space_mb_for_logs',
				last_updated => '2015-12-10 15:43:47',
				value => 'INT 25000',
				config_file => 'records.config',
			},
	},
	## id => 210
	'209' => {
			new => 'Parameter',
			using => {
				name => 'CONFIG proxy.config.log.max_space_mb_for_orphan_logs',
				value => 'INT 25',
				config_file => 'records.config',
				last_updated => '2015-12-10 15:43:46',
			},
	},
	## id => 211
	'210' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.max_space_mb_headroom',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 1000',
			config_file => 'records.config',
		},
	},
	## id => 212
	'211' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.rolling_enabled',
			value => 'INT 1',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 213
	'212' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.rolling_interval_sec',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 86400',
		},
	},
	## id => 214
	'213' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.rolling_offset_hr',
			value => 'INT 0',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 215
	'214' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.rolling_size_mb',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 10',
			config_file => 'records.config',
		},
	},## id => 216
	'215' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.sampling_frequency',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 1',
			config_file => 'records.config',
		},
	},
	## id => 217
	'216' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.separate_host_logs',
			value => 'INT 0',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 218
	'217' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.separate_icp_logs',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 0',
			config_file => 'records.config',
		},
	},
	## id => 219
	'218' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.squid_log_enabled',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 0',
		},
	},
	## id => 220
	'219' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.squid_log_header',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'STRING NULL',
		},
	},
	## id => 221
	'220' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.squid_log_is_ascii',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 0',
		},
	},
	## id => 222
	'221' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.squid_log_name',
			value => 'STRING squid',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 223
	'222' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.log.xml_config_file',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'STRING logs_xml.config',
		},
	},
	## id => 224
	'223' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.mlock_enabled',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 2',
		},
	},
	## id => 225
	'224' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.net.connections_throttle',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 500000',
			config_file => 'records.config',
		},
	},
	## id => 226
	'225' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.net.defer_accept',
			value => 'INT 45',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 227
	'226' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.net.sock_recv_buffer_size_in',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 0',
		},
	},
	## id => 228
	'227' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.net.sock_recv_buffer_size_out',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 0',
			config_file => 'records.config',
		},
	},
	## id => 229
	'228' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.net.sock_send_buffer_size_in',
			value => 'INT 262144',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 230
	'229' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.net.sock_send_buffer_size_out',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 0',
			config_file => 'records.config',
		},
	},
	## id => 231
	'230' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.output.logfile',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'STRING traffic.out',
		},
	},
	## id => 232
	'231' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.process_manager.mgmt_port',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 8084',
			config_file => 'records.config',
		},
	},
	## id => 233
	'232' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.proxy_binary_opts',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'STRING -M',
		},
	},
	## id => 234
	'233' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.proxy_name',
			value => 'STRING __HOSTNAME__',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 235
	'234' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.reverse_proxy.enabled',
			value => 'INT 1',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 236
	'235' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.snapshot_dir',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'STRING snapshots',
		},
	},
	## id => 237
	'236' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.CA.cert.filename',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'STRING NULL',
		},
	},
	## id => 238
	'237' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.CA.cert.path',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'STRING etc/trafficserver',
		},
	},
	## id => 239
	'238' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.CA.cert.path',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'STRING etc/trafficserver/ssl',
		},
	},
	## id => 240
	'239' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.client.CA.cert.filename',
			last_updated => '2015-12-10 15:43:47',
			value => 'STRING NULL',
			config_file => 'records.config',
		},
	},
	## id => 241
	'240' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.client.CA.cert.path',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'STRING etc/trafficserver',
		},
	},
	## id => 242
	'241' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.client.CA.cert.path',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'STRING etc/trafficserver/ssl',
		},
	},
	## id => 243
	'242' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.client.cert.filename',
			value => 'STRING NULL',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 244
	'243' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.client.cert.path',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'STRING etc/trafficserver',
		},
	},
	## id => 245
	'244' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.client.cert.path',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'STRING etc/trafficserver/ssl',
		},
	},
	## id => 246
	'245' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.client.certification_level',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 0',
		},
	},
	## id => 247
	'246' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.client.private_key.filename',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'STRING NULL',
		},
	},
	## id => 248
	'247' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.client.private_key.path',
			value => 'STRING etc/trafficserver',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 249
	'248' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.client.private_key.path',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'STRING etc/trafficserver/ssl',
		},
	},
	## id => 250
	'249' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.client.verify.server',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 0',
			config_file => 'records.config',
		},
	},
	## id => 251
	'250' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.compression',
			value => 'INT 1',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 252
	'251' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.number.threads',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 0',
		},
	},
	## id => 253
	'252' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.server.cert.path',
			last_updated => '2015-12-10 15:43:48',
			value => 'STRING etc/trafficserver',
			config_file => 'records.config',
		},
	},
	## id => 254
	'253' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.server.cert.path',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'STRING etc/trafficserver/ssl',
		},
	},
	## id => 255
	'254' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.server.cert_chain.filename',
			value => 'STRING NULL',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 256
	'255' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.server.cipher_suite',
			value => 'STRING RC4-SHA:AES128-SHA:DES-CBC3-SHA:AES256-SHA:ALL:!aNULL:!EXP:!LOW:!MD5:!SSLV2:!NULL',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 257
	'256' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.server.cipher_suite',
			value => 'STRING ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES128-SHA256:ECDHE-RSA-AES256-SHA384:AES128-GCM-SHA256:AES256-GCM-SHA384:ECDHE-RSA-RC4-SHA:ECDHE-RSA-AES128-SHA:ECDHE-RSA-AES256-SHA:RC4-SHA:RC4-MD5:AES128-SHA:AES256-SHA:DES-CBC3-SHA!SRP:!DSS:!PSK:!aNULL:!eNULL:!SSLv2',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 258
	'257' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.server.honor_cipher_order',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 0',
		},
	},
	## id => 259
	'258' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.server.honor_cipher_order',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 1',
		},
	},
	## id => 260
	'259' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.server.multicert.filename',
			last_updated => '2015-12-10 15:43:47',
			value => 'STRING ssl_multicert.config',
			config_file => 'records.config',
		},
	},
	## id => 261
	'260' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.server.private_key.path',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'STRING etc/trafficserver',
		},
	},
	## id => 262
	'261' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.server.private_key.path',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'STRING etc/trafficserver/ssl',
		},
	},
	## id => 263
	'262' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.SSLv2',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 0',
			config_file => 'records.config',
		},
	},
	## id => 264
	'263' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.SSLv3',
			value => 'INT 1',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 265
	'264' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.SSLv3',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 0',
			config_file => 'records.config',
		},
	},
	## id => 266
	'265' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.ssl.TLSv1',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 1',
		},
	},
	## id => 267
	'266' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.stack_dump_enabled',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 0',
		},
	},
	## id => 268
	'267' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.stack_dump_enabled',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 1',
			config_file => 'records.config',
		},
	},
	## id => 269
	'268' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.syslog_facility',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'STRING LOG_DAEMON',
		},
	},
	## id => 270
	'269' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.system.mmap_max',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 2097152',
		},
	},
	## id => 271
	'270' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.task_threads',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 2',
			config_file => 'records.config',
		},
	},
	## id => 272
	'271' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.temp_dir',
			last_updated => '2015-12-10 15:43:47',
			value => 'STRING /tmp',
			config_file => 'records.config',
		},
	},
	## id => 273
	'272' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.update.concurrent_updates',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 100',
		},
	},
	## id => 274
	'273' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.update.enabled',
			value => 'INT 0',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 275
	'274' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.update.force',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 0',
		},
	},
	## id => 276
	'275' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.update.retry_count',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 10',
		},
	},
	## id => 277
	'276' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.update.retry_interval',
			value => 'INT 2',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 278
	'277' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.url_remap.default_to_server_pac',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'INT 0',
		},
	},
	## id => 279
	'278' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.url_remap.default_to_server_pac_port',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT -1',
		},
	},
	## id => 280
	'279' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.url_remap.filename',
			value => 'STRING remap.config',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 281
	'280' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.url_remap.pristine_host_hdr',
			last_updated => '2015-12-10 15:43:48',
			value => 'INT 0',
			config_file => 'records.config',
		},
	},
	## id => 282
	'281' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.url_remap.remap_required',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 1',
		},
	},
	## id => 283
	'282' => {
		new => 'Parameter',
		using => {
			name => 'CONFIG proxy.config.url_remap.remap_required',
			value => 'INT 0',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 284
	'283' => {
		new => 'Parameter',
		using => {
			name => 'coveragezone.polling.interval',
			last_updated => '2015-12-10 15:43:48',
			value => '86400000',
			config_file => 'CRConfig.json',
		},
	},
	## id => 285
	'284' => {
		new => 'Parameter',
		using => {
			name => 'coveragezone.polling.url',
			config_file => 'CRConfig.json',
			last_updated => '2015-12-10 15:43:46',
			value => 'http://staging.cdnlab.kabletown.net/ipcdn/CZF/current/kabletown_ipcdn_czf-current.json',
		},
	},
	## id => 286
	'285' => {
		new => 'Parameter',
		using => {
			name => 'CoverageZoneJsonURL',
			value => 'http://staging.cdnlab.kabletown.net/ipcdn/CZF/current/kabletown_ipcdn_czf-current.json',
			config_file => 'CRConfig.xml',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 287
	'286' => {
		new => 'Parameter',
		using => {
			name => 'CoverageZoneMapRefreshPeriodHours',
			value => '24',
			config_file => 'CRConfig.xml',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 288
	'287' => {
		new => 'Parameter',
		using => {
			name => 'CoverageZoneMapURL',
			last_updated => '2015-12-10 15:43:48',
			value => 'http://aux.cdnlab.kabletown.net/logs/production/reports/czf/current/kabletown_cdn_czf.xml',
			config_file => 'CRConfig.xml',
		},
	},
	## id => 289
	'288' => {
		new => 'Parameter',
		using => {
			name => 'cron_syncds',
			last_updated => '2015-12-10 15:43:47',
			value => '*/15 * * * * /opt/ort/ipcdn_install_ort.pl syncds error &amp;gt; /tmp/ort/syncds.log 2&amp;gt;&amp;amp;1',
			config_file => 'crontab_root',
		},
	},
	## id => 290
	'289' => {
		new => 'Parameter',
		using => {
			name => 'Disk_Volume',
			config_file => 'storage.config',
			last_updated => '2015-12-10 15:43:47',
			value => '1',
		},
	},
	## id => 291
	'290' => {
		new => 'Parameter',
		using => {
			name => 'NODNAME',
			config_file => 'CRConfig.xml',
			last_updated => '2015-12-10 15:43:46',
			value => 'cdn1.kabletown.net',
		},
	},
	## id => 292
	'291' => {
		new => 'Parameter',
		using => {
			name => 'NODNAME',
			config_file => 'CRConfig.xml',
			last_updated => '2015-12-10 15:43:46',
			value => 'cdn2.kabletown.net',
		},
	},
	## id => 293
	'292' => {
		new => 'Parameter',
		using => {
			name => 'NODNAME',
			last_updated => '2015-12-10 15:43:47',
			value => 'cdn1.kabletown.net',
			config_file => 'CRConfig.json',
		},
	},
	## id => 294
	'293' => {
		new => 'Parameter',
		using => {
			name => 'NODNAME',
			value => 'cdn2.kabletown.net',
			config_file => 'CRConfig.json',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 295
	'294' => {
		new => 'Parameter',
		using => {
			name => 'Drive_Letters',
			last_updated => '2015-12-10 15:43:47',
			value => '0,1,2,3,4,5,6',
			config_file => 'storage.config',
		},
	},
	## id => 296
	'295' => {
		new => 'Parameter',
		using => {
			name => 'Drive_Letters',
			last_updated => '2015-12-10 15:43:48',
			value => 'b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q,r,s,t,u,v,w,x,y',
			config_file => 'storage.config',
		},
	},
	## id => 297
	'296' => {
		new => 'Parameter',
		using => {
			name => 'Drive_Letters',
			config_file => 'storage.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'b,c,d,e,f,g,h,i,j,k,l,m,n,o',
		},
	},
	## id => 298
	'297' => {
		new => 'Parameter',
		using => {
			name => 'Drive_Prefix',
			value => '/dev/ram',
			config_file => 'storage.config',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 299
	'298' => {
		new => 'Parameter',
		using => {
			name => 'Drive_Prefix',
			config_file => 'storage.config',
			last_updated => '2015-12-10 15:43:46',
			value => '/dev/sd',
		},
	},
	## id => 300
	'299' => {
		new => 'Parameter',
		using => {
			name => 'dscp_remap',
			value => '1.0-3.el6.x86_64',
			config_file => 'package',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 301
	'300' => {
		new => 'Parameter',
		using => {
			name => 'dscp_remap',
			last_updated => '2015-12-10 15:43:46',
			value => '3.2.0-4613.el6.x86_64',
			config_file => 'package',
		},
	},
	## id => 302
	'301' => {
		new => 'Parameter',
		using => {
			name => 'dscp_remap',
			value => '1.0-1.el6.x86_64',
			config_file => 'package',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 303
	'302' => {
		new => 'Parameter',
		using => {
			name => 'DsStats',
			config_file => 'traffic_stats.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'tps_3xx',
		},
	},
	## id => 304
	'303' => {
		new => 'Parameter',
		using => {
			name => 'DsStats',
			value => 'tps_4xx',
			config_file => 'traffic_stats.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 305
	'304' => {
		new => 'Parameter',
		using => {
			name => 'DsStats',
			config_file => 'traffic_stats.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'tps_5xx',
		},
	},
	## id => 306
	'305' => {
		new => 'Parameter',
		using => {
			name => 'DsStats',
			config_file => 'traffic_stats.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'tps_total',
		},
	},
	## id => 307
	'306' => {
		new => 'Parameter',
		using => {
			name => 'DsStats',
			config_file => 'traffic_stats.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'kbps',
		},
	},
	## id => 308
	'307' => {
		new => 'Parameter',
		using => {
			name => 'DsStats',
			config_file => 'traffic_stats.config',
			last_updated => '2015-12-10 15:43:48',
			value => 'tps_2xx',
		},
	},
	## id => 309
	'308' => {
		new => 'Parameter',
		using => {
			name => 'DsStats',
			value => 'status_3xx',
			config_file => 'traffic_stats.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 310
	'309' => {
		new => 'Parameter',
		using => {
			name => 'DsStats',
			config_file => 'traffic_stats.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'status_4xx',
		},
	},
	## id => 311
	'310' => {
		new => 'Parameter',
		using => {
			name => 'DsStats',
			config_file => 'traffic_stats.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'status_5xx',
		},
	},
	## id => 312
	'311' => {
		new => 'Parameter',
		using => {
			name => 'error_url',
			config_file => 'url_sig_cdl-c2.config',
			last_updated => '2015-12-10 15:43:46',
			value => '403',
		},
	},
	## id => 313
	'312' => {
		new => 'Parameter',
		using => {
			name => 'geolocation.polling.interval',
			value => '86400000',
			config_file => 'CRConfig.json',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 314
	'313' => {
		new => 'Parameter',
		using => {
			name => 'geolocation.polling.url',
			last_updated => '2015-12-10 15:43:46',
			value => 'https://tm.kabletown.net/MaxMind/GeoLiteCity.dat.gz',
			config_file => 'CRConfig.json',
		},
	},
	## id => 315
	'314' => {
		new => 'Parameter',
		using => {
			name => 'geolocation6.polling.url',
	 		config_file => 'CRConfig.json',
			last_updated => '2015-12-10 15:43:47',
			value => 'https://tm.kabletown.net/MaxMind/GeoLiteCityv6.dat.gz',
		},
	},
	## id => 316
	'315' => {
		new => 'Parameter',
		using => {
			name => 'geolocation6.polling.interval',
	 		config_file => 'CRConfig.json',
			last_updated => '2015-12-10 15:43:48',
			value => '86400000',
		},
	},
	## id => 317
	'316' => {
		new => 'Parameter',
		using => {
			name => 'GeolocationURL',
			value => 'http://aux.cdnlab.kabletown.net:8080/GeoLiteCity.dat.gz',
			config_file => 'CRConfig.xml',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 318
	'317' => {
		new => 'Parameter',
		using => {
			name => 'hack.ttl',
			config_file => 'rascal-config.txt',
			last_updated => '2015-12-10 15:43:47',
			value => '30',
		},
	},
	## id => 319
	'318' => {
		new => 'Parameter',
		using => {
			name => 'header_rewrite',
			config_file => 'package',
			last_updated => '2015-12-10 15:43:46',
			value => '4.0.2-1.el6.x86_64',
		},
	},
	## id => 320
	'319' => {
		new => 'Parameter',
		using => {
			name => 'health.connection.timeout',
			last_updated => '2015-12-10 15:43:48',
			value => '2000',
			config_file => 'rascal.properties',
		},
	},
	## id => 321
	'320' => {
		new => 'Parameter',
		using => {
			name => 'health.event-count',
			config_file => 'rascal-config.txt',
			last_updated => '2015-12-10 15:43:46',
			value => '200',
		},
	},
	## id => 322
	'321' => {
		new => 'Parameter',
		using => {
			name => 'health.polling.interval',
			config_file => 'rascal-config.txt',
			last_updated => '2015-12-10 15:43:46',
			value => '8000',
		},
	},
	## id => 323
	'322' => {
		new => 'Parameter',
		using => {
			name => 'health.polling.url',
			value => 'http://${hostname}/_astats?application=&inf.name=${interface_name}',
			config_file => 'rascal.properties',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 324
	'323' => {
		new => 'Parameter',
		using => {
			name => 'health.threadPool',
			config_file => 'rascal-config.txt',
			last_updated => '2015-12-10 15:43:46',
			value => '4',
		},
	},
	## id => 325
	'324' => {
		new => 'Parameter',
		using => {
			name => 'health.threshold.availableBandwidthInKbps',
			value => '1062500',
			config_file => 'rascal.properties',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 326
	'325' => {
		new => 'Parameter',
		using => {
			name => 'health.threshold.availableBandwidthInKbps',
			last_updated => '2015-12-10 15:43:47',
			value => '1062501',
			config_file => 'rascal.properties',
		},
	},
	## id => 327
	'326' => {
		new => 'Parameter',
		using => {
			name => 'health.threshold.availableBandwidthInKbps',
			last_updated => '2015-12-10 15:43:47',
			value => '1062502',
			config_file => 'rascal.properties',
		},
	},
	## id => 328
	'327' => {
		new => 'Parameter',
		using => {
			name => 'health.threshold.availableBandwidthInKbps',
			value => '>11500000',
			config_file => 'rascal.properties',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 329
	'328' => {
		new => 'Parameter',
		using => {
			name => 'health.threshold.availableBandwidthInKbps',
			config_file => 'rascal.properties',
			last_updated => '2015-12-10 15:43:46',
			value => '>1750000',
		},
	},
	## id => 330
	'329' => {
		new => 'Parameter',
		using => {
			name => 'health.threshold.loadavg',
			last_updated => '2015-12-10 15:43:48',
			value => '25.0',
			config_file => 'rascal.properties',
		},
	},
	## id => 331
	'330' => {
		new => 'Parameter',
		using => {
			name => 'health.threshold.queryTime',
			last_updated => '2015-12-10 15:43:47',
			value => '1000',
			config_file => 'rascal.properties',
		},
	},
	## id => 332
	'331' => {
		new => 'Parameter',
		using => {
			name => 'health.timepad',
			config_file => 'rascal-config.txt',
			last_updated => '2015-12-10 15:43:47',
			value => '30',
		},
	},
	## id => 333
	'332' => {
		new => 'Parameter',
		using => {
			name => 'history.count',
			last_updated => '2015-12-10 15:43:46',
			value => '30',
			config_file => 'rascal.properties',
		},
	},
	## id => 334
	'333' => {
		new => 'Parameter',
		using => {
			name => 'key0',
			last_updated => '2015-12-10 15:43:47',
			value => 'HOOJ3Ghq1x4gChp3iQkqVTcPlOj8UCi3',
			config_file => 'url_sig_cdl-c2.config',
		},
	},
	## id => 335
	'334' => {
		new => 'Parameter',
		using => {
			name => 'key1',
	 		config_file => 'url_sig_cdl-c2.config',
			last_updated => '2015-12-10 15:43:47',
			value => '_9LZYkRnfCS0rCBF7fTQzM9Scwlp2FhO',
		},
	},
	## id => 336
	'335' => {
		new => 'Parameter',
		using => {
			name => 'key2',
			last_updated => '2015-12-10 15:43:47',
			value => 'AFpkxfc4oTiyFSqtY6_ohjt3V80aAIxS',
			config_file => 'url_sig_cdl-c2.config',
		},
	},
	## id => 337
	'336' => {
		new => 'Parameter',
		using => {
			name => 'key3',
			value => 'AL9kzs_SXaRZjPWH8G5e2m4ByTTzkzlc',
			config_file => 'url_sig_cdl-c2.config',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 338
	'337' => {
		new => 'Parameter',
		using => {
			name => 'key4',
	 		config_file => 'url_sig_cdl-c2.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'poP3n3szbD1U4vx1xQXV65BvkVgWzfN8',
		},
	},
	## id => 339
	'338' => {
		new => 'Parameter',
		using => {
			name => 'key5',
			last_updated => '2015-12-10 15:43:47',
			value => '1ir32ng4C4w137p5oq72kd2wqmIZUrya',
			config_file => 'url_sig_cdl-c2.config',
		},
	},
	## id => 340
	'339' => {
		new => 'Parameter',
		using => {
			name => 'key6',
			last_updated => '2015-12-10 15:43:47',
			value => 'B1qLptn2T1b_iXeTCWDcVuYvANtH139f',
			config_file => 'url_sig_cdl-c2.config',
		},
	},
	## id => 341
	'340' => {
		new => 'Parameter',
		using => {
			name => 'key7',
			value => 'PiCV_5OODMzBbsNFMWsBxcQ8v1sK0TYE',
			config_file => 'url_sig_cdl-c2.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 342
	'341' => {
		new => 'Parameter',
		using => {
			name => 'key8',
	 		config_file => 'url_sig_cdl-c2.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'Ggpv6DqXDvt2s1CETPBpNKwaLk4fTM9l',
		},
	},
	## id => 343
	'342' => {
		new => 'Parameter',
		using => {
			name => 'key9',
			value => 'qPlVT_s6kL37aqb6hipDm4Bt55S72mI7',
			config_file => 'url_sig_cdl-c2.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 344
	'343' => {
		new => 'Parameter',
		using => {
			name => 'key10',
			last_updated => '2015-12-10 15:43:46',
			value => 'BsI5A9EmWrobIS1FeuOs1z9fm2t2WSBe',
			config_file => 'url_sig_cdl-c2.config',
		},
	},
	## id => 345
	'344' => {
		new => 'Parameter',
		using => {
			name => 'key11',
	 		config_file => 'url_sig_cdl-c2.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'A54y66NCIj897GjS4yA9RrsSPtCUnQXP',
		},
	},
	## id => 346
	'345' => {
		new => 'Parameter',
		using => {
			name => 'key12',
			last_updated => '2015-12-10 15:43:46',
			value => '2jZH0NDPSJttIr4c2KP510f47EKqTQAu',
			config_file => 'url_sig_cdl-c2.config',
		},
	},
	## id => 347
	'346' => {
		new => 'Parameter',
		using => {
			name => 'key13',
			value => 'XduT2FBjBmmVID5JRB5LEf9oR5QDtBgC',
			config_file => 'url_sig_cdl-c2.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 348
	'347' => {
		new => 'Parameter',
		using => {
			name => 'key14',
	 		config_file => 'url_sig_cdl-c2.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'D9nH0SvK_0kP5w8QNd1UFJ28ulFkFKPn',
		},
	},
	## id => 349
	'348' => {
		new => 'Parameter',
		using => {
			name => 'key15',
	 		config_file => 'url_sig_cdl-c2.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'udKXWYNwbXXweaaLzaKDGl57OixnIIcm',
		},
	},
	## id => 350
	'349' => {
		new => 'Parameter',
		using => {
			name => 'LOCAL proxy.config.cache.interim.storage',
			config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'STRING NULL',
		},
	},
	## id => 351
	'350' => {
		new => 'Parameter',
		using => {
			name => 'LOCAL proxy.local.cluster.type',
	 		config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 3',
		},
	},
	## id => 352
	'351' => {
		new => 'Parameter',
		using => {
			name => 'LOCAL proxy.local.log.collation_mode',
	 		config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => 'INT 0',
		},
	},
	## id => 353
	'352' => {
		new => 'Parameter',
		using => {
			name => 'location',
	 		config_file => 'cache.config',
			last_updated => '2015-12-10 15:43:46',
			value => '/opt/trafficserver/etc/trafficserver/',
		},
	},
	## id => 354
	'353' => {
		new => 'Parameter',
		using => {
			name => 'location',
			last_updated => '2015-12-10 15:43:48',
			value => '/opt/trafficserver/etc/trafficserver/',
			config_file => 'hosting.config',
		},
	},
	## id => 355
	'354' => {
		new => 'Parameter',
		using => {
			name => 'location',
			value => '/opt/trafficserver/etc/trafficserver/',
			config_file => 'parent.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 356
	'355' => {
		new => 'Parameter',
		using => {
			name => 'location',
	 		config_file => 'plugin.config',
			last_updated => '2015-12-10 15:43:47',
			value => '/opt/trafficserver/etc/trafficserver/',
		},
	},
	## id => 357
	'356' => {
		new => 'Parameter',
		using => {
			name => 'location',
	 		config_file => 'records.config',
			last_updated => '2015-12-10 15:43:47',
			value => '/opt/trafficserver/etc/trafficserver/',
		},
	},
	## id => 358
	'357' => {
		new => 'Parameter',
		using => {
			name => 'location',
	 		config_file => 'remap.config',
			last_updated => '2015-12-10 15:43:47',
			value => '/opt/trafficserver/etc/trafficserver/',
		},
	},
	## id => 359
	'358' => {
		new => 'Parameter',
		using => {
			name => 'location',
			value => '/opt/trafficserver/etc/trafficserver/',
			config_file => 'storage.config',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 360
	'359' => {
		new => 'Parameter',
		using => {
			name => 'location',
			last_updated => '2015-12-10 15:43:47',
			value => '/opt/trafficserver/etc/trafficserver/',
			config_file => 'volume.config',
		},
	},
	## id => 361
	'360' => {
		new => 'Parameter',
		using => {
			name => 'location',
	 		config_file => '50-ats.rules',
			last_updated => '2015-12-10 15:43:47',
			value => '/etc/udev/rules.d/',
		},
	},
	## id => 362
	'361' => {
		new => 'Parameter',
		using => {
			name => 'location',
			last_updated => '2015-12-10 15:43:46',
			value => 'XMPP CRConfig node',
			config_file => 'CRConfig.xml',
		},
	},
	## id => 363
	'362' => {
		new => 'Parameter',
		using => {
			name => 'location',
			value => '/etc/kabletown/zones/<zonename>.info',
			config_file => 'dns.zone',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 364
	'363' => {
		new => 'Parameter',
		using => {
			name => 'location',
	 		config_file => 'http-log4j.properties',
			last_updated => '2015-12-10 15:43:46',
			value => '/etc/kabletown',
		},
	},
	## id => 365
	'364' => {
		new => 'Parameter',
		using => {
			name => 'location',
	 		config_file => 'dns-log4j.properties',
			last_updated => '2015-12-10 15:43:48',
			value => '/etc/kabletown',
		},
	},
	## id => 366
	'365' => {
		new => 'Parameter',
		using => {
			name => 'location',
			last_updated => '2015-12-10 15:43:47',
			value => '/etc/kabletown',
			config_file => 'geolocation.properties',
		},
	},
	## id => 367
	'366' => {
		new => 'Parameter',
		using => {
			name => 'location',
			last_updated => '2015-12-10 15:43:47',
			value => '/opt/trafficserver/etc/trafficserver',
			config_file => 'ip_allow.config',
		},
	},
	## id => 368
	'367' => {
		new => 'Parameter',
		using => {
			name => 'location',
			last_updated => '2015-12-10 15:43:46',
			value => '/var/spool/cron',
			config_file => 'crontab_root',
		},
	},
	## id => 369
	'368' => {
		new => 'Parameter',
		using => {
			name => 'location',
			value => 'XMPP CRConfigOTT node',
			config_file => 'CRConfig.xml',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 370
	'369' => {
		new => 'Parameter',
		using => {
			name => 'location',
			last_updated => '2015-12-10 15:43:46',
			value => '/opt/trafficserver/etc/trafficserver',
			config_file => 'hdr_rw_mid_movies-c1.config',
		},
	},
	## id => 371
	'370' => {
		new => 'Parameter',
		using => {
			name => 'location',
			config_file => 'cacheurl.config',
			last_updated => '2015-12-10 15:43:48',
			value => '/opt/trafficserver/etc/trafficserver/',
		},
	},
	## id => 372
	'371' => {
		new => 'Parameter',
		using => {
			name => 'location',
			config_file => 'logs_xml.config',
			last_updated => '2015-12-10 15:43:47',
			value => '/opt/trafficserver/etc/trafficserver',
		},
	},
	## id => 373
	'372' => {
		new => 'Parameter',
		using => {
			name => 'location',
			last_updated => '2015-12-10 15:43:47',
			value => '/opt/traffic_monitor/conf',
			config_file => 'rascal-config.txt',
		},
	},
	## id => 374
	'373' => {
		new => 'Parameter',
		using => {
			name => 'location',
			config_file => '12M_facts',
			last_updated => '2015-12-10 15:43:47',
			value => '/opt/ort',
		},
	},
	## id => 375
	'374' => {
		new => 'Parameter',
		using => {
			name => 'location',
			value => '/opt/trafficserver/etc/trafficserver',
			config_file => 'url_sig_cdl-c2.config',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 376
	'375' => {
		new => 'Parameter',
		using => {
			name => 'location',
			value => '/opt/trafficserver/etc/trafficserver',
			config_file => 'regex_revalidate.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 377
	'376' => {
		new => 'Parameter',
		using => {
			name => 'location',
			value => '/opt/trafficserver/etc/trafficserver',
			config_file => 'drop_qstring.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 378
	'377' => {
		new => 'Parameter',
		using => {
			name => 'location',
			config_file => 'astats.config',
			last_updated => '2015-12-10 15:43:46',
			value => '/opt/trafficserver/etc/trafficserver',
		},
	},
	## id => 379
	'378' => {
		new => 'Parameter',
		using => {
			name => 'location',
			value => '/opt/trafficserver/etc/trafficserver',
			config_file => 'hdr_rw_cdl-c2.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 380
	'379' => {
		new => 'Parameter',
		using => {
			name => 'location',
			config_file => 'hdr_rw_movies-c1.config',
			last_updated => '2015-12-10 15:43:46',
			value => '/opt/trafficserver/etc/trafficserver',
		},
	},
	## id => 381
	'380' => {
		new => 'Parameter',
		using => {
			name => 'location',
			config_file => 'hdr_rw_images-c1.config',
			last_updated => '2015-12-10 15:43:47',
			value => '/opt/trafficserver/etc/trafficserver',
		},
	},
	## id => 382
	'381' => {
		new => 'Parameter',
		using => {
			name => 'location',
			value => '/opt/trafficserver/etc/trafficserver',
			config_file => 'hdr_rw_games-c1.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 383
	'382' => {
		new => 'Parameter',
		using => {
			name => 'LogFormat.Format',
			last_updated => '2015-12-10 15:43:48',
			value => '%<chi> %<caun> [%<cqtq>] "%<cqtx>" %<pssc> %<pscl> %<sssc> %<sscl> %<cqbl> %<pqbl> %<cqhl> %<pshl> %<ttms> %<pqhl> %<sshl> %<phr> %<cfsc> %<pfsc> %<crc> "%<{User-Agent}cqh>"',
			config_file => 'logs_xml.config',
		},
	},
	## id => 384
	'383' => {
		new => 'Parameter',
		using => {
			name => 'LogFormat.Format',
			last_updated => '2015-12-10 15:43:47',
			value => '%<cqtq> chi=%<chi> phn=%<phn> shn=%<shn> url=%<cquuc> cqhm=%<cqhm> cqhv=%<cqhv> pssc=%<pssc> ttms=%<ttms> b=%<pscl> sssc=%<sssc> sscl=%<sscl> cfsc=%<cfsc> pfsc=%<pfsc> crc=%<crc> phr=%<phr> uas="%<{User-Agent}cqh>"',
			config_file => 'logs_xml.config',
		},
	},
	## id => 385
	'384' => {
		new => 'Parameter',
		using => {
			name => 'LogFormat.Name',
			last_updated => '2015-12-10 15:43:48',
			value => 'custom_ats_2',
			config_file => 'logs_xml.config',
		},
	},
	## id => 386
	'385' => {
		new => 'Parameter',
		using => {
			name => 'LogFormat.Name',
			last_updated => '2015-12-10 15:43:48',
			value => 'custom_ats_3',
			config_file => 'logs_xml.config',
		},
	},
	## id => 387
	'386' => {
		new => 'Parameter',
		using => {
			name => 'LogObject.Filename',
			last_updated => '2015-12-10 15:43:48',
			value => 'custom_ats_2',
			config_file => 'logs_xml.config',
		},
	},
	## id => 388
	'387' => {
		new => 'Parameter',
		using => {
			name => 'LogObject.Filename',
			config_file => 'logs_xml.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'custom_ats_3',
		},
	},
	## id => 389
	'388' => {
		new => 'Parameter',
		using => {
			name => 'LogObject.Format',
			value => 'custom_ats_2',
			config_file => 'logs_xml.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 390
	'389' => {
		new => 'Parameter',
		using => {
			name => 'LogObject.Format',
			value => 'custom_ats_3',
			config_file => 'logs_xml.config',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 391
	'390' => {
		new => 'Parameter',
		using => {
			name => 'LogObject.RollingEnabled',
			last_updated => '2015-12-10 15:43:47',
			value => '3',
			config_file => 'logs_xml.config',
		},
	},
	## id => 392
	'391' => {
		new => 'Parameter',
		using => {
			name => 'LogObject.RollingIntervalSec',
			config_file => 'logs_xml.config',
			last_updated => '2015-12-10 15:43:47',
			value => '86400',
		},
	},
	## id => 393
	'392' => {
		new => 'Parameter',
		using => {
			name => 'LogObject.RollingOffsetHr',
			config_file => 'logs_xml.config',
			last_updated => '2015-12-10 15:43:47',
			value => '11',
		},
	},
	## id => 394
	'393' => {
		new => 'Parameter',
		using => {
			name => 'LogObject.RollingSizeMb',
			config_file => 'logs_xml.config',
			last_updated => '2015-12-10 15:43:47',
			value => '1024',
		},
	},
	## id => 395
	'394' => {
		new => 'Parameter',
		using => {
			name => 'maxRevalDurationDays',
	 		config_file => 'regex_revalidate.config',
			last_updated => '2015-12-10 15:43:47',
			value => '90',
		},
	},
	## id => 396
	'395' => {
		new => 'Parameter',
		using => {
			name => 'maxRevalDurationDays',
			value => '60',
			config_file => 'regex_revalidate.config',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 397
	'396' => {
		new => 'Parameter',
		using => {
			name => 'monitor:///opt/tomcat/logs/access.log',
	 		config_file => 'inputs.conf',
			last_updated => '2015-12-10 15:43:46',
			value => 'index=index_odol_test;sourcetype=access_ccr',
		},
	},
	## id => 398
	'397' => {
		new => 'Parameter',
		using => {
			name => 'path',
			config_file => 'astats.config',
			last_updated => '2015-12-10 15:43:47',
			value => '_astats',
		},
	},
	## id => 399
	'398' => {
		new => 'Parameter',
		using => {
			name => 'purge_allow_ip',
	 		config_file => 'ip_allow.config',
			last_updated => '2015-12-10 15:43:47',
			value => '33.101.99.100',
		},
	},
	## id => 400
	'399' => {
		new => 'Parameter',
		using => {
			name => 'RAM_Drive_Letters',
			last_updated => '2015-12-10 15:43:47',
			value => '0,1,2,3,4,5,6,7',
			config_file => 'storage.config',
		},
	},
	## id => 401
	'400' => {
		new => 'Parameter',
		using => {
			name => 'RAM_Drive_Prefix',
			value => '/dev/ram',
			config_file => 'storage.config',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 402
	'401' => {
		new => 'Parameter',
		using => {
			name => 'RAM_Volume',
	 		config_file => 'storage.config',
			last_updated => '2015-12-10 15:43:46',
			value => '2',
		},
	},
	## id => 403
	'402' => {
		new => 'Parameter',
		using => {
			name => 'ramdisk_size',
			config_file => 'grub.conf',
			last_updated => '2015-12-10 15:43:47',
			value => 'ramdisk_size=16777216',
		},
	},
	## id => 404
	'403' => {
		new => 'Parameter',
		using => {
			name => 'record_types',
			config_file => 'astats.config',
			last_updated => '2015-12-10 15:43:48',
			value => '144',
		},
	},
	## id => 405
	'404' => {
		new => 'Parameter',
		using => {
			name => 'regex_revalidate',
			config_file => 'package',
			last_updated => '2015-12-10 15:43:46',
			value => '1.0-1.el6.x86_64',
		},
	},
	## id => 406
	'405' => {
		new => 'Parameter',
		using => {
			name => 'regex_revalidate',
			value => '3.2.0-5695.el6.x86_64',
			config_file => 'package',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 407
	'406' => {
		new => 'Parameter',
		using => {
			name => 'regex_revalidate',
			value => '1.0-4.el6.x86_64',
			config_file => 'package',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 408
	'407' => {
		new => 'Parameter',
		using => {
			name => 'regex_revalidate.so',
			last_updated => '2015-12-10 15:43:46',
			value => '--config regex_revalidate.config',
			config_file => 'plugin.config',
		},
	},
	## id => 409
	'408' => {
		new => 'Parameter',
		using => {
			name => 'remap_stats',
			last_updated => '2015-12-10 15:43:47',
			value => '1.0-1.el6.x86_64',
			config_file => 'package',
		},
	},
	## id => 410
	'409' => {
		new => 'Parameter',
		using => {
			name => 'remap_stats',
			config_file => 'package',
			last_updated => '2015-12-10 15:43:48',
			value => '3.2.0-2.el6.x86_64',
		},
	},
	## id => 411
	'410' => {
		new => 'Parameter',
		using => {
			name => 'remap_stats',
			config_file => 'package',
			last_updated => '2015-12-10 15:43:46',
			value => '1.0-4.el6.x86_64',
		},
	},
	## id => 412
	'411' => {
		new => 'Parameter',
		using => {
			name => 'remap_stats',
			config_file => 'package',
			last_updated => '2015-12-10 15:43:47',
			value => '3.2.0-3.el6.x86_64',
		},
	},
	## id => 413
	'412' => {
		new => 'Parameter',
		using => {
			name => 'remap_stats.so',
			config_file => 'plugin.config',
			last_updated => '2015-12-10 15:43:46',
			value => 'unknown',
		},
	},
	## id => 414
	'413' => {
		new => 'Parameter',
		using => {
			name => 'remap_stats.so',
			last_updated => '2015-12-10 15:43:47',
			value => '',
			config_file => 'plugin.config',
		},
	},
	## id => 416
	'415' => {
			new => 'Parameter',
			using => {
				name => 'tld.soa.admin',
				value => 'traffic_ops',
				config_file => 'CRConfig.json',
				last_updated => '2015-12-10 15:43:47',
			},
		},
	## id => 417
	'416' => {
		new => 'Parameter',
		using => {
			name => 'tld.soa.expire',
			last_updated => '2015-12-10 15:43:48',
			value => '604800',
			config_file => 'CRConfig.json',
		},
	},
	## id => 418
	'417' => {
		new => 'Parameter',
		using => {
			name => 'tld.soa.minimum',
			last_updated => '2015-12-10 15:43:46',
			value => '86400',
			config_file => 'CRConfig.json',
		},
	},
	## id => 419
	'418' => {
		new => 'Parameter',
		using => {
			name => 'tld.soa.refresh',
			config_file => 'CRConfig.json',
			last_updated => '2015-12-10 15:43:48',
			value => '28800',
		},
	},
	## id => 420
	'419' => {
		new => 'Parameter',
		using => {
			name => 'tld.soa.retry',
			config_file => 'CRConfig.json',
			last_updated => '2015-12-10 15:43:47',
			value => '7200',
		},
	},
	## id => 421
	'420' => {
		new => 'Parameter',
		using => {
			name => 'tld.ttls.A',
			config_file => 'CRConfig.json',
			last_updated => '2015-12-10 15:43:48',
			value => '3600',
		},
	},
	## id => 422
	'421' => {
		new => 'Parameter',
		using => {
			name => 'tld.ttls.AAAA',
			value => '3600',
			config_file => 'CRConfig.json',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 423
	'422' => {
		new => 'Parameter',
		using => {
			name => 'tld.ttls.NS',
			value => '3600',
			config_file => 'CRConfig.json',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 424
	'423' => {
		new => 'Parameter',
		using => {
			name => 'tld.ttls.NS',
			last_updated => '2015-12-10 15:43:46',
			value => '7200',
			config_file => 'CRConfig.json',
		},
	},
	## id => 425
	'424' => {
		new => 'Parameter',
		using => {
			name => 'tld.ttls.SOA',
			config_file => 'CRConfig.json',
			last_updated => '2015-12-10 15:43:48',
			value => '172800',
		},
	},
	## id => 426
	'425' => {
		new => 'Parameter',
		using => {
			name => 'tld.ttls.SOA',
			config_file => 'CRConfig.json',
			last_updated => '2015-12-10 15:43:47',
			value => '86400',
		},
	},
	## id => 427
	'426' => {
		new => 'Parameter',
		using => {
			name => 'tm.crConfig.polling.url',
			last_updated => '2015-12-10 15:43:47',
			value => 'https://${tmHostname}/CRConfig-Snapshots/${cdnName}/CRConfig.xml',
			config_file => 'rascal-config.txt',
		},
	},
	## id => 428
	'427' => {
		new => 'Parameter',
		using => {
			name => 'tm.dataServer.polling.url',
			value => 'https://${tmHostname}/dataserver/orderby/id',
			config_file => 'rascal-config.txt',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 429
	'428' => {
		new => 'Parameter',
		using => {
			name => 'tm.healthParams.polling.url',
			config_file => 'rascal-config.txt',
			last_updated => '2015-12-10 15:43:46',
			value => 'https://${tmHostname}/health/${cdnName}',
		},
	},
	## id => 430
	'429' => {
		new => 'Parameter',
		using => {
			name => 'tm.infourl',
			config_file => 'global',
			last_updated => '2015-12-10 15:43:46',
			value => 'http://staging.cdnlab.kabletown.net/tm/info',
		},
	},
	## id => 431
	'430' => {
		new => 'Parameter',
		using => {
			name => 'tm.instance_name',
			last_updated => '2015-12-10 15:43:47',
			value => 'Kabletown CDN',
			config_file => 'global',
		},
	},
	## id => 432
	'431' => {
		new => 'Parameter',
		using => {
			name => 'tm.logourl',
			last_updated => '2015-12-10 15:43:48',
			value => '/images/tc_logo.png',
			config_file => 'global',
		},
	},
	## id => 433
	'432' => {
		new => 'Parameter',
		using => {
			name => 'tm.polling.interval',
			value => '60000',
			config_file => 'rascal-config.txt',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 434
	'433' => {
		new => 'Parameter',
		using => {
			name => 'tm.toolname',
			config_file => 'global',
			last_updated => '2015-12-10 15:43:47',
			value => 'Traffic Ops',
		},
	},
	## id => 435
	'434' => {
		new => 'Parameter',
		using => {
			name => 'tm.url',
			last_updated => '2015-12-10 15:43:46',
			value => 'https://tm.kabletown.net/',
			config_file => 'global',
		},
	},
	## id => 436
	'435' => {
		new => 'Parameter',
		using => {
			name => 'trafficserver',
			last_updated => '2015-12-10 15:43:47',
			value => '0:off	1:off	2:on	3:on	4:on	5:on	6:off',
			config_file => 'chkconfig',
		},
	},
	## id => 437
	'436' => {
		new => 'Parameter',
		using => {
			name => 'trafficserver',
			value => '5.3.2-761.el6.x86_64',
			config_file => 'package',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 438
	'437' => {
		new => 'Parameter',
		using => {
			name => 'trafficserver',
			value => '6.2.1-45.el7.x86_64',
			config_file => 'package',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 439
	'438' => {
		new => 'Parameter',
		using => {
			name => 'trafficserver',
			config_file => 'package',
			last_updated => '2015-12-10 15:43:46',
			value => '6.2.1-48.el7.x86_64',
		},
	},
	## id => 440
	'439' => {
		new => 'Parameter',
		using => {
			name => 'traffic_mon_fwd_proxy',
			last_updated => '2015-12-10 15:43:47',
			value => 'http://proxy.kabletown.net:81',
			config_file => 'global',
		},
	},
	## id => 441
	'440' => {
		new => 'Parameter',
		using => {
			name => 'traffic_rtr_fwd_proxy',
			last_updated => '2015-12-10 15:43:48',
			value => 'http://proxy.kabletown.net:81',
			config_file => 'global',
		},
	},
	## id => 442
	'441' => {
		new => 'Parameter',
		using => {
			name => 'url_sig',
			last_updated => '2015-12-10 15:43:47',
			value => '1.0-3.el6.x86_64',
			config_file => 'package',
		},
	},
	## id => 443
	'442' => {
		new => 'Parameter',
		using => {
			name => 'url_sign',
			last_updated => '2015-12-10 15:43:46',
			value => '1.0-1.el6.x86_64',
			config_file => 'package',
		},
	},
	## id => 444
	'443' => {
		new => 'Parameter',
		using => {
			name => 'url_sign',
			value => '3.2.0-4130.el6.x86_64',
			config_file => 'package',
			last_updated => '2015-12-10 15:43:47',
		},
	},
	## id => 445
	'444' => {
		new => 'Parameter',
		using => {
			name => 'weight',
			value => '1.0',
			config_file => 'parent.config',
			last_updated => '2015-12-10 15:43:46',
		},
	},
	## id => 446
	'445' => {
		new   => 'Parameter',
		using => {
			name        => 'use_tenancy',
			config_file => 'global',
			value       => '1',
		},
	},
);

sub name {
		return "Parameter";
}

sub get_definition {
		my ( $self,
			$name ) = @_;
		return $definition_for{$name};
}

sub all_fixture_names {
	# sort by db name to guarantee insertion order
	return (sort { lc($definition_for{$a}{using}{name}) cmp lc($definition_for{$b}{using}{name}) } keys %definition_for);
}

__PACKAGE__->meta->make_immutable;
1;
