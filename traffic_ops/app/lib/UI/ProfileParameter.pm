package UI::ProfileParameter;
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

# Read
sub read {
	my $self = shift;
	my @data;
	my $orderby = $self->param('orderby') || "profile";
	my $rs_data = $self->db->resultset("ProfileParameter")->search( undef, { prefetch => [ 'profile', 'parameter' ], order_by => $orderby } );
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"profile"      => $row->profile->name,
				"parameter"    => $row->parameter->id,
				"last_updated" => $row->last_updated,
			}
		);
	}
	$self->render( json => \@data );
}

# Delete
sub delete {
	my $self      = shift;
	my $profile   = $self->param('profile');
	my $parameter = $self->param('parameter');

	if ( !&is_oper($self) ) {
		$self->flash( alertmsg => "No can do. Get more privs." );
	}
	else {
		my $prof_name = $self->db->resultset('Profile')->search( { id => $profile } )->get_column('name')->single();
		my $parameter_name = $self->db->resultset('Parameter')->search( { id => $parameter } )->get_column('name')->single();
		my $delete = $self->db->resultset('ProfileParameter')->search( { parameter => $parameter, profile => $profile } );
		$delete->delete();
		$self->flash( alertmsg => 'Success!' );
		&log( $self, "Delete profile parameter " . $prof_name . " <-> " . $parameter_name, "UICHANGE" );
	}

	my $referer = $self->req->headers->header('referer');
	$referer = "/" if ( !defined $referer );
	return if !defined($profile);
	return if !defined($parameter);
	return $self->redirect_to($referer);
}

# Update
# Update not needed for linking table.

# Create
sub create {
	my $self      = shift;
	my $parameter = $self->param('parameter');
	my $profile   = $self->param('profile');
	if ( !&is_oper($self) ) {
		$self->flash( alertmsg => "No can do. Get more privs." );
	}
	else {
		my $prof_name = $self->db->resultset('Profile')->search( { id => $profile } )->get_column('name')->single();
		my $parameter_name = $self->db->resultset('Parameter')->search( { id => $parameter } )->get_column('name')->single();
		my $insert = $self->db->resultset('ProfileParameter')->create( { parameter => $parameter, profile => $profile } )->insert();
		$self->flash( message => "Success!" );
		$self->flash( message => 'Success!' );
		&log( $self, "Create profile parameter " . $prof_name . " <-> " . $parameter_name, "UICHANGE" );
	}
	my $referer = $self->req->headers->header('referer');
	return $self->redirect_to($referer);
}

1;
