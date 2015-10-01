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
		{ prefetch => [ { 'cdn' => undef } ], order_by => $orderby } );
	while ( my $row = $rs_data->next ) {
		my $cdn_name;
		if ( defined $row->cdn ) {
			$cdn_name = $row->cdn->name;
		}

		push(
			@data,
			{   "id"          => $row->id,
				"name"        => $row->name,
				"cdnName"     => $cdn_name,
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
