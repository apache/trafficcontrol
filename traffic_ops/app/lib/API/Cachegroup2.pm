package API::Cachegroup2;
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

sub get_cachegroup_by_id {
    my $self = shift;
    my $id = shift;
    my $row;

    eval {
        $row = $self->db->resultset('Cachegroup')->find( { id => $id }, { prefetch => [ { 'type' => undef, } ]});
    };
    if ($@) {
        $self->app->log->error( "Failed to get cachegroup id = $id: $@" );
        return (undef, "Failed to get cachegroup id = $id: $@")
    }

    my $r;
    eval {
        $r = $self->db->resultset('Cachegroup')->find( { id => $row->parent_cachegroup_id } );
    };
    if ($@) {
        $self->app->log->error( "Failed to get cachegroup id = $id: $@" );
        return (undef, "Failed to get cachegroup id = $id: $@")
    }
    my $parentCachegroup = defined($r) ? $r->name : "";
    eval {
        $r = $self->db->resultset('Cachegroup')->find( { id => $row->secondary_parent_cachegroup_id } );
    };
    if ($@) {
        $self->app->log->error( "Failed to get cachegroup id = $id: $@" );
        return (undef, "Failed to get cachegroup id = $id: $@")
    }
    my $secondaryParentCachegroup = defined($r) ? $r->name : "";

    my $data = {
        "id"     => $row->id,
        "name"   => $row->name,
        "shortName"  => $row->short_name,
        "latitude"    => $row->latitude,
        "longitude"   => $row->longitude,
        "parentCachegroup" => $parentCachegroup,
        "parentCachegroupId" => $row->parent_cachegroup_id,
        "secondaryParentCachegroup" => $secondaryParentCachegroup,
        "secondaryParentCachegroupId" => $row->secondary_parent_cachegroup_id,
        "typeName"        => $row->type->name,
        "lastUpdated" => $row->last_updated,
    };
    return ($data, undef);
}

sub isValidCachegroup {
    my $self = shift;
    my $params = shift;
    my %errFields = ();

    if (!defined($params)) {
        return "parameters must be in JSON format,  please check!";
    }

    if (!defined($params->{'name'})) {
        $errFields{'name'} = 'is required';
    }
    if (!defined($params->{'shortName'})) {
        $errFields{'shorName'} = 'is required';
    }
    if (!defined($params->{typeName})){
        $errFields{'typeName'} = 'is required';
    }
    if (%errFields) {
        return \%errFields;
    }

    if (!($params->{'name'} =~ /^[0-9a-zA-Z_\.\-]+$/)) {
        return "Invalid name. Use alphanumeric . or _ .";
    }
    if (!($params->{'shortName'} =~ /^[0-9a-zA-Z_\.\-]+$/)) {
        return "Invalid shortName. Use alphanumeric . or _ .";
    }
    my $typeName = $params->{typeName};
    my $type_id = $self->get_typeId($typeName);
    if (!defined($type_id)) {
        return "Type ". $typeName . " is not a valid Cache Group type";
    }
    if (defined($params->{'latitude'})) {
        if(!($params->{'latitude'} =~ /^[-]*[0-9]+[.]*[0-9]*/)) {
            return "Invalid latitude entered. Must be a float number.";
        }
        if ( abs $params->{'latitude'} > 90 ) {
            return "Invalid latitude entered. May not exceed +- 90.0.";
        }
    }
    if (defined($params->{'longitude'})) {
        if(!($params->{'longitude'} =~ /^[-]*[0-9]+[.]*[0-9]*/)) {
            return "Invalid longitude entered. Must be a float number.";
        }
        if ( abs $params->{'longitude'} > 180 ) {
            return "Invalid longitude entered. May not exceed +- 180.0.";
        }
    }

    return undef;
}

sub create{
    my $self = shift;
    my $params = $self->req->json;
    if ( !&is_oper($self) ) {
        return $self->forbidden();
    }

    my $err = $self->isValidCachegroup($params);
    if (defined($err)) {
        return $self->alert($err);
    }

    my $cachegroups = $self->get_cachegroups();
    my $name    = $params->{name};
    my $shortName    = $params->{shortName};
    my $parentCachegroup = $params->{parentCachegroup};
    my $secondaryParentCachegroup = $params->{secondaryParentCachegroup};
    my $typeName = $params->{typeName};
    my $type_id = $self->get_typeId($typeName);

    if (exists $cachegroups->{'cachegroups'}->{$name}) {
        return $self->internal_server_error("cache_group_name[".$name."] already exists.");
    }
    if (exists $cachegroups->{'short_names'}->{$shortName}) {
        return $self->internal_server_error("cache_group_shortname[".$shortName."] already exists.");
    }

    my $parentCachegroupId = $cachegroups->{'cachegroups'}->{$parentCachegroup};
    $self->app->log->debug("parentCachegroup[". $parentCachegroup . "]");
    if ( $parentCachegroup ne ""  && !defined($parentCachegroupId) ) {
        return $self->alert("parentCachegroup ". $parentCachegroup . " does not exist.");
    }
    my $secondaryParentCachegroupId = $cachegroups->{'cachegroups'}->{$secondaryParentCachegroup};
    if ( $secondaryParentCachegroup ne ""  && !defined($secondaryParentCachegroupId) ) {
        return $self->alert("secondaryParentCachegroup ". $secondaryParentCachegroup . " does not exist.");
    }
    my $insert = $self->db->resultset('Cachegroup')->create(
        {
            name        => $name,
            short_name  => $shortName,
            latitude    => $params->{latitude},
            longitude  => $params->{longitude},
            parent_cachegroup_id => $parentCachegroupId,
            secondary_parent_cachegroup_id => $secondaryParentCachegroupId,
            type        => $type_id,
        }
    );
    $insert->insert();

    &log( $self, "Create cachegroup with name:" . $name, "APICHANGE" );

    my ($response, $err1) = $self->get_cachegroup_by_id($insert->id);
    if( defined($err1) ) {
        return $self->alert(
            { Error => $err1 }
        );
    }
    return $self->success($response, "Cachegroup successfully created: " . $name);
}

