package Helper::Stats;
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

sub new {
	my $class = shift;
	my $self = bless { c => $class, }, $class;
	return $self;
}

sub series_name {
	my $self            = shift;
	my $cdn_name        = shift;
	my $cachegroup_name = shift;
	my $cache_name      = shift;
	my $metric_type     = shift;

	# 'series' section
	my $delim = ":";

	# Example: <cdn_name>:<cachegroup_name>:<cache_name>:<metric_type>
	return sprintf( "%s$delim%s$delim%s$delim%s", $cdn_name, $cachegroup_name, $cache_name, $metric_type );
}

sub build_summary_query {
	my $self        = shift;
	my $series_name = shift;
	my $start_date  = shift;
	my $end_date    = shift;
	my $interval    = shift;    # Valid interval examples 10m (minutes), 10s (seconds), 1h (hour)
	my $limit       = shift;

	#'summary' section
	return sprintf( '%s "%s" %s',
		"select mean(value), percentile(value, 5), percentile(value, 95), percentile(value, 98), min(value), max(value), sum(value), count(value) from ",
		$series_name, "where time > '$start_date' and time < '$end_date'" );
}

sub build_series_query {
	my $self        = shift;
	my $series_name = shift;
	my $start_date  = shift;
	my $end_date    = shift;
	my $interval    = shift;    # Valid interval examples 10m (minutes), 10s (seconds), 1h (hour)
	my $limit       = shift;

	return sprintf( '%s "%s" %s', "select value from ", $series_name, "where time > '$start_date' and time < '$end_date'" );
}

sub build_summary {
	my $self            = shift;
	my $summary_content = shift;    # in perl hash form

	my $results = $summary_content->{results}[0];
	my $values  = $results->{series}[0]{values}[0];

	my $values_size;
	if ( defined($values) ) {
		$values_size = keys $values;
	}
	my $summary      = ();
	my $series_count = 0;

	if ( defined($values_size) & ( $values_size > 0 ) ) {
		my $avg = $summary_content->{results}[0]{series}[0]{values}[0][1];
		my $average = nearest( .001, $avg );
		$average =~ /([\d\.]+)/;
		$summary->{average}                = $average;
		$summary->{fifthPercentile}        = $summary_content->{results}[0]{series}[0]{values}[0][2];
		$summary->{ninetyFifthPercentile}  = $summary_content->{results}[0]{series}[0]{values}[0][3];
		$summary->{ninetyEighthPercentile} = $summary_content->{results}[0]{series}[0]{values}[0][4];
		$summary->{min}                    = $summary_content->{results}[0]{series}[0]{values}[0][5];
		$summary->{max}                    = $summary_content->{results}[0]{series}[0]{values}[0][6];
		$summary->{total}                  = $summary_content->{results}[0]{series}[0]{values}[0][7];
		$series_count                      = $summary_content->{results}[0]{series}[0]{values}[0][8];
	}
	else {
		$summary->{average}     = 0;
		$summary->{ninetyFifth} = 0;
		$summary->{min}         = 0;
		$summary->{max}         = 0;
		$summary->{total}       = 0;
	}

	return ( $summary, $series_count );

}

sub build_series {
	my $self           = shift;
	my $series_content = shift;
	my $results        = $series_content->{results}[0];
	my $series         = $results->{series}[0];
	my $values         = $series->{values};
	my $values_size;

	if ( defined($values) ) {
		$values_size = keys $values;
	}

	if ( defined($values_size) & ( $values_size > 0 ) ) {
		$values_size = keys $values;
		return $series;
	}
	else {
		return [];
	}
}

1;
