#
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
package UI::PhysLocation;

use UI::Utils;

use Mojo::Base 'Mojolicious::Controller';
use Email::Valid;
use Data::Dumper;

my $finfo = __FILE__ . ":";

# Table view
sub index {
	my $self = shift;

	&navbarpage($self);
}

sub readphys_location {
	my $self = shift;
	my @data;
	my $rs_data = $self->db->resultset("PhysLocation")->search( undef, { prefetch => ['region'], order_by => 'me.name' } );
	while ( my $row = $rs_data->next ) {

		next if $row->short_name eq 'UNDEF';

		push(
			@data, {
				"id"         => $row->id,
				"name"       => $row->name,
				"short_name" => $row->short_name,
				"address"    => $row->address,
				"city"       => $row->city,
				"state"      => $row->state,
				"zip"        => $row->zip,
				"poc"        => $row->poc,
				"phone"      => $row->phone,
				"email"      => $row->email,
				"comments"   => $row->comments,
				"region"     => $row->region->name,
			}
		);
	}
	$self->render( json => \@data );
}

sub readphys_locationtrimmed {
	my $self = shift;
	my @data;
	my $rs_data = $self->db->resultset("PhysLocation")->search( undef, { order_by => 'me.name' } );
	while ( my $row = $rs_data->next ) {

		next if $row->short_name eq 'UNDEF';

		push(
			@data, {
				"name" => $row->name,
			}
		);
	}
	$self->render( json => \@data );
}

sub readregion {
	my $self = shift;
	my @data;
	my $rs_data = $self->db->resultset("Region")->search( undef, { order_by => 'me.name' } );
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"id"   => $row->id,
				"name" => $row->name,
			}
		);
	}
	$self->render( json => \@data );
}

sub edit {
	my $self = shift;
	my $id   = $self->param('id');

	&stash_role($self);
	my $data = $self->db->resultset('PhysLocation')->search( { id => $id } )->single();

	print "data -> $data";
	$self->stash( location => $data, regionname => $data->region->name );
	$self->stash( fbox_layout => 1 );
	$self->render( template => 'phys_location/view' );

	if ( $self->stash('priv_level') > 20 ) {
		$self->render( template => 'phys_location/edit' );
	}
	else {
		$self->render( template => 'phys_location/view' );
	}
}

sub delete {
	my $self = shift;
	my $id   = $self->param('id');

	if ( !&is_oper($self) ) {
		$self->flash( alertmsg => "No can do. Get more privs." );
	}
	else {
		my $name   = $self->db->resultset('PhysLocation')->search( { id => $id } )->get_column('name')->single();
		my $delete = $self->db->resultset('PhysLocation')->search( { id => $id } );
		$delete->delete();
		&log( $self, "Delete phys_location " . $name, "UICHANGE" );
	}
	return $self->redirect_to('/close_fancybox.html');
}

sub check_phys_location_input {
	my $self = shift;

	my $sep = "__NEWLINE__";    # the line separator sub that with \n in the .ep javascript
	my $err = undef;

	# First, check permissions
	if ( !&is_oper($self) ) {
		$err .= "You do not have enough privileges to modify this." . $sep;
	}
	if ( length( $self->param('location.state') ) != 2 ) {
		$err .= "Use 2 char state abbreviation." . $sep;
	}

	return $err;
}

sub update {
	my $self       = shift;
	my $regioninfo = $self->db->resultset('Region')->search( { id => $self->param('location.region') } )->single();
	$self->stash(
		fbox_layout => 1,
		regionname  => $regioninfo->name,
		location    => {
			name       => $self->param('location.name'),
			short_name => $self->param('location.short_name'),
			address    => $self->param('location.address'),
			city       => $self->param('location.city'),
			state      => $self->param('location.state'),
			zip        => $self->param('location.zip'),
			phone      => $self->param('location.phone'),
			poc        => $self->param('location.poc'),
			email      => $self->param('location.email'),
			comments   => $self->param('location.comments'),
			region     => $self->param('location.region')
		}
	);

	if ( $self->is_valid() ) {
		my $update = $self->db->resultset('PhysLocation')->find( { id => $self->param('id') } );
		$update->name( $self->param('location.name') );
		$update->short_name( $self->param('location.short_name') );
		$update->address( $self->param('location.address') );
		$update->city( $self->param('location.city') );
		$update->state( $self->param('location.state') );
		$update->zip( $self->param('location.zip') );
		$update->phone( $self->param('location.phone') );
		$update->poc( $self->param('location.poc') );
		$update->email( $self->param('location.email') );
		$update->comments( $self->param('location.comments') );
		$update->region( $self->param('location.region') );
		$update->update();

		my $name = defined( $self->param('location.name') ) ? $self->param('location.name') : "undef";
		my $id   = defined( $self->param('id') )   ? $self->param('id')   : "undef";

		# if the update has failed, we don't even get here, we go to the exception page.
		&log( $self, "Update phys_location with name " . $name . " and id " . $id, "UICHANGE" );
		$self->flash( message => "Success!" );
		return $self->redirect_to( '/phys_location/' . $id . '/edit' );
	}
	else {
		return $self->render('phys_location/edit');
	}

	my $referer = $self->req->headers->header('referer');
	return $self->redirect_to($referer);
}

