package UI::Federation;
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

# JvD Note: you always want to put Utils as the first use. Sh*t don't work if it's after the Mojo lines.
use UI::Utils;
use List::MoreUtils qw(uniq);

use Mojo::Base 'Mojolicious::Controller';
use Digest::SHA1 qw(sha1_hex);
use Mojolicious::Validator;
use Mojolicious::Validator::Validation;
use Email::Valid;
use Data::GUID;
use Data::Dumper;

# List of Federation Mappings
sub index {
	my $self = shift;
	&navbarpage($self);
}

# NOTE: Do NOT attempt to call this method 'new' or 'init'
#      because Mojo will death spiral.
# Setup a New user for "Add User".
sub add {
	my $self = shift;

	my $current_username = $self->current_user()->{username};
	my $dbh              = $self->db->resultset('TmUser')->search( { username => $current_username } );
	my $tm_user          = $dbh->single;
	&stash_role($self);

	#TODO: drichardson - remove hard coded DS
	my $delivery_services = get_delivery_services( $self, 1 );

	$self->stash(
		tm_user           => $tm_user,
		federation        => {},
		delivery_services => $delivery_services,
		fbox_layout       => 1,
		selected_role_id  => 7,
		mode              => 'add'
	);
}

sub edit {
	my $self   = shift;
	my $fed_id = $self->param('federation_id');

	my $federation;
	my $selected_ds_id;
	my $feds = $self->db->resultset('Federation')->search( { 'id' => $fed_id } );
	while ( my $f = $feds->next ) {
		$federation = $f;
		my $fed_id = $f->id;
		my $federation_deliveryservices =
			$self->db->resultset('FederationDeliveryservice')->search( { federation => $fed_id }, { prefetch => [ 'federation', 'deliveryservice' ] } );
		while ( my $fd = $federation_deliveryservices->next ) {
			$selected_ds_id = $fd->deliveryservice->id;
		}
	}

	my $selected_role_id;
	my $ftusers =
		$self->db->resultset('FederationTmuser')->search( { federation => $fed_id }, { prefetch => [ 'federation', 'tm_user' ] } );
	while ( my $ft = $ftusers->next ) {
		$selected_role_id = $ft->role->id;
	}
	$self->app->log->debug( "selected_role_id #-> " . $selected_role_id );

	my $resolvers = $self->db->resultset('FederationResolver')
		->search( { 'federation_federation_resolvers.federation_resolver' => $fed_id }, { prefetch => 'federation_federation_resolvers' } );
	while ( my $row = $resolvers->next ) {
		my $line = [ $row->id ];
	}

	my $current_username = $self->current_user()->{username};
	my $dbh              = $self->db->resultset('TmUser')->search( { username => $current_username } );
	my $tm_user          = $dbh->single;
	&stash_role($self);

	#TODO: drichardson - remove harded DS
	my $delivery_services = get_delivery_services( $self, 1 );
	$self->app->log->debug( "delivery_services #-> " . Dumper($delivery_services) );
	$self->stash(
		tm_user           => $tm_user,
		selected_ds_id    => $selected_ds_id,
		selected_role_id  => $selected_role_id,
		federation        => $federation,
		mode              => 'edit',
		fbox_layout       => 1,
		delivery_services => $delivery_services
	);
	return $self->render('federation/edit');
}

sub get_delivery_services {
	my $self   = shift;
	my $id     = shift;
	my @ds_ids = $self->db->resultset('Deliveryservice')->search( undef, { orderby => "xml_id" } )->get_column('id')->all;

	my $delivery_services;
	for my $ds_id ( uniq(@ds_ids) ) {
		$self->app->log->debug( "looking for ds_id #-> " . Dumper($ds_id) );
		my $desc = $self->db->resultset('Deliveryservice')->search( { id => $ds_id } )->get_column('xml_id')->single;
		$delivery_services->{$ds_id} = $desc;
	}
	$self->app->log->debug( "delivery_services #-> " . Dumper($delivery_services) );
	return $delivery_services;
}

