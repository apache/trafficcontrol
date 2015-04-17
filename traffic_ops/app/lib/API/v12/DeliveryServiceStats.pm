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
use Utils::Helper;
use Helper::DeliveryServiceStats;
use JSON;
use constant DB_NAME => "deliveryservice_stats";

sub stats2 {
	my $self            = shift;
	my $dsid            = $self->param('ds');
	my $cachegroup_name = $self->param('cacheGroupName');
	my $metric          = $self->param('metricType');
	my $start_date      = $self->param('startDate');
	my $end_date        = $self->param('endDate');
	my $interval        = $self->param('interval') || "1m";         # Valid interval examples 10m (minutes), 10s (seconds), 1h (hour)
	my $limit           = $self->param('limit') || 50;              # How many "series" to limit by
	my $helper          = new Utils::Helper( { mojo => $self } );
	if ( $helper->is_valid_delivery_service($dsid) ) {

		if ( $helper->is_delivery_service_assigned($dsid) ) {
			return $self->v12_deliveryservice_stats( $dsid, $cachegroup_name, $metric, $start_date, $end_date, $interval );
		}
		else {
			return $self->forbidden();
		}
	}
	else {
		$self->success( {} );
	}

}

sub index {
	my $self            = shift;
	my $dsid            = $self->param('dsid');
	my $cachegroup_name = $self->param('cacheGroupName');
	my $metric_type     = $self->param('metricType');
	my $start_date      = $self->param('startDate');
	my $end_date        = $self->param('endDate');
	my $interval        = $self->param('interval');         # Valid interval examples 10m (minutes), 10s (seconds), 1h (hour)
	my $limit           = $self->param('limit');

	my ( $cdn_name, $ds_name ) = $self->deliveryservice_lookup_cdn_name_and_ds_name($dsid);

	my $series_name = Helper::DeliveryServiceStats->series_name( $cdn_name, $ds_name, $cachegroup_name, $metric_type );

	# Build the summary section
	my $summary_query = Helper::DeliveryServiceStats->build_summary_query( $series_name, $start_date, $end_date, $interval, $limit );
	my ( $summary, $series_count ) = $self->get_summary($summary_query);

	# Build the series section
	my $series_query = Helper::DeliveryServiceStats->build_series_query( $series_name, $start_date, $end_date, $interval, $limit );
	my $series = $self->get_series($series_query);
	if ( defined($summary) && defined($series) ) {

		my $parent_node = "stats";
		my $result      = ();
		$result->{$parent_node}{series}               = $series;
		$result->{$parent_node}{seriesCount}          = $series_count;
		$result->{$parent_node}{cdnName}              = $cdn_name;
		$result->{$parent_node}{deliveryServiceName}  = $ds_name;
		$result->{$parent_node}{cacheGroupName}       = $cachegroup_name;
		$result->{$parent_node}{startDate}            = $start_date;
		$result->{$parent_node}{endDate}              = $end_date;
		$result->{$parent_node}{interval}             = int($interval);
		$result->{$parent_node}{metricType}           = $metric_type;
		$result->{$parent_node}{influxdbDatabaseName} = DB_NAME;
		$result->{$parent_node}{influxdbSeriesQuery}  = $series_query;
		$result->{$parent_node}{influxdbSummaryQuery} = $summary_query;
		$result->{$parent_node}{summary}              = $summary;
		return $self->success($result);
	}
	else {
		return $self->alert("Could not retrieve the summary or the series");
	}

}

sub get_summary {
	my $self          = shift;
	my $summary_query = shift;

	my $response_container = $self->influxdb_query( DB_NAME, $summary_query );
	my $response           = $response_container->{'response'};
	my $content            = $response->{_content};

	my $summary;
	my $series_count;
	if ( $response->is_success() ) {
		my $summary_content = decode_json($content);
		( $summary, $series_count ) = Helper::DeliveryServiceStats->build_summary($summary_content);
		return ( $summary, $series_count );
	}
	else {
		my $rc = $response->{_rc};
		return $self->alert( $content, $rc );
	}
}

sub get_series {
	my $self         = shift;
	my $series_query = shift;

	my $response_container = $self->influxdb_query( DB_NAME, $series_query );
	my $response           = $response_container->{'response'};
	my $content            = $response->{_content};

	if ( $response->is_success() ) {
		my $series_content = decode_json($content);
		return Helper::DeliveryServiceStats->build_series($series_content);
	}
	else {
		my $rc = $response->{_rc};
		return $self->alert( $content, $rc );
	}
}

1;
