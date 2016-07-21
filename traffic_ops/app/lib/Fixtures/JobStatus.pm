package Fixtures::JobStatus;
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
	pending => {
		new   => 'JobStatus',
		using => {
			name        => 'PENDING',
			description => 'Job is queued, but has not been picked up by any agents yet'
		},
	},
	## id => 2
	in_progress => {
		new   => 'JobStatus',
		using => {
			name        => 'IN_PROGRESS',
			description => 'Job is being processed by agents'
		},
	},
	## id => 3
	completed => {
		new   => 'JobStatus',
		using => {
			name        => 'COMPLETED',
			description => 'Job has finished'
		},
	},
	## id => 4
	cancelled => {
		new   => 'JobStatus',
		using => {
			name        => 'CANCELLED',
			description => 'Job was cancelled'
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