# Update
sub update {
	my $self        = shift;
	my $fed_id      = $self->param('federation_id');
	my $cname       = $self->param('federation.cname');
	my $role_id     = $self->param('tm_user.role');
	my $description = $self->param('federation.description');
	my $ttl         = $self->param('federation.ttl');

	my $is_valid = $self->is_valid("edit");
	if ( $self->is_valid("edit") ) {
		my $dbh = $self->db->resultset('Federation')->find( { id => $fed_id } );
		$dbh->cname($cname);
		$dbh->description($description);
		$dbh->ttl($ttl);
		$dbh->update();

		my $ftusers =
			$self->db->resultset('FederationTmuser')->search( { federation => $fed_id }, { prefetch => [ 'federation', 'tm_user' ] } );
		while ( my $ft = $ftusers->next ) {
			my $fid    = $ft->federation->id;
			my $fcname = $ft->federation->cname;
			$self->app->log->debug( "fid #-> " . $fid );
			$ft->role($role_id);
			$ft->update();
		}

		$self->flash( message => "Federation was updated successfully." );
		$self->stash( mode => 'edit' );
		return $self->redirect_to( '/federation/' . $fed_id . '/edit' );
	}
	else {
		$self->edit();
	}
}

sub associated_delivery_services {
	my $self       = shift;
	my $tm_user_id = shift;
	my $ds_ids     = shift;

	my $new_id = -1;

	# Sweep the existing DeliveryserviceTmUser relationships
	my $delete = $self->db->resultset('DeliveryserviceTmuser')->search( { tm_user_id => $tm_user_id } );
	$delete->delete();

	# Attached the saved delivery services
	foreach my $ds_id ( @{$ds_ids} ) {
		my $ds_name = $self->db->resultset('Deliveryservice')->search( { id => $ds_id } )->get_column('xml_id')->single();
		my $insert = $self->db->resultset('DeliveryserviceTmuser')->create( { deliveryservice => $ds_id, tm_user_id => $tm_user_id } );

		$new_id = $insert->tm_user_id;
		$insert->insert();
		&log( $self, "Associated Delivery service " . $ds_name . " <-> with tm_user_id: " . $tm_user_id, "UICHANGE" );
	}

}

# Create
sub create {
	my $self = shift;
	&stash_role($self);
	$self->stash( fbox_layout => 1, mode => 'add', federation => {} );
	if ( $self->is_valid("add") ) {
		my $new_id = $self->create_federation_mapping();
		if ( $new_id != -1 ) {
			$self->flash( message => 'Federation created successfully.' );
			return $self->redirect_to('/federation/add');
		}
	}
	else {
		return $self->render('federation/add');
	}
}

sub is_valid {
	my $self = shift;
	my $mode = shift;

	$self->field('federation.cname')->is_required;
	$self->field('federation.ttl')->is_required;

	return $self->valid;
}

sub create_federation_mapping {
	my $self   = shift;
	my $new_id = -1;
	my $dbh    = $self->db->resultset('Federation')->create(
		{
			cname       => $self->param('federation.cname'),
			description => $self->param('federation.description'),
			ttl         => $self->param('federation.ttl'),
		}
	);
	$new_id = $dbh->insert();

	# if the insert has failed, we don't even get here, we go to the exception page.
	&log( $self, "Create federation with name: " . $self->param('federation.name') . " and cname: " . $self->param('federation.name'), "UICHANGE" );
	return $new_id;

}

# Delete
sub delete {
	my $self   = shift;
	my $fed_id = $self->param('federation_id');
	my $cname  = $self->param('federation.cname');

	if ( !&is_oper($self) ) {
		$self->flash( alertmsg => "No can do. Get more privs." );
	}
	else {
		my $delete = $self->db->resultset('Federation')->search( { id => $fed_id } );
		my $resolvers =
			$self->db->resultset('FederationFederationResolver')
			->search( { federation => $fed_id }, { prefetch => [ 'federation', 'federation_resolver' ] } );
		my $ip_address;
		my $cname;
		while ( my $row = $resolvers->next ) {
			my $id = $row->id;
			$self->app->log->debug( "id #-> " . $id );

			#$ip_address = $row->ip_address;
		}
		$delete->delete();
		&log( $self, "Deleted federation: " . $fed_id . " cname: " . $cname . " ip_address: " . $ip_address, "UICHANGE" );
	}
	return $self->redirect_to('/close_fancybox.html');
}

1;
