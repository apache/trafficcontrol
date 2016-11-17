package Fixtures::DeliveryserviceTmuser;
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
	admin => {
		new   => 'DeliveryserviceTmuser',
		using => {
			deliveryservice => 100,
			tm_user_id      => 100,
		},
	},
	portal_ds1 => {
		new   => 'DeliveryserviceTmuser',
		using => {
			deliveryservice => 100,
			tm_user_id      => 200,
		},
	},
	ds_steering_user1 => {
		new   => 'DeliveryserviceTmuser',
		using => {
			deliveryservice => 700,
			tm_user_id      => 700,
		},
	},
	ds_steering_user2 => {
		new   => 'DeliveryserviceTmuser',
		using => {
			deliveryservice => 800,
			tm_user_id      => 700,
		},
	},
);

sub get_definition {
	my ( $self, $name ) = @_;
	return $definition_for{$name};
}

sub all_fixture_names {
	# sort by db deliveryservice to guarantee insertion order
	return (sort { $definition_for{$a}{using}{deliveryservice} cmp $definition_for{$b}{using}{deliveryservice} } keys %definition_for);
}
__PACKAGE__->meta->make_immutable;

1;
