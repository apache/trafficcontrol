package UI::Update;
#
# Copyright 2015 Comcast Cable Communications Management, LLC
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
use Data::Dumper;
use Mojo::Base 'Mojolicious::Controller';

sub readupdate {
	my $self = shift;
	my @data;
	my $host_name = $self->param("host_name");

	my $rs_servers;
	if ( $host_name =~ m/^all$/ ) {
		$rs_servers = $self->db->resultset("Server")->search(undef);
	}
	else {
		$rs_servers = $self->db->resultset("Server")->search( { host_name => $host_name } );
	}

	while ( my $row = $rs_servers->next ) {
		push( @data, { host_name => $row->host_name, upd_pending => $row->upd_pending, host_id => $row->id, status => $row->status->name } );
	}

	$self->render( json => \@data );
}

sub postupdate {

	my $self      = shift;
	my $updated   = $self->param("updated");
	my $host_name = $self->param("host_name");
	if ( !&is_admin($self) ) {
		$self->render_text( "Unauthorized.", status => 401, layout => undef );
		return;
	}

	if ( !defined($updated) ) {
		$self->render_text( "Failed request.  Must provide updated status", status => 500, layout => undef );
		return;
	}

	# resolve server id
	my $serverid = $self->db->resultset("Server")->search( { host_name => $host_name } )->get_column('id')->single;
	if ( !defined $serverid ) {
		$self->render_text( "Failed request.  Unknown server", status => 500, layout => undef );
		return;
	}

	my $update_server = $self->db->resultset('Server')->search( { id => $serverid } );
	if ( defined $updated ) {
		$update_server->update( { upd_pending => $updated } );
	}

print "YAY\n\n";
	# $self->render_text("Success", layout=>undef);

}

sub postupdatequeue {
	my $self     = shift;
	my $setqueue = $self->param("setqueue");
	my $host     = $self->param("id");

	if ( !&is_admin($self) && !&is_oper($self) ) {
		$self->flash( alertmsg => "No can do. Get more privs." );
		return;
	}

	my $edge_type = &type_id( $self, 'EDGE' );
	if ( $host eq "all" ) {

		# default is only edges
		my $update = $self->db->resultset('Server')->search( { type => $edge_type } );
		$update->update( { upd_pending => $setqueue } );
	}
	elsif ( $host eq "edges" ) {
		my $update = $self->db->resultset('Server')->search( { type => $edge_type  } );
		$update->update( { upd_pending => $setqueue } );
	}
	else {
		my $update = $self->db->resultset('Server')->search( { id => $host, } );
		$update->update( { upd_pending => $setqueue } );
	}

	&log( $self, "Flip Update bit (Queue Updates) for server(s):" . $host, "OPER" );

	# $self->render_text("Success", layout=>undef);
	#return $self->redirect_to('/#tabs=1');
}

1;
