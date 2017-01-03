package API::Profile;
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

sub index {
	my $self = shift;
	my @data;
	my $orderby = $self->param('orderby') || "me.name";
	my $parameter_id = $self->param('param');

	if ( defined $parameter_id ) {
		my $rs = $self->db->resultset('ProfileParameter')->search( { parameter => $parameter_id },  { order_by => $orderby }  );
		while ( my $row = $rs->next ) {
			push(
				@data, {
					"id" => $row->profile->id,
					"name" => $row->profile->name,
					"description" => $row->profile->description,
					"cdn" => $row->profile->cdn,
					"type" => $row->profile->type,
					"lastUpdated" => $row->profile->last_updated
				}
			);
		}
	} else {
		my $rs_data = $self->db->resultset("Profile")->search( undef, { order_by => $orderby } );
		while ( my $row = $rs_data->next ) {
			push(
				@data, {
					"id"          => $row->id,
					"name"        => $row->name,
					"description" => $row->description,
					"cdn"         => defined($row->cdn) ? $row->cdn->name : "-",
					"type"        => $row->type,
					"lastUpdated" => $row->last_updated
				}
			);
		}
	}
	$self->success( \@data );
}

sub index_trimmed {
	my $self = shift;
	my @data;
	my $orderby = $self->param('orderby') || "name";
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

sub get_profiles_by_paramId {
	my $self    	= shift;
	my $param_id	= $self->param('id');

	my $param_profiles = $self->db->resultset('ProfileParameter')->search( { parameter => $param_id } );

	my $profiles = $self->db->resultset('Profile')->search(
		{ 'me.id' => { -in => $param_profiles->get_column('profile')->as_query } }
	);

	my @data;
	if ( defined($profiles) ) {
		while ( my $row = $profiles->next ) {
			push(
				@data, {
					"id"          => $row->id,
					"name"        => $row->name,
					"description" => $row->description,
					"lastUpdated" => $row->last_updated
				}
			);
		}
	}

	return $self->success( \@data );
}

sub show {
	my $self = shift;
	my $id   = $self->param('id');

	my $rs_data = $self->db->resultset("Profile")->search( { id => $id } );
	my @data = ();
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"id"          => $row->id,
				"name"        => $row->name,
				"description" => $row->description,
				"lastUpdated" => $row->last_updated
			}
		);
	}
	$self->success( \@data );
}

sub create {
	my $self   = shift;
	print "KK:\n";
	my $params = $self->req->json;
	if ( !defined($params) ) {
		return $self->alert("parameters must be in JSON format,  please check!");
	}

	if ( !&is_oper($self) ) {
		return $self->alert( { Error => " - You must be an admin or oper to perform this operation!" } );
	}

	my $name = $params->{name};
	if ( !defined($name) || $name eq "" || $name =~ /\s/ ) {
		return $self->alert("profile 'name' is required and cannot contain spaces.");
	}

	my $description = $params->{description};
	if ( !defined($description) || $description eq "" ) {
		return $self->alert("profile 'description' is required.");
	}

	my $existing_profile = $self->db->resultset('Profile')->search( { name => $name } )->get_column('name')->single();
	if ( $existing_profile && $name eq $existing_profile ) {
		return $self->alert("profile with name $name already exists.");
	}

	my $existing_desc = $self->db->resultset('Profile')->find( { description => $description } );
	if ($existing_desc) {
		return $self->alert("a profile with the exact same description already exists.");
	}

	my $cdn = $params->{cdn};
	my $type = $params->{type};
	my $insert = $self->db->resultset('Profile')->create(
		{
			name        => $name,
			description => $description,
			cdn         => $cdn,
			type        => $type,
		}
	);
	$insert->insert();
	my $new_id = $insert->id;

	&log( $self, "Created profile with id: " . $new_id . " and name: " . $name, "APICHANGE" );

	my $response;
	$response->{id}          = $new_id;
	$response->{name}        = $name;
	$response->{description} = $description;
	$response->{cdn}         = $cdn;
	$response->{type}        = $type;
	return $self->success($response);
}

