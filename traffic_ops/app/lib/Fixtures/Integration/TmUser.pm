package Fixtures::Integration::TmUser;

# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
# 
#   http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.


# Do not edit! Generated code.
# See https://github.com/Comcast/traffic_control/wiki/The%20Kabletown%20example

use Moose;
extends 'DBIx::Class::EasyFixture';
use namespace::autoclean;

my %definition_for = (
	## id => 1
	'0' => {
		new => 'TmUser',
		using => {
		  username => 'admin',
		  tenant_id => undef,
		  local_passwd => '5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8',
		  uid => '1',
		  role => '1',
		  address_line1 => 'address_line1',
		  token => '',
		  full_name => 'The Admin User',
		  city => 'city',
		  postal_code => '80122',
		  phone_number => '111-111-1111',
		  address_line2 => 'address_line2',
		  registration_sent => '1999-01-01 00:00:00',
		  new_user => '1',
		  last_updated => '2015-12-10 15:43:45',
		  country => 'United States',
		  state_or_province => 'state_or_province',
		  email => 'admin@cable.comcast.com',
		  gid => '1',
		  company => undef,
		  confirm_local_passwd => '5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8',
		},
	},
	## id => 2
  '1' => {
		new => 'TmUser',
		using => {
			company => undef,
			confirm_local_passwd => '5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8',
			username => 'codebig',
			tenant_id => undef,
			gid => '1',
			country => 'United States',
			state_or_province => 'state_or_province',
			email => 'test3@email.com',
			last_updated => '2015-12-10 15:43:45',
			new_user => '1',
			address_line2 => 'address_line8',
			registration_sent => '1999-01-01 00:00:00',
			postal_code => '80124',
			phone_number => '444-444-4444',
			city => 'city',
			full_name => 'The Codebig User',
			token => '',
			uid => '1',
			role => '5',
			address_line1 => 'address_line7',
			local_passwd => '5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8',
		},
  },
	## id => 3
  '2' => {
		new => 'TmUser',
		using => {
			company => undef,
			confirm_local_passwd => '5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8',
			gid => '1',
			username => 'extension',
			tenant_id => undef,
			country => 'United States',
			email => 'plugin@email.com',
			state_or_province => 'state_or_province',
			last_updated => '2015-12-10 15:43:45',
			new_user => '1',
			registration_sent => '1999-01-01 00:00:00',
			address_line2 => 'address_line8',
			postal_code => '80124',
			phone_number => '444-444-4444',
			full_name => 'The Traffic Ops Extension User -- DO NOT REMOVE',
			token => '91504CE6-8E4A-46B2-9F9F-FE7C15228498',
			city => 'city',
			uid => '1',
			role => '5',
			address_line1 => 'address_line7',
			local_passwd => '5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8',
		},
	},
	## id => 4
	'3' => {
		new => 'TmUser',
		using => {
			confirm_local_passwd => '5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8',
			company => undef,
			username => 'migration',
			tenant_id => undef,
			gid => '1',
			state_or_province => undef,
			email => undef,
			country => undef,
			last_updated => '2015-12-10 15:43:45',
			new_user => '1',
			address_line2 => undef,
			registration_sent => '1999-01-01 00:00:00',
			phone_number => undef,
			postal_code => undef,
			city => undef,
			full_name => 'Migration User -- DO NOT REMOVE',
			token => undef,
			uid => '1',
			address_line1 => undef,
			role => '3',
			local_passwd => '5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8',
		},
	},
	## id => 5
	'4' => {
		new => 'TmUser',
		using => {
			last_updated => '2015-12-10 15:43:45',
			email => 'test1@email.com',
			state_or_province => 'state_or_province',
			country => 'United States',
			gid => '1',
			username => 'portal',
			tenant_id => undef,
			confirm_local_passwd => '5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8',
			company => undef,
			local_passwd => '5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8',
			address_line1 => 'address_line3',
			uid => '1',
			role => '5',
			full_name => 'The Portal User',
			city => 'city',
			token => '',
			phone_number => '222-222-2222',
			postal_code => '80122',
			address_line2 => 'address_line4',
			registration_sent => '1999-01-01 00:00:00',
			new_user => '1',
		},
	},
	## id => 6
	'5' => {
		new => 'TmUser',
		using => {
			country => 'United States',
			email => 'test2@email.com',
			state_or_province => 'state_or_province',
			gid => '1',
			username => 'testuser',
			tenant_id => undef,
			last_updated => '2015-12-10 15:43:45',
			company => undef,
			confirm_local_passwd => '5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8',
			city => 'city',
			token => '',
			full_name => 'The Test User',
			postal_code => '80123',
			phone_number => '333-333-3333',
			local_passwd => '5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8',
			address_line1 => 'address_line5',
			uid => '1',
			role => '4',
			new_user => '1',
			address_line2 => 'address_line4',
			registration_sent => '1999-01-01 00:00:00',
		},
	}
);

sub name {
		return "TmUser";
}

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
