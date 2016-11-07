package Fixtures::SteeringUsers;

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

use strict;
use warnings FATAL => 'all';

use Moose;
extends 'DBIx::Class::EasyFixture';
use namespace::autoclean;
use Digest::SHA1 qw(sha1_hex);

my $local_passwd   = sha1_hex('password');
my %definition_for = (
    steering1 => {
        new   => 'TmUser',
        using => {
            id                   => 101,
            username             => 'steering1',
            role                 => 8,
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
            id                   => 102,
            username             => 'steering2',
            role                 => 8,
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
