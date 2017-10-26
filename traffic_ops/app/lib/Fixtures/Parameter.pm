package Fixtures::Parameter;
#
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
	health_threadhold_loadavg => {
		new   => 'Parameter',
		using => {
			id          => 4,
			name        => 'health.threshold.loadavg',
			value       => '25.0',
			config_file => 'rascal.properties',
		},
	},
	health_threadhold_available_bandwidth_in_kbps => {
		new   => 'Parameter',
		using => {
			id          => 5,
			name        => 'health.threshold.availableBandwidthInKbps',
			value       => '>1750000',
			config_file => 'rascal.properties',
		},
	},
	history_count => {
		new   => 'Parameter',
		using => {
			id          => 6,
			name        => 'history.count',
			value       => '30',
			config_file => 'rascal.properties',
		},
	},
	'key0' => {
		new   => 'Parameter',
		using => {
			id          => 7,
			name        => 'key0',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'HOOJ3Ghq1x4gChp3iQkqVTcPlOj8UCi3',
		},
	},
	'key1' => {
		new   => 'Parameter',
		using => {
			id          => 8,
			name        => 'key1',
			config_file => 'url_sig_cdl-c2.config',
			value       => '_9LZYkRnfCS0rCBF7fTQzM9Scwlp2FhO',
		},
	},
	'key2' => {
		new   => 'Parameter',
		using => {
			id          => 9,
			name        => 'key2',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'AFpkxfc4oTiyFSqtY6_ohjt3V80aAIxS',
		},
	},
	'key3' => {
		new   => 'Parameter',
		using => {
			id          => 10,
			name        => 'key3',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'AL9kzs_SXaRZjPWH8G5e2m4ByTTzkzlc',
		},
	},
	'key4' => {
		new   => 'Parameter',
		using => {
			id          => 11,
			name        => 'key4',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'poP3n3szbD1U4vx1xQXV65BvkVgWzfN8',
		},
	},
	'key5' => {
		new   => 'Parameter',
		using => {
			id          => 12,
			name        => 'key5',
			config_file => 'url_sig_cdl-c2.config',
			value       => '1ir32ng4C4w137p5oq72kd2wqmIZUrya',
		},
	},
	'key6' => {
		new   => 'Parameter',
		using => {
			id          => 13,
			name        => 'key6',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'B1qLptn2T1b_iXeTCWDcVuYvANtH139f',
		},
	},
	'key7' => {
		new   => 'Parameter',
		using => {
			id          => 14,
			name        => 'key7',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'PiCV_5OODMzBbsNFMWsBxcQ8v1sK0TYE',
		},
	},
	'key8' => {
		new   => 'Parameter',
		using => {
			id          => 15,
			name        => 'key8',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'Ggpv6DqXDvt2s1CETPBpNKwaLk4fTM9l',
		},
	},
	'key9' => {
		new   => 'Parameter',
		using => {
			id          => 16,
			name        => 'key9',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'qPlVT_s6kL37aqb6hipDm4Bt55S72mI7',
		},
	},
	'key10' => {
		new   => 'Parameter',
		using => {
			id          => 17,
			name        => 'key10',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'BsI5A9EmWrobIS1FeuOs1z9fm2t2WSBe',
		},
	},
	'key11' => {
		new   => 'Parameter',
		using => {
			id          => 18,
			name        => 'key11',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'A54y66NCIj897GjS4yA9RrsSPtCUnQXP',
		},
	},
	'key12' => {
		new   => 'Parameter',
		using => {
			id          => 19,
			name        => 'key12',
			config_file => 'url_sig_cdl-c2.config',
			value       => '2jZH0NDPSJttIr4c2KP510f47EKqTQAu',
		},
	},
	'key13' => {
		new   => 'Parameter',
		using => {
			id          => 20,
			name        => 'key13',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'XduT2FBjBmmVID5JRB5LEf9oR5QDtBgC',
		},
	},
	'key14' => {
		new   => 'Parameter',
		using => {
			id          => 21,
			name        => 'key14',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'D9nH0SvK_0kP5w8QNd1UFJ28ulFkFKPn',
		},
	},
	'key15' => {
		new   => 'Parameter',
		using => {
			id          => 22,
			name        => 'key15',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'udKXWYNwbXXweaaLzaKDGl57OixnIIcm',
		},
	},
	'url_sig_cdl-c2.config_location' => {
		new   => 'Parameter',
		using => {
			id          => 23,
			name        => 'location',
			config_file => 'url_sig_cdl-c2.config',
			value       => '/opt/trafficserver/etc/trafficserver',
		},
	},
	'error_url' => {
		new   => 'Parameter',
		using => {
			id          => 24,
			name        => 'error_url',
			config_file => 'url_sig_cdl-c2.config',
			value       => '403',
		},
	},
	'CONFIG-proxy.config.allocator.debug_filter' => {
		new   => 'Parameter',
		using => {
			id          => 25,
			name        => 'CONFIG proxy.config.allocator.debug_filter',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'CONFIG-proxy.config.allocator.enable_reclaim' => {
		new   => 'Parameter',
		using => {
			id          => 26,
			name        => 'CONFIG proxy.config.allocator.enable_reclaim',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'CONFIG-proxy.config.allocator.max_overage' => {
		new   => 'Parameter',
		using => {
			id          => 27,
			name        => 'CONFIG proxy.config.allocator.max_overage',
			config_file => 'records.config',
			value       => 'INT 3',
		},
	},
	'CONFIG-proxy.config.diags.show_location' => {
		new   => 'Parameter',
		using => {
			id          => 28,
			name        => 'CONFIG proxy.config.diags.show_location',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'CONFIG-proxy.config.http.cache.allow_empty_doc' => {
		new   => 'Parameter',
		using => {
			id          => 29,
			name        => 'CONFIG proxy.config.http.cache.allow_empty_doc',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	'LOCAL-proxy.config.cache.interim.storage' => {
		new   => 'Parameter',
		using => {
			id          => 30,
			name        => 'LOCAL proxy.config.cache.interim.storage',
			config_file => 'records.config',
			value       => 'STRING NULL',
		},
	},
	'CONFIG-proxy.config.http.parent_proxy.file' => {
		new   => 'Parameter',
		using => {
			id          => 31,
			name        => 'CONFIG proxy.config.http.parent_proxy.file',
			config_file => 'records.config',
			value       => 'STRING parent.config',
		},
	},
	'12M_location' => {
		new   => 'Parameter',
		using => {
			id          => 32,
			name        => 'location',
			config_file => '12M_facts',
			value       => '/opt/ort',
		},
	},
	'cacheurl_location' => {
		new   => 'Parameter',
		using => {
			id          => 33,
			name        => 'location',
			config_file => 'cacheurl.config',
			value       => '/opt/trafficserver/etc/trafficserver/',
		},
	},
	'ip_allow_location' => {
		new   => 'Parameter',
		using => {
			id          => 34,
			name        => 'location',
			config_file => 'ip_allow.config',
			value       => '/opt/trafficserver/etc/trafficserver',
		},
	},
	'astats_over_http.so' => {
		new   => 'Parameter',
		using => {
			id          => 35,
			name        => 'astats_over_http.so',
			config_file => 'plugin.config',
			value       => '_astats 33.101.99.100,172.39.19.39,172.39.19.49,172.39.19.49,172.39.29.49',
		},
	},
	'crontab_root_location' => {
		new   => 'Parameter',
		using => {
			id          => 36,
			name        => 'location',
			config_file => 'crontab_root',
			value       => '/var/spool/cron',
		},
	},
	'hdr_rw_cdl-c2.config_location' => {
		new   => 'Parameter',
		using => {
			id          => 37,
			name        => 'location',
			config_file => 'hdr_rw_cdl-c2.config',
			value       => '/opt/trafficserver/etc/trafficserver',
		},
	},
	'50-ats.rules_location' => {
		new   => 'Parameter',
		using => {
			id          => 38,
			name        => 'location',
			config_file => '50-ats.rules',
			value       => '/etc/udev/rules.d/',
		},
	},
	'parent.config_location' => {
		new   => 'Parameter',
		using => {
			id          => 39,
			name        => 'location',
			config_file => 'parent.config',
			value       => '/opt/trafficserver/etc/trafficserver/',
		},
	},
	'remap.config_location' => {
		new   => 'Parameter',
		using => {
			id          => 40,
			name        => 'location',
			config_file => 'remap.config',
			value       => '/opt/trafficserver/etc/trafficserver/',
		},
	},
	'drop_qstring.config_location' => {
		new   => 'Parameter',
		using => {
			id          => 41,
			name        => 'location',
			config_file => 'drop_qstring.config',
			value       => '/opt/trafficserver/etc/trafficserver',
		},
	},
	'LogFormat.Format' => {
		new   => 'Parameter',
		using => {
			id          => 42,
			name        => 'LogFormat.Format',
			config_file => 'logs_xml.config',
			value =>
				'%<cqtq> chi=%<chi> phn=%<phn> shn=%<shn> url=%<cquuc> cqhm=%<cqhm> cqhv=%<cqhv> pssc=%<pssc> ttms=%<ttms> b=%<pscl> sssc=%<sssc> sscl=%<sscl> cfsc=%<cfsc> pfsc=%<pfsc> crc=%<crc> phr=%<phr> uas="%<{User-Agent}cqh>"',
		},
	},
	'LogFormat.Name' => {
		new   => 'Parameter',
		using => {
			id          => 43,
			name        => 'LogFormat.Name',
			config_file => 'logs_xml.config',
			value       => 'custom_ats_2',
		},
	},
	'LogObject.Format' => {
		new   => 'Parameter',
		using => {
			id          => 44,
			name        => 'LogObject.Format',
			config_file => 'logs_xml.config',
			value       => 'custom_ats_2',
		},
	},
	'LogObject.Filename' => {
		new   => 'Parameter',
		using => {
			id          => 45,
			name        => 'LogObject.Filename',
			config_file => 'logs_xml.config',
			value       => 'custom_ats_2',
		},
	},
	'cache.config_location' => {
		new   => 'Parameter',
		using => {
			id          => 46,
			name        => 'location',
			config_file => 'cache.config',
			value       => '/opt/trafficserver/etc/trafficserver/',
		},
	},
	'CONFIG-proxy.config.cache.control.filename' => {
		new   => 'Parameter',
		using => {
			id          => 47,
			name        => 'CONFIG proxy.config.cache.control.filename',
			config_file => 'records.config',
			value       => 'STRING cache.config',
		},
	},
	'regex_revalidate.so' => {
		new   => 'Parameter',
		using => {
			id          => 48,
			name        => 'regex_revalidate.so',
			config_file => 'plugin.config',
			value       => '--config regex_revalidate.config',
		},
	},
	'regex_revalidate.config_location' => {
		new   => 'Parameter',
		using => {
			id          => 49,
			name        => 'location',
			config_file => 'regex_revalidate.config',
			value       => '/opt/trafficserver/etc/trafficserver',
		},
	},
	'hosting.config_location' => {
		new   => 'Parameter',
		using => {
			id          => 50,
			name        => 'location',
			config_file => 'hosting.config',
			value       => '/opt/trafficserver/etc/trafficserver/',
		},
	},
	'volume.config_location' => {
		new   => 'Parameter',
		using => {
			id          => 51,
			name        => 'location',
			config_file => 'volume.config',
			value       => '/opt/trafficserver/etc/trafficserver/',
		},
	},
	'allow_ip' => {
		new   => 'Parameter',
		using => {
			id          => 52,
			name        => 'allow_ip',
			config_file => 'astats.config',
			value       => '127.0.0.1,172.39.0.0/16,33.101.99.0/24',
		},
	},
	'allow_ip6' => {
		new   => 'Parameter',
		using => {
			id          => 53,
			name        => 'allow_ip6',
			config_file => 'astats.config',
			value       => '::1,2033:D011:3300::336/64,2033:D011:3300::335/64,2033:D021:3300::333/64,2033:D021:3300::334/64',
		},
	},
	'record_types' => {
		new   => 'Parameter',
		using => {
			id          => 54,
			name        => 'record_types',
			config_file => 'astats.config',
			value       => '144',
		},
	},
	'astats.config_location' => {
		new   => 'Parameter',
		using => {
			id          => 55,
			name        => 'location',
			config_file => 'astats.config',
			value       => '/opt/trafficserver/etc/trafficserver',
		},
	},
	'astats.config_path' => {
		new   => 'Parameter',
		using => {
			id          => 56,
			name        => 'path',
			config_file => 'astats.config',
			value       => '_astats',
		},
	},
	'storage.config_location' => {
		new   => 'Parameter',
		using => {
			id          => 57,
			name        => 'location',
			config_file => 'storage.config',
			value       => '/opt/trafficserver/etc/trafficserver/',
		},
	},
	'Drive_Prefix' => {
		new   => 'Parameter',
		using => {
			id          => 58,
			name        => 'Drive_Prefix',
			config_file => 'storage.config',
			value       => '/dev/sd',
		},
	},
	'Drive_Letters' => {
		new   => 'Parameter',
		using => {
			id          => 59,
			name        => 'Drive_Letters',
			config_file => 'storage.config',
			value       => 'b,c,d,e,f,g,h,i,j,k,l,m,n,o',
		},
	},
	'Disk_Volume' => {
		new   => 'Parameter',
		using => {
			id          => 60,
			name        => 'Disk_Volume',
			config_file => 'storage.config',
			value       => '1',
		},
	},
	'CONFIG-proxy.config.hostdb.storage_size' => {
		new   => 'Parameter',
		using => {
			id          => 61,
			name        => 'CONFIG proxy.config.hostdb.storage_size',
			config_file => 'records.config',
			value       => 'INT 33554432',
		},
	},
	'regex_revalidate.config_max_days' => {
		new   => 'Parameter',
		using => {
			id          => 63,
			name        => 'maxRevalDurationDays',
			config_file => 'regex_revalidate.config',
			value       => 3,
		},
	},
	'regex_revalidate.config_maxRevalDurationDays' => {
		new   => 'Parameter',
		using => {
			id          => 64,
			name        => 'maxRevalDurationDays',
			config_file => 'regex_revalidate.config',
			value       => 90,
		},
	},
	'unassigned_parameter_1' => {
		new   => 'Parameter',
		using => {
			id          => 65,
			name        => 'unassigned_parameter_1',
			config_file => 'whaterver.config',
			value       => 852,
		},
	},
	'package_trafficserver' => {
		new   => 'Parameter',
		using => {
			id          => 66,
			name        => 'trafficserver',
			config_file => 'package',
			value       => '5.3.2-765.f4354b9.el7.centos.x86_64',
		},
	},
	'use_tenancy' => {
		new   => 'Parameter',
		using => {
			id          => 67,
			name        => 'use_tenancy',
			config_file => 'global',
			value       => '1',
		},
	},
);

sub get_definition {
	my ( $self, $name ) = @_;
	return $definition_for{$name};
}

sub all_fixture_names {
	# sort by db name to guarantee insertion order
	return (sort { lc($definition_for{$a}{using}{id}) cmp lc($definition_for{$b}{using}{id}) } keys %definition_for);
}

__PACKAGE__->meta->make_immutable;

1;
