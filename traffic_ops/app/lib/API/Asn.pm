package API::Asn;
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

# Index
sub index {
	my $self = shift;
	my @data;
	my $orderby = $self->param('orderby') || "asn";
	my $rs_data = $self->db->resultset("Asn")->search( undef, { prefetch => [ { 'cachegroup' => undef } ], order_by => "me." . $orderby } );
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"id"          => $row->id,
				"asn"         => $row->asn,
				"cachegroup"  => $row->cachegroup->name,
				"lastUpdated" => $row->last_updated,
			}
		);
	}
	$self->success( \@data );
}

# Show
sub show {
	my $self = shift;
	my $id   = $self->param('id');

	my $rs_data = $self->db->resultset("Asn")->search( { 'me.id' => $id }, { prefetch => [ 'cachegroup' ] } );
	my @data = ();
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"id"          => $row->id,
				"asn"         => $row->asn,
				"lastUpdated" => $row->last_updated,
				"cachegroup"  => {
					"id" => $row->cachegroup->id,
					"name" => $row->cachegroup->name
				}
			}
		);
	}
	$self->success( \@data );
}

# Index
sub v11_index {
	my $self = shift;
	my @data;
	my $orderby = $self->param('orderby') || "asn";
	my $rs_data = $self->db->resultset("Asn")->search( undef, { prefetch => [ { 'cachegroup' => undef } ], order_by => "me." . $orderby } );
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"id"          => $row->id,
				"asn"         => $row->asn,
				"cachegroup"  => $row->cachegroup->name,
				"lastUpdated" => $row->last_updated,
			}
		);
	}
	$self->success( { "asns" => \@data } );
}
1;
