package UI::Region;
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
	my $self      = shift;
	my %divisions = get_divisions($self);
	$self->stash( region => {}, divisions => \%divisions, selected_division => "default", div_id => "", fbox_layout => 1 );
}

sub get_divisions {
	my $self = shift;
	my %divisions;
	my $d_rs = $self->db->resultset("Division")->search( undef, { order_by => "name" } );
	while ( my $division = $d_rs->next ) {
		$divisions{ $division->name } = $division->id;
	}
	return %divisions;
}

# Create
sub create {
	my $self = shift;
	if ( &is_valid( $self, "add" ) ) {
		my $insert = $self->db->resultset('Region')->create(
			{
				name     => $self->param('region.name'),
				division => $self->param('region.division_id')
			}
		);
		$insert->insert();
		my $id = $insert->id;

		# if the insert has failed, we don't even get here, we go to the exception page.
		&log( $self, "Create Region with name:" . $self->param('region.name'), "UICHANGE" );
		$self->flash( message => "Successfully added Region!" );
		return $self->redirect_to("/region/$id/edit");
	}
	else {
		&stash_role($self);
		my %divisions = get_divisions($self);
		$self->stash( region => {}, divisions => \%divisions, selected_division => "default", fbox_layout => 1 );
		$self->render('region/add');
	}
}

sub is_valid {
	my $self = shift;
	my $mode = shift;
	my $name = $self->param('region.name');

	#Check required fields
	$self->field('region.name')->is_required;
	$self->field('region.division_id')->is_required;
	if ( $mode eq 'add' ) {
		my $existing_region = $self->db->resultset('Region')->search( { name => $name } )->get_column('name')->single();
		if ($existing_region) {
			$self->field('region.name')->is_equal( "", "A Region with name \"$name\" already exists." );
		}
	}
	if ( $mode eq 'edit' ) {
		my $id = $self->param('id');

		#get original name
		my $region_rs = $self->db->resultset('Region');
		my $orig_name = $region_rs->search( { id => $id } )->get_column('name')->single();
		if ( $name ne $orig_name ) {
			my $regions = $region_rs->search( { id => { '!=' => $id } } )->get_column('name');
			while ( my $db_name = $regions->next ) {
				if ( $db_name eq $name ) {
					$self->field('region.name')->is_equal( "", "Region with name \"$name\" already exists." );
				}
			}
		}
	}
	return $self->valid;
}

sub edit {
	my $self              = shift;
	my $id                = $self->param('id');
	my $cursor            = $self->db->resultset('Region')->search( { id => $id } );
	my $rname             = $cursor->get_column('name')->single();
	my $div_id            = $cursor->get_column('division')->single();
	my $selected_division = $self->db->resultset('Division')->search( { id => $div_id } )->get_column('name')->single();
	my %divisions         = get_divisions($self);
	&stash_role($self);
	$self->stash(
		region            => { name => $rname },
		selected_division => $selected_division,
		div_id            => $div_id,
		id                => $id,
		divisions         => \%divisions,
		fbox_layout       => 1
	);
	return $self->render('region/edit');
}

# Update
sub update {
	my $self        = shift;
	my $id          = $self->param('id');
	my $name        = $self->param('region.name');
	my $division_id = $self->param('region.division_id');

	if ( &is_valid( $self, "edit" ) ) {
		my $update = $self->db->resultset('Region')->find( { id => $self->param('id') } );
		$update->name($name);
		$update->division($division_id);
		$update->update();

		# if the update has failed, we don't even get here, we go to the exception page.
		&log( $self, "Update Region with name " . $name . " and id " . $self->param('id'), "UICHANGE" );

		$self->flash( message => "Successfully Updated Region!" );
		return $self->redirect_to("/region/$id/edit");
	}
	else {
		my $selected_division = $self->db->resultset('Division')->search( { id => $division_id } )->get_column('name')->single();
		my %divisions = get_divisions($self);
		&stash_role($self);
		$self->stash(
			region            => {},
			fbox_layout       => 1,
			divisions         => \%divisions,
			selected_division => $selected_division,
			div_id            => $division_id
		);
		$self->render('region/edit');
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
		my $name   = $self->db->resultset('Region')->search( { id => $id } )->get_column('name')->single();
		my $delete = $self->db->resultset('Region')->search( { id => $id } );
		$delete->delete();
		&log( $self, "Delete region " . $name, "UICHANGE" );
	}
	return $self->redirect_to('/close_fancybox.html');
}

1;
