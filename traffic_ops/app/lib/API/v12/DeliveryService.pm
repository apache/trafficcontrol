package API::v12::DeliveryService;
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
use JSON;

sub stats {
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

1;
