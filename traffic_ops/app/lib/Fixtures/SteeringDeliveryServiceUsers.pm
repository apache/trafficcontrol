package Fixtures::SteeringDeliveryServiceUsers;
use strict;
use warnings FATAL => 'all';

use Moose;
extends 'DBIx::Class::EasyFixture';
use namespace::autoclean;
use Digest::SHA1 qw(sha1_hex);

my %definition_for = (
    ds_steering_user1 => {
        new   => 'DeliveryserviceTmuser',
        using => {
            deliveryservice => 10001,
            tm_user_id      => 101,
        },
    },
        ds_steering_user2 => {
        new   => 'DeliveryserviceTmuser',
        using => {
            deliveryservice => 10002,
            tm_user_id      => 102,
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