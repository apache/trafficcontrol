package MojoPlugins::DeliveryService;
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

sub register {
	my ( $self, $app, $conf ) = @_;

	$app->renderer->add_helper(
		hr_string_to_mbps => sub {
			my $self = shift;
			my $inp  = shift;

			if    ( !defined($inp) )     { return 0; }                  # default is 0
			elsif ( $inp =~ /^(\d+)T$/ ) { return $1 * 1000000; }
			elsif ( $inp =~ /^(\d+)G$/ ) { return $1 * 1000; }
			elsif ( $inp =~ /^(\d+)M$/ ) { return $1; }
			elsif ( $inp =~ /^(\d+)k$/ ) { return int( $1 / 1000 ); }
			elsif ( $inp =~ /^\d+$/ )    { return $1; }
			else                         { return -1; }

		}
	);

	$app->renderer->add_helper(
		get_daily_usage => sub {
			my $self            = shift;
			my $dsid            = shift;
			my $cachegroup_name = shift;
			my $peak_usage_type = shift;
			my $start           = shift;
			my $end             = shift;
			my $interval        = shift;

			my ( $cdn_name, $ds_name ) = $self->get_cdn_name_ds_name($dsid);

			my $dh = new Utils::Helper::DateHelper();
			( $start, $end ) = $dh->translate_dates( $start, $end );

			my $j = $self->daily_summary( $cdn_name, $ds_name, $cachegroup_name );

			$self->success($j);
		}
	);

	$app->renderer->add_helper(
		deliveryservice_usage => sub {
			my $self            = shift;
			my $dsid            = shift;
			my $cachegroup_name = shift;
			my $metric_type     = shift;
			my $start           = shift;
			my $end             = shift;
			my $interval        = shift;

			my ( $cdn_name, $ds_name ) = $self->get_cdn_name_ds_name($dsid);

			my $dh = new Utils::Helper::DateHelper();
			( $start, $end ) = $dh->translate_dates( $start, $end );
			my $match = $self->build_match( $cdn_name, $ds_name, $cachegroup_name, $metric_type );
			my $j = $self->stats_data( $match, $start, $end, $interval );
			if ( %{$j} ) {
				$j->{deliveryServiceId} = $dsid;    # add dsId to data structure
			}

			$self->success($j);
		}
	);

	$app->renderer->add_helper(
		v12_deliveryservice_usage => sub {
			my $self            = shift;
			my $dsid            = shift;
			my $cachegroup_name = shift;
			my $metric_type     = shift;
			my $start_date      = shift;
			my $end_date        = shift;
			my $interval        = shift;

			my ( $cdn_name, $ds_name ) = $self->get_cdn_name_ds_name($dsid);

			# 'series' section
			my $delim   = ":";
			my $db_name = "deliveryservice_stats";

			# over-the-top:pixl-tv-linear:us-fl-sarasota:tps_4xx
			my $series_name = sprintf( "%s$delim%s$delim%s$delim%s", $cdn_name, $ds_name, $cachegroup_name, $metric_type );
			$self->app->log->debug( "series_name #-> " . $series_name );

			my $result = ();

			# 'summary' section
			my $summary_query = sprintf( '%s "%s" %s',
				"select mean(value), percentile(value, 95), min(value), max(value), sum(value), count(value) from ",
				$series_name, "where time > '$start_date' and time < '$end_date'" );
			my $response_container = $self->influxdb_query( $db_name, $summary_query );
			my $response           = $response_container->{response};
			my $summary_content    = decode_json( $response->{_content} );
			my $summary            = ();
			$summary->{average}     = $summary_content->{results}[0]{series}[0]{values}[0][1];
			$summary->{ninetyFifth} = $summary_content->{results}[0]{series}[0]{values}[0][2];
			$summary->{min}         = $summary_content->{results}[0]{series}[0]{values}[0][3];
			$summary->{max}         = $summary_content->{results}[0]{series}[0]{values}[0][4];
			$summary->{total}       = $summary_content->{results}[0]{series}[0]{values}[0][5];
			$self->app->log->debug( "summary_content #-> " . Dumper($summary_content) );

			my $series_query = sprintf( '%s "%s" %s', "select value from ", $series_name, "where time > '$start_date' and time < '$end_date'" );
			$response_container = $self->influxdb_query( $db_name, $series_query );
			$response = $response_container->{response};
			my $series_content = decode_json( $response->{_content} );
			my $series         = $series_content->{results}[0]{series};
			$result->{usage}{series}      = $series;
			$result->{usage}{seriesCount} = $summary_content->{results}[0]{series}[0]{values}[0][6];

			if ( %{$result} ) {
				$result->{usage}{cdnName}              = $cdn_name;
				$result->{usage}{deliveryServiceId}    = $dsid;
				$result->{usage}{cacheGroupName}       = $cachegroup_name;
				$result->{usage}{startDate}            = $start_date;
				$result->{usage}{endDate}              = $end_date;
				$result->{usage}{interval}             = $interval;
				$result->{usage}{metricType}           = $metric_type;
				$result->{usage}{influxdbName}         = $db_name;
				$result->{usage}{influxdbSeriesQuery}  = $series_query;
				$result->{usage}{influxdbSummaryQuery} = $summary_query;
				$result->{usage}{summary}              = $summary;
			}
			$self->app->log->debug( "result #-> " . Dumper($result) );

			$self->success($result);
		}
	);

	$app->renderer->add_helper(
		get_cdn_name_ds_name => sub {
			my $self = shift;
			my $dsid = shift;

			my $cdn_name = "all";
			my $ds_name  = "all";
			if ( $dsid ne "all" ) {
				my $ds = $self->db->resultset('Deliveryservice')->search( { id => $dsid }, {} )->single();
				$ds_name = $ds->xml_id;
				my $param =
					$self->db->resultset('ProfileParameter')
					->search( { -and => [ profile => $ds->profile->id, 'parameter.name' => 'CDN_name' ] }, { prefetch => [ 'parameter', 'profile' ] } )
					->single();
				$cdn_name = $param->parameter->value;
			}
			return ( $cdn_name, $ds_name );
		}
	);
	$app->renderer->add_helper(
		build_match => sub {
			my $self            = shift;
			my $cdn_name        = shift;
			my $ds_name         = shift;
			my $cachegroup_name = shift;
			my $peak_usage_type = shift;
			return $cdn_name . ":" . $ds_name . ":" . $cachegroup_name . ":all:" . $peak_usage_type;
		}
	);

}

1;
