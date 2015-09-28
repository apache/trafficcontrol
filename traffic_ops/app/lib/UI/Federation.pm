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

	&stash_role($self);
	$self->stash( federation => {}, fbox_layout => 1, mode => 'add' );
}

# Read
sub read {
	my $self = shift;

	my @data;
	my $orderby = "name";
	$orderby = $self->param('orderby') if ( defined $self->param('orderby') );
	my $dbh = $self->db->resultset("Federation")->search( undef, { prefetch => [ { 'role' => undef } ], order_by => 'me.' . $orderby } );
	while ( my $row = $dbh->next ) {
		push(
			@data, {
				"id"          => $row->id,
				"name"        => $row->name,
				"description" => $row->description,
				"cname"       => $row->cname,
				"ttl"         => $row->ttl,
				"role"        => $row->role->id,
			}
		);
	}
	$self->render( json => \@data );
}

sub edit {
	my $self = shift;
	my $id   = $self->param('id');
	my $dbh  = $self->db->resultset('Federation')->search( { id => $id } );
	my $data = $dbh->single;
	&stash_role($self);

	my %delivery_services = get_delivery_services( $self, $id );
	$self->stash(
		federation        => $data,
		mode              => 'edit',
		fbox_layout       => 1,
		delivery_services => \%delivery_services
	);
	return $self->render('federation/edit');
}

sub get_delivery_services {
	my $self   = shift;
	my $id     = shift;
	my @ds_ids = $self->db->resultset('DeliveryserviceTmuser')->search( { tm_user_id => $id } )->get_column('deliveryservice')->all;

	my %delivery_services;
	for my $ds_id (@ds_ids) {
		my $desc = $self->db->resultset('Deliveryservice')->search( { id => $ds_id } )->get_column('xml_id')->single;
		$delivery_services{$ds_id} = $desc;
	}
	return %delivery_services;
}

# Update
sub update {
	my $self       = shift;
	my $tm_user_id = $self->param('id');
	my @ds_ids     = $self->param('deliveryservices');

	$self->associated_delivery_services( $tm_user_id, \@ds_ids );

	# Prevent these from getting updated
	# Do not modify the local_passwd if it comes across as blank.
	my $local_passwd         = $self->param("tm_user.local_passwd");
	my $confirm_local_passwd = $self->param("tm_user.confirm_local_passwd");

	if ( $self->is_valid("edit") ) {
		my $dbh = $self->db->resultset('TmUser')->find( { id => $tm_user_id } );
		$dbh->username( $self->param('tm_user.username') );
		$dbh->full_name( $self->param('tm_user.full_name') );
		$dbh->role( $self->param('tm_user.role') );
		$dbh->uid(0);
		$dbh->gid(0);

		# ignore the local_passwd and confirm_local_passwd if it comes across as blank (or it didn't change)
		if ( defined($local_passwd) && $local_passwd ne '' ) {
			$dbh->local_passwd( sha1_hex( $self->param('tm_user.local_passwd') ) );
		}
		if ( defined($confirm_local_passwd) && $confirm_local_passwd ne '' ) {
			$dbh->confirm_local_passwd( sha1_hex( $self->param('tm_user.confirm_local_passwd') ) );
		}

		$dbh->company( $self->param('tm_user.company') );
		$dbh->email( $self->param('tm_user.email') );
		$dbh->full_name( $self->param('tm_user.full_name') );
		$dbh->address_line1( $self->param('tm_user.address_line1') );
		$dbh->address_line2( $self->param('tm_user.address_line2') );
		$dbh->city( $self->param('tm_user.city') );
		$dbh->state_or_province( $self->param('tm_user.state_or_province') );
		$dbh->phone_number( $self->param('tm_user.phone_number') );
		$dbh->postal_code( $self->param('tm_user.postal_code') );
		$dbh->country( $self->param('tm_user.country') );
		$dbh->update();
		$self->flash( message => "User was updated successfully." );
		$self->stash( mode => 'edit' );
		return $self->redirect_to( '/federation/' . $tm_user_id . '/edit' );
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
			return $self->redirect_to('/close_fancybox.html');
		}
	}
	else {
		return $self->render('federation/add');
	}
}

sub is_valid {
	my $self = shift;
	my $mode = shift;

	$self->field('federation.name')->is_required;
	$self->field('federation.cname')->is_required;
	$self->field('federation.ttl')->is_required;

	return $self->valid;
}

sub create_federation_mapping {
	my $self   = shift;
	my $new_id = -1;
	my $dbh    = $self->db->resultset('Federation')->create(
		{
			name        => $self->param('federation.name'),
			description => $self->param('federation.description'),
			cname       => $self->param('federation.cname'),
			ttl         => $self->param('federation.ttl'),
			type        => $self->param('federation.type'),
		}
	);
	$new_id = $dbh->insert();

	# if the insert has failed, we don't even get here, we go to the exception page.
	&log( $self, "Create federation with name: " . $self->param('federation.name') . " and cname: " . $self->param('federation.name'), "UICHANGE" );
	return $new_id;

}

1;
