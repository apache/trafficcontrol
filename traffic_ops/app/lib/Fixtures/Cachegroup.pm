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
	## id => 1
	mid_northeast => {
		new   => 'Cachegroup',
		using => {
			name                 => 'cg1-mid-northeast',
			short_name           => 'cg1',
			type                 => 18,
			latitude             => 120,
			longitude            => 120,
			parent_cachegroup_id => undef,
		},
	},
	## id => 2
	mid_northwest => {
		new   => 'Cachegroup',
		using => {
			name                 => 'cg2-mid-northwest',
			short_name           => 'cg2',
			type                 => 18,
			latitude             => 100,
			longitude            => 100,
			parent_cachegroup_id => 1,
		},
	},
	## id => 3
	mid_cg3 => {
		new   => 'Cachegroup',
		using => {
			name                 => 'cg3-mid-south',
			short_name           => 'cg3',
			type                 => 19,
			latitude             => 100,
			longitude            => 100,
			parent_cachegroup_id => undef,
		},
	},
	## id => 4
	edge_cg4 => {
		new   => 'Cachegroup',
		using => {
			name                 => 'cg4-edge-southcentral',
			short_name           => 'cg4',
			type                 => 10,
			latitude             => 100,
			longitude            => 100,
			parent_cachegroup_id => 3,
		},
	},
	## id => 5
	edge_atl => {
		new   => 'Cachegroup',
		using => {
			name                           => 'cg5-edge_atl_group',
			short_name                     => 'cg5',
			type                           => 10,
			latitude                       => 120,
			longitude                      => 120,
			parent_cachegroup_id           => 1,
			secondary_parent_cachegroup_id => 2,
		},
	},
);

sub get_definition {
	my ( $self, $name ) = @_;
	return $definition_for{$name};
}

sub all_fixture_names {

	# sort by db name to guarantee insertion order
	return (sort { $definition_for{$a}{using}{name} cmp $definition_for{$b}{using}{name} } keys %definition_for);
}

# sub all_fixture_names {
# 	return keys %definition_for;
# }

__PACKAGE__->meta->make_immutable;

1;
