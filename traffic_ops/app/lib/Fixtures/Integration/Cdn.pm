package Fixtures::Integration::Cdn;

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
	'cd1' => {
		new => 'Cdn',

		using => {
			id => 1,
			name => 'CDN1',
			dnssec_enabled => '0',
			domain_name => 'cdn1.kabletown.net',
			last_updated => '2015-12-10 15:43:45',
		},
	},
	'cdn2' => {
		new => 'Cdn',
		using => {
			id => 2,
			name => 'CDN2',
			dnssec_enabled => '0',
			domain_name => 'cdn2.kabletown.net',
			last_updated => '2015-12-10 15:43:45',
		},
	},
);

sub name {
		return "Cdn";
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
