package Fixtures::Integration::Regex;

# Do not edit! Generated code.
# See https://github.com/Comcast/traffic_control/wiki/The%20Kabletown%20example

use Moose;
extends 'DBIx::Class::EasyFixture';
use namespace::autoclean;

my %definition_for = (
'0' => { new => 'Regex', => using => { pattern => '.*\.movies\..*', type => '18', id => '1', last_updated => '2015-12-10 15:43:45', }, }, 
'1' => { new => 'Regex', => using => { id => '2', last_updated => '2015-12-10 15:43:45', pattern => '.*\.images\..*', type => '18', }, }, 
'2' => { new => 'Regex', => using => { id => '3', last_updated => '2015-12-10 15:43:45', pattern => '.*\.games\..*', type => '18', }, }, 
'3' => { new => 'Regex', => using => { id => '4', last_updated => '2015-12-10 15:43:45', pattern => '.*\.tv\..*', type => '18', }, }, 
'4' => { new => 'Regex', => using => { id => '11', last_updated => '2015-12-10 15:43:45', pattern => '.*\.movies\..*', type => '18', }, }, 
'5' => { new => 'Regex', => using => { id => '12', last_updated => '2015-12-10 15:43:45', pattern => '.*\.images\..*', type => '18', }, }, 
'6' => { new => 'Regex', => using => { id => '13', last_updated => '2015-12-10 15:43:45', pattern => '.*\.games\..*', type => '18', }, }, 
'7' => { new => 'Regex', => using => { id => '14', last_updated => '2015-12-10 15:43:45', pattern => '.*\.tv\..*', type => '18', }, }, 
); 

sub name {
		return "Regex";
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