sub copy {
	my $self = shift;

	if ( !&is_oper($self) ) {
		return $self->alert( { Error => " - You must be an admin or oper to perform this operation!" } );
	}

	my $name                   = $self->param('profile_name');
	my $profile_copy_from_name = $self->param('profile_copy_from');
	if ( !defined($name) || $name eq "" || $name =~ /\s/ ) {
		return $self->alert("profile 'name' is required and cannot contain spaces.");
	}
	if ( defined($profile_copy_from_name) and ( $profile_copy_from_name eq "" ) ) {
		return $self->alert("profile name 'profile_copy_from' can't be null.");
	}

	my $existing_profile = $self->db->resultset('Profile')->search( { name => $name } )->get_column('name')->single();
	if ( $existing_profile && $name eq $existing_profile ) {
		return $self->alert("profile with name $name already exists.");
	}

	my $rs = $self->db->resultset('Profile')->search( { name => $profile_copy_from_name } );
	my $row1 = $rs->next;
	if ( !$row1 ) {
		return $self->alert("profile_copy_from $profile_copy_from_name doesn't exist.");
	}
	my $profile_copy_from_id = $row1->id;
	my $description          = $row1->description;

	my $cdn = $row1->cdn;
	my $type = $row1->type;
	my $insert = $self->db->resultset('Profile')->create(
		{
			name        => $name,
			description => $description,
			cdn         => $cdn,
			type        => $type,
		}
	);
	$insert->insert();
	my $new_id = $insert->id;

	if ( defined($profile_copy_from_name) ) {
		my $rs_param =
			$self->db->resultset('ProfileParameter')
			->search( { profile => $profile_copy_from_id }, { prefetch => [ { profile => undef }, { parameter => undef } ] } );
		while ( my $row = $rs_param->next ) {
			my $insert = $self->db->resultset('ProfileParameter')->create(
				{
					profile   => $new_id,
					parameter => $row->parameter->id,
				}
			);
			$insert->insert();
		}
	}

	&log( $self, "Created profile from copy with id: " . $new_id . " and name: " . $name, "APICHANGE" );

	my $response;
	$response->{id}              = $new_id;
	$response->{name}            = $name;
	$response->{description}     = $description;
	$response->{profileCopyFrom} = $profile_copy_from_name;
	$response->{idCopyFrom}      = $profile_copy_from_id;
	return $self->success($response);
}

sub update {
	my $self   = shift;
	my $id     = $self->param('id');
	my $params = $self->req->json;

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	my $profile = $self->db->resultset('Profile')->find( { id => $id } );
	if ( !defined($profile) ) {
		return $self->not_found();
	}

	if ( !defined($params) ) {
		return $self->alert("parameters must be in JSON format.");
	}

	my $name = $params->{name};
	if ( !defined($name) || $name eq "" || $name =~ /\s/ ) {
		return $self->alert("profile 'name' is required and cannot contain spaces.");
	}
	if ( $profile->name ne $name ) {
		my $existing = $self->db->resultset('Profile')->find( { name => $name } );
		if ($existing) {
			return $self->alert( "a profile with name " . $name . " already exists." );
		}
	}

	my $description = $params->{description};
	if ( !defined($description) || $description eq "" ) {
		return $self->alert("profile 'description' is required.");
	}
	if ( $profile->description ne $description ) {
		my $existing = $self->db->resultset('Profile')->find( { description => $description } );
		if ($existing) {
			return $self->alert("a profile with the exact same description already exists.");
		}
	}

	my $cdn = $params->{cdn};
	my $type = $params->{type};
	my $values = {
		name        => $name,
		description => $description,
		cdn         => $cdn,
		type        => $type,
	};

	my $rs = $profile->update($values);
	if ($rs) {
		my $response;
		$response->{id}          = $id;
		$response->{name}        = $name;
		$response->{description} = $description;
		$response->{cdn}         = $cdn;
		$response->{type}        = $type;
		&log( $self, "Update profile with id: " . $id . " and name: " . $name, "APICHANGE" );
		return $self->success( $response, "Profile was updated: " . $id );
	}
	else {
		return $self->alert("Profile update failed.");
	}
}

sub delete {
	my $self = shift;
	my $id   = $self->param('id');

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	my $profile = $self->db->resultset('Profile')->find( { id => $id } );
	if ( !defined($profile) ) {
		return $self->not_found();
	}

	if ($profile->name eq 'GLOBAL') {
		return $self->alert("Cannot delete the GLOBAL profile.");
	}

	my $server = $self->db->resultset('Server')->find( { profile => $profile->id } );
	if ( defined($server) ) {
		return $self->alert("the profile is used by some server(s).");
	}
	my $ds = $self->db->resultset('Deliveryservice')->find( { profile => $profile->id } );
	if ( defined($ds) ) {
		return $self->alert("the profile is used by some deliveryservice(s).");
	}

	my $profile_name = $profile->name;
	$profile->delete();

	&log( $self, "Delete profile with id: " . $id . " and name: " . $profile_name, "APICHANGE" );

	return $self->success_message("Profile was deleted.");
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
			push( @data, { "id" => $row->id, "description" => $row->description } );
		}
	}

	$self->success( \@data );
}

1;
