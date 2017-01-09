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
	a_record_host => {
		new   => 'Staticdnsentry',
		using => {
			id              => 100,
			host            => 'A_RECORD_HOST',
			address         => '127.0.0.1',
			type            => 21,
			deliveryservice => 100,
			cachegroup      => 100,
		},
	},
	aaaa_record_host => {
		new   => 'Staticdnsentry',
		using => {
			id              => 200,
			host            => 'AAAA_RECORD_HOST',
			address         => '127.0.0.1',
			deliveryservice => 100,
			cachegroup      => 100,
			type            => 22,
		},
	},
	cname_host => {
		new   => 'Staticdnsentry',
		using => {
			id              => 300,
			host            => 'CNAME_HOST',
			address         => '127.0.0.1',
			deliveryservice => 200,
			type            => 23,
			cachegroup      => 200,
		},
	},
);


sub get_definition {
	my ( $self, $name ) = @_;
	return $definition_for{$name};
}

sub all_fixture_names {
	# sort by db host to guarantee insertion order
	return (sort { $definition_for{$a}{using}{id} cmp $definition_for{$b}{using}{id} } keys %definition_for);
}

__PACKAGE__->meta->make_immutable;

1;
