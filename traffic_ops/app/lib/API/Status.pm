package API::Status;
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

sub index {
	my $self = shift;
	my @data;
	my $orderby = $self->param('orderby') || "name";
	my $rs_data = $self->db->resultset("Status")->search( undef, { order_by => $orderby } );
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"id"          => $row->id,
				"name"        => $row->name,
				"description" => $row->description,
				"lastUpdated" => $row->last_updated
			}
		);
	}
	$self->success( \@data );
}

sub show {
	my $self = shift;
	my $id   = $self->param('id');

	my $rs_data = $self->db->resultset("Status")->search( { id => $id } );
	my @data = ();
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"id"          => $row->id,
				"name"        => $row->name,
				"description" => $row->description,
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

	my $status = $self->db->resultset('Status')->find( { id => $id } );
	if ( !defined($status) ) {
		return $self->not_found();
	}

	if ( !defined( $params->{name} ) ) {
		return $self->alert("Status name is required.");
	}

	my $values = {
		name 		=> $params->{name},
		description => $params->{description}
	};

	my $rs = $status->update($values);
	if ($rs) {
		my $response;
		$response->{id}          = $rs->id;
		$response->{name}        = $rs->name;
		$response->{description} = $rs->description;
		$response->{lastUpdated} = $rs->last_updated;
		&log( $self, "Updated Status name '" . $rs->name . "' for id: " . $rs->id, "APICHANGE" );
		return $self->success( $response, "Status update was successful." );
	}
	else {
		return $self->alert("Status update failed.");
	}

}

sub create {
	my $self   = shift;
	my $params = $self->req->json;

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	my $name = $params->{name};
	if ( !defined($name) ) {
		return $self->alert("Status name is required.");
	}

	my $values = {
		name 			=> $params->{name} ,
		description 	=> $params->{description}
	};

	my $insert = $self->db->resultset('Status')->create($values);
	my $rs = $insert->insert();
	if ($rs) {
		my $response;
		$response->{id}          	= $rs->id;
		$response->{name}        	= $rs->name;
		$response->{description}    = $rs->description;
		$response->{lastUpdated} 	= $rs->last_updated;

		&log( $self, "Created Status name '" . $rs->name . "' for id: " . $rs->id, "APICHANGE" );

		return $self->success( $response, "Status create was successful." );
	}
	else {
		return $self->alert("Status create failed.");
	}

}

sub delete {
	my $self = shift;
	my $id     = $self->param('id');

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	my $status = $self->db->resultset('Status')->find( { id => $id } );
	if ( !defined($status) ) {
		return $self->not_found();
	}

	my $servers = $self->db->resultset('Server')->find( { status => $status->id } );
	if ( defined($servers) ) {
		return $self->alert("This status is currently used by servers.");
	}

	my $rs = $status->delete();
	if ($rs) {
		return $self->success_message("Status deleted.");
	} else {
		return $self->alert( "Status delete failed." );
	}
}





1;
