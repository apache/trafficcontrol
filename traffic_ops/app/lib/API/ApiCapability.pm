package API::ApiCapability;
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

use UI::Utils;

use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;

my $finfo = __FILE__ . ":";

my %valid_http_methods = map { $_ => 1 } ( 'GET', 'POST', 'PUT', 'PATCH', 'DELETE' );

sub index {
	my $self       = shift;
	my $capability = $self->param('capability');

	my %criteria;
	if ( defined $capability ) {
		$criteria{'me.capability'} = $capability;
	}
	my @data;
	my $orderby = "capability";
	$orderby = $self->param('orderby') if ( defined $self->param('orderby') );

	my $rs_data = $self->db->resultset("ApiCapability")->search( \%criteria, { prefetch => ['capability'], order_by => $orderby } );
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"id"          => $row->id,
				"httpMethod"  => $row->http_method,
				"httpRoute"   => $row->route,
				"capability"  => $row->capability->name,
				"lastUpdated" => $row->last_updated
			}
		);
	}
	$self->success( \@data );
}

sub show {
	my $self = shift;
	my $id   = $self->param('id');
	my $alt = "GET /api_capabilities";

	my $rs_data = $self->db->resultset("ApiCapability")->search( 'me.id' => $id );
	if ( !defined($rs_data) ) {
		return $self->with_deprecation("Resource not found.", "error", 404, $alt);
	}

	my @data = ();
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"id"          => $row->id,
				"httpMethod"  => $row->http_method,
				"httpRoute"   => $row->route,
				"capability"  => $row->capability->name,
				"lastUpdated" => $row->last_updated
			}
		);
	}
	$self->deprecation(200, $alt, \@data );
}

sub is_mapping_valid {
	my $self        = shift;
	my $id          = shift;
	my $http_method = shift;
	my $http_route  = shift;
	my $capability  = shift;

	if ( !defined($http_method) ) {
		return ( undef, "HTTP method is required." );
	}

	if ( !exists( $valid_http_methods{$http_method} ) ) {
		return ( undef, "HTTP method \'$http_method\' is invalid. Valid values are: " . join( ", ", sort keys %valid_http_methods ) );
	}

	if ( !defined($http_route) or $http_route eq "" ) {
		return ( undef, "Route is required." );
	}

	if ( !defined($capability) or $capability eq "" ) {
		return ( undef, "Capability name is required." );
	}

	# check if capability exists
	my $rs_data = $self->db->resultset("Capability")->search( { 'name' => { 'like', $capability } } )->single();
	if ( !defined($rs_data) ) {
		return ( undef, "Capability '$capability' does not exist." );
	}

	# search a mapping for the same http_method & route
	$rs_data = $self->db->resultset("ApiCapability")->search( { 'route' => { 'like', $http_route } } )->search(
		{
			'http_method' => { '=', $http_method }
		}
	)->single();

	# if adding a new entry, make sure it is unique
	if ( !defined($id) ) {
		if ( defined($rs_data) ) {
			my $allocated_capability = $rs_data->capability->name;
			return ( undef, "HTTP method '$http_method', route '$http_route' are already mapped to capability: $allocated_capability" );
		}
	}
	else {
		if ( defined($rs_data) ) {
			my $lid = $rs_data->id;
			if ( $lid ne $id ) {
				my $allocated_capability = $rs_data->capability->name;
				return ( undef, "HTTP method '$http_method', route '$http_route' are already mapped to capability: $allocated_capability" );
			}
		}
	}

	return ( 1, undef );
}

