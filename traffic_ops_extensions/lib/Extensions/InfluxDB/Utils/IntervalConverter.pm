package Extensions::InfluxDB::Utils::IntervalConverter;
#
# Copyright 2011-2014, Comcast Corporation. This software and its contents are
# Comcast confidential and proprietary. It cannot be used, disclosed, or
# distributed without Comcast's prior written permission. Modification of this
# software is only allowed at the direction of Comcast Corporation. All allowed
# modifications must be provided to Comcast Corporation.
#
use Time::Seconds;
use Data::Dumper;

my $units = {
	ns => 0.0000000001,
	ms => 0.001,
	s  => 1,
	m  => ONE_MINUTE,
	h  => ONE_HOUR,
	d  => ONE_DAY,
	w  => ONE_WEEK,
	mo => ONE_MONTH,
	y  => ONE_YEAR
};

sub new {
	my $self  = {};
	my $class = shift;

	return ( bless( $self, $class ) );
}

sub to_seconds {
	my $self     = shift;
	my $interval = shift;
	my ( $digits, $unit_of_time ) = $self->parse_interval($interval);

	return $digits * $units->{$unit_of_time};
}

sub to_milliseconds {
	my $self     = shift;
	my $interval = shift;
	my ( $digits, $unit_of_time ) = $self->parse_interval($interval);

	return $digits * $units->{$unit_of_time} * 1000;
}

sub to_nanoseconds {
	my $self     = shift;
	my $interval = shift;
	my ( $digits, $unit_of_time ) = $self->parse_interval($interval);

	return $digits * $units->{$unit_of_time} * 1000000000;
}

sub parse_interval {
	my $self     = shift;
	my $interval = shift;
	my ( $digits, $unit_of_time ) = $interval =~ /(\d*)(.*)/;
	return ( $digits, $unit_of_time );
}

1;
