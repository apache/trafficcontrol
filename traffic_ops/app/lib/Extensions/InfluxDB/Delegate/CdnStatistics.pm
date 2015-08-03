package Extensions::InfluxDB::Delegate::CdnStatistics;
#
# Copyright 2011-2014, Comcast Corporation. This software and its contents are
# Comcast confidential and proprietary. It cannot be used, disclosed, or
# distributed without Comcast's prior written permission. Modification of this
# software is only allowed at the direction of Comcast Corporation. All allowed
# modifications must be provided to Comcast Corporation.
#

# JvD Note: you always want to put Utils as the first use.
use UI::Utils;
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use JSON;
use constant SUCCESS => 0;
use constant ERROR   => 1;
use Utils::Helper::Extensions;
use Extensions::InfluxDB::Builder::CacheStatsBuilder;
use Extensions::InfluxDB::Builder::DeliveryServiceStatsBuilder;
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
		source      => "InfluxDB",
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

	$builder = new Extensions::InfluxDB::Builder::CacheStatsBuilder();

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
	$builder = new Extensions::InfluxDB::Builder::DeliveryServiceStatsBuilder();
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
