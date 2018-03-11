package Fixtures::Role;
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
use Digest::SHA1 qw(sha1_hex);

my %definition_for = (
	disallowed => {
		new   => 'Role',
		using => {
			id          => 1,
			name        => 'disallowed',
			description => 'block all access',
			priv_level  => 0,
		},
	},
	read_only => {
		new   => 'Role',
		using => {
			id          => 2,
			name        => 'read-only',
			description => 'block all access',
			priv_level  => 10,
		},
	},
	federation => {
		new   => 'Role',
		using => {
			id          => 7,
			name        => 'federation',
			description => 'Role for Secondary CZF',
			priv_level  => 15,
		},
	},
	operations => {
		new   => 'Role',
		using => {
			id          => 3,
			name        => 'operations',
			description => 'block all access',
			priv_level  => 20,
		},
	},
	admin => {
		new   => 'Role',
		using => {
			id          => 4,
			name        => 'admin',
			description => 'super-user',
			priv_level  => 30,
		},
	},
	portal => {
		new   => 'Role',
		using => {
			id          => 6,
			name        => 'portal',
			description => 'Portal User',
			priv_level  => 15,
		},
	},
	steering => {
		new   => 'Role',
		using => {
			id          => 8,
			name        => 'steering',
			description => 'Role for Steering Delivery Service',
			priv_level  => 11,
		},
	},
);

sub get_definition {
	my ( $self, $name ) = @_;
	return $definition_for{$name};
}

sub all_fixture_names {
	# sort by db name to guarantee insertion order
	return (sort { $definition_for{$a}{using}{id} cmp $definition_for{$b}{using}{id} } keys %definition_for);
}

__PACKAGE__->meta->make_immutable;

1;
