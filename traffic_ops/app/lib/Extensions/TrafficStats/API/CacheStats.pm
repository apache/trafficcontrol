package Extensions::TrafficStats::API::CacheStats;
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
use JSON;
my $builder;
use Extensions::TrafficStats::Delegate::CacheStatistics;
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

1;
