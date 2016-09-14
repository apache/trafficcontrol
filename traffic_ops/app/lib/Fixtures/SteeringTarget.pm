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
            deliveryservice => 1,
            target => 4,
            weight => 1000,
        }
    },
    steering_target_2 => {
        new => 'SteeringTarget',
        using => {
            deliveryservice => 1,
            target => 5,
            weight => 7654,
        }
    },
    steering_target_3 => {
        new => 'SteeringTarget',
        using => {
            deliveryservice => 2,
            target => 6,
            weight => 123,
        }
    },
    steering_target_4 => {
        new => 'SteeringTarget',
        using => {
            deliveryservice => 2,
            target => 7,
            weight => 999,
        }
    },
);

sub get_definition {
    my ( $self, $name ) = @_;
    return $definition_for{$name};
}

sub all_fixture_names {
	# sort by db deliveryservice to guarantee insertion order
	return (sort { $definition_for{$a}{using}{deliveryservice} cmp $definition_for{$b}{using}{deliveryservice} } keys %definition_for);
}

__PACKAGE__->meta->make_immutable;

1;
