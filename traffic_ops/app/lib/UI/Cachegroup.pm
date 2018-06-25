package UI::Cachegroup;

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

# a note about locations and cachegroups. This used to be "Location", before we had physical locations in 12M. Very confusing.
# What used to be called a location is now called a "cachegroup" and location is now a physical address, not a group of caches working together.
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
	$self->stash( fbox_layout => 1, cg_data => {}, coordinate => {} );
	&stash_role($self);
	if ( $self->stash('priv_level') < 30 ) {
		$self->stash( alertmsg => "Insufficient privileges!" );
		$self->redirect_to('/cachegroups');
	}
}

sub view {
	my $self = shift;
	my $mode = $self->param('mode');
	my $id   = $self->param('id');
	$self->stash( cg_data => {} );

	my $rs_param = $self->db->resultset('Cachegroup')->search( { 'me.id' => $id }, { prefetch => 'coordinate' } );
	my $data = $rs_param->single;

	my $type_id = $self->db->resultset('Cachegroup')->search( { id => $id } )->get_column('type')->single();
	my $selected_type = $self->db->resultset('Type')->search( { id => $type_id } )->get_column('name')->single();

	my @param_ids = $self->db->resultset('CachegroupParameter')->search( { cachegroup => $id } )->get_column('parameter')->all();
	$rs_param = $self->db->resultset('Parameter')->search( { id => { -in => \@param_ids } } );
	my @cachegroup_vars;
	while ( my $row = $rs_param->next ) {
		$self->stash( extra_vars => 1 );
		my $var;
		$var->{value} = $row->value;
		$var->{id}    = $row->id;
		$var->{name}  = $row->name;
		my @profiles = $self->db->resultset('ProfileParameter')->search( { parameter => $row->id } )->get_column('profile')->all();
		my @cachegroups = $self->db->resultset('CachegroupParameter')->search( { parameter => $row->id } )->get_column('cachegroup')->all();
		if ( $#profiles > 0 || $#cachegroups > 1 ) {
			$var->{editable} = 0;
		}
		else {
			$var->{editable} = 1;
		}
		push( @cachegroup_vars, $var );
	}

	my $parent_name = "NO_PARENT";
	if ( defined( $data->parent_cachegroup_id ) ) {
		my $parent_id = $data->parent_cachegroup_id;
		my $rs = $self->db->resultset('Cachegroup')->search( { id => $parent_id } )->single;
		$parent_name = $rs->name;
	}

	my $secondary_parent_name = "NO_PARENT";
	if ( defined( $data->secondary_parent_cachegroup_id ) ) {
		my $secondary_parent_id = $data->secondary_parent_cachegroup_id;
		my $rs = $self->db->resultset('Cachegroup')->search( { id => $secondary_parent_id } )->single;
		$secondary_parent_name = $rs->name;
	}

	$self->stash(
		cachegroup_vars                => \@cachegroup_vars,
		parent_cachegroup_id           => $data->parent_cachegroup_id // -1,
		parent_name                    => $parent_name,
		secondary_parent_cachegroup_id => $data->secondary_parent_cachegroup_id // -1,
		secondary_parent_name          => $secondary_parent_name,
	);

	&stash_role($self);
	my $coordinate_data = defined( $data->coordinate ) ? $data->coordinate : {};
	$self->stash( fbox_layout => 1, cg_data => $data, coordinate => $coordinate_data, selected_type => $selected_type );

	if ( $mode eq "edit" and $self->stash('priv_level') > 20 ) {
		$self->render( template => 'cachegroup/edit' );
	}
	else {
		$self->render( template => 'cachegroup/view' );
	}
}

