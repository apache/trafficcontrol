package Extensions::InfluxDB::Utils::InfluxDBDecorator;
#
# Copyright 2015 Comcast Cable Communications Management, LLC
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

use utf8;
use Data::Dumper;
use JSON;
use File::Slurp;
use Math::Round;
use Carp qw(cluck confess);
use Common::ReturnCodes qw(SUCCESS ERROR);
use Extensions::InfluxDB::Utils::IntervalConverter;
use POSIX qw(strftime);

my $args;

sub new {
	my $self  = {};
	my $class = shift;
	$args = shift;

	return ( bless( $self, $class ) );
}

sub to_influx_series_format {
	my $self        = shift;
	my $json        = shift;
	my $metric_name = shift;
	my $r_to_s      = shift;

	my $rc = SUCCESS;
	my @samples;
	my $prev_time = 0;
	my $time_base = 0;

	my $prev_timebase;
	my $prev_interval_time;
	my $interval                 = $args->{interval};
	my $interval_for_metric_type = $args->{interval_for_metric_type};

	my $prev_timebase_iso8601 = 0;
	for my $dp ( @{ $json->[0]->{data} } ) {
		my $time = $dp->[0];

		#print "TOP time #-> (" . $time . ")\n";
		my $value = $dp->[1];

		#print "time: " . gmtime($time) . " - value: $value \n";
		if ( $prev_time == 0 ) {
			$prev_timebase = $time;
		}
		else {
			$prev_timebase += $interval;
		}
		if ( $prev_time != 0 ) {

			# Keep replicating the value until we reach the next interval change
			my $prior_interval    = $prev_time + $interval;
			my $interval_quotient = $interval_for_metric_type / $interval;
			for ( my $i = 0; $i < $interval_quotient; $i++ ) {
				$prev_timebase_iso8601 = strftime( "%Y-%m-%dT%H:%M:%SZ", gmtime($prev_timebase) );
				my $entry = [ $prev_timebase_iso8601, $value ];
				push( @samples, $entry );
				$prev_timebase += $interval;
			}
		}

		if ( exists( $r_to_s->{$metric_name}->{conversion} ) ) {
			$value = $r_to_s->{$metric_name}->{conversion}->($value);
		}

		$prev_time = $time + $interval;
		$time_base = $time if ( !$time_base );
	}

	my $s = gmtime($time_base);    # first timestamp
	my $e = gmtime($prev_time);    # last timestamp

	my $data = {

		#spdbUrl => $url,
		statName            => $metric_name,
		cdnName             => $args->{cdn},
		deliveryServiceName => $args->{ds_name},
		cacheGroupName      => $args->{cache_group_name},
		hostName            => $args->{host_name},
		elapsed             => $args->{time_elapsed},
		start               => $s,
		end                 => $e,
		interval            => $interval_for_metric_type,
		series              => [
			{
				samples  => \@samples,
				timeBase => $time_base
			}
		]
	};

	return ( $rc, $data );

}

1;
