package UI::Profile;
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
use JSON;

# Table view
sub index {
	my $self = shift;

	&navbarpage($self);
	$self->stash( profile => {} );
}

# for the fancybox view
sub add {
	my $self     = shift;
	my %profiles = get_profiles($self);

	$self->stash_cdn_selector();
	$self->stash_profile_type_selector();

	$self->stash( profile => {}, profiles => \%profiles, fbox_layout => 1 );
}

sub edit {
	my $self   = shift;
	my $id     = $self->param('id');
	my $cursor = $self->db->resultset('Profile')->search( { id => $id } );
	my $data   = $cursor->single;

	$self->stash_cdn_selector(defined($data->cdn) ? $data->cdn->id : undef);
	$self->stash_profile_type_selector($data->type);

	&stash_role($self);
	$self->stash( profile => $data, id => $data->id, routing_disabled => $data->routing_disabled, fbox_layout => 1 );
	return $self->render('profile/edit');
}

# for the fancybox view
sub import {
	my $self = shift;
	$self->stash( fbox_layout => 1, msgs => [] );
}

# for the fancybox view
sub view {
	my $self = shift;

	# my $mode = $self->param('mode');
	my $id = $self->param('id');

	my $rs_param = $self->db->resultset('Profile')->search( { id => $id } );
	my $data = $rs_param->single;
	my $param_count = $self->db->resultset('ProfileParameter')->search( { profile => $id } )->count();

	$self->stash( profile     => $data );
	$self->stash( param_count => $param_count );

	&stash_role($self);

	$self->stash( fbox_layout => 1 );

}

# Read
sub readprofile {
	my $self = shift;
	my @data;
	my $orderby = "name";
	$orderby = $self->param('orderby') if ( defined $self->param('orderby') );
	my $rs_data = $self->db->resultset("Profile")->search( undef, { prefetch => ['cdn'], order_by => 'me.' . $orderby } );
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"id"           => $row->id,
				"name"         => $row->name,
				"type"         => $row->type,
				"cdn"          => defined($row->cdn) ? $row->cdn->name : undef,
				"description"  => $row->description,
				"routing_disabled" => $row->routing_disabled,
				"last_updated" => $row->last_updated,
			}
		);
	}
	$self->render( json => \@data );
}

# Read
sub readprofiletrimmed {
	my $self = shift;
	my @data;
	my $orderby = "name";
	$orderby = $self->param('orderby') if ( defined $self->param('orderby') );
	my $rs_data = $self->db->resultset("Profile")->search( undef, { order_by => $orderby } );
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"name" => $row->name,
			}
		);
	}
	$self->render( json => \@data );
}

# Delete
sub delete {
	my $self = shift;
	my $id   = $self->param('id');

	if ( !&is_admin($self) ) {
		$self->flash( message => "You must be an ADMIN to perform this operation!" );
	}
	else {
		my $p_name = $self->db->resultset('Profile')->search( { id => $id } )->get_column('name')->single();
		my $delete = $self->db->resultset('Profile')->search( { id => $id } );
		$delete->delete();
		&log( $self, "Delete profile " . $p_name, "UICHANGE" );
	}
	return $self->redirect_to('/close_fancybox.html');
}

