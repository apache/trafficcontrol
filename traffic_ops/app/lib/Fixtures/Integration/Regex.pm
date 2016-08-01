package Fixtures::Integration::Regex;

# Do not edit! Generated code.
# See https://github.com/Comcast/traffic_control/wiki/The%20Kabletown%20example

use Moose;
extends 'DBIx::Class::EasyFixture';
use namespace::autoclean;

my %definition_for = (
	## id => 1
	'0' => {
		new => 'Regex',
		using => {
			pattern => '.*\.games\..*',
			last_updated => '2015-12-10 15:43:45',
			type => '20',
		},
	},
	## id => 2
	'1' => {
		new => 'Regex',
		using => {
			pattern => '.*\.games\..*',
			last_updated => '2015-12-10 15:43:45',
			type => '20',
		},
	},
	## id => 3
	'2' => {
 		new => 'Regex',
 		using => {
			pattern => '.*\.images\..*',
 			last_updated => '2015-12-10 15:43:45',
 			type => '20',
 		},
	},
	## id => 4
	'3' => {
		new => 'Regex',
		using => {
			pattern => '.*\.images\..*',
			last_updated => '2015-12-10 15:43:45',
			type => '20',
		},
	},
	## id => 5
	'4' => {
		new => 'Regex',
		using => {
			pattern => '.*\.movies\..*',
			type => '20',
			last_updated => '2015-12-10 15:43:45',
		},
	},
	## id => 6
	'5' => {
		new => 'Regex',
		using => {
			pattern => '.*\.movies\..*',
			last_updated => '2015-12-10 15:43:45',
			type => '20',
		},
	},
	## id => 7
	'6' => {
 		new => 'Regex',
 		using => {
			pattern => '.*\.tv\..*',
 			last_updated => '2015-12-10 15:43:45',
 			type => '20',
 		},
	},
	## id => 8
	'7' => {
 		new => 'Regex',
 		using => {
			pattern => '.*\.tv\..*',
 			last_updated => '2015-12-10 15:43:45',
 			type => '20',
 		},
	},
);

sub name {
		return "Regex";
}

sub get_definition {
		my ( $self,
 			$name ) = @_;
		return $definition_for{$name};
}

sub all_fixture_names {
	# sort by db pattern to guarantee insertion order
	return (sort { $definition_for{$a}{using}{pattern} cmp $definition_for{$b}{using}{pattern} } keys %definition_for);
}

__PACKAGE__->meta->make_immutable;
1;
