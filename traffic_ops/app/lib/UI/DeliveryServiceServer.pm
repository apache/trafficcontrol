package UI::DeliveryServiceServer;
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

use UI::DeliveryService;

sub cpdss_iframe {
	my $self    = shift;
	my $mode    = $self->param('mode');
	my $srvr_id = $self->param('id');

	if ( $mode eq "view" ) {
		my $server = $self->db->resultset('Server')->search( { 'me.id' => $srvr_id }, { prefetch => [ 'cachegroup' ] } )->single();

		my @etypeids = &type_ids( $self, 'EDGE%', 'server' );
		my $rs = $self->db->resultset('Server')->search( { 'me.type' => { -in => \@etypeids }, cdn_id => $server->cdn_id }, 
			{ prefetch => 'profile', order_by => 'host_name' } );
		my @from_server_list;
		while ( my $row = $rs->next ) {
			if ( $row->id == $srvr_id ) {
				next;
			}
			# servers in same cachegroup go at the top
			if ( $row->cachegroup->id == $server->cachegroup->id ) {
				unshift( @from_server_list, $row );
			}
			else {
				push( @from_server_list, $row );
			}
		}

		&stash_role($self);
		$self->stash( fbox_layout      => 1 );
		$self->stash( server           => $server );
		$self->stash( from_server_list => \@from_server_list );
	}
}

sub edit {
	my $self = shift;
	my $mode = $self->param('mode');
	my $id   = $self->param('id');

	my $dss_data;
	my $totals;
	&stash_role($self);

	# Get list of server ids associated with ds
	my $assigned_servers;
	my $rsas = $self->db->resultset('DeliveryserviceServer')->search( { deliveryservice => $id }, { prefetch => [ 'server' ] } );
	while ( my $row = $rsas->next ) {
		$assigned_servers->{ $row->server->id } = 1;
	}

	my $ds = $self->db->resultset('Deliveryservice')->search( { 'me.id' => $id } )->single();
	my $valid_profiles;
	my $psas = $self->db->resultset('Server')->search(
		{ cdn_id => $ds->cdn_id },
		{
			select   => 'profile',
			distinct => 1
		}
	)->get_column('profile');
	while ( my $row = $psas->next ) {
		$valid_profiles->{$row} = 1;
	}

	$ds = $self->db->resultset('Deliveryservice')->search( { 'me.id' => $id }, { prefetch => [ 'cdn' ] } )->single();

	my @types;
	push(@types, &type_ids( $self, 'EDGE%', 'server' ) );
	push(@types, &type_id( $self, 'ORG' ) );
	my $rs      = $self->db->resultset('Server')
		->search( { "me.type" => { -in => \@types } }, { prefetch => [ 'cachegroup', 'type', 'profile', 'status' ], } );
	while ( my $row = $rs->next ) {

		# skip profiles that are not associated with the cdn this ds is in
		if ( !defined( $valid_profiles->{ $row->profile->id } ) ) {
			next;
		}
		if ( !defined( $totals->{ $row->profile->name }->{ $row->cachegroup->name }->{assigned} ) ) {
			$totals->{ $row->profile->name }->{ $row->cachegroup->name }->{assigned}     = 0;
			$totals->{ $row->profile->name }->{ $row->cachegroup->name }->{not_assigned} = 0;
		}
		if ( !defined( $totals->{ $row->profile->name }->{assigned} ) ) {
			$totals->{ $row->profile->name }->{assigned}     = 0;
			$totals->{ $row->profile->name }->{not_assigned} = 0;
		}
		$dss_data->{ $row->profile->name }->{ $row->cachegroup->name }->{ $row->host_name }->{id} = $row->id;
		if ( defined( $assigned_servers->{ $row->id } ) ) {
			$dss_data->{ $row->profile->name }->{ $row->cachegroup->name }->{ $row->host_name }->{assigned} = 1;
			$totals->{ $row->profile->name }->{ $row->cachegroup->name }->{assigned}++;
			$totals->{ $row->profile->name }->{assigned}++;
		}
		else {
			$dss_data->{ $row->profile->name }->{ $row->cachegroup->name }->{ $row->host_name }->{assigned} = 0;
			$totals->{ $row->profile->name }->{ $row->cachegroup->name }->{not_assigned}++;
			$totals->{ $row->profile->name }->{not_assigned}++;
		}
	}

	$self->stash( ds_id            => $id );
	$self->stash( assigned_servers => $dss_data );
	$self->stash( ds_name          => $ds->xml_id . ' (' . UI::DeliveryService::compute_org_server_fqdn($self, $ds->id) . ')' );
	$self->stash( fbox_layout      => 1 );
	$self->stash( dss_data         => $dss_data );
	$self->stash( totals           => $totals );
	$self->stash( cdn_name         => $ds->cdn->name );
}

# Read
sub read {
	my $self = shift;
	my @data;
	my $orderby = "deliveryservice";
	my $limit   = 10;
	$orderby = $self->param('orderby') if ( defined $self->param('orderby') );
	$limit   = $self->param('limit')   if ( defined $self->param('limit') );
	my $rs_data = $self->db->resultset("DeliveryserviceServer")->search( undef, { prefetch => [ 'deliveryservice', 'server' ], order_by => $orderby, rows => $limit } );
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"deliveryservice" => $row->deliveryservice->xml_id,
				"server"          => $row->server->id,
				"last_updated"    => $row->last_updated,
			}
		);
	}

	# print "data = " . Dumper(@data);
	$self->render( json => \@data );
}

