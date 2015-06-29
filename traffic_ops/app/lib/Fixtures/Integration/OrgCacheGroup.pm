package Fixtures::Integration::OrgCacheGroup;
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

	# The test CDN has  6 edge cache groups (formerly known as "Cachegroups"): the 5 largest cities in the US, and Denver ;-)
	# The test CDN has 2 mid tier cache groups: east and west
	# The test CDN has 2 "cloud locations", where the TrafficRouter and TrafficMonitor (and TrafficOps) reside
	# The test CDN has 2 "origin locations", where the multi site origins are hosted
	'org-east' => {
		new   => 'Cachegroup',
		using => {
			id         => 101,
			name       => 'origin-east',
			short_name => 'org-east',
			latitude   => '0',
			longitude  => '0',
			type       => '36',
		},
	},
	'org-west' => {
		new   => 'Cachegroup',
		using => {
			id         => 102,
			name       => 'origin-west',
			short_name => 'org-west',
			latitude   => '0',
			longitude  => '0',
			type       => '36',
		},
	},
);

sub name {
	return "OrgCacheGroup";
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
