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
			deliveryservice => 8,
			server          => 1,
		},
	},
	test_ds1_server_edge13 => {
		new   => 'DeliveryserviceServer',
		using => {
			deliveryservice => 8,
			server          => 3,
		},
	},
	test_ds1_server_mid1 => {
		new   => 'DeliveryserviceServer',
		using => {
			deliveryservice => 8,
			server          => 4,
		},
	},
	test_ds2_server_edge1 => {
		new   => 'DeliveryserviceServer',
		using => {
			deliveryservice => 9,
			server          => 2,
		},
	},
	test_ds2_server_mid1 => {
		new   => 'DeliveryserviceServer',
		using => {
			deliveryservice => 9,
			server          => 5,
		},
	},
	test_ds5_server_edge14 => {
		new   => 'DeliveryserviceServer',
		using => {
			deliveryservice => 5,
			server          => 12,
		},
	},
	test_ds5_server_edge15 => {
		new   => 'DeliveryserviceServer',
		using => {
			deliveryservice => 5,
			server          => 15,
		},
	},
	test_ds6_server_edge14 => {
		new   => 'DeliveryserviceServer',
		using => {
			deliveryservice => 6,
			server          => 14,
		},
	},
	test_ds6_server_edge15 => {
		new   => 'DeliveryserviceServer',
		using => {
			deliveryservice => 6,
			server          => 15,
		},
	},
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
