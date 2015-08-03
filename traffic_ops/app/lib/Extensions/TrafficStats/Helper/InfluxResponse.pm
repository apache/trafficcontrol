package Extensions::TrafficStats::Helper::InfluxResponse;
#
# Copyright 2011-2014, Comcast Corporation. This software and its contents are
# Comcast confidential and proprietary. It cannot be used, disclosed, or
# distributed without Comcast's prior written permission. Modification of this
# software is only allowed at the direction of Comcast Corporation. All allowed
# modifications must be provided to Comcast Corporation.
#
use UI::Utils;
use constant FIVE_MINUTES => 5;

use Data::Dumper;
use JSON;
use POSIX qw(strftime);
use POSIX qw(localtime);
use HTTP::Date;
use Common::ReturnCodes qw(SUCCESS ERROR);

my $args = shift;

sub new {
	my $self  = {};
	my $class = shift;
	$args = shift;

	return ( bless( $self, $class ) );
}

sub parse_retention_period_in_seconds {
	my $self             = shift;
	my $retention_period = shift;

	undef $/;

	my ( $hour, $minutes, $seconds ) = $retention_period =~ /(\d*)h(\d*)m(\d*)s/ms;

	my $hour_in_seconds    = $hour * 60 * 60;
	my $minutes_in_seconds = $minutes * 60;
	my $total_seconds      = $hour_in_seconds + $minutes_in_seconds + $seconds;

	return ($total_seconds);
}

1;
