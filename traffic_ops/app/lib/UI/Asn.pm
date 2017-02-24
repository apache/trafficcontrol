package UI::Asn;
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

sub add {
	my $self = shift;

	print __FILE__ . ":" . __LINE__ . ": asn.add()\n";

	$self->stash( fbox_layout => 1, asn_data => {} );
	&stash_role($self);
	if ( $self->stash('priv_level') < 30 ) {
		$self->stash( alertmsg => "Insufficient privileges!" );
		$self->redirect_to('/asns');
	}
}

# Delete
sub delete {
	my $self = shift;
	my $id   = $self->param('id');

	if ( !&is_admin($self) ) {
		$self->flash( alertmsg => "You must be an ADMIN to perform this operation!" );
	}
	else {
		my $asn    = $self->db->resultset('Asn')->search( { id => $id } )->get_column('asn')->single();
		my $delete = $self->db->resultset('Asn')->search( { id => $id } );
		$delete->delete();
		&log( $self, "Delete asn " . $asn, "UICHANGE" );
	}
	return $self->redirect_to('/asns');
}

sub check_asn_input {
	my $self = shift;

	my $sep = "__NEWLINE__";    # the line separator sub that with \n in the .ep javascript
	my $err = undef;

	# First, check permissions
	if ( !&is_admin($self) ) {
		$err .= "You must be an ADMIN to perform this operation!" . $sep;
		return $err;
	}

	return $err;
}

# Update
sub update {
	my $self       = shift;
	my $id         = $self->param('id');
	my $asn        = $self->param('asn_data.asn');
	my $cachegroup = $self->param('asn_data.cachegroup');

	my $err = &check_asn_input($self);
	if ( defined($err) ) {
		$self->flash( alertmsg => $err );
	}
	else {
		my $update = $self->db->resultset('Asn')->find( { id => $self->param('id') } );
		$update->asn($asn);
		$update->cachegroup($cachegroup);
		$update->update();

		my $name = defined( $self->param('name') ) ? $self->param('name') : "undef";
		my $id   = defined( $self->param('id') )   ? $self->param('id')   : "undef";

		# if the update has failed, we don't even get here, we go to the exception page.
		&log( $self, "Update asn with name " . $name . " and id " . $id, "UICHANGE" );
	}

	$self->flash( alertmsg => "Success!" );
	return $self->redirect_to( '/asns/' . $id . '/edit' );
}

# Create
sub create {
	my $self       = shift;
	my $asn        = $self->param('asn_data.asn');
	my $cachegroup = $self->param('asn_data.cachegroup');
	my $data       = $self->get_asns();
	my $asn_data   = $data->{'asns'};
	my $newid;

	if ( exists $asn_data->{$asn} ) {
		$self->field('asn_data.asn')->is_like( qr/^\/(?!$asn\/)/i, "The ASN already exists." );
		$self->stash(
			fbox_layout => 1,
			asn_data    => {
				asn        => $asn,
				cachegroup => $cachegroup
			}
		);
		return $self->render('asn/add');
	}
	if ( !defined $cachegroup || $cachegroup == -1 ) {
		$self->field('asn_data.cachegroup')->is_like( qr/^\/(?!$cachegroup\/)/i, "Select a cachegroup." );
		$self->stash(
			fbox_layout => 1,
			asn_data    => {
				asn        => $asn,
				cachegroup => $cachegroup
			}
		);
		return $self->render('asn/add');
	}

	my $err = &check_asn_input($self);
	if ( defined($err) ) {
		$self->flash( alertmsg => $err );
	}
	else {
		my $insert = $self->db->resultset('Asn')->create(
			{
				asn        => $self->param('asn_data.asn'),
				cachegroup => $self->param('asn_data.cachegroup'),

			}
		);
		$insert->insert();
		$newid = $insert->id;
		&log( $self, "Create asn with name " . $asn . " and id " . $newid, "UICHANGE" );
	}
	$self->flash( alertmsg => "Success!" );
	return $self->redirect_to( '/asns/' . $newid . "/edit" );
}

sub view {
	my $self = shift;
	my $mode = $self->param('mode');
	my $id   = $self->param('id');

	&stash_role($self);
	$self->stash( asn_data => {}, cg_data => {} );

	my $rs               = $self->db->resultset('Asn')->search( { 'me.id' => $id }, { prefetch => [ 'cachegroup' ] } );
	my $data             = $rs->single;
	my $cache_group_name = $data->cachegroup->name;

	$self->stash( asn_data    => $data );
	$self->stash( fbox_layout => 1 );

	if ( $mode eq "edit" and $self->stash('priv_level') > 20 ) {
		if ( defined $cache_group_name ) {
			$self->stash( cache_group_name => $cache_group_name );
		}
		else {
			$self->stash( cache_group_name => 'NO_CACHEGROUP' );
		}
		$self->render( template => 'asn/edit' );
	}
	else {
		$self->render( template => 'asn/view' );
	}
}

sub get_asns {
	my $self = shift;

	my %data;
	my %asns;
	my $rs = $self->db->resultset('Asn');
	while ( my $asn = $rs->next ) {
		$asns{ $asn->asn } = $asn->asn;
	}
	%data = ( asns => \%asns );

	return \%data;
}

1;
