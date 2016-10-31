package Utils::Helper::DateHelper;
#
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
#
#

use Carp qw(cluck confess);
use Data::Dumper;
use Time::Local;

sub new {
	my $self  = {};
	my $class = shift;
	my $args  = shift;

	return ( bless( $self, $class ) );
}

sub date_to_epoch {
	my $self     = shift;
	my $datetime = shift;
	my $unixtime;
	eval {
		my ( $datestr, $timestr ) = split( /\s/, $datetime );
		my ( $hour, $min, $sec ) = split( /:/, $timestr );
		my ( $year, $mon, $day ) = split( /-/, $datestr );
		$mon = $mon - 1;                                              # This mod uses ordinal numbers
		$unixtime = timegm( $sec, $min, $hour, $day, $mon, $year );
	};
	if ($@) {
		return -1;
	}
	else {
		return ($unixtime);
	}
}

sub translate_dates {
	my $self       = shift;
	my $start_date = shift;
	my $end_date   = shift;

	if ( $end_date eq "now" && $start_date ne "now" ) {
		$end_date = time();
	}
	if ( $start_date eq "yesterday" ) {
		$start_date = $end_date - 86400;
	}
	return ( $start_date, $end_date );
}

1;