# Read
# Note: This sub will be removed soon. It's just here for backward compatibility with current users of /datalocation that still use
# location as the term for cachegroup
sub read {
	my $self = shift;
	my @data;
	my %idnames;
	my $orderby = "name";
	$orderby = $self->param('orderby') if ( defined $self->param('orderby') );

	# Can't figure out how to do the join on the same table
	my $rs_idnames = $self->db->resultset("Cachegroup")->search( undef, { columns => [qw/id name/] } );
	while ( my $row = $rs_idnames->next ) {
		$idnames{ $row->id } = $row->name;
	}

	my $rs_data = $self->db->resultset("Cachegroup")->search( undef, { prefetch => [ { 'type' => undef, }, 'coordinate' ], order_by => 'me.' . $orderby } );
	while ( my $row = $rs_data->next ) {
		if ( defined $row->parent_cachegroup_id ) {
			push(
				@data, {
					"id"                   => $row->id,
					"name"                 => $row->name,
					"short_name"           => $row->short_name,
					"latitude"             => $row->coordinate->latitude,
					"longitude"            => $row->coordinate->longitude,
					"last_updated"         => $row->last_updated,
					"parent_location_id"   => $row->parent_cachegroup_id,
					"parent_location_name" => $idnames{ $row->parent_cachegroup_id },
					"type_id"              => $row->type->id,
					"type_name"            => $row->type->name,
				}
			);
		}
		else {
			push(
				@data, {
					"id"                   => $row->id,
					"name"                 => $row->name,
					"short_name"           => $row->short_name,
					"latitude"             => $row->coordinate->latitude,
					"longitude"            => $row->coordinate->longitude,
					"last_updated"         => $row->last_updated,
					"parent_location_id"   => $row->parent_cachegroup_id,
					"parent_location_name" => undef,
					"type_id"              => $row->type->id,
					"type_name"            => $row->type->name,
				}
			);
		}
	}
	$self->render( json => \@data );
}

# Read
sub readcachegrouptrimmed {
	my $self = shift;
	my @data;
	my %idnames;
	my $orderby = "name";
	$orderby = $self->param('orderby') if ( defined $self->param('orderby') );

	# Can't figure out how to do the join on the same table
	my $rs_idnames = $self->db->resultset("Cachegroup")->search( undef, { columns => [qw/id name/] } );
	while ( my $row = $rs_idnames->next ) {
		$idnames{ $row->id } = $row->name;
	}
	my $rs_data = $self->db->resultset("Cachegroup")->search( undef, { prefetch => [ { 'type' => undef, } ], order_by => 'me.' . $orderby } );
	while ( my $row = $rs_data->next ) {
		push( @data, { "name" => $row->name, } );
	}
	$self->render( json => \@data );
}

# Delete
sub delete {
	my $self = shift;
	my $id   = $self->param('id');

	if ( !&is_admin($self) ) {
		$self->flash( alertmsg => "You must be an ADMIN to perform this operation!" );
	}
	else {
		my $p_name = $self->db->resultset('Cachegroup')->search( { id => $id } )->get_column('name')->single();
		my $delete = $self->db->resultset('Cachegroup')->search( { 'me.id' => $id }, { prefetch => 'coordinate' } );
		my $coordinate = $delete->single()->coordinate;
		$delete->delete();
		if ( defined( $coordinate ) ) {
			$coordinate->delete();
		}
		&log( $self, "Delete cachegroup " . $p_name, "UICHANGE" );
	}
	return $self->redirect_to('/close_fancybox.html');
}

sub check_cachegroup_input {
	my $self = shift;

	my $sep = "__NEWLINE__";    # the line separator sub that with \n in the .ep javascript
	my $err = undef;

	# First, check permissions
	if ( !&is_oper($self) ) {
		$err .= "You do not have enough privileges to modify this." . $sep;
		return $err;
	}

	return $err;
}

