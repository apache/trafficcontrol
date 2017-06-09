package Extensions::TrafficStats::API::CacheStats;
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

# JvD Note: you always want to put Utils as the first use. Sh*t don't work if it's after the Mojo lines.
use UI::Utils;
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use JSON;
my $builder;
use Extensions::TrafficStats::Delegate::CacheStatistics;
use Utils::Helper;
use Utils::Helper::Extensions;
Utils::Helper::Extensions->use;
use Validate::Tiny ':all';
use Common::ReturnCodes qw(SUCCESS ERROR);

sub index {
	my $self        = shift;
	my $cdn_name    = $self->param('cdnName');
	my $metric_type = $self->param('metricType');
	my $start_date  = $self->param('startDate');
	my $end_date    = $self->param('endDate');

	my $query_parameters = {
		cdnName    => $cdn_name,
		metricType => $metric_type,
		startDate  => $start_date,
		endDate    => $end_date
	};

	my ( $is_valid, $result ) = $self->is_valid($query_parameters);

	if ($is_valid) {
		my $cstats = new Extensions::TrafficStats::Delegate::CacheStatistics( $self, $self->get_db_name() );
		my ( $rc, $result ) = $cstats->get_stats();

		if ( $rc == SUCCESS ) {
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

sub is_valid {
	my $self             = shift;
	my $query_parameters = shift;

	my $rules = {
		fields => [qw/cdnName metricType startDate endDate/],

		# Checks to perform on all fields
		checks => [

			# All of these are required
			[qw/cdnName metricType startDate endDate/] => is_required("query parameter is required"),

		]
	};

	# Validate the input against the rules
	my $result = validate( $query_parameters, $rules );

	if ( $result->{success} ) {
		return ( 1, $result->{data} );
	}
	else {
		return ( 0, $result->{error} );
	}
}

sub get_db_name {
	my $self      = shift;
	my $mode      = $self->app->mode;
	my $conf_file = MojoPlugins::InfluxDB->INFLUXDB_CONF_FILE_NAME;
	my $conf      = Utils::JsonConfig->load_conf( $mode, $conf_file );
	return $conf->{cache_stats_db_name};
}

sub get_stat {
	my $self = shift;
	my $database = shift;
	my $query = shift;

	my $response_container = $self->influxdb_query($database, $query);
	my $response           = $response_container->{'response'};
	my $content            = $response->{_content};
	my $summary_content;

	if ( $response->is_success() ) {
		$summary_content   = decode_json($content);
		return $summary_content->{results}[0]{series}[0]{values}[0][1];
	}

	return "";
}

sub current_stats {
	my $self = shift;
	my @stats;
	my $current_bw = $self->get_current_bandwidth();
	my $conns = $self->get_current_connections();
	my $capacity = $self->get_current_capacity();
	my $rs = $self->db->resultset('Cdn');
	while ( my $cdn = $rs->next ) {
		my $cdn_name = $cdn->name;
		my $bw = $current_bw->{$cdn_name};
		my $conn = $conns->{$cdn_name};
		my $cap = $capacity->{$cdn_name};
		push(@stats, ({cdn => $cdn_name, bandwidth => $bw, connections => $conn, capacity => $cap}));
	}
	push(@stats, ({cdn => "total", bandwidth => $current_bw->{"total"}, connections => $conns->{"total"}}));
	return $self->success({"currentStats" => \@stats});
}

sub get_current_bandwidth {
	my $self = shift;
	my $bw;
	my $total_bw = 0;
	my $rs = $self->db->resultset('Cdn');
	while ( my $cdn = $rs->next ) {
		my $cdn_name = $cdn->name;
		my $escaped_cdn_name = Utils::Helper::escape_influxql($cdn_name);
		my $query = "SELECT last(value) FROM \"monthly\".\"bandwidth.cdn.1min\" WHERE cdn = \'$escaped_cdn_name\'";
		my $bandwidth = $self->get_stat("cache_stats", $query);
		if ($bandwidth) {
			$bw->{$cdn_name} = $bandwidth/1000000;
			$total_bw += $bandwidth;
		}
	}
	$bw->{"total"} = $total_bw/1000000;
	return $bw;
}

sub get_current_connections {
	my $self = shift;
	my $conn;
	my $total_conn = 0;
	my $rs = $self->db->resultset('Cdn');
	while ( my $cdn = $rs->next ) {
		my $cdn_name = $cdn->name;
		my $escaped_cdn_name = Utils::Helper::escape_influxql($cdn_name);
		my $query = "select last(value) from \"monthly\".\"connections.cdn.1min\" where cdn = \'$escaped_cdn_name\'";
		my $connections = $self->get_stat("cache_stats", $query);
		if ($connections) {
			$conn->{$cdn_name} = $connections;
			$total_conn += $connections;
		}
	}
	$conn->{"total"} = $total_conn;
	return $conn;
}

sub get_current_capacity {
	my $self = shift;
	my $cap;
	my $rs = $self->db->resultset('Cdn');
	while ( my $cdn = $rs->next ) {
		my $cdn_name = $cdn->name;
		my $escaped_cdn_name = Utils::Helper::escape_influxql($cdn_name);
		my $query = "select last(value) from \"monthly\".\"maxkbps.cdn.1min\" where cdn = \'$escaped_cdn_name\'";
		my $capacity = $self->get_stat("cache_stats", $query);
		if ($capacity) {
		$capacity = $capacity/1000000; #convert to Gbps
		$capacity = $capacity * 0.85;    # need a better way to figure out percentage of max besides hard-coding
		$cap->{$cdn_name} = $capacity;
		}
	}
	return $cap;
}

sub daily_summary {
	my $self = shift;
	my $query = "";
	my $database = "daily_stats";
	my $total_bytesserved = 0;

	my $daily_stats;
	my @max_gbps;
	my @pb_served;
	#get cdns
	my @cdn_names = $self->db->resultset('Server')->search({ 'type.name' => { -like => 'EDGE%' } }, { prefetch => [ 'cdn', 'type' ], group_by => 'cdn.name' } )->get_column('cdn.name')->all();
	foreach my $cdn (@cdn_names) {
		my $bytes_served;
		my $max;
		#get max bw
		$max->{"cdn"} = $cdn;
		$bytes_served->{"cdn"} = $cdn;
		my $escaped_cdn = Utils::Helper::escape_influxql($cdn);
		my $max_bw = $self->get_stat($database, "select max(value) from \"daily_maxgbps\" where cdn = \'$escaped_cdn\'");
		$max->{"highest"} = $max_bw;
		#get last bw
		my $last_bw = $self->get_stat($database, "select last(value) from \"daily_maxgbps\" where cdn = \'$escaped_cdn\'");
		$max->{"yesterday"} = $last_bw;
		push(@max_gbps, $max);
		#get bytesserved
		my $bytesserved = $self->get_stat($database, "select sum(value) from \"daily_bytesserved\" where cdn = \'$escaped_cdn\'");
		$bytes_served->{"bytesServed"} = $bytesserved/1000;
		push(@pb_served, $bytes_served);
		$total_bytesserved += $bytesserved;
	}
	push(@pb_served, ({cdn => "total", bytesServed => $total_bytesserved/1000}));
	$daily_stats->{"maxGbps"} = \@max_gbps;
	$daily_stats->{"petaBytesServed"} = \@pb_served;
	return $self->success({%$daily_stats});
}

1;
