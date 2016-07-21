package Fixtures::Type;
#
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
use Digest::SHA1 qw(sha1_hex);

my %definition_for = (
	## id => 1
	EDGE => {
		new   => 'Type',
		using => {
			name         => 'EDGE',
			description  => 'Edge Cache',
			use_in_table => 'server',
		},
	},
	## id => 2
	MID => {
		new   => 'Type',
		using => {
			name         => 'MID',
			description  => 'Mid Tier Cache',
			use_in_table => 'server',
		},
	},
	## id => 3
	ORG => {
		new   => 'Type',
		using => {
			name         => 'ORG',
			description  => 'Origin',
			use_in_table => 'server',
		},
	},
	## id => 4
	CCR => {
		new   => 'Type',
		using => {
			name         => 'CCR',
			description  => 'Kabletown Content Router',
			use_in_table => 'server',
		},
	},
	## id => 5
	EDGE_LOC => {
		new   => 'Type',
		using => {
			name         => 'EDGE_LOC',
			description  => 'Edge Cachegroup',
			use_in_table => 'cachegroup',
		},
	},
	## id => 6
	MID_LOC => {
		new   => 'Type',
		using => {
			name         => 'MID_LOC',
			description  => 'Mid Cachegroup',
			use_in_table => 'cachegroup',
		},
	},
	## id => 7
	DNS => {
		new   => 'Type',
		using => {
			name         => 'DNS',
			description  => 'DNS Content Routing',
			use_in_table => 'deliveryservice',
		},
	},
	## id => 8
	OTHER_CDN => {
		new   => 'Type',
		using => {
			name         => 'OTHER_CDN',
			description  => 'Other CDN (CDS-IS, Akamai, etc)',
			use_in_table => 'server',
		},
	},
	## id => 9
	HTTP_NO_CACHE => {
		new   => 'Type',
		using => {
			name         => 'HTTP_NO_CACHE',
			description  => 'HTTP Content Routing, no caching',
			use_in_table => 'deliveryservice',
		},
	},
	## id => 10
	HTTP_LIVE => {
		new   => 'Type',
		using => {
			name         => 'HTTP_LIVE',
			description  => 'HTTP Content routing cache in RAM ',
			use_in_table => 'deliveryservice',
		},
	},
	## id => 11
	RASCAL => {
		new   => 'Type',
		using => {
			name         => 'RASCAL',
			description  => 'Rascal health polling & reporting',
			use_in_table => 'server',
		},
	},
	## id => 12
	HOST_REGEXP => {
		new   => 'Type',
		using => {
			name         => 'HOST_REGEXP',
			description  => 'Host header regular expression',
			use_in_table => 'regex',
		},
	},
	## id => 13
	PATH_REGEXP => {
		new   => 'Type',
		using => {
			name         => 'PATH_REGEXP',
			description  => 'Path regular expression',
			use_in_table => 'regex',
		},
	},
	## id => 14
	A_RECORD => {
		new   => 'Type',
		using => {
			name         => 'A_RECORD',
			description  => 'Static DNS A entry',
			use_in_table => 'staticdnsentry',
		}
	},
	## id => 15
	AAAA_RECORD => {
		new   => 'Type',
		using => {
			name         => 'AAAA_RECORD',
			description  => 'Static DNS AAAA entry',
			use_in_table => 'staticdnsentry',
		}
	},
	## id => 16
	CNAME_RECORD => {
		new   => 'Type',
		using => {
			name         => 'CNAME_RECORD',
			description  => 'Static DNS CNAME entry',
			use_in_table => 'staticdnsentry',
		}
	},
	## id => 17
	HTTP_LIVE_NATNL => {
		new   => 'Type',
		using => {
			name         => 'HTTP_LIVE_NATNL',
			description  => 'HTTP Content routing, RAM cache, National',
			use_in_table => 'deliveryservice',
		}
	},
	## id => 18
	DNS_LIVE_NATNL => {
		new   => 'Type',
		using => {
			name         => 'DNS_LIVE_NATNL',
			description  => 'DNS Content routing, RAM cache, National',
			use_in_table => 'deliveryservice',
		}
	},
	## id => 19
	LOCAL => {
		new   => 'Type',
		using => {
			name         => 'LOCAL',
			description  => 'Local User',
			use_in_table => 'tm_user',
		}
	},
	## id => 20
	ACTIVE_DIRECTORY => {
		new   => 'Type',
		using => {
			name         => 'ACTIVE_DIRECTORY',
			description  => 'Active Directory User',
			use_in_table => 'tm_user',
		}
	},
	## id => 21
	TOOLS_SERVER => {
		new   => 'Type',
		using => {
			name         => 'TOOLS_SERVER',
			description  => 'Ops hosts for management',
			use_in_table => 'server',
		}
	},
	## id => 22
	RIAK => {
		new   => 'Type',
		using => {
			name         => 'RIAK',
			description  => 'riak type',
			use_in_table => 'server',
		}
	},
	## id => 23
	INFLUXDB => {
		new   => 'Type',
		using => {
			name         => 'INFLUXDB',
			description  => 'influxdb type',
			use_in_table => 'server',
		}
	},
	## id => 24
	RESOLVE4 => {
		new   => 'Type',
		using => {
			name         => 'RESOLVE4',
			description  => 'federation type resolve4',
			use_in_table => 'federation',
		}
	},
	## id => 25
	RESOLVE6 => {
		new   => 'Type',
		using => {
			name         => 'RESOLVE6',
			description  => 'federation type resolve6',
			use_in_table => 'federation',
		},
	},
	## id => 26
	ANY_MAP => {
		new   => 'Type',
		using => {
			name         => 'ANY_MAP',
			description  => 'any_map type',
			use_in_table => 'deliveryservice',
		}
	},
	## id => 27
	HTTP => {
		new   => 'Type',
		using => {
			name         => 'HTTP',
			description  => 'HTTP Content routing cache ',
			use_in_table => 'deliveryservice',
		},
	},
	## id => 28
	STEERING => {
		new   => 'Type',
		using => {
			name         => 'STEERING',
			description  => 'Steering Delivery Service',
			use_in_table => 'deliveryservice',
		}
	}
);

sub get_definition {
	my ( $self, $name ) = @_;
	return $definition_for{$name};
}

sub all_fixture_names {
	return keys %definition_for;
}

__PACKAGE__->meta->make_immutable;

1;
