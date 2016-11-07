package Fixtures::Job;
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
use POSIX qw(strftime);

my $now = strftime( "%Y-%m-%d %H:%M:%S", gmtime() );
my %definition_for = (
	admin_job1 => {
		new   => 'Job',
		using => {
			agent               => 1,
			keyword             => 'PURGE',
			parameters          => 'TTL:48h',
			asset_url           => 'http://cdn2.edge/foo1/.*',
			asset_type          => 'file',
			status              => 1,
			start_time          => $now,
			job_user            => 1,
			job_deliveryservice => 2,
			entered_time        => $now
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

__PACKAGE__->meta->make_immutable;

1;
