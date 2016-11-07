package Fixtures::Integration::Profile;

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
'0' => { new => 'Profile', => using => { description => 'HP DL380 Edge', id => '2', last_updated => '2015-12-10 15:43:48', name => 'EDGE2_CDN1', }, }, 
'1' => { new => 'Profile', => using => { description => 'HP DL380 Mid', id => '4', last_updated => '2015-12-10 15:43:48', name => 'MID2_CDN1', }, }, 
'2' => { new => 'Profile', => using => { description => 'Comcast Content Router for cdn1.cdn.net', id => '5', last_updated => '2015-12-10 15:43:48', name => 'CCR_CDN1', }, }, 
'3' => { new => 'Profile', => using => { description => 'GLOBAL Traffic Ops Profile -- DO NOT DELETE', id => '6', last_updated => '2015-12-10 15:43:48', name => 'GLOBAL', }, }, 
'4' => { new => 'Profile', => using => { description => 'Comcast Content Router for cdn2.comcast.net', id => '8', last_updated => '2015-12-10 15:43:48', name => 'CCR_CDN2', }, }, 
'5' => { new => 'Profile', => using => { description => 'TrafficMonitor for CDN1', id => '11', last_updated => '2015-12-10 15:43:48', name => 'RASCAL_CDN1', }, }, 
'6' => { new => 'Profile', => using => { id => '12', last_updated => '2015-12-10 15:43:48', name => 'RASCAL_CDN2', description => 'TrafficMonitor for CDN2 ', }, }, 
'7' => { new => 'Profile', => using => { last_updated => '2015-12-10 15:43:48', name => 'EDGE1_CDN2_402', description => 'Dell R720xd, Edge, CDN2 CDN, ATS v4.0.2', id => '16', }, }, 
'8' => { new => 'Profile', => using => { description => 'Dell R720xd, Edge, CDN1 CDN, ATS v4.0.2', id => '19', last_updated => '2015-12-10 15:43:48', name => 'EDGE1_CDN1_402', }, }, 
'9' => { new => 'Profile', => using => { description => 'Dell R720xd, Mid, CDN2 CDN, new vol config, ATS v4.0.x', id => '20', last_updated => '2015-12-10 15:43:48', name => 'MID1_CDN2_402', }, }, 
'10' => { new => 'Profile', => using => { description => 'HP DL380, Edge, CDN1 CDN, ATS v4.0.x', id => '21', last_updated => '2015-12-10 15:43:48', name => 'EDGE2_CDN1_402', }, }, 
'11' => { new => 'Profile', => using => { description => 'HP DL380, Edge, CDN2 CDN, ATS v4.0.x', id => '23', last_updated => '2015-12-10 15:43:48', name => 'EDGE2_CDN2_402', }, }, 
'12' => { new => 'Profile', => using => { description => 'Dell R720xd, Edge, CDN2 CDN, ATS v4.2.1, Consistent Parent', id => '26', last_updated => '2015-12-10 15:43:48', name => 'EDGE1_CDN2_421', }, }, 
'13' => { new => 'Profile', => using => { description => 'Dell R720xd, Edge, CDN1 CDN, ATS v4.2.1, Consistent Parent', id => '27', last_updated => '2015-12-10 15:43:48', name => 'EDGE1_CDN1_421', }, }, 
'14' => { new => 'Profile', => using => { description => 'Dell R720xd, Mid, CDN2 CDN, ATS v4.2.1', id => '30', last_updated => '2015-12-10 15:43:48', name => 'MID1_CDN2_421', }, }, 
'15' => { new => 'Profile', => using => { name => 'MID1_CDN1_421', description => 'Dell R720xd, Mid, CDN1 CDN, ATS v4.2.1', id => '31', last_updated => '2015-12-10 15:43:48', }, }, 
'16' => { new => 'Profile', => using => { id => '34', last_updated => '2015-12-10 15:43:48', name => 'TRAFFIC_STATS', description => 'Traffic Stats profile for all CDNs', }, }, 
'17' => { new => 'Profile', => using => { description => 'HP DL380, Edge, CDN2 CDN, ATS v4.2.1, Consistent Parent', id => '37', last_updated => '2015-12-10 15:43:48', name => 'EDGE2_CDN2_421', }, }, 
'18' => { new => 'Profile', => using => { description => 'Dell r720xd, Edge, CDN1 CDN, ATS v4.2.1, SSL enabled', id => '45', last_updated => '2015-12-10 15:43:48', name => 'EDGE1_CDN1_421_SSL', }, }, 
'19' => { new => 'Profile', => using => { description => 'Riak profile for all CDNs', id => '47', last_updated => '2015-12-10 15:43:48', name => 'RIAK_ALL', }, }, 
'20' => { new => 'Profile', => using => { last_updated => '2015-12-10 15:43:48', name => 'ORG1_CDN1', description => 'Multi site origin profile 1', id => '48', }, }, 
'21' => { new => 'Profile', => using => { name => 'ORG2_CDN1', description => 'Multi site origin profile 2', id => '49', last_updated => '2015-12-10 15:43:48', }, }, 
); 

sub name {
		return "Profile";
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
