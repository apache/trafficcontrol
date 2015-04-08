package Fixtures::Integration::TmUser;
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

my $local_passwd   = sha1_hex('password');
my %definition_for = (
	admin => {
		new   => 'TmUser',
		using => {
			username             => 'admin',
			role                 => 4,
			uid                  => '1',
			gid                  => '1',
			local_passwd         => $local_passwd,
			confirm_local_passwd => $local_passwd,
			full_name            => 'The Admin User',
			email                => 'admin@cable.comcast.com',
			new_user             => '1',
			address_line1        => 'address_line1',
			address_line2        => 'address_line2',
			city                 => 'city',
			state_or_province    => 'state_or_province',
			phone_number         => '111-111-1111',
			postal_code          => '80122',
			country              => 'United States',
			local_user           => '1',
			token                => '',
		},
	},
	migration => {
		new   => 'TmUser',
		using => {
			username             => 'migration',
			role                 => 5,
			uid                  => '1',
			gid                  => '1',
			local_passwd         => $local_passwd,
			confirm_local_passwd => $local_passwd,
			full_name            => 'Migration User -- DO NOT REMOVE',
		},
	},
	portal => {
		new   => 'TmUser',
		using => {
			username             => 'portal',
			role                 => 6,
			uid                  => '1',
			gid                  => '1',
			local_passwd         => $local_passwd,
			confirm_local_passwd => $local_passwd,
			full_name            => 'The Portal User',
			email                => 'test1@email.com',
			new_user             => '1',
			address_line1        => 'address_line3',
			address_line2        => 'address_line4',
			city                 => 'city',
			state_or_province    => 'state_or_province',
			phone_number         => '222-222-2222',
			postal_code          => '80122',
			country              => 'United States',
			local_user           => '1',
			token                => '',
		},
	},
	testuser => {
		new   => 'TmUser',
		using => {
			username             => 'testuser',
			role                 => 3,
			uid                  => '1',
			gid                  => '1',
			local_passwd         => $local_passwd,
			confirm_local_passwd => $local_passwd,
			full_name            => 'The Test User',
			email                => 'test2@email.com',
			new_user             => '1',
			address_line1        => 'address_line5',
			address_line2        => 'address_line4',
			city                 => 'city',
			state_or_province    => 'state_or_province',
			phone_number         => '333-333-3333',
			postal_code          => '80123',
			country              => 'United States',
			local_user           => '1',
			token                => '',
		},
	},
	codebig => {
		new   => 'TmUser',
		using => {
			username             => 'codebig',
			role                 => 6,
			uid                  => '1',
			gid                  => '1',
			local_passwd         => $local_passwd,
			confirm_local_passwd => $local_passwd,
			full_name            => 'The Codebig User',
			email                => 'test3@email.com',
			new_user             => '1',
			address_line1        => 'address_line7',
			address_line2        => 'address_line8',
			city                 => 'city',
			state_or_province    => 'state_or_province',
			phone_number         => '444-444-4444',
			postal_code          => '80124',
			country              => 'United States',
			local_user           => '1',
			token                => '',
		},
	},
	plugin => {
		new   => 'TmUser',
		using => {
			username             => 'extension',
			role                 => 6,
			uid                  => '1',
			gid                  => '1',
			local_passwd         => $local_passwd,
			confirm_local_passwd => $local_passwd,
			full_name            => 'The Traffic Ops Extension User -- DO NOT REMOVE',
			email                => 'plugin@email.com',
			new_user             => '1',
			address_line1        => 'address_line7',
			address_line2        => 'address_line8',
			city                 => 'city',
			state_or_province    => 'state_or_province',
			phone_number         => '444-444-4444',
			postal_code          => '80124',
			country              => 'United States',
			local_user           => '1',
			token                => '91504CE6-8E4A-46B2-9F9F-FE7C15228498',
		},
	},
);

sub name {
	return "TmUser";
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
