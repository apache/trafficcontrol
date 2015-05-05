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
use Data::Dumper;
use Builder::InfluxdbQuery;
use JSON;
my $iq;

sub index {
	my $self            = shift;
	my $cdn_name        = $self->param('cdnName');
	my $ds_name         = $self->param('deliveryServiceName');
	my $cachegroup_name = $self->param('cacheGroupName');
	my $metric_type     = $self->param('metricType');
	my $server_type     = $self->param('serverType');
	my $start_date      = $self->param('startDate');
	my $end_date        = $self->param('endDate');
	my $interval        = $self->param('interval') || "1m";      # Valid interval examples 10m (minutes), 10s (seconds), 1h (hour)
	my $exclude         = $self->param('exclude');
	my $limit           = $self->param('limit');
	my $offset          = $self->param('offset');

	if ( $self->is_valid_delivery_service_name($ds_name) ) {
		if ( $self->is_delivery_service_name_assigned($ds_name) ) {

			# Build the summary section
			$iq = new Builder::InfluxdbQuery(
				{
					cdn_name        => $cdn_name,
					series_name     => $metric_type,
					ds_name         => $ds_name,
					cachegroup_name => $cachegroup_name,
					start_date      => $start_date,
					end_date        => $end_date,
					interval        => $interval
				}
			);

			my $result = ();
			my $summary_query;

			my $include_summary = ( defined($exclude) && $exclude =~ /summary/ ) ? 0 : 1;
			if ($include_summary) {
				( $result, $summary_query ) = $self->build_summary($result);
			}

			my $include_series = ( defined($exclude) && $exclude =~ /series/ ) ? 0 : 1;
			my $series_query;
			if ($include_series) {
				( $result, $series_query ) = $self->build_series($result);
			}

			$result = $self->build_parameters( $result, $summary_query, $series_query );

			return $self->success($result);
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

	my $summary_query = $iq->summary_query();

	my $response_container = $self->influxdb_query( $self->get_db_name(), $summary_query );
	my $response           = $response_container->{'response'};
	my $content            = $response->{_content};

	my $summary;
	my $summary_content;
	my $series_count = 0;
	if ( $response->is_success() ) {
		$summary_content   = decode_json($content);
		$summary           = $iq->summary_response($summary_content);
		$result->{summary} = $summary;
	}
	else {
		return $self->alert( { error_message => $content } );
	}
	return ( $result, $summary_query );
}

sub build_series {
	my $self   = shift;
	my $result = shift;

	my $series_query       = $iq->series_query();
	my $response_container = $self->influxdb_query( $self->get_db_name(), $series_query );
	my $response           = $response_container->{'response'};
	my $content            = $response->{_content};

	my $series;
	if ( $response->is_success() ) {
		my $series_content = decode_json($content);
		$series = $iq->series_response($series_content);
	}
	else {
		return $self->alert( { error_message => $content } );
	}
	my $series_node = "series";
	if ( defined($series) && ( keys $series ) ) {
		$result->{$series_node}{data} = $series;
		my @series_values = $series->{values};
		my $series_count  = $#{ $series_values[0] };
		$result->{$series_node}{count} = $series_count;
	}
	return ( $result, $series_query );
}

sub build_parameters {
	my $self          = shift;
	my $result        = shift;
	my $summary_query = shift;
	my $series_query  = shift;

	my $cdn_name        = $self->param('cdnName');
	my $ds_name         = $self->param('deliveryServiceName');
	my $cachegroup_name = $self->param('cacheGroupName');
	my $metric_type     = $self->param('metricType');
	my $server_type     = $self->param('serverType');
	my $start_date      = $self->param('startDate');
	my $end_date        = $self->param('endDate');
	my $interval        = $self->param('interval') || "1m";      # Valid interval examples 10m (minutes), 10s (seconds), 1h (hour)
	my $exclude         = $self->param('exclude');
	my $limit           = $self->param('limit');
	my $offset          = $self->param('offset');

	my $parent_node = "query";

	my $parameters_node = "parameters";
	$result->{$parent_node}{$parameters_node}{cdnName}             = $cdn_name;
	$result->{$parent_node}{$parameters_node}{deliveryServiceName} = $ds_name;
	$result->{$parent_node}{$parameters_node}{cacheGroupName}      = $cachegroup_name;
	$result->{$parent_node}{$parameters_node}{startDate}           = $start_date;
	$result->{$parent_node}{$parameters_node}{endDate}             = $end_date;
	$result->{$parent_node}{$parameters_node}{interval}            = $interval;
	$result->{$parent_node}{$parameters_node}{metricType}          = $metric_type;
	my $queries_node = "sql";
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
