package API::v12::DeliveryServiceStats;
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

# JvD Note: you always want to put Utils as the first use. Sh*t don't work if it's after the Mojo lines.
use UI::Utils;
use Mojo::Base 'Mojolicious::Controller';
use Builder::InfluxdbBuilder;
use Data::Dumper;
use Builder::DeliveryServiceStatsBuilder;
use JSON;
use HTTP::Date;
use Extensions::Delegate::Statistics;
Utils::Helper::Datasource->load_extensions;

use constant SUCCESS => 1;
use constant ERROR   => 0;
my $builder;

#TODO: drichardson
#      - Add required fields validation see lib/API/User.pm based on Validate::Tiny
#      - Verify how much can be refactored after cache_stats value grouping is complete.
sub index {
	my $self        = shift;
	my $ds_name     = $self->param('deliveryServiceName');
	my $metric_type = $self->param('metricType');
	my $server_type = $self->param('serverType');
	my $start_date  = $self->param('startDate');
	my $end_date    = $self->param('endDate');
	my $interval    = $self->param('interval') || "60s";     # Valid interval examples 10m (minutes), 10s (seconds), 1h (hour)
	my $exclude     = $self->param('exclude');
	my $orderby     = $self->param('orderby');
	my $limit       = $self->param('limit');
	my $offset      = $self->param('offset');

	my $stats = new Extensions::Delegate::Statistics();
	if ( $self->is_valid_delivery_service_name($ds_name) ) {
		if ( $self->is_delivery_service_name_assigned($ds_name) ) {

			# Build the summary section
			$builder = new Builder::DeliveryServiceStatsBuilder(
				{
					series_name => $metric_type,
					ds_name     => $ds_name,
					start_date  => $start_date,
					end_date    => $end_date,
					interval    => $interval,
					orderby     => $orderby,
					limit       => $limit,
					offset      => $offset
				}
			);

			my $result      = ();
			my $start_epoch = str2time($start_date);
			my $end_epoch   = str2time($start_date);
			my $rc          = SUCCESS;
			my $formatted_response;

			#> show retention policies deliveryservice_stats;
			#name	duration	replicaN	default
			#default	0		1		false
			#weekly	120h0m0s	3		true
			#TODO: drichardson - use influx retention period.
			my $retention_period = 2;

			# numeric start/end only which should be done upstream but let's be extra cautious
			if ( ( $start_epoch =~ /^\d+$/ && $end_epoch =~ /^\d+$/ ) && $start_epoch < ( time() - $retention_period - 60 ) )
			{    # -60 for diff between client and our time
				$self->app->log->debug("Retrieving 'long term' stats...");
				( $rc, $formatted_response ) = $stats->long_term( $self, $start_epoch, $end_epoch, $interval );

				#$self->app->log->debug( "formatted_response #-> " . Dumper($formatted_response) );
			}
			else {
				$self->app->log->debug("Retrieving 'short term' stats...");

				( $rc, $formatted_response ) = $stats->short_term( $self, $start_date, $end_date, $interval );
			}

			my $summary_query;

			my $include_summary = ( defined($exclude) && $exclude =~ /summary/ ) ? 0 : 1;
			if ($include_summary) {
				( $rc, $result, $summary_query ) = $self->build_summary($result);
			}

			if ($rc) {
				my $include_series = ( defined($exclude) && $exclude =~ /series/ ) ? 0 : 1;
				my $series_query;
				if ($include_series) {
					( $rc, $result, $series_query ) = $self->build_series($result);
				}
				if ($rc) {
					$result = $self->build_parameters( $result, $summary_query, $series_query );
					return $self->success($result);
				}
				else {
					return $self->alert($result);
				}
			}
			else {
				return $self->alert($result);
			}

		}
		else {
			return $self->forbidden();
		}
	}
	else {
		$self->success( {} );
	}

}

sub build_summary {
	my $self   = shift;
	my $result = shift;

	my $summary_query = $builder->summary_query();
	$self->app->log->debug( "summary_query #-> " . Dumper($summary_query) );

	my $response_container = $self->influxdb_query( $self->get_db_name(), $summary_query );
	my $response           = $response_container->{'response'};
	my $content            = $response->{_content};

	my $summary;
	my $summary_content;
	my $series_count = 0;
	if ( $response->is_success() ) {
		$summary_content   = decode_json($content);
		$summary           = Builder::InfluxdbBuilder->summary_response($summary_content);
		$result->{summary} = $summary;
		return ( SUCCESS, $result, $summary_query );
	}
	else {
		return ( ERROR, $content, undef );
	}
}

sub build_series {
	my $self   = shift;
	my $result = shift;

	my $series_query = $builder->series_query();
	$self->app->log->debug( "series_query #-> " . Dumper($series_query) );
	my $response_container = $self->influxdb_query( $self->get_db_name(), $series_query, "pretty" );
	my $response           = $response_container->{'response'};
	my $content            = $response->{_content};

	my $series;
	if ( $response->is_success() ) {
		my $series_content = decode_json($content);
		$series = Builder::InfluxdbBuilder->series_response($series_content);
		my $series_node = "series";
		if ( defined($series) && ( keys $series ) ) {
			$result->{$series_node} = $series;
			my @series_values = $series->{values};
			my $series_count  = $#{ $series_values[0] };
			$result->{$series_node}{count} = $series_count;
		}
		return ( SUCCESS, $result, $series_query );
	}

	else {
		return ( ERROR, $content, undef );
	}
}

#TODO: drichardson -refactor the common fields
sub build_parameters {
	my $self          = shift;
	my $result        = shift;
	my $summary_query = shift;
	my $series_query  = shift;

	my $ds_name     = $self->param('deliveryServiceName');
	my $metric_type = $self->param('metricType');
	my $server_type = $self->param('serverType');
	my $start_date  = $self->param('startDate');
	my $end_date    = $self->param('endDate');
	my $interval    = $self->param('interval') || "1m";      # Valid interval examples 10m (minutes), 10s (seconds), 1h (hour)
	my $exclude     = $self->param('exclude');
	my $limit       = $self->param('limit');
	my $offset      = $self->param('offset');
	my $orderby     = $self->param('orderby');

	my $parent_node     = "query";
	my $parameters_node = "parameters";
	$result->{$parent_node}{$parameters_node}{deliveryServiceName} = $ds_name;
	$result->{$parent_node}{$parameters_node}{startDate}           = $start_date;
	$result->{$parent_node}{$parameters_node}{endDate}             = $end_date;
	$result->{$parent_node}{$parameters_node}{interval}            = $interval;
	$result->{$parent_node}{$parameters_node}{metricType}          = $metric_type;
	$result->{$parent_node}{$parameters_node}{orderby}             = $orderby;
	$result->{$parent_node}{$parameters_node}{limit}               = $limit;
	$result->{$parent_node}{$parameters_node}{offset}              = $offset;

	my $queries_node = "language";
	$result->{$parent_node}{$queries_node}{influxdbDatabaseName} = $self->get_db_name();
	$result->{$parent_node}{$queries_node}{influxdbSeriesQuery}  = $series_query;
	$result->{$parent_node}{$queries_node}{influxdbSummaryQuery} = $summary_query;

	return $result;
}

sub get_db_name {
	my $self      = shift;
	my $mode      = $self->app->mode;
	my $conf_file = MojoPlugins::InfluxDB->INFLUXDB_CONF_FILE_NAME;
	my $conf      = Utils::JsonConfig->load_conf( $mode, $conf_file );
	return $conf->{deliveryservice_stats_db_name};
}

1;
