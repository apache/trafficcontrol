package Fixtures::DeliveryserviceServer;
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
	test_ds1_server_edge1 => {
		new   => 'DeliveryserviceServer',
		using => {
			deliveryservice => 100,
			server          => 100,
		},
	},
	test_ds1_server_edge13 => {
		new   => 'DeliveryserviceServer',
		using => {
			deliveryservice => 100,
			server          => 1300,
		},
	},
	test_ds1_server_mid1 => {
		new   => 'DeliveryserviceServer',
		using => {
			deliveryservice => 100,
			server          => 300,
		},
	},
	test_ds2_server_edge1 => {
		new   => 'DeliveryserviceServer',
		using => {
			deliveryservice => 200,
			server          => 700,
		},
	},
	test_ds2_server_mid1 => {
		new   => 'DeliveryserviceServer',
		using => {
			deliveryservice => 200,
			server          => 800,
		},
	},
	test_ds5_server_edge14 => {
		new   => 'DeliveryserviceServer',
		using => {
			deliveryservice => 500,
			server          => 1400,
		},
	},
	test_ds5_server_edge15 => {
		new   => 'DeliveryserviceServer',
		using => {
			deliveryservice => 500,
			server          => 1500,
		},
	},
	test_ds6_server_edge14 => {
		new   => 'DeliveryserviceServer',
		using => {
			deliveryservice => 600,
			server          => 1400,
		},
	},
	test_ds6_server_edge15 => {
		new   => 'DeliveryserviceServer',
		using => {
			deliveryservice => 600,
			server          => 1500,
		},
	},
	test_steering_ds1 => {
		new   => 'DeliveryserviceServer',
		using => {
			deliveryservice => 700,
			server          => 900,
		},
	},
);

sub get_definition {
	my ( $self, $name ) = @_;
	return $definition_for{$name};
}

sub all_fixture_names {
	# sort by db name to guarantee insertion order
	return (sort { $definition_for{$a}{using}{deliveryservice} cmp $definition_for{$b}{using}{deliveryservice} } keys %definition_for);
}

__PACKAGE__->meta->make_immutable;

1;
