package Fixtures::SteeringTarget;

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


my %definition_for = (
    steering_target_1 => {
        new => 'SteeringTarget',
        using => {
            deliveryservice => 10001,
            target => 20001,
            weight => 1000,
        }
    },
    steering_target_2 => {
        new => 'SteeringTarget',
        using => {
            deliveryservice => 10001,
            target => 20002,
            weight => 7654,
        }
    },
    steering_target_3 => {
        new => 'SteeringTarget',
        using => {
            deliveryservice => 10002,
            target => 20003,
            weight => 123,
        }
    },
    steering_target_4 => {
        new => 'SteeringTarget',
        using => {
            deliveryservice => 10002,
            target => 20004,
            weight => 999,
        }
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
