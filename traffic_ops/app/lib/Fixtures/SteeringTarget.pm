package Fixtures::SteeringTarget;
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