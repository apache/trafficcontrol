package UI::Types;
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

# Table view
sub index {
	my $self = shift;

	&navbarpage($self);
}

sub add {
	my $self = shift;
	$self->stash( fbox_layout => 1, type_data => {} );
	&stash_role($self);
	if ( $self->stash('priv_level') < 30 ) {
		$self->stash( alertmsg => "Insufficient privileges!" );
		$self->redirect_to('/types');
	}
}

sub view {
	my $self = shift;
	my $mode = $self->param('mode');
	my $id   = $self->param('id');

	&stash_role($self);

	$self->stash( type_data => {} );

	my $rs = $self->db->resultset('Type')->search( { id => $id } );
	my $data = $rs->single;

	$self->stash( type_data   => $data );
	$self->stash( fbox_layout => 1 );

	if ( $mode eq "edit" and $self->stash('priv_level') > 20 ) {
		$self->render( template => 'types/edit' );
	}
	else {
		$self->render( template => 'types/view' );
	}
}

# Read
sub readtype {
	my $self = shift;
	my @data;
	my $orderby = "name";
	$orderby = $self->param('orderby') if ( defined $self->param('orderby') );
	my $rs_data = $self->db->resultset("Type")->search( undef, { order_by => $orderby } );
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"id"           => $row->id,
				"name"         => $row->name,
				"description"  => $row->description,
				"use_in_table" => $row->use_in_table,
				"last_updated" => $row->last_updated,
			}
		);
	}
	$self->render( json => \@data );
}

# Read
sub readtypetrimmed {
	my $self = shift;
	my @data;
	my $orderby = "name";
	$orderby = $self->param('orderby') if ( defined $self->param('orderby') );
	my $rs_data = $self->db->resultset("Type")->search( undef, { order_by => $orderby } );
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"name" => $row->name,
			}
		);
	}
	$self->render( json => \@data );
}

# Delete
sub delete {
	my $self = shift;
	my $id   = $self->param('id');

	if ( !&is_admin($self) ) {
		$self->flash( alertmsg => "You must be an ADMIN to perform this operation!" );
	}
	else {
		my $p_name = $self->db->resultset('Type')->search( { id => $id } )->get_column('name')->single();
		my $delete = $self->db->resultset('Type')->search( { id => $id } );
		$delete->delete();
		&log( $self, "Delete type " . $p_name, "UICHANGE" );
	}
	return $self->redirect_to('/misc');
}

sub check_type_input {
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

# Update
sub update {
	my $self = shift;
	my $id   = $self->param('id');
	$self->stash( fbox_layout => 1 );

	my $err = &check_type_input($self);
	if ( defined($err) ) {
		$self->flash( message => $err );
		return $self->redirect_to( '/types/' . $id . '/edit' );
	}
	else {
		my $update = $self->db->resultset('Type')->find( { id => $self->param('id') } );
		$update->name( $self->param('type_data.name') );
		$update->description( $self->param('type_data.description') );
		$update->use_in_table( $self->param('type_data.use_in_table') );
		$update->update();

		# if the update has failed, we don't even get here, we go to the exception page.
		&log( $self, "Update type with name " . $self->param('type_data.name') . " and id " . $self->param('id'), "UICHANGE" );
	}

	$self->flash( message => "Successsfully updated " . $self->param('type_data.name') );
	return $self->redirect_to( '/types/' . $id . '/edit' );
}

# Create
sub create {
	my $self         = shift;
	my $name         = $self->param('type_data.name');
	my $description  = $self->param('type_data.description');
	my $use_in_table = $self->param('type_data.use_in_table');

	my $err = &check_type_input($self);
	my $p_name = $self->db->resultset('Type')->search( { name => $name } )->get_column('name')->single();

	if ( defined($err) ) {
		$self->flash( message => $err );
		return $self->redirect_to('/types/');
	}

	if ( defined $p_name ) {
		$self->field('type_data.name')->is_like( qr/^\/(?!$name\/)/i, "The name chosen exists." );
		$self->stash(
			fbox_layout => 1,
			type_data   => {
				name         => $name,
				description  => $description,
				use_in_table => $use_in_table
			}
		);
		return $self->render( template => 'types/add' );
	}
	my $new_id = -1;
	if ( defined($err) ) {
		$self->flash( alertmsg => $err );
		$self->redirect_to('/types');
	}
	else {
		my $insert = $self->db->resultset('Type')->create(
			{
				name         => $name,
				description  => $description,
				use_in_table => $use_in_table,
			}
		);
		$insert->insert();
		$new_id = $insert->id;
	}
	$self->flash( message => "Successsfully updated " . $name );

	return $self->redirect_to( '/types/' . $new_id . '/edit' );
}

1;
