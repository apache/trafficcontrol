package MojoPlugins::v12::Cache;
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

use Mojo::Base 'Mojolicious::Plugin';
use Carp qw(cluck confess);
use Data::Dumper;
use Utils::Helper::DateHelper;
use JSON;
use HTTP::Date;

#TODO: drichardson - pull this from the 'Parameters';
use constant DB_NAME => "cache_stats";

sub register {
	my ( $self, $app, $conf ) = @_;

	$app->renderer->add_helper(
		v12_cache_series_name => sub {
			my $self            = shift;
			my $cdn_name        = shift;
			my $ds_name         = shift;
			my $cachegroup_name = shift;
			my $metric_type     = shift;

			# 'series' section
			my $delim = ":";

			# over-the-top:pixl-tv-linear:us-fl-sarasota:tps_4xx
			return sprintf( "%s$delim%s$delim%s$delim%s", $cdn_name, $ds_name, $cachegroup_name, $metric_type );
		}
	);

	$app->renderer->add_helper(
		cache_stats => sub {
			my $self            = shift;
			my $dsid            = shift;
			my $cachegroup_name = shift;
			my $metric_type     = shift;
			my $start_date      = shift;
			my $end_date        = shift;
			my $interval        = shift;    # Valid interval examples 10m (minutes), 10s (seconds), 1h (hour)
			my $limit           = shift;

			my ( $cdn_name, $ds_name ) = $self->lookup_cdn_name_and_ds_name($dsid);

			my $series_name = $self->v12_cache_series_name( $cdn_name, $ds_name, $cachegroup_name, $metric_type );

			#'summary' section
			my $summary_query = sprintf( '%s "%s" %s',
				"select mean(value), percentile(value, 95), min(value), max(value), sum(value), count(value) from ",
				$series_name, "where time > '$start_date' and time < '$end_date' limit 20" );
			my ( $summary_content, $series_count ) = $self->v12_build_cache_summary( $series_name, $summary_query );

			my $series_query = sprintf( '%s "%s" %s', "select value from ", $series_name, "where time > '$start_date' and time < '$end_date' limit 20" );

			#'series' section
			my $series_content = $self->v12_build_cache_series( $series_name, $series_query );

			my $parent_node = "stats";
			my $result      = ();
			$result->{$parent_node}{series}      = $series_content;
			$result->{$parent_node}{seriesCount} = $series_count;

			if ( %{$result} ) {
				$result->{$parent_node}{cdnName}              = $cdn_name;
				$result->{$parent_node}{deliveryServiceId}    = $dsid;
				$result->{$parent_node}{cacheGroupName}       = $cachegroup_name;
				$result->{$parent_node}{startDate}            = $start_date;
				$result->{$parent_node}{endDate}              = $end_date;
				$result->{$parent_node}{interval}             = int($interval);
				$result->{$parent_node}{metricType}           = $metric_type;
				$result->{$parent_node}{influxdbDatabaseName} = DB_NAME;
				$result->{$parent_node}{influxdbSeriesQuery}  = $series_query;
				$result->{$parent_node}{influxdbSummaryQuery} = $summary_query;
				$result->{$parent_node}{summary}              = $summary_content;
			}

			$self->success($result);
		}
	);

	$app->renderer->add_helper(
		v12_build_cache_summary => sub {
			my $self        = shift;
			my $series_name = shift;
			my $query       = shift;

			my $response_container = $self->influxdb_query( DB_NAME, $query );
			my $response = $response_container->{response};

			if ( $response->is_success ) {
				my $content         = $response->{_content};
				my $summary_content = decode_json($content);
				my $results         = $summary_content->{results}[0];
				my $values          = $results->{series}[0]{values}[0];
				my $values_size     = keys $values;
				my $summary         = ();
				my $series_count    = 0;

				if ( $values_size > 1 ) {
					$summary->{average}     = $summary_content->{results}[0]{series}[0]{values}[0][1];
					$summary->{ninetyFifth} = $summary_content->{results}[0]{series}[0]{values}[0][2];
					$summary->{min}         = $summary_content->{results}[0]{series}[0]{values}[0][3];
					$summary->{max}         = $summary_content->{results}[0]{series}[0]{values}[0][4];
					$summary->{total}       = $summary_content->{results}[0]{series}[0]{values}[0][5];
					$series_count           = $summary_content->{results}[0]{series}[0]{values}[0][6];
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
			else {
				$self->internal_server("Could not return deliveryservice stats 'summary'");
			}

		}
	);

	$app->renderer->add_helper(
		v12_build_cache_series => sub {
			my $self        = shift;
			my $series_name = shift;
			my $query       = shift;

			my $response_container = $self->influxdb_query( DB_NAME, $query );
			my $response = $response_container->{response};

			if ( $response->is_success ) {
				my $content     = decode_json( $response->{_content} );
				my $results     = $content->{results}[0];
				my $series      = $results->{series}[0];
				my $values      = $series->{values};
				my $values_size = keys $values;
				if ( $values_size > 0 ) {
					return $series;
				}
				else {
					return [];
				}
			}
			else {
				$self->internal_server("Could not return deliveryservice stats 'series'");
			}

		}
	);

}

1;