sub check_profile_input {
	my $self        = shift;
	my $mode        = shift;
	my $name        = $self->param('profile.name');
	my $description = $self->param('profile.description');
	my $routing_disabled = $self->param('profile.routing_disabled');

	#Check required fields
	$self->field('profile.name')->is_required;
	$self->field('profile.description')->is_required;
	$self->field('profile.type')->is_required;

	$self->field('profile.name')->is_like( qr/^\S+$/, "Profile name cannot contain space(s)." );

	if ( $mode eq 'add' ) {

		#Check for duplicate profile name and description for NEW
		my $existing_profile = $self->db->resultset('Profile')->search( { name        => $name } )->get_column('name')->single();
		my $existing_desc    = $self->db->resultset('Profile')->search( { description => $description } )->get_column('description')->single();
		if ( $existing_profile && $name eq $existing_profile ) {
			$self->field('profile.name')->is_equal( "", "Profile with name \"$name\" already exists." );
		}

		if ( $existing_desc && $description eq $existing_desc ) {
			$self->field('profile.description')->is_equal( "", "Profile with the exact same description already exists" );
		}
	}
	if ( $mode eq 'edit' ) {

		#make sure user didnt enter a name that is already used by another profile.
		my $id = $self->param('id');

		#get original name
		my $profile_rs = $self->db->resultset('Profile');
		my $orig_name = $profile_rs->search( { id => $id } )->get_column('name')->single();
		if ( $name ne $orig_name ) {
			my $profiles = $profile_rs->search( { id => { '!=' => $id } } )->get_column('name');
			while ( my $db_name = $profiles->next ) {
				if ( $db_name eq $name ) {
					$self->field('profile.name')->is_equal( "", "Profile with name \"$name\" already exists." );
				}
			}
		}

		#get original desc
		my $orig_desc = $profile_rs->search( { id => $id } )->get_column('description')->single();
		if ( $description ne $orig_desc ) {

			#get all other descriptions
			my $profiles = $profile_rs->search( { id => { '!=' => $id } } )->get_column('description');
			while ( my $db_desc = $profiles->next ) {
				if ( $db_desc eq $description ) {
					$self->field('profile.description')->is_equal( "", "A profile with the exact same description already exists!" );
				}
			}
		}

		#make sure the CDN matches servers already assigned to the profile
		my $profile = $self->db->resultset('Profile')->search( { 'me.id' => $id}, { prefetch => ['servers'] } )->first();
		my $cdn = $self->param('profile.cdn');
		my $ex_server = $profile->servers->first;
		if ( defined $ex_server ) {
			if ( $cdn != $ex_server->cdn_id ) {
				$self->field('profile.cdn')->is_equal( "", "The assigned CDN does not match the CDN assigned to servers with this profile!" );
			}
		}
	}
	return $self->valid;
}

# Update
sub update {
	my $self        = shift;
	my $id          = $self->param('id');
	my $name        = $self->param('profile.name');
	my $description = $self->param('profile.description');
	my $cdn         = $self->param('profile.cdn');
	my $type        = $self->param('profile.type');
	my $routing_disabled = $self->param('profile.routing_disabled');

	if ( $self->check_profile_input("edit") ) {

		my $update = $self->db->resultset('Profile')->find( { id => $id } );
		$update->name($name);
		$update->description($description);
		$update->cdn($cdn);
		$update->type($type);
		$update->routing_disabled($routing_disabled);
		$update->update();

		# if the update has failed, we don't even get here, we go to the exception page.
		&log( $self, "Update profile with name: $name", "UICHANGE" );

		$self->flash( message => "Success!" );
		return $self->redirect_to("/profile/$id/view");
	}
	else {
		&stash_role($self);

		my $cursor = $self->db->resultset('Profile')->search( { id => $id } );
		my $data   = $cursor->single;

		$self->stash_cdn_selector(defined($data->cdn) ? $data->cdn->id : undef);
		$self->stash( profile => {}, fbox_layout => 1 );
		$self->stash_profile_type_selector($data->type);
		$self->render('profile/edit');
	}

}

sub create {
	my $self   = shift;
	my $new_id = -1;
	my $p_name = $self->param('profile.name');
	my $p_desc = $self->param('profile.description');
	my $p_cdn         = $self->param('profile.cdn');
	my $p_type        = $self->param('profile.type');
	my $routing_disabled = $self->param('profile.routing_disabled');

	print ">>> cdn: $p_cdn t: $p_type \n";
	if ( !&is_admin($self) ) {
		my $err = "You must be an ADMIN to perform this operation!" . "__NEWLINE__";
		return $self->flash( message => $err );
	}
	if ( $self->check_profile_input("add") ) {
		my $insert = $self->db->resultset('Profile')->create(
			{
				name        => $p_name,
				description => $p_desc,
				cdn         => $p_cdn,
				type        => $p_type,
				routing_disabled => $routing_disabled,
			}
		);
		$insert->insert();
		$new_id = $insert->id;

		# if the insert has failed, we don't even get here, we go to the exception page.
		&log( $self, "Create profile with name:" . $self->param('profile.name'), "UICHANGE" );

		if ( defined( $self->param('copy_from_id') ) ) {
			my $cp_id = $self->param('copy_from_id');
			my $rs_param =
				$self->db->resultset('ProfileParameter')->search( { profile => $cp_id }, { prefetch => [ { profile => undef }, { parameter => undef } ] } );
			my $p_name = "";
			while ( my $row = $rs_param->next ) {
				my $insert = $self->db->resultset('ProfileParameter')->create(
					{
						profile   => $new_id,
						parameter => $row->parameter->id,
					}
				);
				$insert->insert();
				$p_name = $row->profile->name;
			}
			&log( $self, "Copy parameter assignments from " . $p_name . " to " . $self->param('name'), "UICHANGE" );
		}
		$self->flash( message => "Success!" );
		return $self->redirect_to("/profile/$new_id/view");
	}
	else {
		&stash_role($self);
		my %profiles = &get_profiles($self);
		$self->stash( profile => {}, profiles => \%profiles, fbox_layout => 1 );
		$self->render('profile/add');
	}
}

