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
	## id => 1
	disallowed => {
		new   => 'Role',
		using => {
			name        => 'disallowed',
			description => 'block all access',
			priv_level  => 0,
		},
	},
	## id => 2
	read_only => {
		new   => 'Role',
		using => {
			name        => 'read-only user',
			description => 'block all access',
			priv_level  => 10,
		},
	},
	## id => 3
	federation => {
		new   => 'Role',
		using => {
			name        => 'federation',
			description => 'Role for Secondary CZF',
			priv_level  => 11,
		},
	},
	## id => 4
	operations => {
		new   => 'Role',
		using => {
			name        => 'operations',
			description => 'block all access',
			priv_level  => 20,
		},
	},
	## id => 5
	admin => {
		new   => 'Role',
		using => {
			name        => 'admin',
			description => 'super-user',
			priv_level  => 30,
		},
	},
	## id => 6
	migrations => {
		new   => 'Role',
		using => {
			name        => 'migrations',
			description => 'database migrations user - DO NOT REMOVE',
			priv_level  => 20,
		},
	},
	## id => 7
	portal => {
		new   => 'Role',
		using => {
			name        => 'portal',
			description => 'Portal User',
			priv_level  => 2,
		},
	},
	## id => 8
	steering => {
		new   => 'Role',
		using => {
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
	return keys %definition_for;
}

__PACKAGE__->meta->make_immutable;

1;
