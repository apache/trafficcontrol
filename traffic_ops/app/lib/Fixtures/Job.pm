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

my %definition_for = (
	job1 => {
		new   => 'Job',
		using => {
			id					=> 100,
			agent               => 1,
			keyword             => 'PURGE',
			parameters          => 'TTL:48h',
			asset_url           => 'http://cdn2.edge/job1/.*',
			asset_type          => 'file',
			status              => 1,
			start_time          => strftime( "%Y-%m-%d %H:%M:%S", gmtime( time() - 1000 ) ),
			job_user            => 100,
			job_deliveryservice => 100,
			entered_time        => strftime( "%Y-%m-%d %H:%M:%S", gmtime( time() - 1000 ) ),
		},
	},
	job2 => {
		new   => 'Job',
		using => {
			id					=> 200,
			agent               => 1,
			keyword             => 'PURGE',
			parameters          => 'TTL:48h',
			asset_url           => 'http://cdn2.edge/job2/.*',
			asset_type          => 'file',
			status              => 1,
			start_time          => strftime( "%Y-%m-%d %H:%M:%S", gmtime( time() - 500 ) ),
			job_user            => 100,
			job_deliveryservice => 200,
			entered_time        => strftime( "%Y-%m-%d %H:%M:%S", gmtime( time() - 500 ) ),
		},
	},
	job3 => {
		new   => 'Job',
		using => {
			id					=> 300,
			agent               => 1,
			keyword             => 'PURGE',
			parameters          => 'TTL:48h',
			asset_url           => 'http://cdn2.edge/job3/.*',
			asset_type          => 'file',
			status              => 1,
			start_time          => strftime( "%Y-%m-%d %H:%M:%S", gmtime( time() - 200 ) ),
			job_user            => 200,
			job_deliveryservice => 100,
			entered_time        => strftime( "%Y-%m-%d %H:%M:%S", gmtime( time() - 200 ) ),
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
