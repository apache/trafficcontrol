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
	## id => 1
	'0' => {
		new => 'Profile',
		using => {
			name => 'CCR_CDN1',
			description => 'Comcast Content Router for cdn1.cdn.net',
			cdn => 1,
			type => 'TR_PROFILE',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 2
	'1' => {
		new => 'Profile',
		using => {
			name => 'CCR_CDN2',
			description => 'Comcast Content Router for cdn2.comcast.net',
			cdn => 2,
			type => 'TR_PROFILE',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 3
	'2' => {
		new => 'Profile',
		using => {
			name => 'EDGE1_CDN1_402',
			description => 'Dell R720xd, Edge, CDN1 CDN, ATS v4.0.2',
			cdn => 1,
			type => 'ATS_PROFILE',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 4
	'3' => {
		new => 'Profile',
		using => {
			name => 'EDGE1_CDN1_421',
			description => 'Dell R720xd, Edge, CDN1 CDN, ATS v4.2.1, Consistent Parent',
			cdn => 1,
			type => 'ATS_PROFILE',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 5
	'4' => {
		new => 'Profile',
		using => {
			name => 'EDGE1_CDN1_421_SSL',
			description => 'Dell r720xd, Edge, CDN1 CDN, ATS v4.2.1, SSL enabled',
			cdn => 1,
			type => 'ATS_PROFILE',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 6
	'5' => {
		new => 'Profile',
		using => {
			name => 'EDGE1_CDN2_402',
			last_updated => '2015-12-10 15:43:48',
			cdn => 2,
			type => 'ATS_PROFILE',
			description => 'Dell R720xd, Edge, CDN2 CDN, ATS v4.0.2',
		},
	},
	## id => 7
	'6' => {
		new => 'Profile',
		using => {
			name => 'EDGE1_CDN2_421',
			type => 'ATS_PROFILE',
			description => 'Dell R720xd, Edge, CDN2 CDN, ATS v4.2.1, Consistent Parent',
			cdn => 2,
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 8
	'7' => {
		new => 'Profile',
		using => {
			name => 'EDGE2_CDN1',
			description => 'HP DL380 Edge',
			cdn => 1,
			type => 'ATS_PROFILE',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 9
	'8' => {
		new => 'Profile',
		using => {
			name => 'EDGE2_CDN1_402',
			description => 'HP DL380, Edge, CDN1 CDN, ATS v4.0.x',
			cdn => 1,
			type => 'ATS_PROFILE',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 10
	'9' => {
		new => 'Profile',
		using => {
			name => 'EDGE2_CDN2_402',
			description => 'HP DL380, Edge, CDN2 CDN, ATS v4.0.x',
			cdn => 2,
			type => 'ATS_PROFILE',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 11
	'10' => {
		new => 'Profile',
		using => {
			name => 'EDGE2_CDN2_421',
			description => 'HP DL380, Edge, CDN2 CDN, ATS v4.2.1, Consistent Parent',
			cdn => 2,
			type => 'ATS_PROFILE',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 12
	'11' => {
		new => 'Profile',
		using => {
			name => 'GLOBAL',
			description => 'GLOBAL Traffic Ops Profile -- DO NOT DELETE',
			type => 'UNK_PROFILE',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 13
	'12' => {
		new => 'Profile',
		using => {
			name => 'MID1_CDN1_421',
			description => 'Dell R720xd, Mid, CDN1 CDN, ATS v4.2.1',
			cdn => 1,
			type => 'ATS_PROFILE',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 14
	'13' => {
		new => 'Profile',
		using => {
			name => 'MID1_CDN2_402',
			description => 'Dell R720xd, Mid, CDN2 CDN, new vol config, ATS v4.0.x',
			cdn => 2,
			type => 'ATS_PROFILE',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 15
	'14' => {
		new => 'Profile',
		using => {
			name => 'MID1_CDN2_421',
			description => 'Dell R720xd, Mid, CDN2 CDN, ATS v4.2.1',
			cdn => 2,
			type => 'ATS_PROFILE',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 16
	'15' => {
		new => 'Profile',
		using => {
			name => 'MID2_CDN1',
			description => 'HP DL380 Mid',
			cdn => 1,
			type => 'ATS_PROFILE',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 17
	'16' => {
		new => 'Profile',
		using => {
			name => 'ORG1_CDN1',
			last_updated => '2015-12-10 15:43:48',
			cdn => 1,
			type => 'ORG_PROFILE',
			description => 'Multi site origin profile 1',
		},
	},
	## id => 18
	'17' => {
		new => 'Profile',
		using => {
			name => 'ORG2_CDN1',
			description => 'Multi site origin profile 2',
			cdn => 1,
			type => 'ORG_PROFILE',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 19
	'18' => {
		new => 'Profile',
		using => {
			name => 'RASCAL_CDN1',
			description => 'TrafficMonitor for CDN1',
			cdn => 1,
			type => 'TM_PROFILE',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 20
	'19' => {
		new => 'Profile',
		using => {
			name => 'RASCAL_CDN2',
			cdn => 2,
			type => 'TM_PROFILE',
			last_updated => '2015-12-10 15:43:48',
			description => 'TrafficMonitor for CDN2 ',
		},
	},
	## id => 21
	'20' => {
		new => 'Profile',
		using => {
			name => 'RIAK_ALL',
			description => 'Riak profile for all CDNs',
			type => 'RIAK_PROFILE',
			last_updated => '2015-12-10 15:43:48',
		},
	},
	## id => 22
	'21' => {
		new => 'Profile',
		using => {
			name => 'TRAFFIC_STATS',
			last_updated => '2015-12-10 15:43:48',
			description => 'Traffic Stats profile for all CDNs',
			type => 'TS_PROFILE',
		},
	},
);

sub name {
		return "Profile";
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
