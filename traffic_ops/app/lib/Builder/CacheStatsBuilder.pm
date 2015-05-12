package Builder::CacheStatsBuilder;
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
use Builder::InfluxdbBuilder;
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
	return Builder::InfluxdbBuilder->validate_keys( $args, $valid_keys );
}

sub summary_query {
	my $self = shift;
	if ( $self->validate_keys() ) {

		#'summary' section
		my $query = sprintf(
			'%s "%s" %s',
			"SELECT mean(value), percentile(value, 5), percentile(value, 95), percentile(value, 98), min(value), max(value), sum(value), count(value) FROM",
			$args->{series_name}, "WHERE time > '$args->{start_date}' AND
		                                         time < '$args->{end_date}' AND
		                                         cdn = '$args->{cdn_name}'
		                                         GROUP BY time($args->{interval}), cdn"
		);

		# cleanup whitespace
		$query =~ s/\\n//g;
		$query =~ s/\s+/ /g;
		return $query;

	}
}

sub series_query {
	my $self = shift;

	# TODO: drichardson - make the sum more dynamic based upon the interval
	my $query = sprintf(
		'%s "%s" %s',
		"SELECT sum(value)*1000/6 FROM",
		$args->{series_name}, "WHERE time > '$args->{start_date}' AND 
                               time < '$args->{end_date}' AND 
							   cdn = '$args->{cdn_name}'
                               GROUP BY time($args->{interval}),
							   cdn ORDER BY asc"
	);

	# cleanup whitespace
	$query =~ s/\\n//g;
	$query =~ s/\s+/ /g;
	return $query;
}

1;
