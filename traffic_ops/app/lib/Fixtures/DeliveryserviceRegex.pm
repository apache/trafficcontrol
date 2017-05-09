package Fixtures::DeliveryserviceRegex;
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
	regex1 => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 100,
			regex           => 200,
			set_number      => 0,
		},
	},
	regex2 => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 200,
			regex           => 100,
			set_number      => 0,
		},
	},
	target_r1_filter => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 400,
			regex           => 100,
			set_number      => 0,
		},
	},
	target_r2_filter => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 400,
			regex           => 200,
			set_number      => 0,
		},
	},
	target_r4_filter => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 700,
			regex           => 300,
			set_number      => 0,
		},
	},
	target_r3_filter => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 600,
			regex           => 400,
			set_number      => 0,
		},
	},
	new_steering => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 300,
			regex           => 600,
			set_number      => 0,
		},
	},
	steering_1 => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 100,
			regex           => 800,
			set_number      => 1,
		},
	},
	steering_2 => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 200,
			regex           => 900,
			set_number      => 0,
		},
	},
	target_1 => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 400,
			regex           => 1000,
			set_number      => 0,
		},
	},
	target_2 => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 500,
			regex           => 1100,
			set_number      => 0,
		},
	},
	target_3 => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 600,
			regex           => 1200,
			set_number      => 0,
		},
	},
	target_4 => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 700,
			regex           => 1300,
			set_number      => 0,
		},
	},
	target_5 => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 800,
			regex           => 1400,
			set_number      => 0,
		},
	},
	target_6 => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 900,
			regex           => 1500,
			set_number      => 0,
		},
	},
	target_7 => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 1000,
			regex           => 1600,
			set_number      => 0,
		},
	},
	target_8 => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 1100,
			regex           => 1700,
			set_number      => 0,
		},
	},
	target_9 => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 1200,
			regex           => 1800,
			set_number      => 0,
		},
	},
	target_10 => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 1300,
			regex           => 1900,
			set_number      => 0,
		},
	},

);

sub get_definition {
	my ( $self, $name ) = @_;
	return $definition_for{$name};
}

sub all_fixture_names {
	# sort by db regex to guarantee insertion order
	return (sort { $definition_for{$a}{using}{regex} cmp $definition_for{$b}{using}{regex} } keys %definition_for);
}

__PACKAGE__->meta->make_immutable;

1;
