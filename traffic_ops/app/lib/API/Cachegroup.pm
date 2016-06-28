package API::Cachegroup;
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
# a note about locations and cachegroups. This used to be "Location", before we had physical locations in 12M. Very confusing.
# What used to be called a location is now called a "cache group" and location is now a physical address, not a group of caches working together.
#

# JvD Note: you always want to put Utils as the first use. Sh*t don't work if it's after the Mojo lines.
use UI::Utils;

use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use JSON;
use MojoPlugins::Response;

# Read
sub index {
	my $self = shift;
	my @data;
	my %idnames;
	my $orderby = $self->param('orderby') || "name";

	# Can't figure out how to do the join on the same table
	my $rs_idnames = $self->db->resultset("Cachegroup")->search( undef, { columns => [qw/id name/] } );
	while ( my $row = $rs_idnames->next ) {
		$idnames{ $row->id } = $row->name;
	}

	my $rs_data = $self->db->resultset("Cachegroup")->search( undef, { prefetch => [ { 'type' => undef, } ], order_by => 'me.' . $orderby } );
	while ( my $row = $rs_data->next ) {
		if ( defined $row->parent_cachegroup_id ) {
			push(
				@data, {
					"id"                   => $row->id,
					"name"                 => $row->name,
					"shortName"            => $row->short_name,
					"latitude"             => $row->latitude,
					"longitude"            => $row->longitude,
					"lastUpdated"          => $row->last_updated,
					"parentCachegroupId"   => $row->parent_cachegroup_id,
					"parentCachegroupName" => $idnames{ $row->parent_cachegroup_id },
					"typeId"               => $row->type->id,
					"typeName"             => $row->type->name,
				}
			);
		}
		else {
			push(
				@data, {
					"id"                   => $row->id,
					"name"                 => $row->name,
					"shortName"            => $row->short_name,
					"latitude"             => $row->latitude,
					"longitude"            => $row->longitude,
					"lastUpdated"          => $row->last_updated,
					"parentCachegroupId"   => $row->parent_cachegroup_id,
					"parentCachegroupName" => undef,
					"typeId"               => $row->type->id,
					"typeName"             => $row->type->name,
				}
			);
		}
	}
	$self->success( \@data );
}

# Read
sub index_trimmed {
	my $self = shift;
	my @data;
	my %idnames;
	my $orderby = $self->param('orderby') || "name";

	my $rs_data = $self->db->resultset("Cachegroup")->search( undef, { prefetch => [ { 'type' => undef, } ], order_by => 'me.' . $orderby } );
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"name" => $row->name,
			}
		);
	}
	$self->success( \@data );
}

sub by_parameter_id {
	my $self    = shift;
	my $paramid = $self->param('paramid');

	my @data;
	my %dsids;
	my %in_use;

	# Get a list of all cachegroup id's associated with this param id
	my $rs_in_use = $self->db->resultset("CachegroupParameter")->search( { 'parameter' => $paramid } );
	while ( my $row = $rs_in_use->next ) {
		$in_use{ $row->cachegroup->id } = 1;
	}

	# Add remaining cachegroup ids to @data
	my $rs_links = $self->db->resultset("Cachegroup")->search( undef, { order_by => "name" } );
	while ( my $row = $rs_links->next ) {
		if ( !defined( $in_use{ $row->id } ) ) {
			push( @data, { "id" => $row->id, "name" => $row->name } );
		}
	}

	$self->success( { cachegroups => \@data } );
}

