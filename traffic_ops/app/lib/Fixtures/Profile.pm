package Fixtures::Profile;
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
	EDGE1 => {
		new   => 'Profile',
		using => {
			id          => 100,
			name        => 'EDGE1',
			description => 'edge description',
			cdn         => 100,
			type        => 'ATS_PROFILE',
		},
	},
	MID1 => {
		new   => 'Profile',
		using => {
			id          => 200,
			name        => 'MID1',
			description => 'mid description',
			cdn         => 100,
			type        => 'ATS_PROFILE',
		},
	},
	CCR1 => {
		new   => 'Profile',
		using => {
			id          => 300,
			name        => 'CCR1',
			description => 'ccr description',
			cdn         => 100,
			type        => 'TR_PROFILE',
		},
	},
	CCR2 => {
		new   => 'Profile',
		using => {
			id          => 301,
			name        => 'CCR2',
			description => 'ccr description',
			cdn         => 200,
			type        => 'TR_PROFILE',
		},
	},
	RIAK1 => {
		new   => 'Profile',
		using => {
			id          => 500,
			name        => 'RIAK1',
			description => 'riak description',
			cdn         => 100,
			type        => 'RIAK_PROFILE',
		},
	},
	RASCAL1 => {
		new   => 'Profile',
		using => {
			id          => 600,
			name        => 'RASCAL1',
			description => 'rascal description',
			cdn         => 100,
			type        => 'TM_PROFILE',
		},
	},
	RASCAL2 => {
		new   => 'Profile',
		using => {
			id          => 700,
			name        => 'RASCAL2',
			description => 'rascal2 description',
			cdn         => 200,
			type        => 'TM_PROFILE',
		},
	},
	MISC => {
		new   => 'Profile',
		using => {
			id          => 8,
			name        => 'MISC',
			description => 'misc profile description',
			type        => 'UNK_PROFILE',
		},
	},
	EDGE2 => {
		new   => 'Profile',
		using => {
			id          => 900,
			name        => 'EDGE2',
			description => 'edge description',
			cdn         => 200,
			type        => 'ATS_PROFILE',
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
