package Fixtures::Integration::Cdn;

# Do not edit! Generated code.
# See https://github.com/Comcast/traffic_control/wiki/The%20Kabletown%20example

use Moose;
extends 'DBIx::Class::EasyFixture';
use namespace::autoclean;

my %definition_for = (
'0' => { new => 'Cdn', => using => { dnssec_enabled => '0', id => '1', last_updated => '2015-12-10 15:43:45', name => 'cdn_number_1', }, }, 
'1' => { new => 'Cdn', => using => { name => 'cdn_number_2', dnssec_enabled => '0', id => '2', last_updated => '2015-12-10 15:43:45', }, }, 
); 

sub name {
		return "Cdn";
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
