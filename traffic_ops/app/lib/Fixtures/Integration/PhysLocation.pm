package Fixtures::Integration::PhysLocation;

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
'0' => { new => 'PhysLocation', => using => { id => '1', region => '1', short_name => 'nyc-1', phone => undef, poc => undef, address => '1 Main Street', city => 'nyc', comments => undef, email => undef, last_updated => '2015-12-10 15:43:45', name => 'plocation-nyc-1', state => 'NY', zip => '12345', }, }, 
'1' => { new => 'PhysLocation', => using => { city => 'nyc', name => 'plocation-nyc-2', poc => undef, last_updated => '2015-12-10 15:43:45', phone => undef, region => '1', short_name => 'nyc-2', address => '2 Broadway', comments => undef, email => undef, id => '2', state => 'NY', zip => '12345', }, }, 
'2' => { new => 'PhysLocation', => using => { comments => undef, id => '3', name => 'plocation-lax-1', phone => undef, region => '2', short_name => 'lax-1', state => 'CA', address => '3 Main Street', city => 'lax', email => undef, last_updated => '2015-12-10 15:43:45', poc => undef, zip => '12345', }, }, 
'3' => { new => 'PhysLocation', => using => { address => '4 Broadway', comments => undef, email => undef, id => '4', last_updated => '2015-12-10 15:43:45', phone => undef, zip => '12345', city => 'lax', name => 'plocation-lax-2', poc => undef, region => '2', short_name => 'lax-2', state => 'CA', }, }, 
'4' => { new => 'PhysLocation', => using => { state => 'IL', zip => '12345', address => '5 Main Street', comments => undef, email => undef, last_updated => '2015-12-10 15:43:45', name => 'plocation-chi-1', phone => undef, city => 'chi', id => '5', poc => undef, region => '3', short_name => 'chi-1', }, }, 
'5' => { new => 'PhysLocation', => using => { phone => undef, short_name => 'chi-2', state => 'IL', address => '6 Broadway', city => 'chi', comments => undef, id => '6', last_updated => '2015-12-10 15:43:45', zip => '12345', email => undef, name => 'plocation-chi-2', poc => undef, region => '3', }, }, 
'6' => { new => 'PhysLocation', => using => { id => '7', last_updated => '2015-12-10 15:43:45', name => 'plocation-hou-1', phone => undef, poc => undef, address => '7 Main Street', city => 'hou', email => undef, region => '3', state => 'TX', comments => undef, short_name => 'hou-1', zip => '12345', }, }, 
'7' => { new => 'PhysLocation', => using => { phone => undef, poc => undef, state => 'TX', zip => '12345', email => undef, city => 'hou', comments => undef, id => '8', last_updated => '2015-12-10 15:43:45', name => 'plocation-hou-2', region => '3', short_name => 'hou-2', address => '8 Broadway', }, }, 
'8' => { new => 'PhysLocation', => using => { region => '1', address => '9 Main Street', email => undef, name => 'plocation-phl-1', last_updated => '2015-12-10 15:43:45', phone => undef, poc => undef, short_name => 'phl-1', state => 'PA', city => 'phl', comments => undef, id => '9', zip => '12345', }, }, 
'9' => { new => 'PhysLocation', => using => { email => undef, state => 'PA', zip => '12345', comments => undef, city => 'phl', id => '10', last_updated => '2015-12-10 15:43:45', name => 'plocation-phl-2', phone => undef, poc => undef, region => '1', address => '10 Broadway', short_name => 'phl-2', }, }, 
'10' => { new => 'PhysLocation', => using => { name => 'plocation-den-1', region => '2', short_name => 'den-1', zip => '12345', city => 'den', comments => undef, last_updated => '2015-12-10 15:43:45', phone => undef, poc => undef, state => 'CO', address => '11 Main Street', email => undef, id => '11', }, }, 
'11' => { new => 'PhysLocation', => using => { id => '12', state => 'CO', zip => '12345', address => '12 Broadway', comments => undef, email => undef, last_updated => '2015-12-10 15:43:45', name => 'plocation-den-2', phone => undef, poc => undef, region => '2', city => 'den', short_name => 'den-2', }, }, 
'12' => { new => 'PhysLocation', => using => { phone => undef, region => '1', short_name => 'clw', comments => undef, email => undef, id => '100', last_updated => '2015-12-10 15:43:45', state => '-', zip => '-', address => '-', city => '-', name => 'cloud-east', poc => undef, }, }, 
'13' => { new => 'PhysLocation', => using => { region => '2', city => '-', comments => undef, email => undef, id => '101', name => 'cloud-west', phone => undef, address => '-', last_updated => '2015-12-10 15:43:45', poc => undef, short_name => 'cle', state => '-', zip => '-', }, }, 
); 

sub name {
		return "PhysLocation";
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
