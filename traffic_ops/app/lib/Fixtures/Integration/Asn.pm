package Fixtures::Integration::Asn;

# Do not edit! Generated code.
# See https://github.com/Comcast/traffic_control/wiki/The%20Kabletown%20example

use Moose;
extends 'DBIx::Class::EasyFixture';
use namespace::autoclean;

my %definition_for = (
	## id => 1
	'0' => {
		new => 'Asn',
		using => {
			asn => '33',
			cachegroup => '8',
			last_updated => '2015-12-10 15:44:36',
		},
	},
	## id => 2
	'1' => {
		new => 'Asn',
		using => {
			asn => '99',
			cachegroup => '8',
			last_updated => '2015-12-10 15:44:36',
		},
	},
	## id => 3
	'2' => {
		new => 'Asn',
		using => {
			asn => '9839',
			cachegroup => '10',
			last_updated => '2015-12-10 15:44:36',
		},
	},
	## id => 4
	'3' => {
		new => 'Asn',
		using => {
			asn => '9930',
			cachegroup => '9',
			last_updated => '2015-12-10 15:44:36',
		},
	},
	## id => 5
	'4' => {
		new => 'Asn',
		using => {
			asn => '9931',
			cachegroup => '12',
			last_updated => '2015-12-10 15:44:36',
		},
	},
	## id => 6
	'5' => {
		new => 'Asn',
		using => {
			asn => '9932',
			cachegroup => '11',
			last_updated => '2015-12-10 15:44:36',
		},
	},
	## id => 7
	'6' => {
		new => 'Asn',
		using => {
			asn => '9933',
			cachegroup => '7',
			last_updated => '2015-12-10 15:44:36',
		},
	},
	## id => 8
	'7' => {
		new => 'Asn',
		using => {
			asn => '9939',
			cachegroup => '10',
			last_updated => '2015-12-10 15:44:36',
		},
	},
);

sub name {
		return "Asn";
}

sub get_definition {
		my ( $self, $name ) = @_;
		return $definition_for{$name};
}

sub all_fixture_names {
	# sort by db asn to guarantee insertion order
	return (sort { $definition_for{$a}{using}{asn} cmp $definition_for{$b}{using}{asn} } keys %definition_for);
}

__PACKAGE__->meta->make_immutable;
1;
