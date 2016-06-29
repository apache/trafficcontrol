package API::Profile;
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
use Data::Dumper;
use JSON;

sub index {
	my $self = shift;
	my @data;
	my $orderby = $self->param('orderby') || "me.name";
	my $rs_data
		= $self->db->resultset("Profile")
		->search( undef,
		{ order_by => $orderby } );
	while ( my $row = $rs_data->next ) {
		push(
			@data,
			{   "id"          => $row->id,
				"name"        => $row->name,
				"description" => $row->description,
				"lastUpdated" => $row->last_updated,
			}
		);
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

sub create {
  my $self = shift;
  my $params = $self->req->json;
  if ( !defined($params) ) {
      return $self->alert("parameters must be in JSON format,  please check!");
  }

  if ( !&is_oper($self) ) {
      return $self->alert( { Error => " - You must be an admin or oper to perform this operation!" } );
  }

	my $name = $params->{name};
	if ( !defined($name) ) {
		return $self->alert("profile 'name' is not given.");
	}
	if ( $name eq "" ) {
		return $self->alert("profile 'name' can't be null.");
	}

  if ( $name =~ /\s/ ) {
    return $self->alert("Profile name cannot contain space(s).");
  }

	my $description = $params->{description};
	if ( !defined($description) ) {
		return $self->alert("profile 'description' is not given.");
	}
	if ( $description eq "" ) {
		return $self->alert("profile 'description' can't be null.");
	}

	my $existing_profile = $self->db->resultset('Profile')->search( { name        => $name } )->get_column('name')->single();
	if ( $existing_profile && $name eq $existing_profile ) {
		return $self->alert("profile with name $name already exists.");
	}

	my $insert = $self->db->resultset('Profile')->create(
		{
			name        => $name,
			description => $description,
		}
	);
	$insert->insert();
	my $new_id = $insert->id;

	my $response;
	$response->{id} = $new_id;
	$response->{name} = $name;
	$response->{description} = $description;
	return $self->success($response);
}

sub copy {
    my $self = shift;

    if ( !&is_oper($self) ) {
        return $self->alert( { Error => " - You must be an admin or oper to perform this operation!" } );
    }

	my $name = $self->param('profile_name');
	my $profile_copy_from_name = $self->param('profile_copy_from');
    if ( !defined($name) ) {
        return $self->alert("profile 'name' is not given.");
    }
    if ( $name eq "" ) {
        return $self->alert("profile 'name' can't be null.");
    }
    if ( defined($profile_copy_from_name) and ( $profile_copy_from_name eq "" ) ) {
        return $self->alert("profile name 'profile_copy_from' can't be null.");
    }

    my $existing_profile = $self->db->resultset('Profile')->search( { name        => $name } )->get_column('name')->single();
    if ( $existing_profile && $name eq $existing_profile ) {
        return $self->alert("profile with name $name already exists.");
    }

    my $rs = $self->db->resultset('Profile')->search( { name => $profile_copy_from_name } );
    my $row1 = $rs->next;
    if ( !$row1 ) {
        return $self->alert("profile_copy_from $profile_copy_from_name doesn't exist.");
    }
    my $profile_copy_from_id = $row1->id;
    my $description = $row1->description;

    my $insert = $self->db->resultset('Profile')->create(
        {
            name        => $name,
            description => $description,
        }
    );
    $insert->insert();
    my $new_id = $insert->id;

    if ( defined($profile_copy_from_name) ) {
        my $rs_param =
        $self->db->resultset('ProfileParameter')->search( { profile => $profile_copy_from_id }, { prefetch => [ { profile => undef }, { parameter => undef } ] } );
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

    my $response;
    $response->{id} = $new_id;
    $response->{name} = $name;
    $response->{description} = $description;
    $response->{profileCopyFrom} = $profile_copy_from_name;
    $response->{idCopyFrom} = $profile_copy_from_id;
    return $self->success($response);
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
