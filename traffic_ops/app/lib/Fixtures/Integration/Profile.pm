package Fixtures::Integration::Profile;
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
	'EDGE2_CDN1' => {
		new   => 'Profile',
		using => {
			id          => 2,
			name        => 'EDGE2_CDN1',
			description => 'HP DL380 Edge',
		},
	},
	'MID2_CDN1' => {
		new   => 'Profile',
		using => {
			id          => 4,
			name        => 'MID2_CDN1',
			description => 'HP DL380 Mid',
		},
	},
	'CCR_CDN1' => {
		new   => 'Profile',
		using => {
			id          => 5,
			name        => 'CCR_CDN1',
			description => 'Comcast Content Router for cdn1.cdn.net',
		},
	},
	'GLOBAL' => {
		new   => 'Profile',
		using => {
			id          => 6,
			name        => 'GLOBAL',
			description => 'GLOBAL Traffic Ops Profile -- DO NOT DELETE',
		},
	},
	'CCR_CDN2' => {
		new   => 'Profile',
		using => {
			id          => 8,
			name        => 'CCR_CDN2',
			description => 'Comcast Content Router for cdn2.comcast.net',
		},
	},
	'TRMON_CDN1' => {
		new   => 'Profile',
		using => {
			id          => 11,
			name        => 'RASCAL_CDN1',
			description => 'TrafficMonitor for CDN1',
		},
	},
	'TRMON_CDN2' => {
		new   => 'Profile',
		using => {
			id          => 12,
			name        => 'RASCAL_CDN2',
			description => 'TrafficMonitor for CDN2 ',
		},
	},
	'EDGE1_CDN2_402' => {
		new   => 'Profile',
		using => {
			id          => 16,
			name        => 'EDGE1_CDN2_402',
			description => 'Dell R720xd, Edge, CDN2 CDN, ATS v4.0.2',
		},
	},
	'EDGE1_CDN1_402' => {
		new   => 'Profile',
		using => {
			id          => 19,
			name        => 'EDGE1_CDN1_402',
			description => 'Dell R720xd, Edge, CDN1 CDN, ATS v4.0.2',
		},
	},
	'MID1_CDN2_402' => {
		new   => 'Profile',
		using => {
			id          => 20,
			name        => 'MID1_CDN2_402',
			description => 'Dell R720xd, Mid, CDN2 CDN, new vol config, ATS v4.0.x',
		},
	},
	'EDGE2_CDN1_402' => {
		new   => 'Profile',
		using => {
			id          => 21,
			name        => 'EDGE2_CDN1_402',
			description => 'HP DL380, Edge, CDN1 CDN, ATS v4.0.x',
		},
	},
	'EDGE2_CDN2_402' => {
		new   => 'Profile',
		using => {
			id          => 23,
			name        => 'EDGE2_CDN2_402',
			description => 'HP DL380, Edge, CDN2 CDN, ATS v4.0.x',
		},
	},
	'EDGE1_CDN2_421' => {
		new   => 'Profile',
		using => {
			id          => 26,
			name        => 'EDGE1_CDN2_421',
			description => 'Dell R720xd, Edge, CDN2 CDN, ATS v4.2.1, Consistent Parent',
		},
	},
	'EDGE1_CDN1_421' => {
		new   => 'Profile',
		using => {
			id          => 27,
			name        => 'EDGE1_CDN1_421',
			description => 'Dell R720xd, Edge, CDN1 CDN, ATS v4.2.1, Consistent Parent',
		},
	},
	'MID1_CDN2_421' => {
		new   => 'Profile',
		using => {
			id          => 30,
			name        => 'MID1_CDN2_421',
			description => 'Dell R720xd, Mid, CDN2 CDN, ATS v4.2.1',
		},
	},
	'MID1_CDN1_421' => {
		new   => 'Profile',
		using => {
			id          => 31,
			name        => 'MID1_CDN1_421',
			description => 'Dell R720xd, Mid, CDN1 CDN, ATS v4.2.1',
		},
	},
	'TRSTATS_ALL' => {
		new   => 'Profile',
		using => {
			id          => 34,
			name        => 'TRSTATS_ALL',
			description => 'TRSTATS (Redis) profile for all CDNs',
		},
	},
	'EDGE2_CDN2_421' => {
		new   => 'Profile',
		using => {
			id          => 37,
			name        => 'EDGE2_CDN2_421',
			description => 'HP DL380, Edge, CDN2 CDN, ATS v4.2.1, Consistent Parent',
		},
	},
	'EDGE1_CDN1_421_SSL' => {
		new   => 'Profile',
		using => {
			id          => 45,
			name        => 'EDGE1_CDN1_421_SSL',
			description => 'Dell r720xd, Edge, CDN1 CDN, ATS v4.2.1, SSL enabled',
		},
	},
	'RIAK_ALL' => {
		new   => 'Profile',
		using => {
			id          => 47,
			name        => 'RIAK_ALL',
			description => 'Riak profile for all CDNs',
		},
	},
	ORG1 => {
		new   => 'Profile',
		using => {
			id          => 48,
			name        => 'ORG1_CDN1',
			description => 'Multi site origin profile 1',
		},
	},
	ORG2 => {
		new   => 'Profile',
		using => {
			id          => 49,
			name        => 'ORG2_CDN1',
			description => 'Multi site origin profile 2',
		},
	},
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
