package Fixtures::Integration::Type;

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
'0' => { new => 'Type', => using => { last_updated => '2015-12-10 15:43:45', name => 'EDGE', use_in_table => 'server', description => 'Edge Cache', id => '1', }, }, 
'1' => { new => 'Type', => using => { use_in_table => 'server', description => 'Mid Tier Cache', id => '2', last_updated => '2015-12-10 15:43:45', name => 'MID', }, }, 
'2' => { new => 'Type', => using => { description => 'Origin', id => '3', last_updated => '2015-12-10 15:43:45', name => 'ORG', use_in_table => 'server', }, }, 
'3' => { new => 'Type', => using => { description => 'Comcast Content Router (aka Traffic Router)', id => '4', last_updated => '2015-12-10 15:43:45', name => 'CCR', use_in_table => 'server', }, }, 
'4' => { new => 'Type', => using => { description => 'Edge Cachegroup', id => '6', last_updated => '2015-12-10 15:43:45', name => 'EDGE_LOC', use_in_table => 'cachegroup', }, }, 
'5' => { new => 'Type', => using => { name => 'MID_LOC', use_in_table => 'cachegroup', description => 'Mid Cachegroup', id => '7', last_updated => '2015-12-10 15:43:45', }, }, 
'6' => { new => 'Type', => using => { id => '8', last_updated => '2015-12-10 15:43:45', name => 'HTTP', use_in_table => 'deliveryservice', description => 'HTTP Content Routing', }, }, 
'7' => { new => 'Type', => using => { last_updated => '2015-12-10 15:43:45', name => 'DNS', use_in_table => 'deliveryservice', description => 'DNS Content Routing', id => '9', }, }, 
'8' => { new => 'Type', => using => { use_in_table => 'server', description => 'Riak keystore', id => '10', last_updated => '2015-12-10 15:43:45', name => 'RIAK', }, }, 
'9' => { new => 'Type', => using => { description => 'HTTP Content Routing, no caching', id => '11', last_updated => '2015-12-10 15:43:45', name => 'HTTP_NO_CACHE', use_in_table => 'deliveryservice', }, }, 
'10' => { new => 'Type', => using => { description => 'HTTP Content routing cache in RAM', id => '13', last_updated => '2015-12-10 15:43:45', name => 'HTTP_LIVE', use_in_table => 'deliveryservice', }, }, 
'11' => { new => 'Type', => using => { description => 'Rascal (aka Traffic Monitor) server', id => '15', last_updated => '2015-12-10 15:43:45', name => 'RASCAL', use_in_table => 'server', }, }, 
'12' => { new => 'Type', => using => { description => 'Host header regular expression', id => '18', last_updated => '2015-12-10 15:43:45', name => 'HOST_REGEXP', use_in_table => 'regex', }, }, 
'13' => { new => 'Type', => using => { description => 'URL path regular expression', id => '19', last_updated => '2015-12-10 15:43:45', name => 'PATH_REGEXP', use_in_table => 'regex', }, }, 
'14' => { new => 'Type', => using => { description => 'HTTP header regular expression', id => '20', last_updated => '2015-12-10 15:43:45', name => 'HEADER_REGEXP', use_in_table => 'regex', }, }, 
'15' => { new => 'Type', => using => { name => 'A_RECORD', use_in_table => 'staticdnsentry', description => 'Static DNS A entry', id => '21', last_updated => '2015-12-10 15:43:45', }, }, 
'16' => { new => 'Type', => using => { id => '22', last_updated => '2015-12-10 15:43:45', name => 'AAAA_RECORD', use_in_table => 'staticdnsentry', description => 'Static DNS AAAA entry', }, }, 
'17' => { new => 'Type', => using => { description => 'Static DNS CNAME entry', id => '23', last_updated => '2015-12-10 15:43:45', name => 'CNAME_RECORD', use_in_table => 'staticdnsentry', }, }, 
'18' => { new => 'Type', => using => { last_updated => '2015-12-10 15:43:45', name => 'HTTP_LIVE_NATNL', use_in_table => 'deliveryservice', description => 'HTTP Content routing, RAM cache, National', id => '24', }, }, 
'19' => { new => 'Type', => using => { use_in_table => 'server', description => 'traffic stats server', id => '25', last_updated => '2015-12-10 15:43:45', name => 'TRAFFIC_STATS', }, }, 
'20' => { new => 'Type', => using => { description => 'DNS Content routing, RAM cache, National', id => '26', last_updated => '2015-12-10 15:43:45', name => 'DNS_LIVE_NATNL', use_in_table => 'deliveryservice', }, }, 
'21' => { new => 'Type', => using => { description => 'DNS Content routing, RAM cache, Lo
cal', id => '27', last_updated => '2015-12-10 15:43:45', name => 'DNS_LIVE', use_in_table => 'deliveryservice', }, }, 
'22' => { new => 'Type', => using => { description => 'Local User', id => '28', last_updated => '2015-12-10 15:43:45', name => 'LOCAL', use_in_table => 'tm_user', }, }, 
'23' => { new => 'Type', => using => { description => 'Active Directory User', id => '29', last_updated => '2015-12-10 15:43:45', name => 'ACTIVE_DIRECTORY', use_in_table => 'tm_user', }, }, 
'24' => { new => 'Type', => using => { name => 'TOOLS_SERVER', use_in_table => 'server', description => 'Ops hosts for managment ', id => '30', last_updated => '2015-12-10 15:43:45', }, }, 
'25' => { new => 'Type', => using => { id => '31', last_updated => '2015-12-10 15:43:45', name => 'CHECK_EXTENSION_BOOL', use_in_table => 'to_extension', description => 'TO Extension for checkmark in Server Check', }, }, 
'26' => { new => 'Type', => using => { last_updated => '2015-12-10 15:43:45', name => 'CHECK_EXTENSION_NUM', use_in_table => 'to_extension', description => 'TO Extenstion for int value in Server Check', id => '32', }, }, 
'27' => { new => 'Type', => using => { name => 'CHECK_EXTENSION_OPEN_SLOT', use_in_table => 'to_extension', description => 'Open slot for check in Server Status', id => '33', last_updated => '2015-12-10 15:43:45', }, }, 
'28' => { new => 'Type', => using => { description => 'Extension for additional configuration file', id => '34', last_updated => '2015-12-10 15:43:45', name => 'CONFIG_EXTENSION', use_in_table => 'to_extension', }, }, 
'29' => { new => 'Type', => using => { use_in_table => 'to_extension', description => 'Extension source for 12M graphs', id => '35', last_updated => '2015-12-10 15:43:45', name => 'STATISTIC_EXTENSION', }, }, 
'30' => { new => 'Type', => using => { description => 'Multi Site Origin "Cachegroup"', id => '36', last_updated => '2015-12-10 15:43:45', name => 'ORG_LOC', use_in_table => 'cachegroup', }, }, 
); 

sub name {
		return "Type";
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