# Update
sub update {
	my $self                           = shift;
	my $id                             = $self->param('id');
	my $priv_level                     = $self->stash('priv_level');
	my $cachegroup_vars                = $self->stash('cachegroup_vars');
	my $extra_vars                     = $self->stash('extra_vars');
	my $parent_cachegroup_id           = $self->param('cg_data.parent_cachegroup_id') // -1;
	my $secondary_parent_cachegroup_id = $self->param('cg_data.secondary_parent_cachegroup_id') // -1;

	$self->stash(
		id              => $id,
		fbox_layout     => 1,
		priv_level      => $priv_level,
		cachegroup_vars => $cachegroup_vars,
		extra_vars      => $extra_vars,
		cg_data         => {
			id                             => $id,
			name                           => $self->param('cg_data.name'),
			short_name                     => $self->param('cg_data.short_name'),
			parent_cachegroup_id           => $parent_cachegroup_id,
			secondary_parent_cachegroup_id => $secondary_parent_cachegroup_id,
			type                           => $self->param('cg_data.type')
		},
		coordinate => {
			latitude => $self->param('coordinate.latitude'),
			longitude => $self->param('coordinate.longitude')
		}
	);

	if ( !$self->isValidCachegroup() ) {
		return $self->render( template => 'cachegroup/edit' );
	}

	# JvD Note: the "cachegroup" parameter in this $self->param is actually the parent_cachegroup_id, because, i'm re-using the
	# cachegroup.js functions.
	my $err = &check_cachegroup_input($self);
	if ( defined($err) ) {
		$self->flash( alertmsg => $err );
	}
	else {
		my $update = $self->db->resultset('Cachegroup')->find( { id => $self->param('id') }, { prefetch => 'coordinate' } );
		$update->name( $self->param('cg_data.name') );
		$update->short_name( $self->param('cg_data.short_name') );
		$update->coordinate->name( 'from_cachegroup_' . $self->param('cg_data.name') );
		$update->coordinate->latitude( $self->param('coordinate.latitude') );
		$update->coordinate->longitude( $self->param('coordinate.longitude') );
		if ( $parent_cachegroup_id != -1 ) {
			$update->parent_cachegroup_id( $self->param('cg_data.parent_cachegroup_id') );
		}
		else {
			$update->parent_cachegroup_id( undef );
		}
		if ( $secondary_parent_cachegroup_id != -1 ) {
			$update->secondary_parent_cachegroup_id( $self->param('cg_data.secondary_parent_cachegroup_id') );
		}
                else {
                        $update->secondary_parent_cachegroup_id( undef );
                }
		$update->type( $self->param('cg_data.type') );
		$update->update();
		$update->coordinate->update();

		foreach my $param ( $self->param ) {
			next unless $param =~ /^param:/;
			my $param_id = $param;
			$param_id =~ s/^param://;
			$update = $self->db->resultset('Parameter')->find( { id => $param_id } );
			$update->value( $self->param($param) );
			$update->update();
		}

		# if the update has failed, we don't even get here, we go to the exception page.
	}

	&log( $self, "Update cachegroup with name:" . $self->param('cg_data.name'), "UICHANGE" );
	$self->flash( message => "Successfully updated Cache Group." );
	return $self->redirect_to( '/cachegroup/edit/' . $id );
}

