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
use JSON;
use MojoPlugins::Response;

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

sub create{
    my $self = shift;
    my $region_name = $self->param('region_name');
    my $params = $self->req->json;
    if (!defined($params)) {
        return $self->alert("parameters must be in JSON format,  please check!");
    }
    if ( !&is_oper($self) ) {
        return $self->alert("You must be an ADMIN or OPER to perform this operation!");
    }

    my $existing_physlocation = $self->db->resultset('PhysLocation')->search( { name => $params->{name} } )->get_column('name')->single();
    if (defined($existing_physlocation)){
        return $self->alert("physical location[". $params->{name} . "] already exists.");
    }
    $existing_physlocation = $self->db->resultset('PhysLocation')->search( { name => $params->{shortName} } )->get_column('name')->single();
    if (defined($existing_physlocation)){
        return $self->alert("physical location with shortName[". $params->{shortName} . "] already exists.");
    }
    my $region_id = $self->db->resultset('Region')->search( { name => $region_name } )->get_column('id')->single();
    if (!defined($region_id)) {
        return $self->alert("region[". $region_name . "] does not exist.");
    }

    my $insert = $self->db->resultset('PhysLocation')->create(
        {
            name     => $params->{name},
            short_name     => $params->{shortName},
            region     => $region_id,
            address     => $self->undef_to_default($params->{address}, ""),
            city     => $self->undef_to_default($params->{city}, ""),
            state     => $self->undef_to_default($params->{state}, ""),
            zip     => $self->undef_to_default($params->{zip}, ""),
            phone     => $self->undef_to_default($params->{phone}, ""),
            poc     => $self->undef_to_default($params->{poc}, ""),
            email     => $self->undef_to_default($params->{email}, ""),
            comments  => $self->undef_to_default($params->{comments}, ""),
        } );
    $insert->insert();

    my $response;
    my $rs = $self->db->resultset('PhysLocation')->find( { id => $insert->id } );
    if (defined($rs)) {
        $response->{id}     = $rs->id;
        $response->{name}   = $rs->name;
        $response->{shortName}   = $rs->short_name;
        $response->{regionName}   = $region_name;
        $response->{regionId}   = $rs->region->id;
        $response->{address}   = $rs->address;
        $response->{city}   = $rs->city;
        $response->{state}   = $rs->state;
        $response->{zip}   = $rs->zip;
        $response->{phone}   = $rs->phone;
        $response->{poc}   = $rs->poc;
        $response->{email}   = $rs->email;
        $response->{comments}   = $rs->comments;
        return $self->success($response);
    }
    return $self->alert("create region ". $params->{name}." failed.");
}

sub undef_to_default {
    my $self    = shift;
    my $v       = shift;
    my $default = shift;

    return $v || $default;
}

1;
