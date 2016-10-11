package API::Types;
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

sub index {
	my $self         = shift;
	my $use_in_table = $self->param('useInTable');
	my $orderby      = $self->param('orderby') || "name";

	my @data;
	my %criteria;

	if ( defined $use_in_table ) {
		$criteria{'me.use_in_table'} = $use_in_table;
	}

	my $rs_data = $self->db->resultset("Type")->search( \%criteria, { order_by => $orderby } );

	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"id"          => $row->id,
				"name"        => $row->name,
				"description" => $row->description,
				"useInTable"  => $row->use_in_table,
				"lastUpdated" => $row->last_updated
			}
		);
	}
	$self->success( \@data );
}

sub index_trimmed {
	my $self = shift;
	my @data;
	my $orderby = $self->param('orderby') || "name";
	my $rs_data = $self->db->resultset("Type")->search( undef, { order_by => $orderby } );
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"name" => $row->name
			}
		);
	}
	$self->success( \@data );
}

sub show {
	my $self = shift;
	my $id   = $self->param('id');

	my $rs_data = $self->db->resultset("Type")->search( { id => $id } );
	my @data = ();
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"id"          => $row->id,
				"name"        => $row->name,
				"description" => $row->description,
				"useInTable"  => $row->use_in_table,
				"lastUpdated" => $row->last_updated
			}
		);
	}
	$self->success( \@data );
}

sub update {
	my $self   = shift;
	my $id     = $self->param('id');
	my $params = $self->req->json;

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	my $type = $self->db->resultset('Type')->find( { id => $id } );
	if ( !defined($type) ) {
		return $self->not_found();
	}

	if ( !defined( $params->{name} ) ) {
		return $self->alert("Type name is required.");
	}

	my $values = {
		name 			=> $params->{name},
		description 	=> $params->{description},
		use_in_table 	=> $params->{useInTable}
	};

	my $rs = $type->update($values);
	if ($rs) {
		my $response;
		$response->{id}          = $rs->id;
		$response->{name}        = $rs->name;
		$response->{description} = $rs->description;
		$response->{useInTable} = $rs->description;
		$response->{lastUpdated} = $rs->use_in_table;

		&log( $self, "Updated Type name '" . $rs->name . "' for id: " . $rs->id, "APICHANGE" );

		return $self->success( $response, "Type update was successful." );
	}
	else {
		return $self->alert("Type update failed.");
	}

}


1;
