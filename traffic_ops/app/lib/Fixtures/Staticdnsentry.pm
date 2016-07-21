package Fixtures::Staticdnsentry;
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
	a_record_host => {
		new   => 'Staticdnsentry',
		using => {
			host            => 'A_RECORD_HOST',
			address         => '127.0.0.1',
			type            => 21,
			deliveryservice => 1,
			cachegroup      => 1,
		},
	},
	## id => 2
	aaaa_record_host => {
		new   => 'Staticdnsentry',
		using => {
			host            => 'AAAA_RECORD_HOST',
			address         => '127.0.0.1',
			deliveryservice => 1,
			cachegroup      => 1,
			type            => 22,
		},
	},
	## id => 3
	cname_host => {
		new   => 'Staticdnsentry',
		using => {
			host            => 'CNAME_HOST',
			address         => '127.0.0.1',
			deliveryservice => 2,
			type            => 23,
			cachegroup      => 2,
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
