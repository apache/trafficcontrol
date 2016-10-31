package Extensions::TrafficStats::API::DeliveryServiceStats;
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
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use Extensions::TrafficStats::Delegate::Statistics;
use Common::ReturnCodes qw(SUCCESS ERROR);
use Validate::Tiny ':all';
use UI::Utils;

my $builder;

sub index {
	my $self             = shift;
	my $ds_name          = $self->param('deliveryServiceName');
	my $metric_type      = $self->param('metricType');
	my $start_date       = $self->param('startDate');
	my $end_date         = $self->param('endDate');
	my $query_parameters = { deliveryServiceName => $ds_name, metricType => $metric_type, startDate => $start_date, endDate => $end_date };

	my ( $is_valid, $result ) = $self->is_valid($query_parameters);
	if ($is_valid) {
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
				return $self->forbidden("Forbidden. Delivery service not assigned to user.");
			}
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
		fields => [qw/deliveryServiceName metricType startDate endDate/],

		# Checks to perform on all fields
		checks => [

			# All of these are required
			[qw/deliveryServiceName metricType startDate endDate/] => is_required("query parameter is required"),

		]
	};

	# Validate the input against the rules
	my $result = validate( $query_parameters, $rules );

	if ( $result->{success} ) {

		#print "success: " . dump( $result->{data} );
		return ( 1, $result->{data} );
	}
	else {
		print "failed " . Dumper( $result->{error} );
		return ( 0, $result->{error} );
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
