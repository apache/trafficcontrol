package Fixtures::PhysLocation;
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
	denver => {
		new   => 'PhysLocation',
		using => {
			id         => 100,
			name       => 'Denver',
			short_name => 'denver',
			address    => '1234 mile high circle',
			city       => 'Denver',
			state      => 'CO',
			zip        => '80202',
			poc        => undef,
			phone      => '303-111-1111',
			email      => undef,
			comments   => undef,
			region     => 100,
		},
	},
	boulder => {
		new   => 'PhysLocation',
		using => {
			id         => 200,
			name       => 'Boulder',
			short_name => 'boulder',
			address    => '1234 green way',
			city       => 'Boulder',
			state      => 'CO',
			zip        => '80301',
			poc        => undef,
			phone      => '303-222-2222',
			email      => undef,
			comments   => undef,
			region     => 100,
		},
	},
	atlanta => {
		new   => 'PhysLocation',
		using => {
			id         => 300,
			name       => 'HotAtlanta',
			short_name => 'atlanta',
			address    => '1234 southern way',
			city       => 'Atlanta',
			state      => 'GA',
			zip        => '30301',
			poc        => undef,
			phone      => '404-222-2222',
			email      => undef,
			comments   => undef,
			region     => 100,
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
