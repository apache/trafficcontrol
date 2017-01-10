package Fixtures::Federation;
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
#
#
use Moose;
extends 'DBIx::Class::EasyFixture';
use namespace::autoclean;
use Digest::SHA1 qw(sha1_hex);

my %definition_for = (
	name1 => {
		new   => 'Federation',
		using => {
			id          => 1,
			cname       => 'cname1.',
			description => 'resolver4 type',
			ttl         => 86400,
		},
	},
	name2 => {
		new   => 'Federation',
		using => {
			id          => 2,
			cname       => 'cname2.',
			description => 'resolver4 type',
			ttl         => 86400,
		},
	},
	name3 => {
		new   => 'Federation',
		using => {
			id          => 3,
			cname       => 'cname3.',
			description => 'resolver4 type',
			ttl         => 86400,
		},
	},
	name4 => {
		new   => 'Federation',
		using => {
			id          => 4,
			cname       => 'cname4.',
			description => 'resolver4 type',
			ttl         => 86400,
		},
	},
);

sub get_definition {
	my ( $self, $name ) = @_;
	return $definition_for{$name};
}

sub all_fixture_names {
	# sort by db cname to guarantee insertion order
	return (sort { $definition_for{$a}{using}{id} cmp $definition_for{$b}{using}{id} } keys %definition_for);
}

__PACKAGE__->meta->make_immutable;

1;