# Create
sub create {
	my $self        = shift;
	my $name        = $self->param('cg_data.name');
	my $short_name  = $self->param('cg_data.short_name');
	my $latitude    = $self->param('coordinate.latitude');
	my $longitude   = $self->param('coordinate.longitude');
	my $cachegroup  = $self->param('cg_data.parent_cachegroup_id');
	my $type        = $self->param('cg_data.type');
	my $data        = $self->get_cachegroups();
	my $cachegroups = $data->{'cachegroups'};
	my $short_names = $data->{'short_names'};

	if ( !$self->isValidCachegroup() ) {
		$self->stash(
			fbox_layout => 1,
			cg_data     => {
				name       => $name,
				short_name => $short_name,
			},
			coordinate => {
				latitude => $latitude,
				longitude => $longitude
			}
		);
		return $self->render('cachegroup/add');
	}
	if ( exists $cachegroups->{$name} ) {
		$self->field('cg_data.name')->is_like( qr/^\/(?!$name\/)/i, "The name exists." );
		$self->stash(
			fbox_layout => 1,
			cg_data     => {
				name       => $name,
				short_name => $short_name,
			},
			coordinate  => {
				latitude  => $latitude,
				longitude => $longitude
			}
		);
		return $self->render('cachegroup/add');
	}
	if ( exists $short_names->{$short_name} ) {
		$self->field('cg_data.short_name')->is_like( qr/^\/(?!$short_name\/)/i, "The short name exists." );
		$self->stash(
			fbox_layout => 1,
			cg_data     => {
				name       => $name,
				short_name => $short_name,
			},
			coordinate  => {
				latitude  => $latitude,
				longitude => $longitude
			}
		);
		return $self->render('cachegroup/add');
	}

	my $new_id = -1;
	my $err    = &check_cachegroup_input($self);
	if ( defined($err) ) {
		return $self->redirect_to( '/cachegroup/edit/' . $new_id );
	}
	else {

		my $parent_cachegroup_id = $cachegroup;    # sharing the code in JS for create and edit.
		$parent_cachegroup_id = undef if ( $parent_cachegroup_id == -1 );

		my $coordinate = $self->db->resultset('Coordinate')->create(
			{
				name => 'from_cachegroup_' . $name,
				latitude => $latitude,
				longitude => $longitude
			}
		);
		$coordinate->insert();
		my $coordinate_id = $coordinate->id;

		my $insert = $self->db->resultset('Cachegroup')->create(
			{
				name                 => $name,
				short_name           => $short_name,
				coordinate           => $coordinate_id,
				parent_cachegroup_id => $parent_cachegroup_id,
				type                 => $type,
			}
		);
		$insert->insert();
		$new_id = $insert->id;
	}
	if ( $new_id == -1 ) {
		my $referer = $self->req->headers->header('referer');
		return $self->redirect_to($referer);
	}
	else {
		&log( $self, "Create cachegroup with name:" . $self->param('cg_data.name'), "UICHANGE" );
		$self->flash( message => "Successfully updated Cache Group." );
		return $self->redirect_to( '/cachegroup/edit/' . $new_id );
	}
}

sub isValidCachegroup {
	my $self = shift;
	$self->field('cg_data.name')->is_required->is_like( qr/^[0-9a-zA-Z_\.\-]+$/, "Use alphanumeric . or _ ." );
	$self->field('cg_data.short_name')->is_required->is_like( qr/^[0-9a-zA-Z_\.\-]+$/, "Use alphanumeric . or _" );
	$self->field('coordinate.latitude')->is_required->is_like( qr/^[-]*[0-9]+[.]*[0-9]*/, "Invalid latitude entered." );
	$self->field('coordinate.longitude')->is_required->is_like( qr/^[-]*[0-9]+[.]*[0-9]*/, "Invalid latitude entered." );
	my $latitude  = $self->param('coordinate.latitude');
	my $longitude = $self->param('coordinate.longitude');

	if ( abs $latitude > 90 ) {
		$self->field('coordinate.latitude')->is_required->is_like( qr/^\./, "May not exceed +- 90.0." );
	}
	if ( abs $longitude > 180 ) {
		$self->field('coordinate.longitude')->is_required->is_like( qr/^\./, "May not exceed +- 180.0." );
	}

	return $self->valid;
}

sub get_cachegroups {
	my $self = shift;

	my %data;
	my %cachegroups;
	my %short_names;
	my $rs = $self->db->resultset('Cachegroup');
	while ( my $cachegroup = $rs->next ) {
		$cachegroups{ $cachegroup->name }       = $cachegroup->id;
		$short_names{ $cachegroup->short_name } = $cachegroup->id;
	}
	%data = ( cachegroups => \%cachegroups, short_names => \%short_names );

	return \%data;
}

1;
