package Fixtures::Integration::EdgeCacheGroup;
#
# Copyright 2015 Comcast Cable Communications Management, LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
use Moose;
extends 'DBIx::Class::EasyFixture';
use namespace::autoclean;

my %definition_for = (

	# The test CDN has edge cache groups (formerly known as "Cachegroups") in the 5 largest cities in the US, and Denver ;-)
	# and 2 mid tier cache groups: east and west
	'us-bb-newyork' => {
		new   => 'Cachegroup',
		using => {
			id                 => 91,
			name               => 'us-ny-newyork',
			short_name         => 'nyc',
			latitude           => '40.71435',
			longitude          => '-74.00597',
			parent_cachegroup_id => '1',
			type               => '6',
		},
	},
	'us-ca-losangeles' => {
		new   => 'Cachegroup',
		using => {
			id                 => 92,
			name               => 'us-ca-losangeles',
			short_name         => 'lax',
			latitude           => '34.05',
			longitude          => '-118.25',
			parent_cachegroup_id => '2',
			type               => '6',
		},
	},
	'us-il-chicago' => {
		new   => 'Cachegroup',
		using => {
			id                 => 93,
			name               => 'us-il-chicago',
			short_name         => 'chi',
			latitude           => '41.881944',
			longitude          => '-87.627778',
			parent_cachegroup_id => '2',
			type               => '6',
		},
	},
	'us-tx-houston' => {
		new   => 'Cachegroup',
		using => {
			id                 => 94,
			name               => 'us-tx-houston',
			short_name         => 'hou',
			latitude           => '29.762778',
			longitude          => '-95.383056',
			parent_cachegroup_id => '1',
			type               => '6',
		},
	},
	'us-pa-philadelphia' => {
		new   => 'Cachegroup',
		using => {
			id                 => 95,
			name               => 'us-pa-philadelphia',
			short_name         => 'phl',
			latitude           => '39.664722',
			longitude          => '-75.565278',
			parent_cachegroup_id => '1',
			type               => '6',
		},
	},
	'us-co-denver' => {
		new   => 'Cachegroup',
		using => {
			id                 => 96,
			name               => 'us-co-denver',
			short_name         => 'den',
			latitude           => '39.739167',
			longitude          => '-104.984722',
			parent_cachegroup_id => '2',
			type               => '6',
		},
	},
);

sub name {
	return "EdgeCacheGroup";
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
