package Extensions::TrafficStats::Delegate::Statistics;

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

use Data::Dumper;
use Time::HiRes qw(gettimeofday tv_interval);
use Math::Round qw(nearest);
use JSON;
use POSIX qw(strftime);
use Extensions::TrafficStats::Builder::DeliveryServiceStatsBuilder;
use Extensions::TrafficStats::Helper::InfluxResponse;
use HTTP::Date;
use Utils::Helper::DateHelper;
use Carp qw(cluck confess);
use Common::ReturnCodes qw(SUCCESS ERROR);
use Utils::Deliveryservice;
use Time::Seconds;
use Time::Piece;
use DateTime::Format::ISO8601;
use constant ONE_DAY_IN_SECONDS          => 86400;
use constant THREE_DAYS                  => ONE_DAY * 3;
use constant SECONDS_IN_CAPTURE_INTERVAL => 60;

# constants do not interpolate
my $delim = ":";

my $builder;
my $mojo;
my $db_name;

sub new {
	my $self  = {};
	my $class = shift;
	$mojo    = shift;
	$db_name = shift;
	return ( bless( $self, $class ) );
}

sub info {
	return {
		name        => "Statistics",
		version     => "1.2",
		source      => "TrafficStats",
		info_url    => "",
		description => "Statistics Stub",
		isactive    => 1,
		script_file => "",
	};
}

sub set_info {
	my $self   = shift;
	my $result = shift;
	$result->{version} = $self->info()->{version};
	$result->{source}  = $self->info()->{source};

}

# TrafficStats
sub get_stats {
	my $self = shift;

	# version 1.2 parameters
	my $ds_name     = $mojo->param("deliveryServiceName");
	my $metric_type = $mojo->param("metricType");
	my $server_type = $mojo->param("serverType");
	my $start_date  = $mojo->param("startDate");
	my $end_date    = $mojo->param("endDate");
	my $stats_only  = $mojo->param("stats");
	my $data_only   = $mojo->param("data");
	my $type        = $mojo->param("type");
	my $interval    = $mojo->param("interval");
	my $match       = $mojo->param("match");
	my $exclude     = $mojo->param("exclude");
	my $orderby     = $mojo->param("orderby");
	my $limit       = $mojo->param("limit");
	my $offset      = $mojo->param("offset");

	# This parameter allows the API to override the retention period because
	# We can't wait for 30 days for data build up if it hasn't been 30 days yet.
	my $retention_period_in_days = $mojo->param("retentionPeriodInDays");

	# Build the summary section
	$builder = new Extensions::TrafficStats::Builder::DeliveryServiceStatsBuilder(
		{
			deliveryServiceName => $ds_name,
			metricType          => $metric_type,
			startDate           => $start_date,
			endDate             => $end_date,
			dbName              => $db_name,
			interval            => $interval,
			orderby             => $orderby,
			exclude             => $exclude,
			limit               => $limit,
			offset              => $offset,
		}
	);

	my $summary_query;
	my $rc = SUCCESS;
	my $result;

	my $include_summary = ( defined($exclude) && $exclude =~ /summary/ ) ? 0 : 1;
	if ($include_summary) {
		( $rc, $result, $summary_query ) = $self->build_summary( $metric_type, $start_date, $end_date, $result );
	}

	if ( $rc == SUCCESS ) {
		my $include_series = ( defined($exclude) && $exclude =~ /series/ ) ? 0 : 1;
		my $series_query;
		if ($include_series) {
			( $rc, $result, $series_query ) = $self->build_series($result);
		}
		if ( $rc == SUCCESS ) {
			$result = build_parameters( $self, $result, $summary_query, $series_query );
		}
		else {
			return ( ERROR, $result );
		}
	}
	else {
		return ( ERROR, $result );
	}
	$self->set_info($result);
	return ( SUCCESS, $result );
}

