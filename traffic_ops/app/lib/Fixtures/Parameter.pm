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
	## id => 1
	'allow_ip' => {
		new   => 'Parameter',
		using => {
			name        => 'allow_ip',
			config_file => 'astats.config',
			value       => '127.0.0.1,172.39.0.0/16,33.101.99.0/24',
		},
	},
	## id => 2
	'allow_ip6' => {
		new   => 'Parameter',
		using => {
			name        => 'allow_ip6',
			config_file => 'astats.config',
			value       => '::1,2033:D011:3300::336/64,2033:D011:3300::335/64,2033:D021:3300::333/64,2033:D021:3300::334/64',
		},
	},
	## id => 3
	'astats_over_http.so' => {
		new   => 'Parameter',
		using => {
			name        => 'astats_over_http.so',
			config_file => 'plugin.config',
			value       => '_astats 33.101.99.100,172.39.19.39,172.39.19.49,172.39.19.49,172.39.29.49',
		},
	},
	## id => 4
	'CONFIG-proxy.config.allocator.debug_filter' => {
		new   => 'Parameter',
		using => {
			name        => 'CONFIG proxy.config.allocator.debug_filter',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	## id => 5
	'CONFIG-proxy.config.allocator.enable_reclaim' => {
		new   => 'Parameter',
		using => {
			name        => 'CONFIG proxy.config.allocator.enable_reclaim',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	## id => 6
	'CONFIG-proxy.config.allocator.max_overage' => {
		new   => 'Parameter',
		using => {
			name        => 'CONFIG proxy.config.allocator.max_overage',
			config_file => 'records.config',
			value       => 'INT 3',
		},
	},
	## id => 7
	'CONFIG-proxy.config.cache.control.filename' => {
		new   => 'Parameter',
		using => {
			name        => 'CONFIG proxy.config.cache.control.filename',
			config_file => 'records.config',
			value       => 'STRING cache.config',
		},
	},
	## id => 8
	'CONFIG-proxy.config.diags.show_location' => {
		new   => 'Parameter',
		using => {
			name        => 'CONFIG proxy.config.diags.show_location',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	## id => 9
	'CONFIG-proxy.config.http.cache.allow_empty_doc' => {
		new   => 'Parameter',
		using => {
			name        => 'CONFIG proxy.config.http.cache.allow_empty_doc',
			config_file => 'records.config',
			value       => 'INT 0',
		},
	},
	## id => 10
	'CONFIG-proxy.config.http.parent_proxy.file' => {
		new   => 'Parameter',
		using => {
			name        => 'CONFIG proxy.config.http.parent_proxy.file',
			config_file => 'records.config',
			value       => 'STRING parent.config',
		},
	},
	## id => 11
	'CONFIG-proxy.config.hostdb.storage_size' => {
		new   => 'Parameter',
		using => {
			name        => 'CONFIG proxy.config.hostdb.storage_size',
			config_file => 'records.config',
			value       => 'INT 33554432',
		},
	},
	## id => 12
	'Disk_Volume' => {
		new   => 'Parameter',
		using => {
			name        => 'Disk_Volume',
			config_file => 'storage.config',
			value       => '1',
		},
	},
	## id => 13
	domain_name => {
		new   => 'Parameter',
		using => {
			name        => 'domain_name',
			value       => 'foo.com',
			config_file => 'CRConfig.json',
		},
	},
	## id => 14
	'Drive_Letters' => {
		new   => 'Parameter',
		using => {
			name        => 'Drive_Letters',
			config_file => 'storage.config',
			value       => 'b,c,d,e,f,g,h,i,j,k,l,m,n,o',
		},
	},
	## id => 15
	'Drive_Prefix' => {
		new   => 'Parameter',
		using => {
			name        => 'Drive_Prefix',
			config_file => 'storage.config',
			value       => '/dev/sd',
		},
	},
	## id => 16
	'error_url' => {
		new   => 'Parameter',
		using => {
			name        => 'error_url',
			config_file => 'url_sig_cdl-c2.config',
			value       => '403',
		},
	},
	## id => 17
	health_threadhold_available_bandwidth_in_kbps => {
		new   => 'Parameter',
		using => {
			name        => 'health.threshold.availableBandwidthInKbps',
			value       => '>1750000',
			config_file => 'rascal.properties',
		},
	},
	## id => 18
	health_threadhold_loadavg => {
		new   => 'Parameter',
		using => {
			name        => 'health.threshold.loadavg',
			value       => '25.0',
			config_file => 'rascal.properties',
		},
	},
	## id => 19
	history_count => {
		new   => 'Parameter',
		using => {
			name        => 'history.count',
			value       => '30',
			config_file => 'rascal.properties',
		},
	},
	## id => 20
	'key0' => {
		new   => 'Parameter',
		using => {
			name        => 'key0',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'HOOJ3Ghq1x4gChp3iQkqVTcPlOj8UCi3',
		},
	},
	## id => 21
	'key1' => {
		new   => 'Parameter',
		using => {
			name        => 'key1',
			config_file => 'url_sig_cdl-c2.config',
			value       => '_9LZYkRnfCS0rCBF7fTQzM9Scwlp2FhO',
		},
	},
	## id => 22
	'key2' => {
		new   => 'Parameter',
		using => {
			name        => 'key2',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'AFpkxfc4oTiyFSqtY6_ohjt3V80aAIxS',
		},
	},
	## id => 23
	'key3' => {
		new   => 'Parameter',
		using => {
			name        => 'key3',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'AL9kzs_SXaRZjPWH8G5e2m4ByTTzkzlc',
		},
	},
	## id => 24
	'key4' => {
		new   => 'Parameter',
		using => {
			name        => 'key4',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'poP3n3szbD1U4vx1xQXV65BvkVgWzfN8',
		},
	},
	## id => 25
	'key5' => {
		new   => 'Parameter',
		using => {
			name        => 'key5',
			config_file => 'url_sig_cdl-c2.config',
			value       => '1ir32ng4C4w137p5oq72kd2wqmIZUrya',
		},
	},
	## id => 26
	'key6' => {
		new   => 'Parameter',
		using => {
			name        => 'key6',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'B1qLptn2T1b_iXeTCWDcVuYvANtH139f',
		},
	},
	## id => 27
	'key7' => {
		new   => 'Parameter',
		using => {
			name        => 'key7',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'PiCV_5OODMzBbsNFMWsBxcQ8v1sK0TYE',
		},
	},
	## id => 28
	'key8' => {
		new   => 'Parameter',
		using => {
			name        => 'key8',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'Ggpv6DqXDvt2s1CETPBpNKwaLk4fTM9l',
		},
	},
	## id => 29
	'key9' => {
		new   => 'Parameter',
		using => {
			name        => 'key9',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'qPlVT_s6kL37aqb6hipDm4Bt55S72mI7',
		},
	},
	## id => 30
	'key10' => {
		new   => 'Parameter',
		using => {
			name        => 'key10',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'BsI5A9EmWrobIS1FeuOs1z9fm2t2WSBe',
		},
	},
	## id => 31
	'key11' => {
		new   => 'Parameter',
		using => {
			name        => 'key11',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'A54y66NCIj897GjS4yA9RrsSPtCUnQXP',
		},
	},
	## id => 32
	'key12' => {
		new   => 'Parameter',
		using => {
			name        => 'key12',
			config_file => 'url_sig_cdl-c2.config',
			value       => '2jZH0NDPSJttIr4c2KP510f47EKqTQAu',
		},
	},
	## id => 33
	'key13' => {
		new   => 'Parameter',
		using => {
			name        => 'key13',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'XduT2FBjBmmVID5JRB5LEf9oR5QDtBgC',
		},
	},
	## id => 34
	'key14' => {
		new   => 'Parameter',
		using => {
			name        => 'key14',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'D9nH0SvK_0kP5w8QNd1UFJ28ulFkFKPn',
		},
	},
	## id => 35
	'key15' => {
		new   => 'Parameter',
		using => {
			name        => 'key15',
			config_file => 'url_sig_cdl-c2.config',
			value       => 'udKXWYNwbXXweaaLzaKDGl57OixnIIcm',
		},
	},
	## id => 36
	'LOCAL-proxy.config.cache.interim.storage' => {
		new   => 'Parameter',
		using => {
			name        => 'LOCAL proxy.config.cache.interim.storage',
			config_file => 'records.config',
			value       => 'STRING NULL',
		},
	},
	## id => 37
	'url_sig_cdl-c2.config_location' => {
		new   => 'Parameter',
		using => {
			name        => 'location',
			config_file => 'url_sig_cdl-c2.config',
			value       => '/opt/trafficserver/etc/trafficserver',
		},
	},
	## id => 38
	'12M_location' => {
		new   => 'Parameter',
		using => {
			name        => 'location',
			config_file => '12M_facts',
			value       => '/opt/ort',
		},
	},
	## id => 39
	'cacheurl_location' => {
		new   => 'Parameter',
		using => {
			name        => 'location',
			config_file => 'cacheurl.config',
			value       => '/opt/trafficserver/etc/trafficserver/',
		},
	},
	## id => 40
	'ip_allow_location' => {
		new   => 'Parameter',
		using => {
			name        => 'location',
			config_file => 'ip_allow.config',
			value       => '/opt/trafficserver/etc/trafficserver',
		},
	},
	## id => 41
	'crontab_root_location' => {
		new   => 'Parameter',
		using => {
			name        => 'location',
			config_file => 'crontab_root',
			value       => '/var/spool/cron',
		},
	},
	## id => 42
	'hdr_rw_cdl-c2.config_location' => {
		new   => 'Parameter',
		using => {
			name        => 'location',
			config_file => 'hdr_rw_cdl-c2.config',
			value       => '/opt/trafficserver/etc/trafficserver',
		},
	},
	## id => 43
	'50-ats.rules_location' => {
		new   => 'Parameter',
		using => {
			name        => 'location',
			config_file => '50-ats.rules',
			value       => '/etc/udev/rules.d/',
		},
	},
	## id => 44
	'parent.config_location' => {
		new   => 'Parameter',
		using => {
			name        => 'location',
			config_file => 'parent.config',
			value       => '/opt/trafficserver/etc/trafficserver/',
		},
	},
	## id => 45
	'remap.config_location' => {
		new   => 'Parameter',
		using => {
			name        => 'location',
			config_file => 'remap.config',
			value       => '/opt/trafficserver/etc/trafficserver/',
		},
	},
	## id => 46
	'drop_qstring.config_location' => {
		new   => 'Parameter',
		using => {
			name        => 'location',
			config_file => 'drop_qstring.config',
			value       => '/opt/trafficserver/etc/trafficserver',
		},
	},
	## id => 47
	'cache.config_location' => {
		new   => 'Parameter',
		using => {
			name        => 'location',
			config_file => 'cache.config',
			value       => '/opt/trafficserver/etc/trafficserver/',
		},
	},
	## id => 48
	'regex_revalidate.config_location' => {
		new   => 'Parameter',
		using => {
			name        => 'location',
			config_file => 'regex_revalidate.config',
			value       => '/opt/trafficserver/etc/trafficserver',
		},
	},
	## id => 49
	'hosting.config_location' => {
		new   => 'Parameter',
		using => {
			name        => 'location',
			config_file => 'hosting.config',
			value       => '/opt/trafficserver/etc/trafficserver/',
		},
	},
	## id => 50
	'volume.config_location' => {
		new   => 'Parameter',
		using => {
			name        => 'location',
			config_file => 'volume.config',
			value       => '/opt/trafficserver/etc/trafficserver/',
		},
	},
	## id => 51
	'astats.config_location' => {
		new   => 'Parameter',
		using => {
			name        => 'location',
			config_file => 'astats.config',
			value       => '/opt/trafficserver/etc/trafficserver',
		},
	},
	## id => 52
	'storage.config_location' => {
		new   => 'Parameter',
		using => {
			name        => 'location',
			config_file => 'storage.config',
			value       => '/opt/trafficserver/etc/trafficserver/',
		},
	},
	## id => 53
	'LogFormat.Format' => {
		new   => 'Parameter',
		using => {
			name        => 'LogFormat.Format',
			config_file => 'logs_xml.config',
			value =>
				'%<cqtq> chi=%<chi> phn=%<phn> shn=%<shn> url=%<cquuc> cqhm=%<cqhm> cqhv=%<cqhv> pssc=%<pssc> ttms=%<ttms> b=%<pscl> sssc=%<sssc> sscl=%<sscl> cfsc=%<cfsc> pfsc=%<pfsc> crc=%<crc> phr=%<phr> uas="%<{User-Agent}cqh>"',
		},
	},
	## id => 54
	'LogFormat.Name' => {
		new   => 'Parameter',
		using => {
			name        => 'LogFormat.Name',
			config_file => 'logs_xml.config',
			value       => 'custom_ats_2',
		},
	},
	## id => 55
	'LogObject.Format' => {
		new   => 'Parameter',
		using => {
			name        => 'LogObject.Format',
			config_file => 'logs_xml.config',
			value       => 'custom_ats_2',
		},
	},
	## id => 56
	'LogObject.Filename' => {
		new   => 'Parameter',
		using => {
			name        => 'LogObject.Filename',
			config_file => 'logs_xml.config',
			value       => 'custom_ats_2',
		},
	},
	## id => 57
	'regex_revalidate.config_max_days' => {
		new   => 'Parameter',
		using => {
			name        => 'maxRevalDurationDays',
			config_file => 'regex_revalidate.config',
			value       => 90,
		},
	},
	## id => 58
	'regex_revalidate.config_maxRevalDurationDays' => {
		new   => 'Parameter',
		using => {
			name        => 'maxRevalDurationDays',
			config_file => 'regex_revalidate.config',
			value       => 90,
		},
	},
	## id => 59
	'astats.config_path' => {
		new   => 'Parameter',
		using => {
			name        => 'path',
			config_file => 'astats.config',
			value       => '_astats',
		},
	},
	## id => 60
	'record_types' => {
		new   => 'Parameter',
		using => {
			name        => 'record_types',
			config_file => 'astats.config',
			value       => '144',
		},
	},
	## id => 61
	'regex_revalidate.so' => {
		new   => 'Parameter',
		using => {
			name        => 'regex_revalidate.so',
			config_file => 'plugin.config',
			value       => '--config regex_revalidate.config',
		},
	},
	## id => 62
	'regex_revalidate.config_snapshot_dir' => {
		new   => 'Parameter',
		using => {
			name        => 'snapshot_dir',
			config_file => 'regex_revalidate.config',
			value       => 'public/Trafficserver-Snapshots/',
		},
	}
);

sub get_definition {
	my ( $self, $name ) = @_;
	return $definition_for{$name};
}

sub all_fixture_names {
	# sort by db name to guarantee insertion order
	return (sort { lc($definition_for{$a}{using}{name}) cmp lc($definition_for{$b}{using}{name}) } keys %definition_for);
}

__PACKAGE__->meta->make_immutable;

1;
