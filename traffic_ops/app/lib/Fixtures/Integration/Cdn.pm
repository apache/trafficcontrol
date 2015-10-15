package Fixtures::Integration::Cdn;
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
use Moose;
extends 'DBIx::Class::EasyFixture';
use namespace::autoclean;
use Digest::SHA1 qw(sha1_hex);

my %definition_for = (
	cdn_number_1 => {
		new   => 'Cdn',
		using => {
			id          => 1,
			name        => 'cdn_number_1',
		},
	},
	cdn_number_2 => {
		new   => 'Cdn',
		using => {
			id          => 2,
			name        => 'cdn_number_2',
		},
	},
	cdn_number_3 => {
		new   => 'Cdn',
		using => {
			id          => 3,
			name        => 'ALL',
		},
	},
);

sub name {
	return "Cdn";
}

sub get_definition {
	my ( $self, $name ) = @_;
	return $definition_for{$name};
}

sub all_fixture_names {
	return keys %definition_for;
}

__PACKAGE__->meta->make_immutable;

1;
