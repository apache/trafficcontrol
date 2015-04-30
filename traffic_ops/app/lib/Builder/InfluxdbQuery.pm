package Builder::InfluxdbQuery;
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

my $args;

sub new {
	my $self  = {};
	my $class = shift;
	$args = shift;

	return ( bless( $self, $class ) );
}

sub now {
	my $self      = shift;
	my $now_param = shift;
	if ( defined($now_param) ) {
		return "now()";
	}

}

sub valid_keys {
	my $self       = shift;
	my $valid      = 1;
	my $valid_keys = {
		cdn_name        => 1,
		ds_name         => 1,
		cachegroup_name => 1,
		start_date      => 1,
		end_date        => 1,
		series_name     => 1,
		interval        => 1,
		limit           => 1
	};
	foreach my $k ( keys $args ) {
		unless ( defined( $valid_keys->{$k} ) ) {
			confess("'$k' is not a valid key constructor key.");
			$valid = 0;
		}
	}
	return $valid;
}

sub summary_query {
	my $self = shift;
	if ( valid_keys() ) {

		#'summary' section
		return sprintf(
			'%s %s %s',
			"SELECT mean(value), percentile(value, 5), percentile(value, 95), percentile(value, 98), min(value), max(value), sum(value), count(value) FROM",
			$args->{series_name}, "WHERE time > '$args->{start_date}' AND 
                                         time < '$args->{end_date}' AND 
                                         cdn = '$args->{cdn_name}' AND
                                         cachegroup = '$args->{cachegroup_name}' 
                                         GROUP BY time($args->{interval}), cdn, cachegroup, deliveryservice"
		);

	}
}

sub series_query {
	my $self = shift;
	return sprintf(
		'%s %s %s',
		"SELECT value FROM",
		$args->{series_name}, "WHERE time > '$args->{start_date}' AND time < '$args->{end_date}' AND cdn = '$args->{cdn_name}' AND cachegroup =
				'$args->{cachegroup_name}'"
	);
}

sub summary_response {
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

sub series_response {
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
