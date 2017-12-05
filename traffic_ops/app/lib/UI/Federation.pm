package UI::Federation;
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
use List::MoreUtils qw(uniq);

use Mojo::Base 'Mojolicious::Controller';
use Mojolicious::Validator;
use Mojolicious::Validator::Validation;
use Email::Valid;
use Data::GUID;
use Data::Dumper;
use constant FEDERATION_ROLE => 'federation';

# List of Federation Mappings
sub index {
	my $self = shift;
	&navbarpage($self);
}

sub add {
	my $self = shift;

	my $current_username = $self->current_user()->{username};
	my $dbh              = $self->db->resultset('TmUser')->search( { username => $current_username } );
	my $tm_user          = $dbh->single;
	&stash_role($self);

	# default the ds_id to 0 because _form.html.ep expects it to be there
	$self->stash(
		tm_user     => $tm_user,
		ds_id       => 0,
		user_id     => 0,
		federation  => {},
		fbox_layout => 1,
		role_name   => FEDERATION_ROLE,    # the federation role
		mode        => 'add'
	);
}

sub edit {
	my $self   = shift;
	my $fed_id = $self->param('federation_id');

	my $federation;
	my $ds_id;
	my $feds =
		$self->db->resultset('Federation')->search( { 'id' => $fed_id } );
	my $fed_count = $feds->count();
	if ( $fed_count > 0 ) {
		while ( my $f = $feds->next ) {
			$federation = $f;
			my $fed_id = $f->id;
			my $federation_deliveryservices =
				$self->db->resultset('FederationDeliveryservice')->search( { federation => $fed_id }, { prefetch => [ 'federation', 'deliveryservice' ] } );
			while ( my $fd = $federation_deliveryservices->next ) {
				$ds_id = $fd->deliveryservice->id;
			}
		}

		my $role_name;
		my $user_id;
		my $ftusers = $self->db->resultset('FederationTmuser')->search( { federation => $fed_id }, { prefetch => [ 'federation', 'tm_user', 'role' ] } );
		while ( my $ft = $ftusers->next ) {
			$user_id   = $ft->tm_user->id;
			$role_name = $ft->role->name;
		}

		my $current_username = $self->current_user()->{username};
		my $dbh              = $self->db->resultset('TmUser')->search( { username => $current_username } );
		my $tm_user          = $dbh->single;
		&stash_role($self);

		my $delivery_services = get_delivery_services( $self, $ds_id );
		$self->stash(
			tm_user           => $tm_user,
			ds_id             => $ds_id,
			user_id           => $user_id,
			role_name         => $role_name,
			federation        => $federation,
			mode              => 'edit',
			fbox_layout       => 1,
			delivery_services => $delivery_services
		);
		return $self->render('federation/edit');
	}
	else {
		return $self->not_found();
	}
}

# .json format for the jqTree widge
sub users {
	my $self = shift;
	my $data;
	my $federation_role_id = $self->db->resultset('Role')->search( { name => FEDERATION_ROLE }, undef )->get_column('id')->single();

	my $fed_users = $self->db->resultset('TmUser')->search( { role => $federation_role_id }, { order_by => 'full_name' } );
	while ( my $row = $fed_users->next ) {
		push(
			@$data, {
				id       => $row->id,
				username => $row->username,
				fullname => $row->full_name,
				tenant   => $row->company
			}
		);
	}
	return $self->render( json => $data );
}

# .json format for the jqTree widge
sub resolvers {
	my $self   = shift;
	my $fed_id = $self->param('federation_id');

	my $data;
	my $fed_fed_resolvers =
		$self->db->resultset('FederationFederationResolver')->search( { federation => $fed_id }, { prefetch => ['federation_resolver'] } );
	my $nodes;
	my $resolvers = $self->group_resolvers($fed_id);

	for my $r ( sort( keys(%$resolvers) ) ) {
		my $children;
		my $ip_addresses = $resolvers->{$r};
		foreach my $ip_addr (@$ip_addresses) {
			my $resolver_node = { label => $ip_addr };
			push( @{$children}, $resolver_node );
		}

		$nodes = { label => $r, children => $children };
		push( @{$data}, $nodes );
	}
	return $self->render( json => $data );
}

#  Groups the cname ip addresses together by resolver type.
sub group_resolvers {
	my $self   = shift;
	my $fed_id = shift;

	my $data;
	my $fed_fed_resolvers =
		$self->db->resultset('FederationFederationResolver')->search( { federation => $fed_id }, { prefetch => [ { 'federation_resolver' => 'type' } ] } );
	my $resolvers;
	while ( my $row = $fed_fed_resolvers->next ) {
		my $fed_resolver    = $row->federation_resolver;
		my $fed_resolver_id = $row->federation_resolver->id;
		my $ip_address      = $row->federation_resolver->ip_address;
		my $type_name       = lc $row->federation_resolver->type->name;

		if ( !defined $resolvers->{$type_name} ) {
			$resolvers->{$type_name} = [];
		}
		push( @{ $resolvers->{$type_name} }, $ip_address );
	}
	return $resolvers;
}

sub get_delivery_services {
	my $self   = shift;
	my $id     = shift;
	my @ds_ids = $self->db->resultset('Type')->search( { name => { -like => 'DNS%' } } )->get_column('id')->all;

	my $delivery_services;
	for my $ds_id ( uniq(@ds_ids) ) {
		my $desc = $self->db->resultset('Deliveryservice')->search( { id => $ds_id } )->get_column('xml_id')->single;
		$delivery_services->{$ds_id} = $desc;
	}
	return $delivery_services;
}

