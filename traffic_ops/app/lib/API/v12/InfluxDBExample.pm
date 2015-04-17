package API::v12::InfluxDBExample;
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
use JSON;
use Data::Dumper;
use Helper::DeliveryStats;

sub query {
	my $self  = shift;
	my $query = $self->param("q");

	my $response_container = $self->influxdb_query( "deliveryservice_stats", $query );
	my $response           = $response_container->{'response'};
	my $content            = $response->{_content};
	if ( $response->is_success() ) {
		return $self->success( decode_json($content) );
	}
	else {
		my $rc = $response->{_rc};
		return $self->alert( $content, $rc );
	}
}

sub write_point {
	my $self = shift;

	my $write_points       = $self->req->json;
	my $response_container = $self->influxdb_write($write_points);
	my $response           = $response_container->{response};
	my $rc                 = $response->{_rc};
	my $content            = $response->{_msg};
	if ( $rc == 200 ) {
		return $self->success($content);
	}
	else {
		return $self->alert( $content, $rc );
	}
}
1;
