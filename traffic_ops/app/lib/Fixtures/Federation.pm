package Fixtures::Federation;
#
# Copyright 2015 Comcast Cable Communications Management, LLC
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
			name        => 'name1',
			description => 'resolver4 type',
			cname       => 'cname1',
			ttl         => 86400,
			type        => 33,
		},
	},
	name2 => {
		new   => 'Federation',
		using => {
			id          => 2,
			name        => 'name2',
			description => 'resolver4 type',
			cname       => 'cname2',
			ttl         => 86400,
			type        => 33,
		},
	},
	name3 => {
		new   => 'Federation',
		using => {
			id          => 3,
			name        => 'name3',
			description => 'resolver6 type',
			cname       => 'cname3',
			ttl         => 86400,
			type        => 34,
		},
	},
	name3 => {
		new   => 'Federation',
		using => {
			id          => 4,
			name        => 'name4',
			description => 'resolver6 type',
			cname       => 'cname4',
			ttl         => 86400,
			type        => 34,
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
