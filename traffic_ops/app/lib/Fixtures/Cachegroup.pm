package Fixtures::Cachegroup;
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
	mid_northeast => {
		new   => 'Cachegroup',
		using => {
			id                   => 100,
			name                 => 'mid-northeast-group',
			short_name           => 'ne',
			type                 => 2,
			latitude             => 120,
			longitude            => 120,
			parent_cachegroup_id => undef,
		},
	},
	mid_northwest => {
		new   => 'Cachegroup',
		using => {
			id                   => 200,
			name                 => 'mid-northwest-group',
			short_name           => 'nw',
			type                 => 2,
			latitude             => 100,
			longitude            => 100,
			parent_cachegroup_id => 100,
		},
	},
	mid_cg3 => {
		new   => 'Cachegroup',
		using => {
			id                   => 800,
			name                 => 'mid_cg3',
			short_name           => 'mid_cg3',
			type                 => 6,
			latitude             => 100,
			longitude            => 100,
			parent_cachegroup_id => undef,
		},
	},
	edge_cg4 => {
		new   => 'Cachegroup',
		using => {
			id                   => 900,
			name                 => 'edge_cg4',
			short_name           => 'edge_cg4',
			type                 => 5,
			latitude             => 100,
			longitude            => 100,
			parent_cachegroup_id => 800,
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

# sub all_fixture_names {
# 	return keys %definition_for;
# }

__PACKAGE__->meta->make_immutable;

1;
