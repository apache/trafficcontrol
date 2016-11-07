package Fixtures::SteeringDeliveryserviceRegex;

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

my %definition_for = (
    target_r1_filter => {
        new   => 'DeliveryserviceRegex',
        using => {
            deliveryservice => 20001,
            regex           => 21001,
            set_number      => 0,
        },
    },
    target_r2_filter => {
        new   => 'DeliveryserviceRegex',
        using => {
            deliveryservice => 20001,
            regex           => 21002,
            set_number      => 0,
        },
    },
    target_r3_filter => {
        new   => 'DeliveryserviceRegex',
        using => {
            deliveryservice => 20003,
            regex           => 21003,
            set_number      => 0,
        },
    },
    target_r4_filter => {
        new   => 'DeliveryserviceRegex',
        using => {
            deliveryservice => 20004,
            regex           => 21004,
            set_number      => 0,
        },
    },
    steering_1 => {
        new   => 'DeliveryserviceRegex',
        using => {
            deliveryservice => 10001,
            regex           => 21101,
            set_number      => 0,
        },
    },
    steering_2 => {
        new   => 'DeliveryserviceRegex',
        using => {
            deliveryservice => 10002,
            regex           => 21102,
            set_number      => 0,
        },
    },
    target_1 => {
        new   => 'DeliveryserviceRegex',
        using => {
            deliveryservice => 20001,
            regex           => 22201,
            set_number      => 0,
        },
    },
    target_2 => {
        new   => 'DeliveryserviceRegex',
        using => {
            deliveryservice => 20002,
            regex           => 22202,
            set_number      => 0,
        },
    },
    target_3 => {
        new   => 'DeliveryserviceRegex',
        using => {
            deliveryservice => 20003,
            regex           => 22203,
            set_number      => 0,
        },
    },
    target_4 => {
        new   => 'DeliveryserviceRegex',
        using => {
            deliveryservice => 20004,
            regex           => 22204,
            set_number      => 0,
        },
    },
    new_steering => {
        new   => 'DeliveryserviceRegex',
        using => {
            deliveryservice => 10003,
            regex           => 21103,
            set_number      => 0,
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
