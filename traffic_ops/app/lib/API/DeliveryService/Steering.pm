package API::DeliveryService::Steering;
#
# Copyright 2016 Comcast Cable Communications Management, LLC
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

use Mojo::Base 'Mojolicious::Controller';
use UI::Utils;
use Data::Dumper;

sub index {
    my $self = shift;

    if (&is_admin( $self )) {
        my $data = $self->find_steering();
        return $self->success($data);
    }

    if (&is_steering( $self )) {
        return $self->success("Yo!");
    }

    return $self->render(json => {"message" => "unauthorized"}, status => 401);
}

sub find_steering {
    my $self = shift;
    my %steering;

    my $rs_data = $self->db->resultset('SteeringView')->search({}, {order_by => ['steering_xml_id', 'target_xml_id']});

    while ( my $row = $rs_data->next ) {
        my $target_id = $row->target_id;
        my $rs_filters = $self->db->resultset('RegexByDeliveryServiceList')->search({'ds_id' => $target_id });

        my $filters = [];
        while (my $r2 = $rs_filters->next) {
            push(@{$filters}, $r2->pattern);
        }

        if (! exists($steering{$row->steering_xml_id})) {
            $steering{$row->steering_xml_id} = {"deliveryService" => $row->steering_xml_id};
        }

        my $steering_entry = $steering{$row->steering_xml_id};

        if (! exists($steering_entry->{"targets"})) {
            $steering_entry->{"targets"} = [];
        }

        my $targets = $steering_entry->{"targets"};

        push(@{$targets},{
            'deliveryService' => $row->target_xml_id,
            'weight' => $row->weight,
            'filters' => $filters,
        });
    }

    my $response = [];
    foreach my $key (sort (keys(%steering))) {
        push(@{$response}, $steering{$key});
    }

    return $response;
}

#sub add() {
#    my $self = shift;
#
#
#    if (!&is_admin($self)) {
#        return $self->render(json => {"message" => "unauthorized"}, status => 401);
#    }
#
#    my $xml_id = $self->req->json->{'id'};
#    my $ds_count = $self->db->resultset('Deliveryservice')->count( { xml_id => $xml_id } )->get_column('id');
#
#    print STDERR "count found $ds_count\n";
#
#    my @steered_delivery_services = $self->req->json->{'steeredDeliveryServices'};
#
#    foreach (@steered_delivery_services)
#}

1;