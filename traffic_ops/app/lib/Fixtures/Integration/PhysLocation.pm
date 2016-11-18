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
	## id => 1
	'0' => {
		new => 'PhysLocation',
		using => {
			name => 'cloud-east',
			phone => undef,
			region => '1',
			short_name => 'clw',
			comments => undef,
			email => undef,
			last_updated => '2015-12-10 15:43:45',
			state => '-',
			zip => '-',
			address => '-',
			city => '-',
			poc => undef,
		},
	},
	## id => 2
	'1' => {
		new => 'PhysLocation',
		using => {
			name => 'cloud-west',
			region => '2',
			city => '-',
			comments => undef,
			email => undef,
			phone => undef,
			address => '-',
			last_updated => '2015-12-10 15:43:45',
			poc => undef,
			short_name => 'cle',
			state => '-',
			zip => '-',
		},
	},
	## id => 3
	'2' => {
		new => 'PhysLocation',
		using => {
			name => 'plocation-chi-1',
			state => 'IL',
			zip => '12345',
			address => '5 Main Street',
			comments => undef,
			email => undef,
			last_updated => '2015-12-10 15:43:45',
			phone => undef,
			city => 'chi',
			poc => undef,
			region => '3',
			short_name => 'chi-1',
		},
	},
	## id => 4
	'3' => {
		new => 'PhysLocation',
		using => {
			name => 'plocation-chi-2',
			phone => undef,
			short_name => 'chi-2',
			state => 'IL',
			address => '6 Broadway',
			city => 'chi',
			comments => undef,
			last_updated => '2015-12-10 15:43:45',
			zip => '12345',
			email => undef,
			poc => undef,
			region => '3',
		},
	},
	## id => 5
	'4' => {
		new => 'PhysLocation',
		using => {
			name => 'plocation-den-1',
			region => '2',
			short_name => 'den-1',
			zip => '12345',
			city => 'den',
			comments => undef,
			last_updated => '2015-12-10 15:43:45',
			phone => undef,
			poc => undef,
			state => 'CO',
			address => '11 Main Street',
			email => undef,
		},
	},
	## id => 6
	'5' => {
		new => 'PhysLocation',
		using => {
			name => 'plocation-den-2',
			state => 'CO',
			zip => '12345',
			address => '12 Broadway',
			comments => undef,
			email => undef,
			last_updated => '2015-12-10 15:43:45',
			phone => undef,
			poc => undef,
			region => '2',
			city => 'den',
			short_name => 'den-2',
		},
	},
	## id => 7
	'6' => {
		new => 'PhysLocation',
		using => {
			name => 'plocation-hou-1',
			last_updated => '2015-12-10 15:43:45',
			phone => undef,
			poc => undef,
			address => '7 Main Street',
			city => 'hou',
			email => undef,
			region => '3',
			state => 'TX',
			comments => undef,
			short_name => 'hou-1',
			zip => '12345',
		},
	},
	## id => 8
	'7' => {
		new => 'PhysLocation',
		using => {
			name => 'plocation-hou-2',
			phone => undef,
			poc => undef,
			state => 'TX',
			zip => '12345',
			email => undef,
			city => 'hou',
			comments => undef,
			last_updated => '2015-12-10 15:43:45',
			region => '3',
			short_name => 'hou-2',
			address => '8 Broadway',
		},
	},
	## id => 9
	'8' => {
		new => 'PhysLocation',
		using => {
			name => 'plocation-lax-1',
			comments => undef,
			phone => undef,
			region => '2',
			short_name => 'lax-1',
			state => 'CA',
			address => '3 Main Street',
			city => 'lax',
			email => undef,
			last_updated => '2015-12-10 15:43:45',
			poc => undef,
			zip => '12345',
		},
	},
	## id => 10
	'9' => {
		new => 'PhysLocation',
		using => {
			name => 'plocation-lax-2',
			address => '4 Broadway',
			comments => undef,
			email => undef,
			last_updated => '2015-12-10 15:43:45',
			phone => undef,
			zip => '12345',
			city => 'lax',
			poc => undef,
			region => '2',
			short_name => 'lax-2',
			state => 'CA',
		},
	},
	## id => 11
	'10' => {
		new => 'PhysLocation',
		using => {
			name => 'plocation-nyc-1',
			region => '1',
			short_name => 'nyc-1',
			phone => undef,
			poc => undef,
			address => '1 Main Street',
			city => 'nyc',
			comments => undef,
			email => undef,
			last_updated => '2015-12-10 15:43:45',
			state => 'NY',
			zip => '12345',
		},
	},
	## id => 12
	'11' => {
		new => 'PhysLocation',
		using => {
			name => 'plocation-nyc-2',
			city => 'nyc',
			poc => undef,
			last_updated => '2015-12-10 15:43:45',
			phone => undef,
			region => '1',
			short_name => 'nyc-2',
			address => '2 Broadway',
			comments => undef,
			email => undef,
			state => 'NY',
			zip => '12345',
		},
	},
	## id => 13
	'12' => {
		new => 'PhysLocation',
		using => {
			name => 'plocation-phl-1',
			region => '1',
			address => '9 Main Street',
			email => undef,
			last_updated => '2015-12-10 15:43:45',
			phone => undef,
			poc => undef,
			short_name => 'phl-1',
			state => 'PA',
			city => 'phl',
			comments => undef,
			zip => '12345',
		},
	},
	## id => 14
	'13' => {
		new => 'PhysLocation',
		using => {
			name => 'plocation-phl-2',
			email => undef,
			state => 'PA',
			zip => '12345',
			comments => undef,
			city => 'phl',
			last_updated => '2015-12-10 15:43:45',
			phone => undef,
			poc => undef,
			region => '1',
			address => '10 Broadway',
			short_name => 'phl-2',
		},
	},
);

sub name {
		return "PhysLocation";
}

sub get_definition {
		my ( $self,
$name ) = @_;
		return $definition_for{$name};
}

sub all_fixture_names {
	# sort by db name to guarantee insertion order
	return (sort { $definition_for{$a}{using}{name} cmp $definition_for{$b}{using}{name} } keys %definition_for);
}

__PACKAGE__->meta->make_immutable;
1;
