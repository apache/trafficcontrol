package Fixtures::Integration::Division;

# Do not edit! Generated code.
# See https://github.com/Comcast/traffic_control/wiki/The%20Kabletown%20example

use Moose;
extends 'DBIx::Class::EasyFixture';
use namespace::autoclean;

my %definition_for = (
'0' => { new => 'Division', => using => { id => '1', last_updated => '2015-12-10 15:43:45', name => 'East', }, }, 
'1' => { new => 'Division', => using => { last_updated => '2015-12-10 15:43:45', name => 'West', id => '2', }, }, 
); 

sub name {
		return "Division";
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
