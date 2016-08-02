package Fixtures::Integration::DeliveryserviceRegex;

# Do not edit! Generated code.
# See https://github.com/Comcast/traffic_control/wiki/The%20Kabletown%20example

use Moose;
extends 'DBIx::Class::EasyFixture';
use namespace::autoclean;

my %definition_for = (
	'0' => {
		new => 'DeliveryserviceRegex',
		using => {
			regex => '5',
			set_number => '0',
			deliveryservice => '4',
		},
	},
	'1' => {
		new => 'DeliveryserviceRegex',
		using => {
			regex => '3',
			deliveryservice => '3',
			set_number => '0',
		},
	},
	'2' => {
		new => 'DeliveryserviceRegex',
		using => {
			regex => '1',
			set_number => '0',
			deliveryservice => '2',
		},
	},
	'3' => {
		new => 'DeliveryserviceRegex',
		using => {
			regex => '7',
			deliveryservice => '5',
			set_number => '0',
		},
	},
	'4' => {
		new => 'DeliveryserviceRegex',
		using => {
			regex => '6',
			deliveryservice => '7',
			set_number => '0',
		},
	},
	'5' => {
		new => 'DeliveryserviceRegex',
		using => {
			regex => '4',
			deliveryservice => '1',
			set_number => '0',
		},
	},
	'6' => {
		new => 'DeliveryserviceRegex',
		using => {
			regex => '2',
			deliveryservice => '6',
			set_number => '0',
		},
	},
	'7' => {
		new => 'DeliveryserviceRegex',
		using => {
			regex => '8',
			deliveryservice => '8',
			set_number => '0',
		},
	},
);

sub name {
		return "DeliveryserviceRegex";
}

sub get_definition {
		my ( $self, $name ) = @_;
		return $definition_for{$name};
}

sub all_fixture_names {
	# sort by db regex to guarantee insertion order
	return (sort { $definition_for{$a}{using}{regex} cmp $definition_for{$b}{using}{regex} } keys %definition_for);
}

__PACKAGE__->meta->make_immutable;
1;
