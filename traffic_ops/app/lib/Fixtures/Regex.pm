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
			id      => 100,
			pattern => '.*\.omg-01\..*',
			type    => 19,
		},
	},
	regex_1 => {
		new   => 'Regex',
		using => {
			id      => 200,
			pattern => '.*\.foo\..*',
			type    => 19,
		},
	},
	target_filter_1 => {
		new => 'Regex',
		using => {
			id      => 300,
			pattern => '.*/force-to-one/.*',
			type => 20,
		},
	},
	target_filter_1_2 => {
		new => 'Regex',
		using => {
			id      => 400,
			pattern => '.*/force-to-one-also/.*',
			type => 20,
		},
	},
	target_filter_4 => {
		new => 'Regex',
		using => {
			id      => 500,
			pattern => '.*/go-to-four/.*',
			type => 20,
		},
	},
	target_filter_3 => {
		new => 'Regex',
		using => {
			id      => 600,
			pattern => '.*/use-three/.*',
			type => 20,
		},
	},
	hr_new_steering => {
		new => 'Regex',
		using => {
			id      => 700,
			pattern => '.*\.new-steering-ds\..*',
			type => 19,
		},
	},
	hr_steering_1 => {
		new => 'Regex',
		using => {
			id      => 800,
			pattern => '.*\.steering-ds1\..*',
			type => 19,
		},
	},
	hr_steering_2 => {
		new => 'Regex',
		using => {
			id      => 900,
			pattern => '.*\.steering-ds2\..*',
			type => 19,
		},
	},
	hr_target_1 => {
		new => 'Regex',
		using => {
			id      => 1000,
			pattern => '.*\.target-ds1\..*',
			type => 19,
		},
	},
	hr_target_2 => {
		new => 'Regex',
		using => {
			id      => 1100,
			pattern => '.*\.target-ds2\..*',
			type => 19,
		},
	},
	hr_target_3 => {
		new => 'Regex',
		using => {
			id      => 1200,
			pattern => '.*\.target-ds3\..*',
			type => 19,
		},
	},
	hr_target_4 => {
		new => 'Regex',
		using => {
			id      => 1300,
			pattern => '.*\.target-ds4\..*',
			type => 19,
		},
	},
	hr_target_5 => {
		new => 'Regex',
		using => {
			id      => 1400,
			pattern => '.*\.target-ds5\..*',
			type => 19,
		},
	},
	hr_target_6 => {
		new => 'Regex',
		using => {
			id      => 1500,
			pattern => '.*\.target-ds6\..*',
			type => 19,
		},
	},
	hr_target_7 => {
		new => 'Regex',
		using => {
			id      => 1600,
			pattern => '.*\.target-ds7\..*',
			type => 19,
		},
	},
	hr_target_8 => {
		new => 'Regex',
		using => {
			id      => 1700,
			pattern => '.*\.target-ds8\..*',
			type => 19,
		},
	},
	hr_target_9 => {
		new => 'Regex',
		using => {
			id      => 1800,
			pattern => '.*\.target-ds9\..*',
			type => 19,
		},
	},
	hr_target_10 => {
		new => 'Regex',
		using => {
			id      => 1900,
			pattern => '.*\.target-ds10\..*',
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