sub doImport {
	my $self             = shift;
	my $new_id           = -1;
	my $in_data          = $self->param('profile_to_import');
	my $f_name           = $in_data->{filename};
	my $data             = JSON->new->utf8->decode( $in_data->asset->{content} );
	my $p_name           = $data->{profile}->{name};
	my $p_desc           = $data->{profile}->{description};
	my $p_type           = $data->{profile}->{type};
	my $existing_profile = $self->db->resultset('Profile')->search( { name => $p_name } )->get_column('name')->single();
	my $existing_desc    = $self->db->resultset('Profile')->search( { description => $p_desc } )->get_column('description')->single();
	my @valid_types      = @{$self->db->source('ProfileTypeValue')->column_info('value')->{extra}->{list}};
	my @msgs;

	if ($existing_profile) {
		push( @msgs, "A profile with the name \"$p_name\" already exists!" );
	}
	if ($existing_desc) {
		push( @msgs, "A profile with the exact same description already exists!" );
	}

	if (! grep(/^$p_type$/, @valid_types )) {
		my $vtypes = join(', ', @valid_types);
		push( @msgs, "Profile contains type \"$p_type\" which is not a valid profile type. Valid types are: $vtypes" );
	}

	my $msgs_size = @msgs;
	if ( $msgs_size > 0 ) {
		&stash_role($self);
		$self->stash( fbox_layout => 1, msgs => \@msgs );
		return $self->render('profile/import');
	}
	else {
		my $insert = $self->db->resultset('Profile')->create(
			{
				name        => $p_name,
				description => $p_desc,
				type        => $p_type,
			}
		);
		$insert->insert();
		$new_id = $insert->id;

		my $new_count      = 0;
		my $existing_count = 0;
		my %done;
		foreach my $param ( @{ $data->{parameters} } ) {
			my $param_name        = $param->{name};
			my $param_config_file = $param->{config_file};
			my $param_value       = $param->{value};
			my $param_id =
				$self->db->resultset('Parameter')
				->search( { -and => [ name => $param_name, value => $param_value, config_file => $param_config_file ] }, { rows => 1 } )->get_column('id')
				->single();
			if ( !defined($param_id) ) {
				my $insert = $self->db->resultset('Parameter')->create(
					{
						name        => $param_name,
						config_file => $param_config_file,
						value       => $param_value,
					}
				);
				$insert->insert();
				$param_id = $insert->id();
				$new_count++;
			}
			else {
				next if defined( $done{$param_id} );    # sometimes the profiles we import have dupes?
				$existing_count++;
			}

			my $link_insert = $self->db->resultset('ProfileParameter')->create(
				{
					parameter => $param_id,
					profile   => $new_id,
				}
			);
			$link_insert->insert();
			$done{$param_id} = $new_id;
		}
		&log( $self, "Import profile " . $p_name . " with " . $new_count . " new and " . $existing_count . " existing parameters.", "UICHANGE" );
		$self->flash( message => => "Success!" );
		return $self->redirect_to("/profile/$new_id/view");
	}
}

sub availableprofile {
	my $self = shift;
	my @data;
	my $paramid = $self->param('paramid');
	my %dsids;
	my %in_use;

	# Get a list of all profile id's associated with this param id
	my $rs_in_use = $self->db->resultset("ProfileParameter")->search( { 'parameter' => $paramid } );
	while ( my $row = $rs_in_use->next ) {
		$in_use{ $row->profile->id } = undef;
	}

	# Add remaining profile ids to @data
	my $rs_links = $self->db->resultset("Profile")->search( undef, { order_by => "description" } );
	while ( my $row = $rs_links->next ) {
		if ( !exists( $in_use{ $row->id } ) ) {
			push( @data, { "id" => $row->id, "name" => $row->name, "description" => $row->description } );
		}
	}

	$self->render( json => \@data );
}

