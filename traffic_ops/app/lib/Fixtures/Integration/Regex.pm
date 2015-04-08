package Fixtures::Integration::Regex;
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
	# HOST_REGEXP => {
	# 	new   => 'Type',
	# 	using => {
	# 		id           => 18,
	# 		name         => 'HOST_REGEXP',
	# 		description  => 'Host header regular expression',
	# 		use_in_table => 'regex',
	# 	},
	# },
	regex_1 => {
		new   => 'Regex',
		using => {
			id      => 1,
			pattern => '.*\.movies\..*',
			type    => 18,
		},
	},
	regex_2 => {
		new   => 'Regex',
		using => {
			id      => 2,
			pattern => '.*\.images\..*',
			type    => 18,
		},
	},
	regex_3 => {
		new   => 'Regex',
		using => {
			id      => 3,
			pattern => '.*\.games\..*',
			type    => 18,
		},
	},
	regex_4 => {
		new   => 'Regex',
		using => {
			id      => 4,
			pattern => '.*\.tv\..*',
			type    => 18,
		},
	},
	regex_11 => {
		new   => 'Regex',
		using => {
			id      => 11,
			pattern => '.*\.movies\..*',
			type    => 18,
		},
	},
	regex_12 => {
		new   => 'Regex',
		using => {
			id      => 12,
			pattern => '.*\.images\..*',
			type    => 18,
		},
	},
	regex_13 => {
		new   => 'Regex',
		using => {
			id      => 13,
			pattern => '.*\.games\..*',
			type    => 18,
		},
	},
	regex_14 => {
		new   => 'Regex',
		using => {
			id      => 14,
			pattern => '.*\.tv\..*',
			type    => 18,
		},
	},

	regex_2 => {
		new   => 'Regex',
		using => {
			id      => 2,
			pattern => '.*\.images\..*',
			type    => 18,
		},
	},
	regex_3 => {
		new   => 'Regex',
		using => {
			id      => 3,
			pattern => '.*\.games\..*',
			type    => 18,
		},
	},
	regex_4 => {
		new   => 'Regex',
		using => {
			id      => 4,
			pattern => '.*\.tv\..*',
			type    => 18,
		},
	},
	regex_11 => {
		new   => 'Regex',
		using => {
			id      => 11,
			pattern => '.*\.movies\..*',
			type    => 18,
		},
	},
	regex_12 => {
		new   => 'Regex',
		using => {
			id      => 12,
			pattern => '.*\.images\..*',
			type    => 18,
		},
	},
	regex_13 => {
		new   => 'Regex',
		using => {
			id      => 13,
			pattern => '.*\.games\..*',
			type    => 18,
		},
	},
	regex_14 => {
		new   => 'Regex',
		using => {
			id      => 14,
			pattern => '.*\.tv\..*',
			type    => 18,
		},
	},
	
);

sub name {
	return "Regex";
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
