package Extensions::TrafficStats::Delegate::CdnStatistics;

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

# JvD Note: you always want to put Utils as the first use.
use UI::Utils;
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use JSON;
use constant SUCCESS => 0;
use constant ERROR   => 1;
use Utils::Helper::Extensions;
use Extensions::TrafficStats::Builder::CacheStatsBuilder;
use Extensions::TrafficStats::Builder::DeliveryServiceStatsBuilder;
Utils::Helper::Extensions->use;

my $builder;
my $mojo;
my $deliveryservice_stats_db_name;
my $cache_stats_db_name;

sub new {
	my $self  = {};
	my $class = shift;
	$mojo                          = shift;
	$cache_stats_db_name           = shift;
	$deliveryservice_stats_db_name = shift;
	return ( bless( $self, $class ) );
}

sub info {
	return {
		name        => "CdnStatistics",
		version     => "1.2",
		info_url    => "",
		source      => "TrafficStats",
		description => "Cdn Statistics Stub",
		isactive    => 1,
		script_file => "Extensions::Delegate::CdnStatistics",
	};
}

sub set_info {
	my $self   = shift;
	my $result = shift;
	$result->{version} = $self->info()->{version};
	$result->{source}  = $self->info()->{source};

}

sub get_usage_overview {
	my $self = shift;

	$builder = new Extensions::TrafficStats::Builder::CacheStatsBuilder();

	# ---------------------
	# maxGbps
	# ---------------------
	my $query = $builder->usage_overview_max_kbps_query();
	my $rc;
	my $result;
	my $response;
	my $stat_value;
	( $rc, $response, $stat_value ) = $self->lookup_stat( $cache_stats_db_name, $query );

	if ( $rc == SUCCESS ) {
		$result->{maxGbps} = $stat_value;
	}
	else {
		return ( ERROR, $response, undef );
	}

	# ---------------------
	# currentGbps
	# ---------------------
	$query = $builder->usage_overview_current_gbps_query();
	( $rc, $response, $stat_value ) = $self->lookup_stat( $cache_stats_db_name, $query );

	if ( $rc == SUCCESS ) {
		$self->set_info($result);
		$result->{currentGbps} = $stat_value;
	}
	else {
		return ( ERROR, $response, undef );
	}

	# ---------------------
	# tps
	# ---------------------
	$query   = $builder->usage_overview_current_gbps_query();
	$builder = new Extensions::TrafficStats::Builder::DeliveryServiceStatsBuilder();
	$query   = $builder->usage_overview_tps_query();
	( $rc, $response, $stat_value ) = $self->lookup_stat( $deliveryservice_stats_db_name, $query );

	if ( $rc == SUCCESS ) {
		$result->{tps} = int($stat_value);
	}
	else {
		return ( ERROR, $response, undef );
	}

	return ( SUCCESS, $result, $query );
}

sub lookup_stat {
	my $self    = shift;
	my $db_name = shift;
	my $query   = shift;

	my $response_container = $mojo->influxdb_query( $db_name, $query );
	my $response           = $response_container->{'response'};
	my $json_content       = $response->{_content};

	my $result;
	my $summary;
	my $stat_value;

	if ( $response->is_success() ) {
		my $content    = decode_json($json_content);
		my $stat_value = $content->{results}[0]{series}[0]->{values}[0][1];
		return ( SUCCESS, $result, $stat_value );
	}
	else {
		return ( ERROR, $json_content, undef );
	}
}

1;