# Update
sub update {
	my $self        = shift;
	my $fed_id      = $self->param('federation_id');
	my $ds_id       = $self->param('ds_id');
	my $user_id     = $self->param('user_id');
	my $cname       = $self->param('federation.cname');
	my $description = $self->param('federation.description');
	my $ttl         = $self->param('federation.ttl');

	my $federation_role_id = $self->db->resultset('Role')->search( { name => FEDERATION_ROLE }, undef )->get_column('id')->single();
	my $is_valid = $self->is_valid();
	if ( $self->is_valid("edit") ) {
		my $dbh =
			$self->db->resultset('Federation')->find( { id => $fed_id } );
		$dbh->cname($cname);
		$dbh->description($description);
		$dbh->ttl($ttl);
		$dbh->update();

		my $ft = $self->db->resultset('FederationTmuser')->find_or_create(
			{
				federation => $fed_id,
				tm_user    => $user_id,
				role       => $federation_role_id,
			}
		);

		if ( defined($ft) ) {
			$ft->federation($fed_id);
			$ft->tm_user($user_id);
			$ft->role($federation_role_id);
			$ft->update();
		}

		my $fdses =
			$self->db->resultset('FederationDeliveryservice')->search( { federation => $fed_id }, { prefetch => [ 'federation', 'deliveryservice' ] } );
		while ( my $fd = $fdses->next ) {
			$fd->deliveryservice($ds_id);
			$fd->update();
		}

		$self->flash( message => "Federation was updated successfully." );
		$self->stash( mode => 'edit' );
		return $self->redirect_to( '/federation/' . $fed_id . '/edit' );
	}
	else {
		$self->edit();
	}
}

# Create
sub create {
	my $self    = shift;
	my $ds_id   = $self->param("ds_id");
	my $user_id = $self->param("user_id");
	my $cname   = $self->param("federation.cname");
	my $desc    = $self->param("federation.description");
	my $ttl     = $self->param("federation.ttl");
	&stash_role($self);
	$self->stash(
		deliveryservice_name => undef,
		user_id              => $user_id,
		ds_id                => $ds_id,
		federation           => {},
		fbox_layout          => 1,
		role_name            => FEDERATION_ROLE,
		mode                 => 'add'
	);

	my $existing_fed = $self->db->resultset('Federation')->search( { cname => $cname } )->get_column('cname')->single();
	if ($existing_fed) {
		$self->field('federation.cname')->is_equal( "", "A Federation with name \"$cname\" already exists." );
	}
	if ( !$existing_fed && $self->is_valid("add") ) {
		my $new_id = $self->create_federation( $ds_id, $user_id, $cname, $desc, $ttl );
		if ( $new_id > 0 ) {
			$self->app->log->debug("redirecting....");

			$self->flash( message => "Successfully added Federation!" );
			return $self->redirect_to("/federation/$new_id/edit");
		}
	}
	else {
		return $self->render('federation/add');
	}
}

sub create_federation {
	my $self = shift;

	my $ds_id         = shift;
	my $user_id       = shift;
	my $cname         = shift;
	my $desc          = shift;
	my $ttl           = shift;
	my $federation_id = -1;
	my $fed           = $self->db->resultset('Federation')->create(
		{
			cname       => $cname,
			description => $desc,
			ttl         => $ttl,
		}
	);
	my $f = $fed->insert();
	$federation_id = $f->id;

	if ( $federation_id > 0 ) {
		my $fed_ds_id = -1;
		my $fed_ds    = $self->db->resultset('FederationDeliveryservice')->create(
			{
				federation      => $federation_id,
				deliveryservice => $ds_id,
			}
		);
		$fed_ds_id = $fed_ds->insert();

		if ( $fed_ds_id > 0 ) {
			my $federation_role_id = $self->db->resultset('Role')->search( { name => FEDERATION_ROLE }, undef )->get_column('id')->single();
			my $ft = $self->db->resultset('FederationTmuser')->create(
				{
					federation => $federation_id,
					tm_user    => $user_id,
					role       => $federation_role_id,
				}
			);
		}
		$fed_ds_id = $fed_ds->insert();

		my $ds = $self->db->resultset('Deliveryservice')->search( { id => $ds_id } )->single();

		# if the insert has failed, we don't even get here, we go to the exception page.
		&log( $self, "Created federation with CNAME: " . $cname . " and Delivery Service:  " . $ds->xml_id, "UICHANGE" );
	}
	return $federation_id;

}

sub is_valid {
	my $self = shift;

	$self->field('federation.cname')->is_required;
	$self->field('federation.cname')->is_like( qr/\.$/, "CNAME must end with a period." );
	$self->field('federation.ttl')->is_required;
	$self->field('ds_id')->is_required;
	$self->field('user_id')->is_required;

	return $self->valid;
}

# Delete
sub delete {
	my $self   = shift;
	my $fed_id = $self->param('federation_id');

	if ( !&is_oper($self) ) {
		$self->flash( alertmsg => "No can do. Get more privs." );
	}
	else {

		my $fed_resolver = $self->db->resultset('FederationResolver')
			->search( { 'federation_federation_resolvers.federation' => $fed_id }, { prefetch => 'federation_federation_resolvers' } );

		if ( defined($fed_resolver) ) {
			$fed_resolver->delete();
		}

		my $federation = $self->db->resultset('Federation')->search( { id => $fed_id } )->single();
		if ( defined($federation) ) {
			my $cname = $federation->cname;
			$federation->delete();
			my $msg = sprintf( "Deleted Federation -- cname: %s", $cname );
			&log( $self, $msg, "UICHANGE" );
		}

	}
	return $self->redirect_to('/close_fancybox.html');
}

1;
