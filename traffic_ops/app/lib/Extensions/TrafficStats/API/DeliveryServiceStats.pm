package Extensions::TrafficStats::API::DeliveryServiceStats;
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
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use Extensions::TrafficStats::Delegate::Statistics;
use Common::ReturnCodes qw(SUCCESS ERROR);

my $builder;

sub index {
	my $self    = shift;
	my $ds_name = $self->param('deliveryServiceName');

	if ( $self->is_valid_delivery_service_name($ds_name) ) {
		if ( $self->is_delivery_service_name_assigned($ds_name) || &is_admin($self) || &is_oper($self) ) {

			my $stats = new Extensions::TrafficStats::Delegate::Statistics( $self, $self->get_db_name() );

			my ( $rc, $result ) = $stats->get_stats();
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
