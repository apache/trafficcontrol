package Extensions::TrafficStats::Builder::BaseBuilder;
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
use Carp qw(cluck confess);

my $mojo;

sub new {
	my $self  = {};
	my $class = shift;
	$mojo = shift;
	return ( bless( $self, $class ) );
}

sub validate_keys {
	my $self       = shift;
	my $args       = shift;
	my $valid_keys = shift;
	my $valid      = 1;
	foreach my $k ( keys %$args ) {
		unless ( defined( $valid_keys->{$k} ) ) {
			confess("'$k' is not a valid key constructor key.");
			$valid = 0;
		}
	}
	return $valid;
}

sub summary_response {
	my $self            = shift;
	my $summary_content = shift;    # in perl hash form

	my $results = $summary_content->{results}[0];
	my $values  = $results->{series}[0]{values}[0];

	my $summary = ();

	my $avg = $summary_content->{results}[0]{series}[0]{values}[0][1];

	my $average = nearest( .001, $avg );
	$average =~ /([\d\.]+)/;
	$summary->{average} = $average;
	my $fifth_percentile = $summary_content->{results}[0]{series}[0]{values}[0][2];
	$summary->{fifthPercentile} = ( defined($fifth_percentile) ) ? $fifth_percentile : 0;

	my $ninety_fifth_percentile = $summary_content->{results}[0]{series}[0]{values}[0][3];
	$summary->{ninetyFifthPercentile} = ( defined($ninety_fifth_percentile) ) ? $ninety_fifth_percentile : 0;

	my $ninety_eighth_percentile = $summary_content->{results}[0]{series}[0]{values}[0][4];
	$summary->{ninetyEighthPercentile} = ( defined($ninety_eighth_percentile) ) ? $ninety_eighth_percentile : 0;

	my $min = $summary_content->{results}[0]{series}[0]{values}[0][5];
	$summary->{min} = ( defined($min) ) ? $min : 0;

	my $max = $summary_content->{results}[0]{series}[0]{values}[0][6];
	$summary->{max} = ( defined($max) ) ? $max : 0;

	my $count = $summary_content->{results}[0]{series}[0]{values}[0][7];
	$summary->{count} = ( defined($count) ) ? $count : 0;

	return $summary;

}

sub series_response {
	my $self           = shift;
	my $series_content = shift;
	my $results        = $series_content->{results}[0];
	my $series         = $results->{series}[0];
	my $values         = $series->{values};
	my $values_size;

	if ( defined($values) ) {
		$values_size = @$values;
	}

	if ( defined($values_size) & ( $values_size > 0 ) ) {
		return $series;
	}
	else {
		return [];
	}
}

sub clean_whitespace {
	my $self  = shift;
	my $query = shift;

	# cleanup whitespace
	$query =~ s/\\n//g;
	$query =~ s/\s+/ /g;
	return $query;
}

sub to_influxdb_date {
	my $self = shift;
	my $date = shift;

	if ( defined($date) && $date eq "now" ) {
		$date = 'now()';
	}
	else {
		$date = "'" . $date . "'";
	}
	return $date;
}

sub append_clauses {
	my $self  = shift;
	my $query = shift;
	my $args  = shift;
	$query = defined( $args->{orderby} ) ? $query .= " ORDER BY " . $args->{orderby} : $query;
	$query = defined( $args->{limit} )   ? $query .= " LIMIT " . $args->{limit}      : $query;
	$query = defined( $args->{offset} )  ? $query .= " OFFSET " . $args->{offset}    : $query;
	return $query;
}

1;