sub available_for_parameter {
	my $self = shift;
	my @data;
	my $paramid = $self->param('paramid');
	my %dsids;
	my %in_use;

	# Get a list of all profile id's associated with this param id
	my $rs_in_use = $self->db->resultset("CachegroupParameter")->search( { 'parameter' => $paramid } );
	while ( my $row = $rs_in_use->next ) {
		$in_use{ $row->cachegroup->id } = 1;
	}

	# Add remaining cachegroup ids to @data
	my $rs_links = $self->db->resultset("Cachegroup")->search( undef, { order_by => "name" } );
	while ( my $row = $rs_links->next ) {
		if ( !defined( $in_use{ $row->id } ) ) {
			push( @data, { "id" => $row->id, "name" => $row->name } );
		}
	}

	$self->success( \@data );
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

sub create{
    my $self = shift;
    my $params = $self->req->json;
    if (!defined($params)) {
        return $self->alert("parameters must be in JSON format,  please check!"); 
    }

    if ( !&is_oper($self) ) {
        return $self->alert("You must be an ADMIN or OPER to perform this operation!");
    }

    my $cachegroups = $self->get_cachegroups();
    my $name    = $params->{name};
    my $short_name    = $params->{short_name};
    my $parent_cachegroup = $params->{parent_cachegroup};
    my $secondary_parent_cachegroup = $params->{secondary_parent_cachegroup};
    my $type_name = $params->{type_name};
    my $type_id = $self->get_typeId($type_name);

    if (!defined($type_id)) {
        return $self->alert("Type ". $type_name . " is not a valid Cache Group type"); 
    }
    if (exists $cachegroups->{'cachegroups'}->{$name}) {
        return $self->internal_server_error("cache_group_name[".$name."] already exists.");
    }
    if (exists $cachegroups->{'short_names'}->{$short_name}) {
        return $self->internal_server_error("cache_group_shortname[".$short_name."] already exists.");
    }

    my $parent_cachegroup_id = $cachegroups->{'cachegroups'}->{$parent_cachegroup};
    $self->app->log->debug("parent_cachegroup[". $parent_cachegroup . "]");
    if ( $parent_cachegroup ne ""  && !defined($parent_cachegroup_id) ) {
        return $self->alert("parent_cachegroup ". $parent_cachegroup . " does not exist."); 
    }
    my $secondary_parent_cachegroup_id = $cachegroups->{'cachegroups'}->{$secondary_parent_cachegroup};
    if ( $secondary_parent_cachegroup ne ""  && !defined($secondary_parent_cachegroup_id) ) {
        return $self->alert("secondary_parent_cachegroup ". $secondary_parent_cachegroup . " does not exist."); 
    }
    my $insert = $self->db->resultset('Cachegroup')->create(
        {
            name        => $name,
            short_name  => $short_name,
            latitude    => $params->{latitude},
            longitude  => $params->{longitude},
            parent_cachegroup_id => $parent_cachegroup_id,
            secondary_parent_cachegroup_id => $secondary_parent_cachegroup_id,
            type        => $type_id,
        }
    );
    $insert->insert();
   
    my $response;
    my $rs = $self->db->resultset('Cachegroup')->find( { id => $insert->id } );
    if (defined($rs)) {
        $response->{id}     = $rs->id;
        $response->{name}   = $rs->name;
        $response->{short_name}  = $rs->short_name;
        $response->{latitude}    = $rs->latitude;
        $response->{longitude}   = $rs->longitude;
        $response->{parent_cachegroup} = $parent_cachegroup;
        $response->{parent_cachegroup_id} = $rs->parent_cachegroup_id;
        $response->{secondary_parent_cachegroup} = $secondary_parent_cachegroup;
        $response->{secondary_parent_cachegroup_id} = $rs->secondary_parent_cachegroup_id;
        $response->{type}        = $rs->type->id;
        $response->{last_updated} = $rs->last_updated;
    }
    return $self->success($response);
}

sub get_typeId {
    my $self      = shift;
    my $type_name = shift;

    my $rs = $self->db->resultset("Type")->find( { name => $type_name } );
    my $type_id;
    if (defined($rs) && ($rs->use_in_table eq "cachegroup")) {
        $type_id = $rs->id;
    }
    return($type_id);
}

sub postupdatequeue {
    my $self       = shift;
    my $params = $self->req->json;
    if ( !&is_oper($self) ) {
        return $self->forbidden("Forbidden. Insufficent privileges.");
    }

    my $name;
    my $id   = $self->param('id');
    $name = $self->db->resultset('Cachegroup')->search( { id => $id } )->get_column('name')->single();

    if (! defined($name)) {
        return $self->alert("cachegroup id[".$id."] does not exist.");
    }

    my $cdn = $params->{cdn};
    my $cdn_id = $self->db->resultset('Cdn')->search( { name => $cdn })->get_column('id')->single();
    if( !defined($cdn_id) ) {
        return $self->alert("cdn " . $cdn . " does not exist.");
    }

    my $setqueue = $params->{action};
    if ( !defined($setqueue)) {
        return $self->alert("action needed, should be queue or dequeue.");
    }
    if ( $setqueue eq "queue") {
        $setqueue = 1
    } elsif	($setqueue eq "dequeue") {
        $setqueue = 0
    } else {
        return $self->alert("action should be queue or dequeue.");
    }

    my @profiles;
    @profiles = $self->db->resultset('Server')->search(
        { 'cdn.name' => $cdn },
        {
            prefetch => 'cdn',
            select   => 'me.profile',
            distinct => 1
        }
    )->get_column('profile')->all();
    my $update = $self->db->resultset('Server')->search(
        {
            -and => [
                cachegroup => $id ,
                profile    => { -in => \@profiles }
            ]
        }
    );

    my $response;
    my @svrs = ();
    if ( $update->count() > 0 ) {
        $update->update( { upd_pending => $setqueue } );
        my @row = $update->get_column('host_name')->all();
        foreach my $svr ( @row ){
            push( @svrs, $svr);
        }
    }

    $response->{serverNames} = \@svrs;
    $response->{action}  = ($setqueue==1)?"queue":"dequeue";
    $response->{cdn}  = $cdn;
    $response->{cachegroupName}  = $name;
    $response->{cachegroupId} = $id;
    return $self->success($response);
}

1;
