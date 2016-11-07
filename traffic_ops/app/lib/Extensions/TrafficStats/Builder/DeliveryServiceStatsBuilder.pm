package Extensions::TrafficStats::Builder::DeliveryServiceStatsBuilder;
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

use utf8;
use Data::Dumper;
use JSON;
use File::Slurp;
use Math::Round;
use Extensions::TrafficStats::Builder::BaseBuilder;
use Carp qw(cluck confess);

my $args;

sub new {
	my $self  = {};
	my $class = shift;
	$args = shift;

	return ( bless( $self, $class ) );
}

sub validate_keys {
	my $self       = shift;
	my $valid_keys = {
		deliveryServiceName => 1,
		metricType          => 1,
		startDate           => 1,
		endDate             => 1,
		dbName              => 1,
		interval            => 1,
		orderby             => 1,
		exclude             => 1,
		limit               => 1,
		offset              => 1
	};
	return Extensions::TrafficStats::Builder::BaseBuilder->validate_keys( $args, $valid_keys );
}

sub summary_query {
	my $self = shift;
	if ( $self->validate_keys() ) {

		my $end_date = Extensions::TrafficStats::Builder::BaseBuilder->to_influxdb_date( $args->{endDate} );

		#'summary' section
		my $query = sprintf(
			'%s %s %s',
			qq[SELECT mean(value), percentile(value, 5), percentile(value, 95), percentile(value, 98), min(value), max(value), count(value)
                FROM "$args->{dbName}"."monthly"."$args->{metricType}.ds.1min"
                WHERE time >= '$args->{startDate}' AND time <= $end_date AND cachegroup = 'total' AND deliveryservice = '$args->{deliveryServiceName}']
		);

		$query = Extensions::TrafficStats::Builder::BaseBuilder->append_clauses( $query, $args );
		return Extensions::TrafficStats::Builder::BaseBuilder->clean_whitespace($query);
	}
}

sub series_query {
	my $self = shift;

	my $end_date = Extensions::TrafficStats::Builder::BaseBuilder->to_influxdb_date( $args->{endDate} );

	my $query = sprintf(
		'%s %s %s',
		qq[SELECT sum(value)/count(value)
		    FROM "$args->{dbName}"."monthly"."$args->{metricType}.ds.1min"
		    WHERE cachegroup = 'total' AND deliveryservice = '$args->{deliveryServiceName}' AND time >='$args->{startDate}' AND time <= $end_date GROUP BY time($args->{interval}), cachegroup]
	);

	$query = Extensions::TrafficStats::Builder::BaseBuilder->append_clauses( $query, $args );

	return Extensions::TrafficStats::Builder::BaseBuilder->clean_whitespace($query);
}

sub usage_overview_tps_query {
	my $self = shift;

	if ( $self->validate_keys() ) {
		my $query = qq[
		    SELECT sum(value)/6
		    FROM "$args->{dbName}"."monthly"."tps.ds.1min"
		    WHERE cachegroup = 'total' and time > now() - 60s
		];
		return Extensions::TrafficStats::Builder::BaseBuilder->clean_whitespace($query);
	}
}

1;
