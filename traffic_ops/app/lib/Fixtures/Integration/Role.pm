package Fixtures::Integration::Role;

# Do not edit! Generated code.
# See https://github.com/Comcast/traffic_control/wiki/The%20Kabletown%20example

use Moose;
extends 'DBIx::Class::EasyFixture';
use namespace::autoclean;

my %definition_for = (
'0' => { new => 'Role', => using => { description => 'block all access', id => '1', name => 'disallowed', priv_level => '0', }, }, 
'1' => { new => 'Role', => using => { description => 'block all access', id => '2', name => 'read-only user', priv_level => '10', }, }, 
'2' => { new => 'Role', => using => { description => 'block all access', id => '3', name => 'operations', priv_level => '20', }, }, 
'3' => { new => 'Role', => using => { description => 'super-user', id => '4', name => 'admin', priv_level => '30', }, }, 
'4' => { new => 'Role', => using => { description => 'database migrations user - DO NOT REMOVE', id => '5', name => 'migrations', priv_level => '20', }, }, 
'5' => { new => 'Role', => using => { description => 'Portal User', id => '6', name => 'portal', priv_level => '2', }, }, 
); 

sub name {
		return "Role";
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
