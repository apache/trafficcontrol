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
'0' => { new => 'TmUser', => using => { address_line2 => 'address_line4', local_passwd => '5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8', phone_number => '333-333-3333', role => '3', state_or_province => 'state_or_province', company => undef, id => '71', new_user => '1', postal_code => '80123', full_name => 'The Test User', gid => '1', registration_sent => '1999-01-01 00:00:00', username => 'testuser', address_line1 => 'address_line5', city => 'city', confirm_local_passwd => '5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8', email => 'test2@email.com', uid => '1', country => 'United States', last_updated => '2015-12-10 15:43:45', token => '', }, }, 
'1' => { new => 'TmUser', => using => { last_updated => '2015-12-10 15:43:45', country => 'United States', email => 'test3@email.com', full_name => 'The Codebig User', id => '72', phone_number => '444-444-4444', registration_sent => '1999-01-01 00:00:00', token => '', username => 'codebig', address_line2 => 'address_line8', city => 'city', local_passwd => '5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8', new_user => '1', role => '6', address_line1 => 'address_line7', company => undef, confirm_local_passwd => '5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8', gid => '1', postal_code => '80124', state_or_province => 'state_or_province', uid => '1', }, }, 
'2' => { new => 'TmUser', => using => { city => 'city', company => undef, country => 'United States', postal_code => '80124', address_line1 => 'address_line7', address_line2 => 'address_line8', role => '6', token => '91504CE6-8E4A-46B2-9F9F-FE7C15228498', uid => '1', username => 'extension', full_name => 'The Traffic Ops Extension User -- DO NOT REMOVE', local_passwd => '5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8', last_updated => '2015-12-10 15:43:45', phone_number => '444-444-4444', registration_sent => '1999-01-01 00:00:00', email => 'plugin@email.com', id => '73', new_user => '1', state_or_province => 'state_or_province', confirm_local_passwd => '5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8', gid => '1', }, }, 
'3' => { new => 'TmUser', => using => { id => '74', last_updated => '2015-12-10 15:43:45', local_passwd => '5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8', registration_sent => '1999-01-01 00:00:00', role => '6', city => 'city', confirm_local_passwd => '5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8', country => 'United States', username => 'portal', postal_code => '80122', token => '', uid => '1', email => 'test1@email.com', full_name => 'The Portal User', phone_number => '222-222-2222', address_line1 => 'address_line3', company => undef, state_or_province => 'state_or_province', new_user => '1', address_line2 => 'address_line4', gid => '1', }, }, 
'4' => { new => 'TmUser', => using => { address_line1 => 'address_line1', postal_code => '80122', token => '', uid => '1', username => 'admin', address_line2 => 'address_line2', city => 'city', company => undef, full_name => 'The Admin User', local_passwd => '5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8', state_or_province => 'state_or_province', confirm_local_passwd => '5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8', country => 'United States', email => 'admin@cable.comcast.com', id => '75', last_updated => '2015-12-10 15:43:45', role => '4', gid => '1', new_user => '1', phone_number => '111-111-1111', registration_sent => '1999-01-01 00:00:00', }, }, 
'5' => { new => 'TmUser', => using => { state_or_province => undef, token => undef, country => undef, email => undef, local_passwd => '5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8', registration_sent => '1999-01-01 00:00:00', uid => '1', username => 'migration', last_updated => '2015-12-10 15:43:45', postal_code => undef, full_name => 'Migration User -- DO NOT REMOVE', new_user => '1', address_line2 => undef, city => undef, confirm_local_passwd => '5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8', id => '76', phone_number => undef, role => '5', address_line1 => undef, company => undef, gid => '1', }, }, 
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
