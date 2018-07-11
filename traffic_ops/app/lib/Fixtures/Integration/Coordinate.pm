package Fixtures::Integration::Coordinate;

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
		new => 'Coordinate',
		using => {
			name => 'dc-cloudeast',
			longitude => '0',
			latitude => '0',
		},
	},
	## id => 2
	'1' => {
		new => 'Coordinate',
		using => {
			name => 'dc-cloudwest',
			latitude => '0',
			longitude => '0',
		},
	},
	## id => 3
	'2' => {
		new => 'Coordinate',
		using => {
			name => 'origin-east',
			latitude => '0',
			longitude => '0',
		},
	},
	## id => 4
	'3' => {
		new => 'Coordinate',
		using => {
			name => 'mid-east',
			latitude => '0',
			longitude => '0',
		},
	},
	## id => 5
	'4' => {
		new => 'Coordinate',
		using => {
			name => 'origin-west',
			latitude => '0',
			longitude => '0',
		},
	},
	## id => 6
	'5' => {
		new => 'Coordinate',
		using => {
			name => 'mid-west',
			latitude => '0',
			longitude => '0',
		},
	},
	## id => 7
	'6' => {
		new => 'Coordinate',
		using => {
			name => 'us-ca-losangeles',
			latitude => '34.05',
			longitude => '-118.25',
		},
	},
	## id => 8
	'7' => {
		new => 'Coordinate',
		using => {
			name => 'us-co-denver',
			latitude => '39.739167',
			longitude => '-104.984722',
		},
	},
	## id => 9
	'8' => {
		new => 'Coordinate',
		using => {
			name => 'us-il-chicago',
			latitude => '41.881944',
			longitude => '-87.627778',
		},
	},
	## id => 10
	'9' => {
		new => 'Coordinate',
		using => {
			name => 'us-ny-newyork',
			latitude => '40.71435',
			longitude => '-74.00597',
		},
	},
	## id => 11
	'10' => {
		new => 'Coordinate',
		using => {
			name => 'us-pa-philadelphia',
			latitude => '39.664722',
			longitude => '-75.565278',
		},
	},
	## id => 12
	'11' => {
		new => 'Coordinate',
		using => {
			name => 'us-tx-houston',
			latitude => '29.762778',
			longitude => '-95.383056',
		},
	},
);

sub name {
		return "Coordinate";
}

sub get_definition {
		my ( $self, $name ) = @_;
		return $definition_for{$name};
}

sub all_fixture_names {
	# sort by db short_name to guarantee insertion order
	return (sort { $definition_for{$a}{using}{name} cmp $definition_for{$b}{using}{name} } keys %definition_for);
}

__PACKAGE__->meta->make_immutable;
1;
