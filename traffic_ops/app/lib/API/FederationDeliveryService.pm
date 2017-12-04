package API::FederationDeliveryService;
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

sub index {
    my $self    = shift;
    my $fed_id  = $self->param('fedId');

    my @data;
    my $rs_data = $self->db->resultset("FederationDeliveryservice")->search( { 'federation' => $fed_id }, { prefetch => [ 'deliveryservice' ] } );
    while ( my $row = $rs_data->next ) {
        push(
            @data, {
                "id"        => $row->deliveryservice->id,
                "cdn"       => $row->deliveryservice->cdn->name,
                "type"      => $row->deliveryservice->type->name,
                "xmlId"     => $row->deliveryservice->xml_id,
            }
        );
    }
    $self->success( \@data );
}

sub assign_dss_to_federation {
    my $self        = shift;
    my $fed_id      = $self->param('fedId');
    my $params      = $self->req->json;
    my $ds_ids      = $params->{dsIds};
    my $replace     = $params->{replace};
    my $count       = 0;

    if ( !&is_admin($self) ) {
        return $self->forbidden();
    }

    my $fed = $self->db->resultset('Federation')->find( { id => $fed_id } );
    if ( !defined($fed) ) {
        return $self->not_found();
    }

    if ( ref($ds_ids) ne 'ARRAY' ) {
        return $self->alert("Delivery Service IDs must be an array");
    }

    if ( $replace ) {
        if (!scalar @{ $ds_ids }) {
            return $self->alert("A federation must have at least one delivery service assigned");
        }
        # start fresh and delete existing fed/ds associations if replace=true
        my $delete = $self->db->resultset('FederationDeliveryservice')->search( { federation => $fed_id } );
        $delete->delete();
    }

    my @values = ( [ qw( federation deliveryservice ) ]); # column names are required for 'populate' function

    foreach my $ds_id (@{ $ds_ids }) {
        push(@values, [ $fed_id, $ds_id ]);
        $count++;
    }

    $self->db->resultset("FederationDeliveryservice")->populate(\@values);

    my $msg = $count . " delivery service(s) were assigned to the " . $fed->cname . " federation";
    &log( $self, $msg, "APICHANGE" );

    my $response = $params;
    return $self->success($response, $msg);
}

sub delete {
    my $self        = shift;
    my $fed_id      = $self->param('fedId');
    my $ds_id       = $self->param('dsId');

    if ( !&is_admin($self) ) {
        return $self->forbidden();
    }

    my $fed_dss = $self->db->resultset("FederationDeliveryservice")->search( { 'federation' => $fed_id } );
    if ( $fed_dss->count() < 2 ) {
        return $self->alert("A federation must have at least one delivery service assigned");
    }

    my $fed_ds = $self->db->resultset("FederationDeliveryservice")->search( { 'federation.id' => $fed_id, 'deliveryservice' => $ds_id }, { prefetch => [ 'federation', 'deliveryservice' ] } );
    if ( !defined($fed_ds) ) {
        return $self->not_found();
    }

    my $row = $fed_ds->next;
    my $rs = $fed_ds->delete();
    if ($rs) {
        my $msg = "Removed delivery service [ " . $row->deliveryservice->xml_id . " ] from federation [ " . $row->federation->cname . " ]";
        &log( $self, $msg, "APICHANGE" );
        return $self->success_message($msg);
    }

    return $self->alert( "Failed to remove delivery service from federation." );
}

1;
