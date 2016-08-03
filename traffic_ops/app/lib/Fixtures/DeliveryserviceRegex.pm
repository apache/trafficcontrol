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
	regex2 => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 9,
			regex           => 7,
			set_number      => 0,
		},
	},
	target_r1_filter => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 4,
			regex           => 1,
			set_number      => 0,
		},
	},
	regex1 => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 8,
			regex           => 5,
			set_number      => 0,
		},
	},
	target_r2_filter => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 4,
			regex           => 2,
			set_number      => 0,
		},
	},
	target_r4_filter => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 7,
			regex           => 3,
			set_number      => 0,
		},
	},
	target_r3_filter => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 6,
			regex           => 4,
			set_number      => 0,
		},
	},
	new_steering => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 3,
			regex           => 6,
			set_number      => 0,
		},
	},
	steering_1 => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 1,
			regex           => 8,
			set_number      => 0,
		},
	},
	steering_2 => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 2,
			regex           => 9,
			set_number      => 0,
		},
	},
	target_1 => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 4,
			regex           => 10,
			set_number      => 0,
		},
	},
	target_2 => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 5,
			regex           => 11,
			set_number      => 0,
		},
	},
	target_3 => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 6,
			regex           => 12,
			set_number      => 0,
		},
	},
		target_4 => {
		new   => 'DeliveryserviceRegex',
		using => {
			deliveryservice => 7,
			regex           => 13,
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
