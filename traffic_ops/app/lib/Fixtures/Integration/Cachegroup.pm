package Fixtures::Integration::Cachegroup;

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
'0' => { new => 'Cachegroup', => using => { name => 'mid-east', parent_cachegroup_id => '101', short_name => 'east', type => '7', id => '1', last_updated => '2015-12-10 15:44:36', latitude => '0', longitude => '0', }, }, 
'1' => { new => 'Cachegroup', => using => { type => '7', id => '2', last_updated => '2015-12-10 15:44:36', latitude => '0', longitude => '0', name => 'mid-west', parent_cachegroup_id => '102', short_name => 'west', }, }, 
'2' => { new => 'Cachegroup', => using => { short_name => 'clw', type => '4', id => '3', last_updated => '2015-12-10 15:44:36', latitude => '0', longitude => '0', name => 'dc-cloudwest', parent_cachegroup_id => undef, }, }, 
'3' => { new => 'Cachegroup', => using => { longitude => '0', name => 'dc-cloudeast', parent_cachegroup_id => undef, short_name => 'cle', type => '4', id => '5', last_updated => '2015-12-10 15:44:36', latitude => '0', }, }, 
'4' => { new => 'Cachegroup', => using => { last_updated => '2015-12-10 15:44:36', latitude => '40.71435', longitude => '-74.00597', name => 'us-ny-newyork', parent_cachegroup_id => '1', short_name => 'nyc', type => '6', id => '91', }, }, 
'5' => { new => 'Cachegroup', => using => { parent_cachegroup_id => '2', short_name => 'lax', type => '6', id => '92', last_updated => '2015-12-10 15:44:36', latitude => '34.05', longitude => '-118.25', name => 'us-ca-losangeles', }, }, 
'6' => { new => 'Cachegroup', => using => { parent_cachegroup_id => '2', short_name => 'chi', type => '6', id => '93', last_updated => '2015-12-10 15:44:36', latitude => '41.881944', longitude => '-87.627778', name => 'us-il-chicago', }, }, 
'7' => { new => 'Cachegroup', => using => { parent_cachegroup_id => '1', short_name => 'hou', type => '6', id => '94', last_updated => '2015-12-10 15:44:36', latitude => '29.762778', longitude => '-95.383056', name => 'us-tx-houston', }, }, 
'8' => { new => 'Cachegroup', => using => { name => 'us-pa-philadelphia', parent_cachegroup_id => '1', short_name => 'phl', type => '6', id => '95', last_updated => '2015-12-10 15:44:36', latitude => '39.664722', longitude => '-75.565278', }, }, 
'9' => { new => 'Cachegroup', => using => { type => '6', id => '96', last_updated => '2015-12-10 15:44:36', latitude => '39.739167', longitude => '-104.984722', name => 'us-co-denver', parent_cachegroup_id => '2', short_name => 'den', }, }, 
'10' => { new => 'Cachegroup', => using => { short_name => 'org-east', type => '36', id => '101', last_updated => '2015-12-10 15:44:36', latitude => '0', longitude => '0', name => 'origin-east', parent_cachegroup_id => undef, }, }, 
'11' => { new => 'Cachegroup', => using => { name => 'origin-west', parent_cachegroup_id => undef, short_name => 'org-west', type => '36', id => '102', last_updated => '2015-12-10 15:44:36', latitude => '0', longitude => '0', }, }, 
); 

sub name {
		return "Cachegroup";
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
