package API::PhysLocation;
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

use UI::Utils;

use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;

my $finfo = __FILE__ . ":";

sub index {
	my $self = shift;
	my @data;
	my $orderby = $self->param('orderby') || "name";
	my $rs_data = $self->db->resultset("PhysLocation")->search( undef, { prefetch => ['region'], order_by => 'me.' . $orderby } );
	while ( my $row = $rs_data->next ) {

		next if $row->short_name eq 'UNDEF';

		push(
			@data, {
				"id"        => $row->id,
				"name"      => $row->name,
				"shortName" => $row->short_name,
				"address"   => $row->address,
				"city"      => $row->city,
				"state"     => $row->state,
				"zip"       => $row->zip,
				"poc"       => $row->poc,
				"phone"     => $row->phone,
				"email"     => $row->email,
				"comments"  => $row->comments,
				"region"    => $row->region->name,
			}
		);
	}
	$self->success( \@data );
}

sub index_trimmed {
	my $self = shift;
	my @data;
	my $orderby = $self->param('orderby') || "name";
	my $rs_data = $self->db->resultset("PhysLocation")->search( undef, { prefetch => ['region'], order_by => 'me.' . $orderby } );
	while ( my $row = $rs_data->next ) {

		next if $row->short_name eq 'UNDEF';

		push(
			@data, {
				"name" => $row->name,
			}
		);
	}
	$self->success( \@data );
}

1;
