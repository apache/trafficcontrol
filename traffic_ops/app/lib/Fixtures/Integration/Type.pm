package Fixtures::Integration::Type;

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
		new => 'Type',
		using => {
			name => 'A_RECORD',
			use_in_table => 'staticdnsentry',
			description => 'Static DNS A entry',
			last_updated => '2015-12-10 15:43:45',
		},
	},
	## id => 2
	'1' => {
		new => 'Type',
		using => {
			last_updated => '2015-12-10 15:43:45',
			name => 'AAAA_RECORD',
			use_in_table => 'staticdnsentry',
			description => 'Static DNS AAAA entry',
		},
	},
	## id => 3
	'2' => {
		new => 'Type',
		using => {
			description => 'Active Directory User',
			last_updated => '2015-12-10 15:43:45',
			name => 'ACTIVE_DIRECTORY',
			use_in_table => 'tm_user',
		},
	},
	## id => 4
	'3' => {
		new => 'Type',
		using => {
			description => 'Comcast Content Router (aka Traffic Router)',
			last_updated => '2015-12-10 15:43:45',
			name => 'CCR',
			use_in_table => 'server',
		},
	},
	## id => 5
	'4' => {
		new => 'Type',
		using => {
			last_updated => '2015-12-10 15:43:45',
			name => 'CHECK_EXTENSION_BOOL',
			use_in_table => 'to_extension',
			description => 'TO Extension for checkmark in Server Check',
		},
	},
	## id => 6
	'5' => {
		new => 'Type',
		using => {
			last_updated => '2015-12-10 15:43:45',
			name => 'CHECK_EXTENSION_NUM',
			use_in_table => 'to_extension',
			description => 'TO Extenstion for int value in Server Check',
		},
	},
	## id => 7
	'6' => {
		new => 'Type',
		using => {
			name => 'CHECK_EXTENSION_OPEN_SLOT',
			use_in_table => 'to_extension',
			description => 'Open slot for check in Server Status',
			last_updated => '2015-12-10 15:43:45',
		},
	},
	## id => 8
	'7' => {
		new => 'Type',
		using => {
			description => 'Static DNS CNAME entry',
			last_updated => '2015-12-10 15:43:45',
			name => 'CNAME_RECORD',
			use_in_table => 'staticdnsentry',
		},
	},
	## id => 9
	'8' => {
		new => 'Type',
		using => {
			description => 'Extension for additional configuration file',
			last_updated => '2015-12-10 15:43:45',
			name => 'CONFIG_EXTENSION',
			use_in_table => 'to_extension',
		},
	},
	## id => 10
	'9' => {
		new => 'Type',
		using => {
			last_updated => '2015-12-10 15:43:45',
			name => 'DNS',
			use_in_table => 'deliveryservice',
			description => 'DNS Content Routing',
		},
	},
	## id => 11
	'10' => {
		new => 'Type',
		using => {
			description => 'DNS Content routing, RAM cache, Local',
			last_updated => '2015-12-10 15:43:45',
			name => 'DNS_LIVE',
			use_in_table => 'deliveryservice',
		},
	},
	## id => 12
	'11' => {
		new => 'Type',
		using => {
			description => 'DNS Content routing, RAM cache, National',
			last_updated => '2015-12-10 15:43:45',
			name => 'DNS_LIVE_NATNL',
			use_in_table => 'deliveryservice',
		},
	},
	## id => 13
	'12' => {
		new => 'Type',
		using => {
			last_updated => '2015-12-10 15:43:45',
			name => 'EDGE',
			use_in_table => 'server',
			description => 'Edge Cache',
		},
	},
	## id => 14
	'13' => {
		new => 'Type',
		using => {
			description => 'Edge Cachegroup',
			last_updated => '2015-12-10 15:43:45',
			name => 'EDGE_LOC',
			use_in_table => 'cachegroup',
		},
	},
	## id => 15
	'14' => {
		new => 'Type',
		using => {
			description => 'HTTP header regular expression',
			last_updated => '2015-12-10 15:43:45',
			name => 'HEADER_REGEXP',
			use_in_table => 'regex',
		},
	},
	## id => 16
	'15' => {
		new => 'Type',
		using => {
			last_updated => '2015-12-10 15:43:45',
			name => 'HTTP',
			use_in_table => 'deliveryservice',
			description => 'HTTP Content Routing',
		},
	},
	## id => 17
	'16' => {
		new => 'Type',
		using => {
			description => 'HTTP Content routing cache in RAM',
			last_updated => '2015-12-10 15:43:45',
			name => 'HTTP_LIVE',
			use_in_table => 'deliveryservice',
		},
	},
	## id => 18
	'17' => {
		new => 'Type',
		using => {
			last_updated => '2015-12-10 15:43:45',
			name => 'HTTP_LIVE_NATNL',
			use_in_table => 'deliveryservice',
			description => 'HTTP Content routing, RAM cache, National',
		},
	},
	## id => 19
	'18' => {
		new => 'Type',
		using => {
			description => 'HTTP Content Routing, no caching',
			last_updated => '2015-12-10 15:43:45',
			name => 'HTTP_NO_CACHE',
			use_in_table => 'deliveryservice',
		},
	},
	## id => 20
	'19' => {
		new => 'Type',
		using => {
			description => 'Host header regular expression',
			last_updated => '2015-12-10 15:43:45',
			name => 'HOST_REGEXP',
			use_in_table => 'regex',
		},
	},
	## id => 21
	'20' => {
		new => 'Type',
		using => {
			description => 'Local User',
			last_updated => '2015-12-10 15:43:45',
			name => 'LOCAL',
			use_in_table => 'tm_user',
		},
	},
	## id => 22
	'21' => {
		new => 'Type',
		using => {
			use_in_table => 'server',
			description => 'Mid Tier Cache',
			last_updated => '2015-12-10 15:43:45',
			name => 'MID',
		},
	},
	## id => 23
	'22' => {
		new => 'Type',
		using => {
			name => 'MID_LOC',
			use_in_table => 'cachegroup',
			description => 'Mid Cachegroup',
			last_updated => '2015-12-10 15:43:45',
		},
	},
	## id => 24
	'23' => {
		new => 'Type',
		using => {
			description => 'Origin',
			last_updated => '2015-12-10 15:43:45',
			name => 'ORG',
			use_in_table => 'server',
		},
	},
	## id => 25
	'24' => {
		new => 'Type',
		using => {
			description => 'Multi Site Origin "Cachegroup"',
			last_updated => '2015-12-10 15:43:45',
			name => 'ORG_LOC',
			use_in_table => 'cachegroup',
		},
	},
	## id => 26
	'25' => {
		new => 'Type',
		using => {
			description => 'URL path regular expression',
			last_updated => '2015-12-10 15:43:45',
			name => 'PATH_REGEXP',
			use_in_table => 'regex',
		},
	},
	## id => 27
	'26' => {
		new => 'Type',
		using => {
			description => 'Rascal (aka Traffic Monitor) server',
			last_updated => '2015-12-10 15:43:45',
			name => 'RASCAL',
			use_in_table => 'server',
		},
	},
	## id => 28
	'27' => {
		new => 'Type',
		using => {
			use_in_table => 'server',
			description => 'Riak keystore',
			last_updated => '2015-12-10 15:43:45',
			name => 'RIAK',
		},
	},
	## id => 29
	'28' => {
		new => 'Type',
		using => {
			use_in_table => 'to_extension',
			description => 'Extension source for 12M graphs',
			last_updated => '2015-12-10 15:43:45',
			name => 'STATISTIC_EXTENSION',
		},
	},
	## id => 30
	'29' => {
		new => 'Type',
		using => {
			name => 'TOOLS_SERVER',
			use_in_table => 'server',
			description => 'Ops hosts for managment ',
			last_updated => '2015-12-10 15:43:45',
		},
	},
	## id => 31
	'30' => {
		new => 'Type',
		using => {
			use_in_table => 'server',
			description => 'traffic stats server',
			last_updated => '2015-12-10 15:43:45',
			name => 'TRAFFIC_STATS',
		},
	},
);

sub name {
		return "Type";
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