sub create {
	my $self   = shift;
	my $params = $self->req->json;

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	if ( !defined($params) ) {
		return $self->alert("Parameters must be in JSON format.");
	}

	my $http_method = $params->{httpMethod} if defined( $params->{httpMethod} );
	my $http_route  = $params->{httpRoute}  if defined( $params->{httpRoute} );
	my $capability  = $params->{capability} if defined( $params->{capability} );
	my $id          = undef;

	my ( $is_valid, $errStr ) = $self->is_mapping_valid( $id, $http_method, $http_route, $capability );
	if ( !$is_valid ) {
		return $self->alert($errStr);
	}

	my $values = {
		http_method => $http_method,
		route       => $http_route,
		capability  => $capability
	};

	my $insert = $self->db->resultset('ApiCapability')->create($values);
	my $rs     = $insert->insert();
	my $alt = "[NO ALTERNATE - See https://traffic-control-cdn.readthedocs.io/en/latest/api/api_capabilities.html#post]";
	if ($rs) {
		my $response;
		$response->{id}          = $rs->id;
		$response->{httpMethod}  = $rs->http_method;
		$response->{httpRoute}   = $rs->route;
		$response->{capability}  = $rs->capability->name;
		$response->{lastUpdated} = $rs->last_updated;

		&log( $self,
			"Created API-Capability mapping: '$response->{httpMethod}', '$response->{httpRoute}', '$response->{capability}' for id: " . $response->{id},
			"APICHANGE" );

		return $self->with_deprecation("API-Capability mapping was created.", "success", 200, $alt);
	}
	else {
		return $self->with_deprecation("API-Capability mapping creation failed.", "error", 500, $alt);
	}
}

sub update {
	my $self   = shift;
	my $id     = $self->param('id');
	my $params = $self->req->json;
	my $alt = "[NO ALTERNATE - See https://traffic-control-cdn.readthedocs.io/en/latest/api/api_capabilities_id.html#put]";

	if ( !&is_oper($self) ) {
		return $self->with_deprecation("Forbidden", "error", 403, $alt);
	}

	if ( !defined($params) ) {
		return $self->with_deprecation("Parameters must be in JSON format.", "error", 400, $alt);
	}

	my $http_method = $params->{httpMethod} if defined( $params->{httpMethod} );
	my $http_route  = $params->{httpRoute}  if defined( $params->{httpRoute} );
	my $capability  = $params->{capability} if defined( $params->{capability} );

	my $mapping = $self->db->resultset('ApiCapability')->find( { id => $id } );
	if ( !defined($mapping) ) {
		return $self->with_deprecation("Resource not found.", "error", 404, $alt);
	}

	my ( $is_valid, $errStr ) = $self->is_mapping_valid( $id, $http_method, $http_route, $capability );
	if ( !$is_valid ) {
		return $self->with_deprecation($errStr, "error", 400, $alt);
	}

	my $values = {
		http_method => $http_method,
		route       => $http_route,
		capability  => $capability
	};

	my $rs = $mapping->update($values);
	if ($rs) {
		my $response;
		$response->{id}          = $rs->id;
		$response->{httpMethod}  = $rs->http_method;
		$response->{httpRoute}   = $rs->route;
		$response->{capability}  = $rs->capability->name;
		$response->{lastUpdated} = $rs->last_updated;

		&log( $self,
			"Updated API-Capability mapping: '$response->{httpMethod}', '$response->{httpRoute}', '$response->{capability}' for id: " . $response->{id},
			"APICHANGE" );
		return $self->with_deprecation("API-Capability mapping was updated.", "success", 200, $alt, $response);
	}
	else {
		return $self->with_deprecation("API-Capability mapping update failed.", "error", 400, $alt);
	}
}

sub delete {
	my $self = shift;
	my $id   = $self->param('id');
	my $alt = "[NO ALTERNATE - See https://traffic-control-cdn.readthedocs.io/en/latest/api/api_capabilities_id.html#delete]";

	if ( !&is_oper($self) ) {
		return $self->with_deprecation("Forbidden", "error", 403, $alt);
	}

	my $mapping = $self->db->resultset('ApiCapability')->find( { id => $id } );
	if ( !defined($mapping) ) {
		return $self->with_deprecation("Resource not found.", "error", 404, $alt);
	}

	my $rs = $mapping->delete();
	if ($rs) {
		return $self->with_deprecation("API-Capability mapping deleted.", "success", 200, $alt);
	}
	else {
		return $self->with_deprecation("API-Capability mapping deletion failed.", "error", 400, $alt);
	}
}

1;
