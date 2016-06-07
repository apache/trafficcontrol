package Fixtures::SteeringType;
use strict;
use warnings FATAL => 'all';

use Moose;
extends 'DBIx::Class::EasyFixture';
use namespace::autoclean;
use Digest::SHA1 qw(sha1_hex);

my %definition_for = (
    STEERING_REGEXP => {
        new   => 'Type',
        using => {
            id           => 987,
            name         => 'STEERING_REGEXP',
            description  => 'Steering target filter regular expression',
            use_in_table => 'regex',
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