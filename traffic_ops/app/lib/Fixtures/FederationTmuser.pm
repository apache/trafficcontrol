package Fixtures::FederationTmuser;
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
	federation_tm_user1 => {
		new   => 'FederationTmuser',
		using => {
			federation => 1,
			tm_user    => 500,
			role       => 7,
		},
	},
	federation_tm_user2 => {
		new   => 'FederationTmuser',
		using => {
			federation => 2,
			tm_user    => 500,
			role       => 7,
		},
	}
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