sub add {
	my $self   = shift;
	my @params = $self->param;

	$self->stash( fbox_layout => 1, );
	$self->stash( location    => {} );
	&stash_role($self);
	if ( $self->stash('priv_level') < 30 ) {
		$self->stash( alertmsg => "Insufficient privileges!" );
		$self->redirect_to('/phys_locations');
	}
	foreach my $field (@params) {
		$self->stash( $field => $self->param($field) );
	}
}

sub create {
	my $self           = shift;
	my $name           = $self->param('location.name'),
		my $short_name = $self->param('location.short_name'),
		my $address    = $self->param('location.address'),
		my $city       = $self->param('location.city'),
		my $state      = $self->param('location.state'),
		my $zip        = $self->param('location.zip'),
		my $phone      = $self->param('location.phone'),
		my $poc        = $self->param('location.poc'),
		my $email      = $self->param('location.email'),
		my $comments   = $self->param('location.comments'),
		my $region     = $self->param('location.region'),
		my $data       = $self->get_phsyical_location_names();

	my $names       = $data->{'names'};
	my $short_names = $data->{'short_names'};
	my $regioninfo  = $self->db->resultset('Region')->search( { id => $self->param('location.region') } )->single();

	if ( exists $names->{$name} ) {
		$self->field('location.name')->is_like( qr/^\/(?!$name\/)/, "The name chosen is already used." );
		$self->stash(
			fbox_layout => 1,
			location    => {
				name       => $name,
				short_name => $short_name,
				address    => $address,
				city       => $city,
				state      => $state,
				zip        => $zip,
				phone      => $phone,
				poc        => $poc,
				email      => $email,
				comments   => $comments,
				region     => $region
			}
		);
		return $self->render('phys_location/add');
	}
	if ( exists $short_names->{$short_name} ) {
		$self->field('location.short_name')->is_like( qr/^\/(?!$short_name\/)/, "The short name chosen is already used." );
		$self->stash(
			fbox_layout => 1,
			location    => {
				name       => $name,
				short_name => $short_name,
				address    => $address,
				city       => $city,
				state      => $state,
				zip        => $zip,
				phone      => $phone,
				poc        => $poc,
				email      => $email,
				comments   => $comments,
				region     => $region
			}
		);
		return $self->render('phys_location/add');
	}
	if ( $phone eq "000-000-0000" ) {
		$phone = "";
	}
	my $new_id = -1;

	$self->stash(
		fbox_layout => 1,
		regionname  => $regioninfo->name,
		location    => {
			name       => $name,
			short_name => $short_name,
			address    => $address,
			city       => $city,
			state      => $state,
			zip        => $zip,
			phone      => $phone,
			poc        => $poc,
			email      => $email,
			comments   => $comments,
			region     => $region
		}
	);
	if ( !$self->is_valid() ) {
		return $self->render('phys_location/add');
	}
	else {
		my $insert = $self->db->resultset('PhysLocation')->create(
			{
				name       => $name,
				short_name => $short_name,
				address    => $address,
				city       => $city,
				state      => $state,
				zip        => $zip,
				phone      => $phone,
				poc        => $poc,
				email      => $email,
				comments   => $comments,
				region     => $region,
			}
		);
		$insert->insert();
		$new_id = $insert->id;
	}

	$self->flash( message => "Success!" );
	return $self->redirect_to( '/phys_location/' . $new_id . '/edit' );
}

sub get_phsyical_location_names {
	my $self = shift;

	my %data;
	my %names;
	my %short_names;
	my $rs = $self->db->resultset('PhysLocation');
	while ( my $row = $rs->next ) {
		$names{ $row->name }             = $row->id;
		$short_names{ $row->short_name } = $row->id;
	}

	%data = ( names => \%names, short_names => \%short_names );

	return \%data;
}

sub is_valid {
	my $self = shift;

	$self->field('location.name')->is_required->is_like( qr/^[a-zA-Z0-9\s]+/, "Use alphanumeric characters." );
	$self->field('location.short_name')->is_required->is_like( qr/^[a-zA-Z0-9\ ]+/, "Use alphanumeric characters." );
	$self->field('location.address')->is_required->is_like( qr/^[a-zA-Z0-9\ \.\-]+/, "Use alphanumeric characters, '.','-', or space." );
	$self->field('location.city')->is_required->is_like( qr/^[a-zA-Z0-9\ \.\-]+/, "Use alphanumeric characters, '.','-', or space." );
	$self->field('location.state')->is_required->is_like( qr/^[A-Z]{2}/, "Uppercase 2 char. state abbreviation." );
	$self->field('location.phone')->is_like( qr/^$|[0-9]{3}\-[0-9]{3}\-[0-9]{4}/, "Phone number format is: ###-###-####" );
	$self->field('location.email')->check( sub {
			my ( $value, $params ) = @_;
			return if Email::Valid->address($value);
			return  "Enter a valid email address.";
			});
	return $self->valid;
}

1;
