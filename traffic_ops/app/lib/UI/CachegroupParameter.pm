package UI::CachegroupParameter;
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

sub index {
	my $self = shift;
	my @data;
	my $orderby = $self->param('orderby') || "cachegroup";
	my $rs_data = $self->db->resultset("CachegroupParameter")->search( undef, { prefetch => [ 'cachegroup', 'parameter' ], order_by => $orderby } );
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"cachegroup"   => $row->cachegroup->name,
				"parameter"    => $row->parameter->id,
				"last_updated" => $row->last_updated,
			}
		);
	}
	$self->render( json => \@data );
}

# Delete
sub delete {
	my $self       = shift;
	my $cachegroup = $self->param('cachegroup');
	my $parameter  = $self->param('parameter');

	if ( !&is_oper($self) ) {
		$self->flash( alertmsg => "No can do. Get more privs." );
	}
	else {
		my $cg_name = $self->db->resultset('Cachegroup')->search( { id => $cachegroup } )->get_column('name')->single();
		my $parameter_name = $self->db->resultset('Parameter')->search( { id => $parameter } )->get_column('name')->single();
		my $delete = $self->db->resultset('CachegroupParameter')->search( { parameter => $parameter, cachegroup => $cachegroup } );
		$delete->delete();
		$self->flash( alertmsg => 'Success!' );
		&log( $self, "Delete cachegroup parameter link " . $cg_name . " <-> " . $parameter_name, "UICHANGE" );
	}

	my $referer = $self->req->headers->header('referer');
	$referer = "/" if ( !defined $referer );
	return if !defined($cachegroup);
	return if !defined($parameter);
	return $self->redirect_to($referer);
}

# Update
# Update not needed for linking table.

# Create
sub create {
	my $self       = shift;
	my $parameter  = $self->param('parameter');
	my $cachegroup = $self->param('cachegroup');

	if ( !&is_oper($self) ) {
		$self->flash( message => "No can do. Get more privs." );
	}
	else {
		my $prof_name = $self->db->resultset('Cachegroup')->search( { id => $cachegroup } )->get_column('name')->single();
		my $parameter_name = $self->db->resultset('Parameter')->search( { id => $parameter } )->get_column('name')->single();
		my $insert = $self->db->resultset('CachegroupParameter')->create( { parameter => $parameter, cachegroup => $cachegroup } )->insert();
		$self->flash( message => "Success!" );
		&log( $self, "Create cachegroup parameter link " . $prof_name . " <-> " . $parameter_name, "UICHANGE" );
	}

	# my $referer = $self->req->headers->header('referer');
	return $self->redirect_to("/parameter/$parameter");
}
1;
