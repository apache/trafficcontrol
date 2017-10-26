package Extensions::TrafficStats::Builder::CacheStatsBuilder;
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
		cdn_name    => 1,
		start_date  => 1,
		end_date    => 1,
		series_name => 1,
		interval    => 1,
		orderby     => 1,
		limit       => 1,
		offset      => 1
	};
	return Extensions::TrafficStats::Builder::BaseBuilder->validate_keys( $args, $valid_keys );
}

sub usage_overview_max_kbps_query {
	my $self = shift;
	if ( $self->validate_keys() ) {
		my $query = "SELECT sum(value)/1000/1000/6 FROM maxKbps WHERE time > now() - 60s";

		#my $query = "SELECT sum(value) FROM maxKbps where time > now() - 1m";
		return Extensions::TrafficStats::Builder::BaseBuilder->clean_whitespace($query);
	}
}

sub usage_overview_current_gbps_query {
	my $self = shift;
	if ( $self->validate_keys() ) {
		my $query = "SELECT sum(value)/1000/1000/6 FROM bandwidth WHERE time > now() - 60s";

		#my $query = "SELECT sum(value)*1000/6 FROM bandwidth where cachegroup = 'total' and time > now() - 10m";
		return Extensions::TrafficStats::Builder::BaseBuilder->clean_whitespace($query);
	}
}

sub summary_query {
	my $self = shift;
	if ( $self->validate_keys() ) {

		#'summary' section
		my $query = qq[SELECT mean(value), percentile(value, 5), percentile(value, 95), percentile(value, 98), min(value), max(value), sum(value), count(value)
				FROM "monthly"."$args->{series_name}.cdn.1min"
				WHERE cdn = '$args->{cdn_name}'
					AND time > '$args->{start_date}'
					AND time < '$args->{end_date}'
					GROUP BY time($args->{interval}), cdn];

		$query = Extensions::TrafficStats::Builder::BaseBuilder->append_clauses( $query, $args );

		return Extensions::TrafficStats::Builder::BaseBuilder->clean_whitespace($query);

	}
}

sub series_query {
	my $self = shift;

	my $query = qq[SELECT sum(value)/count(value)
			FROM "monthly"."$args->{series_name}.cdn.1min"
			WHERE cdn = '$args->{cdn_name}'
				AND time > '$args->{start_date}'
				AND time < '$args->{end_date}'
				GROUP BY time($args->{interval}), cdn
				ORDER BY asc];

	$query = Extensions::TrafficStats::Builder::BaseBuilder->append_clauses( $query, $args );

	return Extensions::TrafficStats::Builder::BaseBuilder->clean_whitespace($query);
}

1;