sub build_summary {
	my $self        = shift;
	my $metric_type = shift;
	my $start_date  = shift;
	my $end_date    = shift;
	my $result      = shift;

	my $summary_query = $builder->summary_query();
	$mojo->app->log->debug( "summary_query #-> " . Dumper($summary_query) );

	my $response_container = $mojo->influxdb_query( $db_name, $summary_query );
	my $response           = $response_container->{'response'};
	my $content            = $response->{_content};

	my $summary;
	my $summary_content;
	my $series_count = 0;
	if ( $response->is_success() ) {
		$summary_content = decode_json($content);

		my $ib = Extensions::TrafficStats::Builder::BaseBuilder->new($mojo);
		$summary = $ib->summary_response($summary_content);
		$result->{summary} = $summary;
		$self->build_totals( $metric_type, $result );
		return ( SUCCESS, $result, $summary_query );
	}
	else {
		return ( ERROR, $content, undef );
	}
}

sub build_totals {
	my $self        = shift;
	my $metric_type = shift;
	my $summary     = shift;
	my $average     = $summary->{summary}{average};
	my $count       = $summary->{summary}{count};

	# Use intervalInSeconds to calculate total for the time period.
	#  Default is 10s, but can be overridden by an extension.
	my $interval_in_sec = $summary->{summary}{intervalInSeconds} // SECONDS_IN_CAPTURE_INTERVAL;

	my $total_tps = ( $count * $interval_in_sec ) * $average;

	if ( $metric_type =~ /kbps/ ) {

		#we divide by 8 bytes for totalBytes
		$summary->{summary}{totalBytes}        = $total_tps / 8;
		$summary->{summary}{totalTransactions} = undef;
	}
	else {
		$summary->{summary}{totalBytes}        = undef;
		$summary->{summary}{totalTransactions} = $total_tps;
	}

}

sub build_series {
	my $self   = shift;
	my $result = shift;

	my $series_query = $builder->series_query();
	$mojo->app->log->debug( "series_query #-> " . Dumper($series_query) );
	my $response_container = $mojo->influxdb_query( $db_name, $series_query, "pretty" );
	my $response           = $response_container->{'response'};
	my $content            = $response->{_content};

	my $series;
	if ( $response->is_success() ) {

		my $series_content = decode_json($content);
		my $ib             = Extensions::TrafficStats::Builder::BaseBuilder->new($mojo);
		$series = $ib->series_response($series_content);
		my $series_node = "series";
		if ( defined($series) && ( ref($series) eq "HASH" ) ) {
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

# Append to the incoming result hash the additional sections.
sub build_parameters {
	my $self            = shift;
	my $result          = shift;
	my $summary_query   = shift;
	my $series_query    = shift;
	my $cachegroup_name = $mojo->param("cacheGroupName");
	my $ds_name         = $mojo->param("deliveryServiceName");
	my $metric_type     = $mojo->param("metricType");
	my $start_date      = $mojo->param("startDate");
	my $end_date        = $mojo->param("endDate");
	my $interval        = $mojo->param("interval");
	my $host_name       = $mojo->param("hostName");
	my $orderby         = $mojo->param("orderby");
	my $limit           = $mojo->param("limit");
	my $exclude         = $mojo->param("exclude");
	my $offset          = $mojo->param("offset");

	my $parent_node     = "query";
	my $parameters_node = "parameters";
	$result->{$parent_node}{$parameters_node}{deliveryServiceName} = $ds_name;
	$result->{$parent_node}{$parameters_node}{startDate}           = $start_date;
	$result->{$parent_node}{$parameters_node}{endDate}             = $end_date;
	$result->{$parent_node}{$parameters_node}{interval}            = $interval;
	$result->{$parent_node}{$parameters_node}{metricType}          = $metric_type;
	$result->{$parent_node}{$parameters_node}{orderby}             = $orderby;
	$result->{$parent_node}{$parameters_node}{limit}               = $limit;
	$result->{$parent_node}{$parameters_node}{exclude}             = $exclude;
	$result->{$parent_node}{$parameters_node}{offset}              = $offset;

	my $queries_node = "language";
	if ( defined($series_query) ) {
		$result->{$parent_node}{$queries_node}{influxdbDatabaseName} = $db_name;
		$result->{$parent_node}{$queries_node}{influxdbSeriesQuery}  = $series_query;
		$result->{$parent_node}{$queries_node}{influxdbSummaryQuery} = $summary_query;
	}

	return $result;
}

1;
