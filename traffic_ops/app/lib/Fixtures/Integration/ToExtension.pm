package Fixtures::Integration::ToExtension;

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


# Do not edit! Generated code.
# See https://github.com/Comcast/traffic_control/wiki/The%20Kabletown%20example

use Moose;
extends 'DBIx::Class::EasyFixture';
use namespace::autoclean;

my %definition_for = (
'0' => { new => 'ToExtension', using => { name => 'ILO_PING', isactive => '1', last_updated => '2015-12-10 15:44:37', servercheck_short_name => 'ILO', type => '5', version => '1.0.0', info_url => 'http://foo.com/bar.html', description => undef, script_file => 'ping', servercheck_column_name => 'aa', additional_config_json => '{ "select": "ilo_ip_address", "cron": "9 * * * *" }', }, },
'1' => { new => 'ToExtension', using => { name => '10G_PING', script_file => 'ping', servercheck_short_name => '10G', version => '1.0.0', info_url => 'http://foo.com/bar.html', isactive => '1', servercheck_column_name => 'ab', type => '5', additional_config_json => '{ "select": "ip_address", "cron": "18 * * * *" }', description => undef, last_updated => '2015-12-10 15:44:37', }, },
'2' => { new => 'ToExtension', using => { name => 'FQDN_PING', isactive => '1', last_updated => '2015-12-10 15:44:37', servercheck_short_name => 'FQDN', type => '5', version => '1.0.0', description => undef, info_url => 'http://foo.com/bar.html', script_file => 'ping', servercheck_column_name => 'ac', additional_config_json => '{ "select": "host_name", "cron": "27 * * * *" }', }, },
'3' => { new => 'ToExtension', using => { name => 'CHECK_DSCP', servercheck_column_name => 'ad', servercheck_short_name => 'DSCP', type => '5', version => '1.0.0', last_updated => '2015-12-10 15:44:37', script_file => 'dscp', additional_config_json => '{ "select": "ilo_ip_address", "cron": "36 * * * *" }', description => undef, info_url => 'http://foo.com/bar.html', isactive => '1', }, },
'4' => { new => 'ToExtension', using => { name => 'OPEN', additional_config_json => '', info_url => 'http://foo.com/bar.html', script_file => 'dscp', version => '1.0.0', description => undef, isactive => '0', last_updated => '2015-12-10 15:44:37', servercheck_column_name => 'ae', servercheck_short_name => '', type => '7', }, },
'5' => { new => 'ToExtension', using => { name => 'OPEN', isactive => '0', servercheck_short_name => '', type => '7', version => '1.0.0', description => undef, info_url => 'http://foo.com/bar.html', last_updated => '2015-12-10 15:44:37', script_file => 'dscp', servercheck_column_name => 'af', additional_config_json => '', }, },
'6' => { new => 'ToExtension', using => { name => 'IPV6_PING', servercheck_short_name => '10G6', version => '1.0.0', additional_config_json => '{ "select": "ip6_address", "cron": "0 * * * *" }', description => undef, script_file => 'ping', info_url => 'http://foo.com/bar.html', isactive => '1', last_updated => '2015-12-10 15:44:37', servercheck_column_name => 'ag', type => '5', }, },
'7' => { new => 'ToExtension', using => { name => 'OPEN', servercheck_short_name => '', isactive => '0', last_updated => '2015-12-10 15:44:37', script_file => 'dscp', servercheck_column_name => 'ah', version => '1.0.0', additional_config_json => '', description => undef, info_url => 'http://foo.com/bar.html', type => '7', }, },
'8' => { new => 'ToExtension', using => { name => 'CHECK_STATS', additional_config_json => '{ "select": "ilo_ip_address", "cron": "54 * * * *" }', description => undef, isactive => '1', script_file => 'ping', servercheck_short_name => 'STAT', type => '5', info_url => 'http://foo.com/bar.html', last_updated => '2015-12-10 15:44:37', servercheck_column_name => 'ai', version => '1.0.0', }, },
'9' => { new => 'ToExtension', using => { name => 'OPEN', info_url => 'http://foo.com/bar.html', script_file => 'dscp', servercheck_column_name => 'aj', type => '7', version => '1.0.0', additional_config_json => '', description => undef, isactive => '0', last_updated => '2015-12-10 15:44:37', servercheck_short_name => '', }, },
'10' => { new => 'ToExtension', using => { name => 'CHECK_MTU', servercheck_column_name => 'ak', servercheck_short_name => 'MTU', version => '1.0.0', additional_config_json => '{ "select": "ip_address", "cron": "45 * * * *" }', description => undef, last_updated => '2015-12-10 15:44:37', script_file => 'ping', type => '5', info_url => 'http://foo.com/bar.html', isactive => '1', }, },
'11' => { new => 'ToExtension', using => { name => 'CHECK_TRAFFIC_ROUTER_STATUS', additional_config_json => '{ "select": "ilo_ip_address", "cron": "10 * * * *" }', description => undef, info_url => 'http://foo.com/bar.html', isactive => '1', last_updated => '2015-12-10 15:44:37', script_file => 'ping', version => '1.0.0', servercheck_column_name => 'al', servercheck_short_name => 'TRTR', type => '5', }, },
'12' => { new => 'ToExtension', using => { name => 'CHECK_TRAFFIC_MONITOR_STATUS', type => '5', version => '1.0.0', description => undef, isactive => '1', last_updated => '2015-12-10 15:44:37', script_file => 'ping', servercheck_short_name => 'TRMO', additional_config_json => '{ "select": "ip_address", "cron": "10 * * * *" }', info_url => 'http://foo.com/bar.html', servercheck_column_name => 'am', }, },
'13' => { new => 'ToExtension', using => { name => 'CACHE_HIT_RATIO_LAST_15', info_url => 'http://foo.com/bar.html', isactive => '1', servercheck_column_name => 'an', version => '1.0.0', additional_config_json => '{ "select": "ilo_ip_address", "cron": "0,15,30,45 * * * *" }', description => undef, servercheck_short_name => 'CHR', type => '6', last_updated => '2015-12-10 15:44:37', script_file => 'ping', }, },
'14' => { new => 'ToExtension', using => { name => 'DISK_UTILIZATION', script_file => 'ping', type => '6', info_url => 'http://foo.com/bar.html', isactive => '1', last_updated => '2015-12-10 15:44:37', servercheck_column_name => 'ao', servercheck_short_name => 'CDU', version => '1.0.0', additional_config_json => '{ "select": "ilo_ip_address", "cron": "20 * * * *" }', description => undef, }, },
'15' => { new => 'ToExtension', using => { name => 'ORT_ERROR_COUNT', description => undef, info_url => 'http://foo.com/bar.html', script_file => 'ping', servercheck_short_name => 'ORT', type => '6', servercheck_column_name => 'ap', version => '1.0.0', additional_config_json => '{ "select": "ilo_ip_address", "cron": "40 * * * *" }', isactive => '1', last_updated => '2015-12-10 15:44:37', }, },
'16' => { new => 'ToExtension', using => { name => 'OPEN', version => '1.0.0', additional_config_json => '', info_url => 'http://foo.com/bar.html', script_file => 'dscp', servercheck_column_name => 'aq', type => '7', description => undef, isactive => '0', last_updated => '2015-12-10 15:44:37', servercheck_short_name => '', }, },
'17' => { new => 'ToExtension', using => { name => 'OPEN', isactive => '0', last_updated => '2015-12-10 15:44:37', script_file => 'dscp', servercheck_short_name => '', additional_config_json => '', description => undef, info_url => 'http://foo.com/bar.html', type => '7', version => '1.0.0', servercheck_column_name => 'ar', }, },
'18' => { new => 'ToExtension', using => { name => 'OPEN', info_url => 'http://foo.com/bar.html', isactive => '0', script_file => 'dscp', servercheck_short_name => '', type => '7', version => '1.0.0', description => undef, last_updated => '2015-12-10 15:44:37', servercheck_column_name => 'bf', additional_config_json => '', }, },
'19' => { new => 'ToExtension', using => { name => 'OPEN', servercheck_short_name => '', version => '1.0.0', description => undef, last_updated => '2015-12-10 15:44:37', script_file => 'dscp', servercheck_column_name => 'at', type => '7', additional_config_json => '', info_url => 'http://foo.com/bar.html', isactive => '0', }, },
'20' => { new => 'ToExtension', using => { name => 'OPEN', description => undef, isactive => '0', last_updated => '2015-12-10 15:44:37', servercheck_column_name => 'au', servercheck_short_name => '', type => '7', additional_config_json => '', version => '1.0.0', script_file => 'dscp', info_url => 'http://foo.com/bar.html', }, },
'21' => { new => 'ToExtension', using => { name => 'OPEN', script_file => 'dscp', description => undef, info_url => 'http://foo.com/bar.html', last_updated => '2015-12-10 15:44:37', version => '1.0.0', additional_config_json => '', isactive => '0', servercheck_column_name => 'av', servercheck_short_name => '', type => '7', }, },
'22' => { new => 'ToExtension', using => { name => 'OPEN', type => '7', version => '1.0.0', isactive => '0', script_file => 'dscp', servercheck_column_name => 'aw', servercheck_short_name => '', additional_config_json => '', description => undef, info_url => 'http://foo.com/bar.html', last_updated => '2015-12-10 15:44:37', }, },
'23' => { new => 'ToExtension', using => { name => 'OPEN', info_url => 'http://foo.com/bar.html', isactive => '0', last_updated => '2015-12-10 15:44:37', script_file => 'dscp', description => undef, servercheck_column_name => 'ax', servercheck_short_name => '', type => '7', version => '1.0.0', additional_config_json => '', }, },
'24' => { new => 'ToExtension', using => { name => 'OPEN', description => undef, isactive => '0', servercheck_short_name => '', type => '7', version => '1.0.0', script_file => 'dscp', servercheck_column_name => 'ay', additional_config_json => '', info_url => 'http://foo.com/bar.html', last_updated => '2015-12-10 15:44:37', }, },
'25' => { new => 'ToExtension', using => { name => 'OPEN', servercheck_column_name => 'az', servercheck_short_name => '', type => '7', additional_config_json => '', info_url => 'http://foo.com/bar.html', last_updated => '2015-12-10 15:44:37', version => '1.0.0', description => undef, isactive => '0', script_file => 'dscp', }, },
'26' => { new => 'ToExtension', using => { name => 'OPEN', additional_config_json => '', isactive => '0', type => '7', version => '1.0.0', servercheck_short_name => '', description => undef, info_url => 'http://foo.com/bar.html', last_updated => '2015-12-10 15:44:37', script_file => 'dscp', servercheck_column_name => 'ba', }, },
'27' => { new => 'ToExtension', using => { name => 'OPEN', script_file => 'dscp', servercheck_column_name => 'bb', servercheck_short_name => '', info_url => 'http://foo.com/bar.html', isactive => '0', last_updated => '2015-12-10 15:44:37', type => '7', version => '1.0.0', additional_config_json => '', description => undef, }, },
'28' => { new => 'ToExtension', using => { name => 'OPEN', description => undef, isactive => '0', servercheck_column_name => 'bc', version => '1.0.0', servercheck_short_name => '', type => '7', additional_config_json => '', info_url => 'http://foo.com/bar.html', last_updated => '2015-12-10 15:44:37', script_file => 'dscp', }, },
'29' => { new => 'ToExtension', using => { name => 'OPEN', info_url => 'http://foo.com/bar.html', script_file => 'dscp', servercheck_column_name => 'bd', servercheck_short_name => '', additional_config_json => '', last_updated => '2015-12-10 15:44:37', type => '7', version => '1.0.0', description => undef, isactive => '0', }, },
'30' => { new => 'ToExtension', using => { name => 'OPEN', description => undef, last_updated => '2015-12-10 15:44:37', servercheck_short_name => '', version => '1.0.0', script_file => 'dscp', servercheck_column_name => 'be', type => '7', additional_config_json => '', info_url => 'http://foo.com/bar.html', isactive => '0', }, },
);

sub name {
		return "ToExtension";
}

sub get_definition {
		my ( $self, $name ) = @_;
		return $definition_for{$name};
}

sub all_fixture_names {
	# sort by db name to guarantee insertion order
	return (sort { $definition_for{$a}{using}{name} cmp $definition_for{$b}{using}{name} } keys %definition_for);
}

__PACKAGE__->meta->make_immutable;
1;