sub update{
    my $self = shift;
    my $params = $self->req->json;
    if ( !&is_oper($self) ) {
        return $self->forbidden();
    }

    my $err = $self->isValidCachegroup($params);
    if (defined($err)) {
        return $self->alert($err);
    }

    my $id = $self->param('id');
    my $update = $self->db->resultset('Cachegroup')->find( { id => $id } );
    if( !defined($update) ) {
        return $self->not_found();
    }

    my $type_id = undef;
    if (defined($params->{typeName})){
        my $typeName = $params->{typeName};
        $type_id = $self->get_typeId($typeName);
    }

    my $cachegroups = $self->get_cachegroups();
    my $parentCachegroupId = undef;
    if (defined($params->{parentCachegroup})){
        my $parentCachegroup = $params->{parentCachegroup};
        $parentCachegroupId = $cachegroups->{'cachegroups'}->{$parentCachegroup};
        if ( $parentCachegroup ne ""  && !defined($parentCachegroupId) ) {
            return $self->alert("parentCachegroup ". $parentCachegroup . " does not exist.");
        }
        if (defined($parentCachegroupId) && $parentCachegroupId == $id) {
            return $self->alert("Could not set the Cache Group itself as parent.");
        }
    }
    my $secondaryParentCachegroupId = undef;
    if (defined($params->{secondaryParentCachegroup})){
        my $secondaryParentCachegroup = $params->{secondaryParentCachegroup};
        $secondaryParentCachegroupId = $cachegroups->{'cachegroups'}->{$secondaryParentCachegroup};
        if ( $secondaryParentCachegroup ne ""  && !defined($secondaryParentCachegroupId) ) {
            return $self->alert("secondaryParentCachegroup ". $secondaryParentCachegroup . " does not exist.");
        }
        if (defined($secondaryParentCachegroupId) && $secondaryParentCachegroupId == $id) {
            return $self->alert("Could not set the Cache Group itself as secondary parent.");
        }
    }

    eval {
        $update->update(
            {
                name        => defined($params->{'name'}) ? $params->{'name'} : $update->name,
                short_name  => defined($params->{'shortName'}) ? $params->{'shortName'} : $update->short_name,
                latitude    => defined($params->{'latitude'}) ? $params->{'latitude'} : $update->latitude,
                longitude   => defined($params->{'longitude'}) ? $params->{'longitude'} : $update->longitude,
                parent_cachegroup_id            => defined($params->{parentCachegroup}) ? $parentCachegroupId : $update->parent_cachegroup_id,
                secondary_parent_cachegroup_id  => defined($params->{secondaryParentCachegroup}) ? $secondaryParentCachegroupId : $update->secondary_parent_cachegroup_id,
                type        => defined($type_id) ? $type_id : $update->type,
            }
        );
    };
    if ($@) {
        $self->app->log->error( "Failed to update cachegroup id = $id: $@" );
        return $self->alert(
            { Error => "Failed to update server: $@" }
        );
    }
    $update->update();

    &log( $self, "Update cachegroup with name:" . $update->name, "APICHANGE" );

    my ($response, $err1) = $self->get_cachegroup_by_id($id);
    if( defined($err1) ) {
        return $self->alert(
            { Error => $err1 }
        );
    }
    return $self->success($response, "Cachegroup was updated: " . $update->name);
}

sub delete{
    my $self = shift;
    my $rs;
    if ( !&is_oper($self) ) {
        return $self->forbidden();
    }

    my $id = $self->param('id');
    my $cg = $self->db->resultset('Cachegroup')->find( { id => $id } );
    if ( !defined($cg) ) {
        return $self->not_found();
    }
    $rs = $self->db->resultset('Cachegroup')->search( { parent_cachegroup_id => $id } );
    if ($rs->count() > 0) {
        $self->app->log->error( "Failed to delete cachegroup id = $id, which has children" );
        return $self->alert("Failed to delete cachegroup id = $id, which has children");
    }
    $rs = $self->db->resultset('Cachegroup')->search( { secondary_parent_cachegroup_id => $id } );
    if ($rs->count() > 0) {
        $self->app->log->error( "Failed to delete cachegroup id = $id, which has children" );
        return $self->alert("Failed to delete cachegroup id = $id, which has children");
    }
    $rs = $self->db->resultset('Server')->search( { cachegroup => $id } );
    if ($rs->count() > 0) {
        $self->app->log->error( "Failed to delete cachegroup id = $id has servers" );
        return $self->alert("Failed to delete cachegroup id = $id has servers");
    }
    my $delete = $self->db->resultset('Cachegroup')->search( { id => $id } );
    my $name = $delete->get_column('name')->single();
    $delete->delete();

    &log( $self, "Delete cachegroup " . $name, "APICHANGE" );

    return $self->success_message("Cachegroup was deleted: ". $name);
}

sub get_typeId {
    my $self      = shift;
    my $typeName = shift;

    my $rs = $self->db->resultset("Type")->find( { name => $typeName } );
    my $type_id;
    if (defined($rs) && ($rs->use_in_table eq "cachegroup")) {
        $type_id = $rs->id;
    }
    return($type_id);
}

1;
