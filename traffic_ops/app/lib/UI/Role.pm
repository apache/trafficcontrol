package UI::Role;
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
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;

# Read
sub read {
	my $self = shift;
	my @data;
	my $orderby = $self->param('orderby') || "name";

	my $rs_data = $self->db->resultset("Role")->search( undef, { order_by => $orderby } );
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"id"          => $row->id,
				"name"        => $row->name,
				"description" => $row->description,
				"priv_level"  => $row->priv_level,
			}
		);
	}
	$self->render( json => \@data );
}

# Delete
sub delete {
	my $self = shift;
	my $id   = $self->param('id');
	return 1 if !defined($id);
	my $delete = $self->db->resultset('Role')->search( { id => $id } );
	$delete->delete();
	return $self->redirect_to('/misc');
}

# Update
sub update {
	my $self = shift;
	my $update = $self->db->resultset('Role')->find( { id => $self->param('id') } );
	$update->name( $self->param('name') );
	$update->description( $self->param('description') );
	$update->update();
	return $self->redirect_to('/misc');
}

# Create
sub create {
	my $self   = shift;
	my $insert = $self->db->resultset('Role')->create(
		{
			name        => $self->param('name'),
			description => $self->param('description'),
			priv_level  => $self->param('priv_level'),
		}
	);
	$insert->insert();
	return $self->redirect_to('/misc');
}

1;
