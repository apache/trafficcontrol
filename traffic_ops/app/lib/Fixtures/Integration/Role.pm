package Fixtures::Integration::Role;

# Do not edit! Generated code.
# See https://github.com/Comcast/traffic_control/wiki/The%20Kabletown%20example

use Moose;
extends 'DBIx::Class::EasyFixture';
use namespace::autoclean;

my %definition_for = (
	## id => 1
	'0' => {
		new => 'Role',
		using => {
			name => 'admin',
			description => 'super-user',
			priv_level => '30',
		},
	},
	## id => 2
	'1' => {
		new => 'Role',
		using => {
			name => 'disallowed',
			description => 'block all access',
			priv_level => '0',
		},
	},
	## id => 3
	'2' => {
		new => 'Role',
		using => {
			name => 'migrations',
			description => 'database migrations user - DO NOT REMOVE',
			priv_level => '20',
		},
	},
	## id => 4
	'3' => {
		new => 'Role',
		using => {
			name => 'operations',
			description => 'block all access',
			priv_level => '20',
		},
	},
	## id => 5
	'4' => {
		new => 'Role',
		using => {
			name => 'portal',
			description => 'Portal User',
			priv_level => '2',
		},
	},
	## id => 6
	'5' => {
		new => 'Role',
		using => {
			name => 'read-only user',
			description => 'block all access',
			priv_level => '10',
		},
	},
);

sub name {
		return "Role";
}

sub get_definition {
		my ( $self,
			$name ) = @_;
		return $definition_for{$name};
}

sub all_fixture_names {
	# sort by db name to guarantee insertion order
	return (sort { $definition_for{$a}{using}{name} cmp $definition_for{$b}{using}{name} } keys %definition_for);
}

__PACKAGE__->meta->make_immutable;
1;
