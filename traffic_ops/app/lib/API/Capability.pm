package API::Capability;
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

sub index {
	my $self = shift;
	my @data;
	my $orderby = "name";
	$orderby = $self->param('orderby') if ( defined $self->param('orderby') );

	my $rs_data = $self->db->resultset("Capability")->search( undef, { order_by => $orderby } );
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"name"        => $row->name,
				"description" => $row->description,
				"lastUpdated" => $row->last_updated
			}
		);
	}
	$self->success( \@data );
}

sub name {
	my $self = shift;
	my $name = $self->param('name');

	my $alt = "GET /capabilities with the 'name' query parameter";

	my $rs_data = $self->db->resultset("Capability")->search( 'me.name' => $name );
	if ( !defined($rs_data) ) {
		return $self->with_deprecation("Resource not found.", "error", 404, $alt);
	}
	my @data = ();
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"name"        => $row->name,
				"description" => $row->description,
				"lastUpdated" => $row->last_updated
			}
		);
	}
	$self->deprecation(200, $alt, \@data );
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

	my $name        = $params->{name}        if defined( $params->{name} );
	my $description = $params->{description} if defined( $params->{description} );

	if ( !defined($name) or $name eq "" ) {
		return $self->alert("Name is required.");
	}

	if ( !defined($description) or $description eq "" ) {
		return $self->alert("Description is required.");
	}

	# check if capability exists
	my $rs_data = $self->db->resultset("Capability")->search( { 'name' => { 'like', $name } } )->single();
	if ( defined($rs_data) ) {
		return $self->alert("Capability '$name' already exists.");
	}

	my $values = {
		name        => $name,
		description => $description
	};

	my $insert = $self->db->resultset('Capability')->create($values);
	my $rs     = $insert->insert();
	if ($rs) {
		my $response;
		$response->{name}        = $rs->name;
		$response->{description} = $rs->description;

		&log( $self, "Created Capability: '$response->{name}', '$response->{description}'", "APICHANGE" );

		return $self->success( $response, "Capability was created." );
	}
	else {
		return $self->alert("Capability creation failed.");
	}
}

sub update {
	my $self   = shift;
	my $name   = $self->param('name');
	my $params = $self->req->json;

	if ( !&is_oper($self) ) {
		return $self->with_deprecation_with_no_alternative("Forbidden", "error", 403);
	}

	if ( !defined($params) ) {
		return $self->with_deprecation_with_no_alternative("Parameters must be in JSON format.", "error", 400);
	}

	my $description = $params->{description} if defined( $params->{description} );

	my $capability = $self->db->resultset('Capability')->find( { name => $name } );
	if ( !defined($capability) ) {
		return $self->with_deprecation_with_no_alternative("Resource not found.", "error", 404);
	}

	if ( !defined($description) or $description eq "" ) {
		return $self->with_deprecation_with_no_alternative("Description is required.", "error", 400);
	}

	my $values = { description => $description };

	my $rs = $capability->update($values);
	if ($rs) {
		my $response;
		$response->{name}        = $rs->name;
		$response->{description} = $rs->description;
		$response->{lastUpdated} = $rs->last_updated;

		&log( $self, "Updated Capability: '$response->{name}', '$response->{description}'", "APICHANGE" );

		return $self->with_deprecation_with_no_alternative("Capability was updated.", "success", 200, $response);
	}
	else {
		return $self->with_deprecation_with_no_alternative("Capability update failed.", "error", 400);
	}
}

sub delete {
	my $self = shift;
	my $name = $self->param('name');

	if ( !&is_oper($self) ) {
		return $self->with_deprecation_with_no_alternative("Forbidden", "error", 403);
	}

	my $capability = $self->db->resultset('Capability')->find( { name => $name } );
	if ( !defined($capability) ) {
		return $self->with_deprecation_with_no_alternative("Resource not found.", "error", 404);
	}

	# make sure no api_capability refers to this capability
	my $rs_data = $self->db->resultset("ApiCapability")->find( { 'me.capability' => $name } );
	if ( defined($rs_data) ) {
		my $reference_id = $rs_data->id;
		return $self->with_deprecation_with_no_alternative("Capability \'$name\' is refered by an api_capability mapping: $reference_id. Deletion failed.", "error", 400);
	}

	my $rs = $capability->delete();
	if ($rs) {
		return $self->with_deprecation_with_no_alternative("Capability deleted.", "success", 200);
	}
	else {
		return $self->with_deprecation_with_no_alternative("Capability deletion failed.", "error", 400);
	}
}

1;