sub export {
	my $self = shift;
	my $id   = $self->param('id');

	my $jdata = {};
	my $pname;
	my $rs = $self->db->resultset('ProfileParameter')->search( { profile => $id }, { prefetch => [ { parameter => undef }, { profile => undef } ] } );
	my $i = 0;
	while ( my $row = $rs->next ) {
		if ( !defined( $jdata->{profile} ) ) {
			$jdata->{profile}->{name}        = $row->profile->name;
			$jdata->{profile}->{description} = $row->profile->description;
			$jdata->{profile}->{type}        = $row->profile->type;
			$pname                           = $row->profile->name;
		}
		$jdata->{parameters}->[$i] = {
			name        => $row->parameter->name,
			config_file => $row->parameter->config_file,
			value       => $row->parameter->value
		};
		$i++;
	}

	my $text = JSON->new->utf8->encode($jdata);

	$self->res->headers->content_type("application/download");
	my $fname = $pname . ".traffic_ops";
	$self->res->headers->content_disposition("attachment; filename=\"$fname\"");
	$self->render( text => $text, format => 'txt', profile => {} );
}

sub compareprofile {
	my $self   = shift;
	my $pid1   = $self->param('profile1');
	my $pid2   = $self->param('profile2');
	my $pname1 = $self->db->resultset('Profile')->search( { id => $pid1 } )->get_column('name')->single();
	my $pname2 = $self->db->resultset('Profile')->search( { id => $pid2 } )->get_column('name')->single();

	&stash_role($self);
	$self->stash( pname1 => $pname1, pname2 => $pname2, pid1 => $pid1, pid2 => $pid2 );
	&navbarpage($self);
}

sub acompareprofile {
	my $self = shift;
	my $pid1 = $self->param('profile1');
	my $pid2 = $self->param('profile2');
	my %data = ( "aaData" => undef );

	my $rs = $self->db->resultset('ProfileParameter')->search( { profile => $pid1 }, { prefetch => [ { parameter => undef }, { profile => undef } ] } );
	my $params1;
	while ( my $row = $rs->next ) {
		$params1->{ $row->parameter->name . ':' . $row->parameter->config_file } =
			{ name => $row->parameter->name, config_file => $row->parameter->config_file, value => $row->parameter->value };
	}

	$rs = $self->db->resultset('ProfileParameter')->search( { profile => $pid2 }, { prefetch => [ { parameter => undef }, { profile => undef } ] } );
	my $params2;
	while ( my $row = $rs->next ) {
		$params2->{ $row->parameter->name . ':' . $row->parameter->config_file } =
			{ name => $row->parameter->name, config_file => $row->parameter->config_file, value => $row->parameter->value };
	}

	foreach my $key ( keys %{$params1} ) {
		if ( !defined( $params2->{$key} ) ) {
			my @line = [ $params1->{$key}->{name}, $params1->{$key}->{config_file}, $params1->{$key}->{value}, "undef" ];
			push( @{ $data{'aaData'} }, @line );
		}
		elsif ( $params1->{$key}->{value} ne $params2->{$key}->{value} ) {
			my @line = [ $params1->{$key}->{name}, $params1->{$key}->{config_file}, $params1->{$key}->{value}, $params2->{$key}->{value} ];
			push( @{ $data{'aaData'} }, @line );
		}
		delete $params2->{$key};
	}

	foreach my $key ( keys %{$params2} ) {
		my @line = [ $params2->{$key}->{name}, $params2->{$key}->{config_file}, "undef", $params2->{$key}->{value} ];
		push( @{ $data{'aaData'} }, @line );
	}

	$self->render( json => \%data );
}

sub get_profiles {
	my $self = shift;
	my %profiles;
	my $p_rs = $self->db->resultset("Profile")->search( undef, { order_by => "name" } );
	while ( my $profile = $p_rs->next ) {
		$profiles{ $profile->name . " - " . $profile->description } = $profile->id;
	}
	return %profiles;
}
1;