sub clone_server {
	my $self = shift;

	my $from_server = $self->param('from_server');
	my $to_server   = $self->param('to_server');

	# foreach my $param ( $self->param ) {
	# 	print $param . " -> " . $self->param($param) . "\n";
	# }

	my @dslist = $self->db->resultset('DeliveryserviceServer')->search( { server => $from_server } )->get_column('deliveryservice')->all();

	# clean up
	my $delete = $self->db->resultset('DeliveryserviceServer')->search( { server => $to_server } );
	$delete->delete();

	my $numlinks = 0;
	foreach my $ds (@dslist) {

		# print ">>> " . $ds . "\n";
		my $insert = $self->db->resultset('DeliveryserviceServer')->create(
			{
				deliveryservice => $ds,
				server          => $to_server,
			}
		);
		$insert->insert();
		$numlinks++;

		my $ds = $self->db->resultset('Deliveryservice')->search( { id => $ds } )->single();
		&UI::DeliveryService::header_rewrite( $self, $ds->id, $ds->profile, $ds->xml_id, $ds->edge_header_rewrite, "edge" );
		&UI::DeliveryService::regex_remap( $self, $ds->id, $ds->profile, $ds->xml_id, $ds->regex_remap );
	        &UI::DeliveryService::cacheurl( $self, $ds->id, $ds->profile, $ds->xml_id, $ds->cacheurl );

	}

	my $from_name = $self->db->resultset('Server')->search( { id => $from_server } )->get_column('host_name')->single();
	my $to_name   = $self->db->resultset('Server')->search( { id => $to_server } )->get_column('host_name')->single();
	&log( $self, "Clone deliveryservice links from " . $from_name . " to " . $to_name . " (" . $numlinks . " links cloned)", "UICHANGE" );

	$self->flash( alertmsg => "Success!" );
	return $self->redirect_to('/close_fancybox.html');
}

sub assign_servers {
	my $self = shift;
	my $dsid     = $self->param('id');

	my @server_ids;
	foreach my $param ( $self->param ) {
		next
			if ( $param eq 'id' || $self->param($param) ne 'on' );    # we only get the 'on', but still.
		my ( $rubbish, $srvr_id ) = split( /_/, $param );
		push( @server_ids, $srvr_id );
	}

	# clean up
	my $delete = $self->db->resultset('DeliveryserviceServer')->search( { deliveryservice => $dsid } );
	$delete->delete();

	# and associate what was checked
	my $numlinks = 0;
	foreach my $s_id (@server_ids) {
		my $insert = $self->db->resultset('DeliveryserviceServer')->create(
			{
				deliveryservice => $dsid,
				server          => $s_id,
			}
		);
		$insert->insert();
		$numlinks++;
	}

	my $ds = $self->db->resultset('Deliveryservice')->search( { id => $dsid } )->single();
	&UI::DeliveryService::header_rewrite( $self, $ds->id, $ds->profile, $ds->xml_id, $ds->edge_header_rewrite, "edge" );
        &UI::DeliveryService::regex_remap( $self, $ds->id, $ds->profile, $ds->xml_id, $ds->regex_remap );
        &UI::DeliveryService::cacheurl( $self, $ds->id, $ds->profile, $ds->xml_id, $ds->cacheurl );

	&log( $self, "Link deliveryservice " . $ds->xml_id . " to " . $numlinks . " servers", "UICHANGE" );

	$self->flash( alertmsg => "Success!" );
	my $referer = $self->req->headers->header('referer');
	return $self->redirect_to($referer);
}

# Create
sub create {
	my $self            = shift;
	my $new_id          = -1;
	my $server          = $self->param('server');
	my $deliveryservice = $self->param('deliveryservice');
	if ( !&is_oper($self) ) {
		$self->flash( alertmsg => "No can do. Get more privs." );
	}
	else {
		my $server_name = $self->db->resultset('Server')->search( { id => $server } )->get_column('host_name')->single();
		my $insert = $self->db->resultset('DeliveryserviceServer')->create( { server => $server, deliveryservice => $deliveryservice } )->insert();
		$new_id = $insert->id;

		my $ds = $self->db->resultset('Deliveryservice')->search( { id => $deliveryservice } )->single();
		&UI::DeliveryService::header_rewrite( $self, $ds->id, $ds->profile, $ds->xml_id, $ds->edge_header_rewrite, "edge" );

		$self->flash( alertmsg => 'Success!' );
		&log( $self, "Create deliveryservice server link " . $ds->xml_id . " <-> " . $server_name, "UICHANGE" );
	}
	if ( $new_id == -1 ) {
		my $referer = $self->req->headers->header('referer');
		if ( defined($referer) ) {
			return $self->redirect_to($referer);
		}
		else {
			return $self->render( text => "ERR = ", layout => undef );    # for testing - $referer is not defined there.
		}
	}
	return $self->redirect_to('/deliveryservices');
}

1;
