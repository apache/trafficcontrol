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

    if (&is_admin($self) || &is_steering($self)) {
        my $data = $self->find_steering();
        return $self->success($data);
    }

    return $self->render(json => {"message" => "unauthorized"}, status => 401);
}

sub find_steering {
    my $self = shift;
    my $steering_filter  = $self->param('xml_id');

    my %steering;

    my $rs_data = $self->db->resultset('SteeringView')->search({}, {order_by => ['steering_xml_id', 'target_xml_id']});

    while ( my $row = $rs_data->next ) {
        if ($steering_filter && $row->steering_xml_id ne $steering_filter) {
            next;
        }

        if (!&is_admin($self)) {
            my $name = $self->current_user()->{username};
            my $user_id = $self->db->resultset('TmUser')->search( { username => $name}, {columns => 'id'} )->single->id;
            my $dsu_row = $self->db->resultset('DeliveryserviceTmuser')->search({tm_user_id => $user_id, deliveryservice => $row->steering_id})->single;

            if (!$dsu_row) {
                next;
            }
        }

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

    if ($steering_filter) {
        return (values %steering)[0];
    }

    my $response = [];
    foreach my $key (sort (keys(%steering))) {
        push(@{$response}, $steering{$key});
    }

    return $response;
}

sub add() {
    my $self = shift;

    if (!&is_admin($self)) {
        return $self->render(json => {"message" => "unauthorized"}, status => 401);
    }

    my $steering_xml_id = $self->req->json->{'deliveryService'};

    if (!$steering_xml_id || !$self->req->json->{'targets'}) {
        return $self->render(json => {"message" => "bad request"}, status => 400);
    }

    if (!(ref($self->req->json->{'targets'}) eq "ARRAY")) {
        return $self->render(json => {"message" => "bad request"}, status => 400);
    }

    my $target_xml_ids = [];
    foreach my $target (@{$self->req->json->{'targets'}}) {
        if (!(ref($target) eq "HASH") || !$target->{'deliveryService'}) {
            return $self->render(json => {"message" => "bad request"}, status => 400);
        }
        push(@{$target_xml_ids}, $target->{'deliveryService'});
    }

    my $ds = $self->db->resultset('Deliveryservice')->search( { xml_id => $steering_xml_id } )->get_column('id')->single();

    if (!$ds) {
        return $self->render(json => {}, status => 409);
    }

    my $rows = [];

    foreach my $xml_id (@{$target_xml_ids}) {
        my $target_ds = $self->db->resultset('Deliveryservice')->search({xml_id => $xml_id})->get_column('id')->single();

        if (!$target_ds) {
            return $self->render(json => {}, status => 409);
        }

        push(@{$rows}, [$ds, $target_ds])
    }

    my $transaction_guard = $self->db->txn_scope_guard;

    for my $row (@{$rows}) {
        $self->db->resultset('SteeringTarget')->create({deliveryservice => $row->[0], target => $row->[1], weight => 0});
    }

    $transaction_guard->commit;

    $self->res->headers->header('Location', '/internal/api/1.2/steering/' . $steering_xml_id . ".json");
    return $self->render(json => {}, status => 201);
}

1;