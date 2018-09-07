package API::FederationResolver;
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

use UI::Utils;

use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use JSON;
use Validate::Tiny ':all';
use Net::CIDR;
use Data::Validate::IP qw(is_ipv4 is_ipv6);

sub index {
    my $self = shift;
    my @data;
    my $rs_data = $self->db->resultset("FederationResolver")->search( undef, { prefetch => ['type'] } );
    while ( my $row = $rs_data->next ) {
        push(
            @data, {
                "id"        => $row->id,
                "ipAddress" => $row->ip_address,
                "type"      => $row->type->name
            }
        );
    }
    $self->success( \@data );
}

sub create {
    my $self        = shift;
    my $params      = $self->req->json;

    if ( !&is_admin($self) ) {
        return $self->forbidden();
    }

    my ( $is_valid, $result ) = $self->is_federation_resolver_valid($params);

    if ( !$is_valid ) {
        return $self->alert($result);
    }

    my $existing = $self->db->resultset('FederationResolver')->find( { ip_address => $params->{ipAddress} } );
    if ( $existing ) {
        return $self->alert("$params->{ipAddress} already in use");
    }

    my $values = {
        ip_address  => $params->{ipAddress},
        type        => $params->{typeId},
    };

    my $insert = $self->db->resultset('FederationResolver')->create($values);
    my $rs = $insert->insert();
    if ($rs) {
        my $response;
        $response->{id}         = $rs->id;
        $response->{ipAddress}  = $rs->ip_address;
        $response->{typeId}     = $rs->type->id;

        my $msg = "Federation Resolver created [ IP = " . $rs->ip_address . " ] with id: " . $rs->id;
        &log( $self, $msg, "APICHANGE" );
        return $self->success( $response, $msg );
    } else {
        return $self->alert("Federation Resolver creation failed.");
    }
}

sub delete {
    my $self	    = shift;
    my $fed_res_id	= $self->param('id');

    if ( !&is_admin($self) ) {
        return $self->forbidden();
    }

    my $fed_res = $self->db->resultset('FederationResolver')->find( { id => $fed_res_id } );
    if ( !defined($fed_res) ) {
        return $self->not_found();
    }

    my $rs = $fed_res->delete();
    if ($rs) {
        my $msg = "Federation resolver deleted [ IP = " . $rs->ip_address . " ] with id: " . $rs->id;
        &log( $self, $msg, "APICHANGE" );
        return $self->success_message($msg);
    } else {
        return $self->alert( "Federation resolver delete failed." );
    }
}

sub is_federation_resolver_valid {
    my $self   = shift;
    my $params = shift;

    my $rules = {
        fields => [ qw/ipAddress typeId/ ],

        # Validation checks to perform
        checks => [
            ipAddress   => [ \&is_valid_ip ],
            typeId      => [ is_required("is required") ],
        ]

    };

    # Validate the input against the rules
    my $result = validate( $params, $rules );

    if ( $result->{success} ) {
        return ( 1, $result->{data} );
    }
    else {
        return ( 0, $result->{error} );
    }
}

sub is_valid_ip {
    my ( $value, $params ) = @_;

    if (!is_ipv4($value) && !is_ipv6($value) && !Net::CIDR::cidr2range($value)) {
        return "invalid. $value is not a valid ip address or range.";
    }

    return undef;
}

1;
