package Fixtures::SteeringDeliveryserviceRegex;
use strict;
use warnings FATAL => 'all';

use Moose;
extends 'DBIx::Class::EasyFixture';
use namespace::autoclean;

my %definition_for = (
    target_r1 => {
        new   => 'DeliveryserviceRegex',
        using => {
            deliveryservice => 20001,
            regex           => 21001,
            set_number      => 0,
        },
    },
    target_r2 => {
        new   => 'DeliveryserviceRegex',
        using => {
            deliveryservice => 20001,
            regex           => 21002,
            set_number      => 0,
        },
    },
    target_r3 => {
        new   => 'DeliveryserviceRegex',
        using => {
            deliveryservice => 20003,
            regex           => 21003,
            set_number      => 0,
        },
    },
    target_r4 => {
        new   => 'DeliveryserviceRegex',
        using => {
            deliveryservice => 20004,
            regex           => 21004,
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