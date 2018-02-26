package Fixtures::TmUser;
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

my $local_passwd   = sha1_hex('password');
my %definition_for = (
	admin => {
		new   => 'TmUser',
		using => {
			id                   => 100,
			username             => 'admin',
			tenant_id            => undef,
			role                 => 4,
			uid                  => '1',
			gid                  => '1',
			local_passwd         => $local_passwd,
			confirm_local_passwd => $local_passwd,
			full_name            => 'The Admin User',
			email                => 'admin@kabletown.com',
			new_user             => '1',
			address_line1        => 'address_line1',
			address_line2        => 'address_line2',
			city                 => 'city',
			state_or_province    => 'state_or_province',
			phone_number         => '111-111-1111',
			postal_code          => '80122',
			country              => 'United States',
			token                => '',
			registration_sent    => '1999-01-01 00:00:00',
		},
	},
	portal => {
		new   => 'TmUser',
		using => {
			id                   => 200,
			username             => 'portal',
			tenant_id            => undef,
			role                 => 6,
			uid                  => '1',
			gid                  => '1',
			local_passwd         => $local_passwd,
			confirm_local_passwd => $local_passwd,
			full_name            => 'The Portal User',
			email                => 'portal@kabletown.com',
			new_user             => '1',
			address_line1        => 'address_line3',
			address_line2        => 'address_line4',
			city                 => 'city',
			state_or_province    => 'state_or_province',
			phone_number         => '222-222-2222',
			postal_code          => '80122',
			country              => 'United States',
			token                => '',
			registration_sent    => '1999-01-01 00:00:00',
		},
	},
	codebig => {
		new   => 'TmUser',
		using => {
			id                   => 300,
			username             => 'codebig',
			tenant_id            => undef,
			role                 => 6,
			uid                  => '1',
			gid                  => '1',
			local_passwd         => $local_passwd,
			confirm_local_passwd => $local_passwd,
			full_name            => 'The Codebig User',
			email                => 'codebig@kabletown.com',
			new_user             => '1',
			address_line1        => 'address_line7',
			address_line2        => 'address_line8',
			city                 => 'city',
			state_or_province    => 'state_or_province',
			phone_number         => '444-444-4444',
			postal_code          => '80124',
			country              => 'United States',
			token                => '',
			registration_sent    => '1999-01-01 00:00:00',
		},
	},
	federation => {
		new   => 'TmUser',
		using => {
			id                   => 500,
			username             => 'federation',
			tenant_id            => undef,
			role                 => 7,
			uid                  => '1',
			gid                  => '1',
			local_passwd         => $local_passwd,
			confirm_local_passwd => $local_passwd,
			full_name            => 'The federations User',
			email                => 'federation@kabletown.com',
			new_user             => '1',
			address_line1        => 'address_line1',
			address_line2        => 'address_line2',
			city                 => 'city',
			state_or_province    => 'state_or_province',
			phone_number         => '333-333-3333',
			postal_code          => '80123',
			country              => 'United States',
			token                => '',
			registration_sent    => '1999-01-01 00:00:00',
		},
	},
	steering1 => {
		new   => 'TmUser',
		using => {
			id                   => 600,
			username             => 'steering1',
			tenant_id            => undef,
			role                 => 7,
			uid                  => '1',
			gid                  => '1',
			local_passwd         => $local_passwd,
			confirm_local_passwd => $local_passwd,
			full_name            => 'The steering User 1',
			email                => 'steering1@kabletown.com',
			new_user             => '1',
			address_line1        => 'address_line1',
			address_line2        => 'address_line2',
			city                 => 'city',
			state_or_province    => 'state_or_province',
			phone_number         => '333-333-3333',
			postal_code          => '80123',
			country              => 'United States',
			token                => '',
			registration_sent    => '1999-01-01 00:00:00',
		},
	},
	steering2 => {
		new   => 'TmUser',
		using => {
			id                   => 700,
			username             => 'steering2',
			tenant_id            => undef,
			role                 => 7,
			uid                  => '1',
			gid                  => '1',
			local_passwd         => $local_passwd,
			confirm_local_passwd => $local_passwd,
			full_name            => 'The steering User 2',
			email                => 'steering2@kabletown.com',
			new_user             => '1',
			address_line1        => 'address_line1',
			address_line2        => 'address_line2',
			city                 => 'city',
			state_or_province    => 'state_or_province',
			phone_number         => '333-333-3333',
			postal_code          => '80123',
			country              => 'United States',
			token                => '',
			registration_sent    => '1999-01-01 00:00:00',
		},
	},
	
	admin_root => {
		new   => 'TmUser',
		using => {
			id                   => 800,
			username             => 'admin-root',
			tenant_id            => 10**9,
			role                 => 4,
			uid                  => '1',
			gid                  => '1',
			local_passwd         => $local_passwd,
			confirm_local_passwd => $local_passwd,
			full_name            => 'The Admin User for the "root" tenant',
			email                => 'admin-root@kabletown.com',
			new_user             => '1',
			address_line1        => 'address_line1',
			address_line2        => 'address_line2',
			city                 => 'city',
			state_or_province    => 'state_or_province',
			phone_number         => '111-111-1111',
			postal_code          => '80122',
			country              => 'United States',
			token                => '',
			registration_sent    => '1999-01-01 00:00:00',
		},
	},
	
	portal_root => {
		new   => 'TmUser',
		using => {
			id                   => 900,
			username             => 'portal-root',
			tenant_id            => 10**9,
			role                 => 6,
			uid                  => '1',
			gid                  => '1',
			local_passwd         => $local_passwd,
			confirm_local_passwd => $local_passwd,
			full_name            => 'The Portal User for the "root" tenant',
			email                => 'portal-root@kabletown.com',
			new_user             => '1',
			address_line1        => 'address_line3',
			address_line2        => 'address_line4',
			city                 => 'city',
			state_or_province    => 'state_or_province',
			phone_number         => '222-222-2222',
			postal_code          => '80122',
			country              => 'United States',
			token                => '',
			registration_sent    => '1999-01-01 00:00:00',
		},
	},

	read_only_root => {
		new   => 'TmUser',
		using => {
			id                   => 1000,
			username             => 'read-only-root',
			tenant_id            => 10**9,
			role                 => 2,
			uid                  => '1',
			gid                  => '1',
			local_passwd         => $local_passwd,
			confirm_local_passwd => $local_passwd,
			full_name            => 'The Read Only User for the "root" tenant',
			email                => 'read-only-root@kabletown.com',
			new_user             => '1',
			address_line1        => 'address_line3',
			address_line2        => 'address_line4',
			city                 => 'city',
			state_or_province    => 'state_or_province',
			phone_number         => '222-222-2222',
			postal_code          => '80122',
			country              => 'United States',
			token                => '',
			registration_sent    => '1999-01-01 00:00:00',
		},
	},

);

sub get_definition {
	my ( $self, $name ) = @_;
	return $definition_for{$name};
}

sub all_fixture_names {
	# sort by db username to guarantee insertion order
	return (sort { $definition_for{$a}{using}{username} cmp $definition_for{$b}{using}{username} } keys %definition_for);
}

__PACKAGE__->meta->make_immutable;

1;
