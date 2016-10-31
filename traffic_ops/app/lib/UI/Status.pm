package UI::Status;
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
	my $orderby = "name";
	$orderby = $self->param('orderby') if ( defined $self->param('orderby') );
	my $rs_data = $self->db->resultset("Status")->search( undef, { order_by => $orderby } );
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"id"           => $row->id,
				"name"         => $row->name,
				"description"  => $row->description,
				"last_updated" => $row->last_updated,
			}
		);
	}
	$self->render( json => \@data );
}

sub delete {
	my $self = shift;
	my $id   = $self->param('id');

	if ( !&is_admin($self) ) {
		$self->flash( alertmsg => "You must be an ADMIN to perform this operation!" );
	}
	else {
		my $p_name = $self->db->resultset('Status')->search( { id => $id } )->get_column('name')->single();
		my $delete = $self->db->resultset('Status')->search( { id => $id } );
		$delete->delete();
		&log( $self, "Delete status " . $p_name, "UICHANGE" );
	}
	return $self->redirect_to('/misc');
}

sub check_status_input {
	my $self = shift;

	my $sep = "__NEWLINE__";    # the line separator sub that with \n in the .ep javascript
	my $err = undef;

	# First, check permissions
	if ( !&is_admin($self) ) {
		$err .= "You must be an ADMIN to perform this operation!" . $sep;
		return $err;
	}

	return $err;
}

sub update {
	my $self = shift;
	my $id   = $self->param('id');

	my $err = &check_status_input($self);
	if ( defined($err) ) {
		$self->flash( alertmsg => $err );
	}
	else {
		my $update = $self->db->resultset('Status')->find( { id => $self->param('id') } );
		$update->name( $self->param('name') );
		$update->description( $self->param('description') );
		$update->update();

		# if the update has failed, we don't even get here, we go to the exception page.
		&log( $self, "Update status with name " . $self->param('name') . " and id " . $self->param('id'), "UICHANGE" );
	}

	$self->flash( alertmsg => "Success!" );
	return $self->redirect_to('/misc');
}

sub create {
	my $self = shift;

	my $err = &check_status_input($self);
	if ( defined($err) ) {
		$self->flash( alertmsg => $err );
	}
	else {
		my $insert = $self->db->resultset('Status')->create(
			{
				name        => $self->param('name'),
				description => $self->param('description'),
			}
		);
		$insert->insert();
	}
	$self->flash( alertmsg => "Success!" );
	return $self->redirect_to('/misc');
}

sub is_valid_status {
	my $self   = shift;
	my $status = shift;
	my $valid  = 0;

	my $row = $self->db->resultset("Status")->search( { name => $status } )->single;

	if ( defined($row) ) {
		return ( $row->id );
	}
	else {
		return (undef);
	}
}

1;
