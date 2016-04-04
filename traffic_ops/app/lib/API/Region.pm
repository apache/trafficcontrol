package API::Region;
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
use JSON;
use MojoPlugins::Response;

my $finfo = __FILE__ . ":";

sub index {
	my $self = shift;
	my @data;
	my $orderby = $self->param('orderby') || "name";
	my $rs_data = $self->db->resultset("Region")->search( undef, { prefetch => ['division'], order_by => 'me.' . $orderby } );
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"id"   => $row->id,
				"name" => $row->name,
			}
		);
	}
	$self->success( \@data );
}

sub create{
    my $self = shift;
    my $division_name = $self->param('division_name');
    my $params = $self->req->json;
    if (!defined($params)) {
        return $self->alert("parameters must be in JSON format,  please check!"); 
    }
    if ( !&is_oper($self) ) {
        return $self->alert("You must be an ADMIN or OPER to perform this operation!");
    }

    my $existing_region = $self->db->resultset('Region')->search( { name => $params->{name} } )->get_column('name')->single();
    if (defined($existing_region)) {
        return $self->alert("region[". $params->{name} . "] already exists."); 
    }

    my $divsion_id = $self->db->resultset('Division')->search( { name => $division_name } )->get_column('id')->single();
    if (!defined($divsion_id)) {
        return $self->alert("division[". $division_name . "] does not exist."); 
    }

    my $insert = $self->db->resultset('Region')->create(
        {
            name     => $params->{name},
            division => $divsion_id
        } );
    $insert->insert();
   
    my $response;
    my $rs = $self->db->resultset('Region')->find( { id => $insert->id } );
    if (defined($rs)) {
        $response->{id}     = $rs->id;
        $response->{name}   = $rs->name;
        $response->{division_name}  = $division_name;
        $response->{divsion_id}    = $rs->division->id;
        return $self->success($response);
    }
    return $self->alert("create region ". $params->{name}." failed.");
}

1;
