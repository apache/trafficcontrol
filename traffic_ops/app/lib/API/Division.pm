package API::Division;

# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.


use UI::Utils;
use UI::Division;
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use JSON;
use MojoPlugins::Response;

sub index {
	my $self = shift;
	my @data;
	my $orderby = $self->param('orderby') || "name";
	my $rs_data = $self->db->resultset("Division")->search( undef, { order_by => 'me.' . $orderby } );
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"id"          => $row->id,
				"name"        => $row->name,
				"lastUpdated" => $row->last_updated
			}
		);
	}
	$self->success( \@data );
}

sub index_by_name {
	my $self = shift;
	my $name   = $self->param('name');

	my $rs_data = $self->db->resultset("Division")->search( { name => $name } );
	my @data = ();
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"id"          => $row->id,
				"name"        => $row->name,
				"lastUpdated" => $row->last_updated
			}
		);
	}
	$self->deprecation(200, "GET /divisions with the 'name' parameter", \@data );
}


sub show {
	my $self = shift;
	my $id   = $self->param('id');

	my $rs_data = $self->db->resultset("Division")->search( { id => $id } );
	my @data = ();
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"id"          => $row->id,
				"name"        => $row->name,
				"lastUpdated" => $row->last_updated
			}
		);
	}
	$self->success( \@data );
}

sub update {
	my $self   = shift;
	my $id     = $self->param('id');
	my $params = $self->req->json;

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	my $division = $self->db->resultset('Division')->find( { id => $id } );
	if ( !defined($division) ) {
		return $self->not_found();
	}

	if ( !defined($params) ) {
		return $self->alert("parameters must be in JSON format.");
	}

	if ( !defined( $params->{name} ) ) {
		return $self->alert("Division name is required.");
	}

	my $values = { name => $params->{name} };

	my $rs = $division->update($values);
	if ($rs) {
		my $response;
		$response->{id}          = $rs->id;
		$response->{name}        = $rs->name;
		$response->{lastUpdated} = $rs->last_updated;
		&log( $self, "Updated Division name '" . $rs->name . "' for id: " . $rs->id, "APICHANGE" );
		return $self->success( $response, "Division update was successful." );
	}
	else {
		return $self->alert("Division update failed.");
	}

}

sub create {
	my $self   = shift;
	my $params = $self->req->json;
	if ( !defined($params) ) {
		return $self->alert("parameters must be in JSON format,  please check!");
	}

	if ( !&is_oper($self) ) {
		return $self->alert( { Error => " - You must be an ADMIN or OPER to perform this operation!" } );
	}

	my $name = $params->{name};
	if ( !defined($name) ) {
		return $self->alert("division 'name' is not given.");
	}

	#Check for duplicate division name
	my $existing_division = $self->db->resultset('Division')->search( { name => $name } )->get_column('name')->single();
	if ($existing_division) {
		return $self->alert("A division with name \"$name\" already exists.");
	}

	my $insert = $self->db->resultset('Division')->create( { name => $name } );
	$insert->insert();

	my $response;
	my $rs = $self->db->resultset('Division')->find( { id => $insert->id } );
	if ( defined($rs) ) {
		$response->{id}   = $rs->id;
		$response->{name} = $rs->name;
		return $self->success($response);
	}
	return $self->alert("create division failed.");
}

sub delete {
	my $self = shift;
	my $id     = $self->param('id');

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	my $division = $self->db->resultset('Division')->find( { id => $id } );
	if ( !defined($division) ) {
		return $self->not_found();
	}

	my $regions = $self->db->resultset('Region')->find( { division => $division->id } );
	if ( defined($regions) ) {
		return $self->alert("This division is currently used by regions.");
	}

	my $rs = $division->delete();
	if ($rs) {
		return $self->success_message("Division deleted.");
	} else {
		return $self->alert( "Division delete failed." );
	}
}

sub delete_by_name {
	my $self = shift;
	my $name = $self->param('name');

	my $alt = "DELETE /divsions/{{ID}}";

	if ( !&is_oper($self) ) {
		return $self->with_deprecation("Forbidden", "error", 403, $alt);
	}

	my $division = $self->db->resultset('Division')->find( { name => $name } );
	if ( !defined($division) ) {
		return $self->with_deprecation("Resource not found.", "error", 404, $alt);
	}

	my $regions = $self->db->resultset('Region')->find( { division => $division->id } );
	if ( defined($regions) ) {
		return $self->with_deprecation("This division is currently used by regions.", "error", 400, $alt);
	}


	my $rs = $division->delete();
	if ($rs) {
		return $self->with_deprecation("Division deleted.", "success", 200, $alt);
	} else {
		return $self->with_deprecation("Division delete failed.", "error", 400, $alt);
	}
}



1;
