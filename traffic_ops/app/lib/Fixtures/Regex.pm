package Fixtures::Regex;
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

my %definition_for = (
	regex_omg01 => {
		new   => 'Regex',
		using => {
			id      => 1,
			pattern => '.*\.omg-01\..*',
			type    => 19,
		},
	},
	regex_1 => {
		new   => 'Regex',
		using => {
			id      => 2,
			pattern => '.*\.foo\..*',
			type    => 19,
		},
	},
	target_filter_1 => {
		new => 'Regex',
		using => {
			id      => 3,
			pattern => '.*/force-to-one/.*',
			type => 20,
		},
	},
	target_filter_1_2 => {
		new => 'Regex',
		using => {
			id      => 4,
			pattern => '.*/force-to-one-also/.*',
			type => 20,
		},
	},
	target_filter_4 => {
		new => 'Regex',
		using => {
			id      => 5,
			pattern => '.*/go-to-four/.*',
			type => 20,
		},
	},
	target_filter_3 => {
		new => 'Regex',
		using => {
			id      => 6,
			pattern => '.*/use-three/.*',
			type => 20,
		},
	},
	hr_new_steering => {
		new => 'Regex',
		using => {
			id      => 7,
			pattern => '.*\.new-steering-ds\..*',
			type => 19,
		},
	},
	hr_steering_1 => {
		new => 'Regex',
		using => {
			id      => 8,
			pattern => '.*\.steering-ds1\..*',
			type => 19,
		},
	},
	hr_steering_2 => {
		new => 'Regex',
		using => {
			id      => 9,
			pattern => '.*\.steering-ds2\..*',
			type => 19,
		},
	},
	hr_target_1 => {
		new => 'Regex',
		using => {
			id      => 10,
			pattern => '.*\.target-ds1\..*',
			type => 19,
		},
	},
	hr_target_2 => {
		new => 'Regex',
		using => {
			id      => 11,
			pattern => '.*\.target-ds2\..*',
			type => 19,
		},
	},
	hr_target_3 => {
		new => 'Regex',
		using => {
			id      => 12,
			pattern => '.*\.target-ds3\..*',
			type => 19,
		},
	},
	hr_target_4 => {
		new => 'Regex',
		using => {
			id      => 13,
			pattern => '.*\.target-ds4\..*',
			type => 19,
		},
	},
);

sub get_definition {
	my ( $self, $name ) = @_;
	return $definition_for{$name};
}

sub all_fixture_names {
	# sort by db pattern to guarantee insertion order
	return (sort { $definition_for{$a}{using}{id} cmp $definition_for{$b}{using}{id} } keys %definition_for);
}

__PACKAGE__->meta->make_immutable;

1;
