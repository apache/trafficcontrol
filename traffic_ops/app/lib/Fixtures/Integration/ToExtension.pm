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
'0' => { new => 'ToExtension', => using => { isactive => '1', last_updated => '2015-12-10 15:44:37', servercheck_short_name => 'ILO', type => '31', version => '1.0.0', info_url => 'http://foo.com/bar.html', description => undef, id => '1', name => 'ILO_PING', script_file => 'ping', servercheck_column_name => 'aa', additional_config_json => '{ "select": "ilo_ip_address", "cron": "9 * * * *" }', }, },
'1' => { new => 'ToExtension', => using => { name => '10G_PING', script_file => 'ping', servercheck_short_name => '10G', version => '1.0.0', id => '2', info_url => 'http://foo.com/bar.html', isactive => '1', servercheck_column_name => 'ab', type => '31', additional_config_json => '{ "select": "ip_address", "cron": "18 * * * *" }', description => undef, last_updated => '2015-12-10 15:44:37', }, },
'2' => { new => 'ToExtension', => using => { isactive => '1', last_updated => '2015-12-10 15:44:37', servercheck_short_name => 'FQDN', type => '31', version => '1.0.0', id => '3', description => undef, info_url => 'http://foo.com/bar.html', name => 'FQDN_PING', script_file => 'ping', servercheck_column_name => 'ac', additional_config_json => '{ "select": "host_name", "cron": "27 * * * *" }', }, },
'3' => { new => 'ToExtension', => using => { servercheck_column_name => 'ad', servercheck_short_name => 'DSCP', type => '31', version => '1.0.0', id => '4', last_updated => '2015-12-10 15:44:37', name => 'CHECK_DSCP', script_file => 'dscp', additional_config_json => '{ "select": "ilo_ip_address", "cron": "36 * * * *" }', description => undef, info_url => 'http://foo.com/bar.html', isactive => '1', }, },
'4' => { new => 'ToExtension', => using => { additional_config_json => '', id => '5', info_url => 'http://foo.com/bar.html', name => 'OPEN', script_file => 'dscp', version => '1.0.0', description => undef, isactive => '0', last_updated => '2015-12-10 15:44:37', servercheck_column_name => 'ae', servercheck_short_name => '', type => '33', }, },
'5' => { new => 'ToExtension', => using => { isactive => '0', name => 'OPEN', servercheck_short_name => '', type => '33', version => '1.0.0', description => undef, info_url => 'http://foo.com/bar.html', last_updated => '2015-12-10 15:44:37', script_file => 'dscp', servercheck_column_name => 'af', additional_config_json => '', id => '6', }, },
'6' => { new => 'ToExtension', => using => { servercheck_short_name => '10G6', version => '1.0.0', additional_config_json => '{ "select": "ip6_address", "cron": "0 * * * *" }', description => undef, id => '7', name => 'IPV6_PING', script_file => 'ping', info_url => 'http://foo.com/bar.html', isactive => '1', last_updated => '2015-12-10 15:44:37', servercheck_column_name => 'ag', type => '31', }, },
'7' => { new => 'ToExtension', => using => { servercheck_short_name => '', id => '8', isactive => '0', last_updated => '2015-12-10 15:44:37', script_file => 'dscp', servercheck_column_name => 'ah', version => '1.0.0', additional_config_json => '', description => undef, info_url => 'http://foo.com/bar.html', name => 'OPEN', type => '33', }, },
'8' => { new => 'ToExtension', => using => { additional_config_json => '{ "select": "ilo_ip_address", "cron": "54 * * * *" }', description => undef, id => '9', isactive => '1', script_file => 'ping', servercheck_short_name => 'STAT', type => '31', info_url => 'http://foo.com/bar.html', last_updated => '2015-12-10 15:44:37', name => 'CHECK_STATS', servercheck_column_name => 'ai', version => '1.0.0', }, },
'9' => { new => 'ToExtension', => using => { id => '10', info_url => 'http://foo.com/bar.html', script_file => 'dscp', servercheck_column_name => 'aj', type => '33', version => '1.0.0', additional_config_json => '', description => undef, isactive => '0', last_updated => '2015-12-10 15:44:37', name => 'OPEN', servercheck_short_name => '', }, },
'10' => { new => 'ToExtension', => using => { servercheck_column_name => 'ak', servercheck_short_name => 'MTU', version => '1.0.0', additional_config_json => '{ "select": "ip_address", "cron": "45 * * * *" }', description => undef, last_updated => '2015-12-10 15:44:37', name => 'CHECK_MTU', script_file => 'ping', type => '31', id => '11', info_url => 'http://foo.com/bar.html', isactive => '1', }, },
'11' => { new => 'ToExtension', => using => { additional_config_json => '{ "select": "ilo_ip_address", "cron": "10 * * * *" }', description => undef, info_url => 'http://foo.com/bar.html', isactive => '1', last_updated => '2015-12-10 15:44:37', name => 'CHECK_TRAFFIC_ROUTER_STATUS', script_file => 'ping', version => '1.0.0', id => '12', servercheck_column_name => 'al', servercheck_short_name => 'TRTR', type => '31', }, },
'12' => { new => 'ToExtension', => using => { type => '31', version => '1.0.0', description => undef, id => '13', isactive => '1', last_updated => '2015-12-10 15:44:37', script_file => 'ping', servercheck_short_name => 'TRMO', additional_config_json => '{ "select": "ip_address", "cron": "10 * * * *" }', info_url => 'http://foo.com/bar.html', name => 'CHECK_TRAFFIC_MONITOR_STATUS', servercheck_column_name => 'am', }, },
'13' => { new => 'ToExtension', => using => { id => '14', info_url => 'http://foo.com/bar.html', isactive => '1', name => 'CACHE_HIT_RATIO_LAST_15', servercheck_column_name => 'an', version => '1.0.0', additional_config_json => '{ "select": "ilo_ip_address", "cron": "0,15,30,45 * * * *" }', description => undef, servercheck_short_name => 'CHR', type => '32', last_updated => '2015-12-10 15:44:37', script_file => 'ping', }, },
'14' => { new => 'ToExtension', => using => { script_file => 'ping', type => '32', info_url => 'http://foo.com/bar.html', isactive => '1', id => '15', last_updated => '2015-12-10 15:44:37', name => 'DISK_UTILIZATION', servercheck_column_name => 'ao', servercheck_short_name => 'CDU', version => '1.0.0', additional_config_json => '{ "select": "ilo_ip_address", "cron": "20 * * * *" }', description => undef, }, },
'15' => { new => 'ToExtension', => using => { description => undef, info_url => 'http://foo.com/bar.html', script_file => 'ping', servercheck_short_name => 'ORT', type => '32', servercheck_column_name => 'ap', version => '1.0.0', additional_config_json => '{ "select": "ilo_ip_address", "cron": "40 * * * *" }', id => '16', isactive => '1', last_updated => '2015-12-10 15:44:37', name => 'ORT_ERROR_COUNT', }, },
'16' => { new => 'ToExtension', => using => { version => '1.0.0', additional_config_json => '', info_url => 'http://foo.com/bar.html', name => 'OPEN', script_file => 'dscp', servercheck_column_name => 'aq', type => '33', description => undef, id => '17', isactive => '0', last_updated => '2015-12-10 15:44:37', servercheck_short_name => '', }, },
'17' => { new => 'ToExtension', => using => { isactive => '0', last_updated => '2015-12-10 15:44:37', script_file => 'dscp', servercheck_short_name => '', additional_config_json => '', description => undef, info_url => 'http://foo.com/bar.html', type => '33', version => '1.0.0', id => '18', name => 'OPEN', servercheck_column_name => 'ar', }, },
'18' => { new => 'ToExtension', => using => { info_url => 'http://foo.com/bar.html', isactive => '0', name => 'OPEN', script_file => 'dscp', servercheck_short_name => '', type => '33', version => '1.0.0', id => '19', description => undef, last_updated => '2015-12-10 15:44:37', servercheck_column_name => 'bf', additional_config_json => '', }, }, 
'19' => { new => 'ToExtension', => using => { servercheck_short_name => '', version => '1.0.0', description => undef, id => '20', last_updated => '2015-12-10 15:44:37', script_file => 'dscp', servercheck_column_name => 'at', type => '33', additional_config_json => '', info_url => 'http://foo.com/bar.html', isactive => '0', name => 'OPEN', }, },
'20' => { new => 'ToExtension', => using => { description => undef, id => '21', isactive => '0', last_updated => '2015-12-10 15:44:37', servercheck_column_name => 'au', servercheck_short_name => '', type => '33', additional_config_json => '', version => '1.0.0', name => 'OPEN', script_file => 'dscp', info_url => 'http://foo.com/bar.html', }, },
'21' => { new => 'ToExtension', => using => { script_file => 'dscp', description => undef, id => '22', info_url => 'http://foo.com/bar.html', last_updated => '2015-12-10 15:44:37', name => 'OPEN', version => '1.0.0', additional_config_json => '', isactive => '0', servercheck_column_name => 'av', servercheck_short_name => '', type => '33', }, },
'22' => { new => 'ToExtension', => using => { type => '33', version => '1.0.0', id => '23', isactive => '0', script_file => 'dscp', servercheck_column_name => 'aw', name => 'OPEN', servercheck_short_name => '', additional_config_json => '', description => undef, info_url => 'http://foo.com/bar.html', last_updated => '2015-12-10 15:44:37', }, },
'23' => { new => 'ToExtension', => using => { id => '24', info_url => 'http://foo.com/bar.html', isactive => '0', last_updated => '2015-12-10 15:44:37', name => 'OPEN', script_file => 'dscp', description => undef, servercheck_column_name => 'ax', servercheck_short_name => '', type => '33', version => '1.0.0', additional_config_json => '', }, },
'24' => { new => 'ToExtension', => using => { description => undef, isactive => '0', servercheck_short_name => '', type => '33', version => '1.0.0', script_file => 'dscp', servercheck_column_name => 'ay', additional_config_json => '', id => '25', info_url => 'http://foo.com/bar.html', last_updated => '2015-12-10 15:44:37', name => 'OPEN', }, },
'25' => { new => 'ToExtension', => using => { servercheck_column_name => 'az', servercheck_short_name => '', type => '33', additional_config_json => '', info_url => 'http://foo.com/bar.html', last_updated => '2015-12-10 15:44:37', name => 'OPEN', version => '1.0.0', description => undef, id => '26', isactive => '0', script_file => 'dscp', }, },
'26' => { new => 'ToExtension', => using => { additional_config_json => '', id => '27', isactive => '0', type => '33', version => '1.0.0', servercheck_short_name => '', description => undef, info_url => 'http://foo.com/bar.html', last_updated => '2015-12-10 15:44:37', name => 'OPEN', script_file => 'dscp', servercheck_column_name => 'ba', }, },
'27' => { new => 'ToExtension', => using => { script_file => 'dscp', servercheck_column_name => 'bb', servercheck_short_name => '', id => '28', info_url => 'http://foo.com/bar.html', isactive => '0', last_updated => '2015-12-10 15:44:37', name => 'OPEN', type => '33', version => '1.0.0', additional_config_json => '', description => undef, }, },
'28' => { new => 'ToExtension', => using => { description => undef, isactive => '0', name => 'OPEN', servercheck_column_name => 'bc', version => '1.0.0', servercheck_short_name => '', type => '33', additional_config_json => '', id => '29', info_url => 'http://foo.com/bar.html', last_updated => '2015-12-10 15:44:37', script_file => 'dscp', }, },
'29' => { new => 'ToExtension', => using => { info_url => 'http://foo.com/bar.html', script_file => 'dscp', servercheck_column_name => 'bd', servercheck_short_name => '', additional_config_json => '', id => '30', last_updated => '2015-12-10 15:44:37', name => 'OPEN', type => '33', version => '1.0.0', description => undef, isactive => '0', }, },
'30' => { new => 'ToExtension', => using => { description => undef, last_updated => '2015-12-10 15:44:37', servercheck_short_name => '', version => '1.0.0', name => 'OPEN', script_file => 'dscp', servercheck_column_name => 'be', type => '33', additional_config_json => '', id => '31', info_url => 'http://foo.com/bar.html', isactive => '0', }, },
);

sub name {
		return "ToExtension";
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
