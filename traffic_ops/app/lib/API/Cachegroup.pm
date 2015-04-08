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

1;
