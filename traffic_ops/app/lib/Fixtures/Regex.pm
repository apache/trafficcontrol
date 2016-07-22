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
	## id => 1
	target_filter_1 => {
		new => 'Regex',
		using => {
			pattern => '.*/force-to-one/.*',
			type => 28,
		},
	},
	## id => 2
	target_filter_1_2 => {
		new => 'Regex',
		using => {
			pattern => '.*/force-to-one-also/.*',
			type => 28,
		},
	},
	## id => 3
	target_filter_4 => {
		new => 'Regex',
		using => {
			pattern => '.*/go-to-four/.*',
			type => 28,
		},
	},
	## id => 4
	target_filter_3 => {
		new => 'Regex',
		using => {
			pattern => '.*/use-three/.*',
			type => 28,
		},
	},
	## id => 5
	regex_1 => {
		new   => 'Regex',
		using => {
			pattern => '.*\.foo\..*',
			type    => 15,
		},
	},
	## id => 6
	hr_new_steering => {
		new => 'Regex',
		using => {
			pattern => '.*\.new-steering-ds\..*',
			type => 15,
		},
	},
	## id => 7
	regex_omg01 => {
		new   => 'Regex',
		using => {
			pattern => '.*\.omg-01\..*',
			type    => 15,
		},
	},
	## id => 8
	hr_steering_1 => {
		new => 'Regex',
		using => {
			pattern => '.*\.steering-ds1\..*',
			type => 15,
		},
	},
	## id => 9
	hr_steering_2 => {
		new => 'Regex',
		using => {
			pattern => '.*\.steering-ds2\..*',
			type => 15,
		},
	},
	## id => 10
	hr_target_1 => {
		new => 'Regex',
		using => {
			pattern => '.*\.target-ds1\..*',
			type => 15,
		},
	},
	## id => 11
	hr_target_2 => {
		new => 'Regex',
		using => {
			pattern => '.*\.target-ds2\..*',
			type => 15,
		},
	},
	## id => 12
	hr_target_3 => {
		new => 'Regex',
		using => {
			pattern => '.*\.target-ds3\..*',
			type => 15,
		},
	},
	## id => 13
	hr_target_4 => {
		new => 'Regex',
		using => {
			pattern => '.*\.target-ds4\..*',
			type => 15,
		},
	},
);

sub get_definition {
	my ( $self, $name ) = @_;
	return $definition_for{$name};
}

sub all_fixture_names {
	# sort by db pattern to guarantee insertion order
	return (sort { $definition_for{$a}{using}{pattern} cmp $definition_for{$b}{using}{pattern} } keys %definition_for);
}

__PACKAGE__->meta->make_immutable;

1;
