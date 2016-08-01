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
	## id => 1
	'0' => {
		new => 'Cachegroup',
		using => {
			name => 'dc-cloudeast',
			longitude => '0',
			parent_cachegroup_id => undef,
			short_name => '0-cle',
			type => '4',
			last_updated => '2015-12-10 15:44:36',
			latitude => '0',
		},
	},
	## id => 2
	'1' => {
		new => 'Cachegroup',
		using => {
			name => 'dc-cloudwest',
			short_name => '1-clw',
			type => '4',
			last_updated => '2015-12-10 15:44:36',
			latitude => '0',
			longitude => '0',
			parent_cachegroup_id => undef,
		},
	},
	## id => 3
	'2' => {
		new => 'Cachegroup',
		using => {
			name => 'origin-east',
			short_name => '2-org-east',
			type => '25',
			last_updated => '2015-12-10 15:44:36',
			latitude => '0',
			longitude => '0',
			parent_cachegroup_id => undef,
		},
	},
	## id => 4
	'3' => {
		new => 'Cachegroup',
		using => {
			name => '3-mid-east',
			parent_cachegroup_id => '5',
			short_name => 'east',
			type => '23',
			last_updated => '2015-12-10 15:44:36',
			latitude => '0',
			longitude => '0',
		},
	},
	## id => 5
	'4' => {
		new => 'Cachegroup',
		using => {
			name => 'origin-west',
			parent_cachegroup_id => undef,
			short_name => '4-org-west',
			type => '25',
			last_updated => '2015-12-10 15:44:36',
			latitude => '0',
			longitude => '0',
		},
	},
	## id => 6
	'5' => {
		new => 'Cachegroup',
		using => {
			name => '5-mid-west',
			type => '23',
			last_updated => '2015-12-10 15:44:36',
			latitude => '0',
			longitude => '0',
			parent_cachegroup_id => '6',
			short_name => 'west',
		},
	},
	## id => 7
	'6' => {
		new => 'Cachegroup',
		using => {
			name => 'us-ca-losangeles',
			parent_cachegroup_id => '4',
			short_name => '6-lax',
			type => '14',
			last_updated => '2015-12-10 15:44:36',
			latitude => '34.05',
			longitude => '-118.25',
		},
	},
	## id => 8
	'7' => {
		new => 'Cachegroup',
		using => {
			name => 'us-co-denver',
			type => '14',
			last_updated => '2015-12-10 15:44:36',
			latitude => '39.739167',
			longitude => '-104.984722',
			parent_cachegroup_id => '4',
			short_name => '7-den',
		},
	},
	## id => 9
	'8' => {
		new => 'Cachegroup',
		using => {
			name => 'us-il-chicago',
			parent_cachegroup_id => '4',
			short_name => '8-chi',
			type => '14',
			last_updated => '2015-12-10 15:44:36',
			latitude => '41.881944',
			longitude => '-87.627778',
		},
	},
	## id => 10
	'9' => {
		new => 'Cachegroup',
		using => {
			name => 'us-ny-newyork',
			last_updated => '2015-12-10 15:44:36',
			latitude => '40.71435',
			longitude => '-74.00597',
			parent_cachegroup_id => '3',
			short_name => '9-nyc',
			type => '14',
		},
	},
	## id => 11
	'10' => {
		new => 'Cachegroup',
		using => {
			name => 'us-pa-philadelphia',
			parent_cachegroup_id => '3',
			short_name => '10-phl',
			type => '14',
			last_updated => '2015-12-10 15:44:36',
			latitude => '39.664722',
			longitude => '-75.565278',
		},
	},
	## id => 12
	'11' => {
		new => 'Cachegroup',
		using => {
			name => 'us-tx-houston',
			parent_cachegroup_id => '3',
			short_name => '11-hou',
			type => '14',
			last_updated => '2015-12-10 15:44:36',
			latitude => '29.762778',
			longitude => '-95.383056',
		},
	},
);

sub name {
		return "Cachegroup";
}

sub get_definition {
		my ( $self, $name ) = @_;
		return $definition_for{$name};
}

sub all_fixture_names {
	# sort by db short_name to guarantee insertion order
	return (sort { $definition_for{$a}{using}{short_name} cmp $definition_for{$b}{using}{short_name} } keys %definition_for);
}

__PACKAGE__->meta->make_immutable;
1;
