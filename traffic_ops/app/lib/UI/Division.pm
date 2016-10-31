package UI::Division;
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

# for the fancybox view
sub add {
	my $self = shift;
	$self->stash( division => {}, fbox_layout => 1 );
}

# Create
sub create {
	my $self = shift;
	if ( $self->is_valid() ) {
		my $insert = $self->db->resultset('Division')->create( { name => $self->param('division.name') } );
		$insert->insert();
		my $id = $insert->id;

		# if the insert has failed, we don't even get here, we go to the exception page.
		&log( $self, "Create division with name:" . $self->param('division.name'), "UICHANGE" );
		$self->flash( message => "Successfully added Division!" );
		return $self->redirect_to("/division/$id/edit");
	}
	else {
		&stash_role($self);
		$self->stash( division => {}, fbox_layout => 1 );
		$self->render('division/add');
	}
}

sub is_valid {
	my $self = shift;
	my $name = $self->param('division.name');

	#Check required fields
	$self->field('division.name')->is_required;

	#Check for duplicate division name
	my $existing_division = $self->db->resultset('Division')->search( { name => $name } )->get_column('name')->single();
	if ($existing_division) {
		$self->field('division.name')->is_equal( "", "A division with name \"$name\" already exists." );
	}
	return $self->valid;
}

sub edit {
	my $self   = shift;
	my $id     = $self->param('id');
	my $cursor = $self->db->resultset('Division')->search( { id => $id } );
	my $data   = $cursor->single;
	&stash_role($self);
	$self->stash( division => $data, id => $data->id, fbox_layout => 1 );
	return $self->render('division/edit');
}

# Update
sub update {
	my $self = shift;
	my $id   = $self->param('id');
	my $name = $self->param('division.name');

	if ( $self->is_valid() ) {
		my $update = $self->db->resultset('Division')->find( { id => $self->param('id') } );
		$update->name($name);
		$update->update();

		# if the update has failed, we don't even get here, we go to the exception page.
		&log( $self, "Update division with name " . $name . " and id " . $self->param('id'), "UICHANGE" );

		$self->flash( message => "Successfully Updated Division!" );
		return $self->redirect_to("/division/$id/edit");
	}
	else {
		&stash_role($self);
		$self->stash( division => {}, fbox_layout => 1 );
		$self->render('division/edit');
	}
}

# Delete
sub delete {
	my $self = shift;
	my $id   = $self->param('id');

	if ( !&is_admin($self) ) {
		$self->flash( message => "You must be an ADMIN to perform this operation!" );
	}
	else {
		my $name   = $self->db->resultset('Division')->search( { id => $id } )->get_column('name')->single();
		my $delete = $self->db->resultset('Division')->search( { id => $id } );
		$delete->delete();
		&log( $self, "Delete division " . $name, "UICHANGE" );
	}
	return $self->redirect_to('/close_fancybox.html');
}

1;
