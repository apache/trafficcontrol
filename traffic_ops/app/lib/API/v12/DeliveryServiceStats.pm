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
use JSON;
use HTTP::Date;
use Extensions::Delegate::Statistics;
use Utils::Helper::Extensions;
Utils::Helper::Extensions->use;
use Common::ReturnCodes qw(SUCCESS ERROR);

my $builder;

#TODO: drichardson
#      - Add required fields validation see lib/API/User.pm based on Validate::Tiny
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

	if ( $self->is_valid_delivery_service_name($ds_name) ) {
		if ( $self->is_delivery_service_name_assigned($ds_name) ) {

			my $stats = new Extensions::Delegate::Statistics(
				$self, {
					deliveryServiceName => $ds_name,
					metricType          => $metric_type,
					startDate           => $start_date,
					endDate             => $end_date,
					interval            => $interval,
					exclude             => $exclude,
					orderby             => $orderby,
					limit               => $limit,
					dbName              => $self->get_db_name(),
					offset              => $offset
				}
			);

			# Extensions Contract:
			#  "$rc": will be either SUCCESS or ERROR (****the implemented Extension should use the same constants for consistency)
			#  "$result": should always come back as hash that will be forwarded to the Client as JSON.
			my ( $rc, $result ) = $stats->get_stats();

			$self->app->log->debug( "top.rc #-> " . Dumper($rc) );

			#$self->app->log->debug( "top.result #-> " . Dumper($result) );

			if ( $rc == SUCCESS ) {
				return $self->success($result);
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

sub get_db_name {
	my $self      = shift;
	my $mode      = $self->app->mode;
	my $conf_file = MojoPlugins::InfluxDB->INFLUXDB_CONF_FILE_NAME;
	my $conf      = Utils::JsonConfig->load_conf( $mode, $conf_file );
	return $conf->{deliveryservice_stats_db_name};
}

1;
