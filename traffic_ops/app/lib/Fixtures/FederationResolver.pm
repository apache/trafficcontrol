package Fixtures::FederationResolver;
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
	ipv4_resolver1 => {
		new   => 'FederationResolver',
		using => {
			id         => 1,
			ip_address => "127.0.0.1/32",
			type       => 33,
		},
	},
	ipv4_resolver2 => {
		new   => 'FederationResolver',
		using => {
			id         => 2,
			ip_address => "127.0.0.2/32",
			type       => 33,
		},
	},
	ipv6_resolver1 => {
		new   => 'FederationResolver',
		using => {
			id         => 3,
			ip_address => "FE80::0202:B3FF:FE1E:8329/128",
			type       => 34,
		},
	},
	ipv6_resolver2 => {
		new   => 'FederationResolver',
		using => {
			id         => 4,
			ip_address => "FE80::0202:B3FF:FE1E:8330/128",
			type       => 34,
		},
	},
);

sub get_definition {
	my ( $self, $name ) = @_;
	return $definition_for{$name};
}

sub all_fixture_names {
	# sort by db ip_address to guarantee insertion order
	return (sort { $definition_for{$a}{using}{id} cmp $definition_for{$b}{using}{id} } keys %definition_for);
}

__PACKAGE__->meta->make_immutable;

1;
