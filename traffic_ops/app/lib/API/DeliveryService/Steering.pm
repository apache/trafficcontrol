package API::DeliveryService::Steering;
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

use Mojo::Base 'Mojolicious::Controller';
use UI::Utils;
use Data::Dumper;


sub index {
    my $self = shift;

    if (&is_admin($self) || &is_steering($self)) {
        my $data = $self->find_steering($self->param('xml_id'));

        if (!$data) {
            return $self->render(json => {}, status => 404);
        }

        return $self->success($data);
    }

    return $self->render(json => {"message" => "unauthorized"}, status => 401);
}

sub find_steering {
    my $self = shift;
    my $steering_xml_id  = shift;

    my %steering;

    my $rs_data = $self->db->resultset('SteeringView')->search({}, {order_by => ['steering_xml_id', 'target_xml_id']});

    while ( my $row = $rs_data->next ) {
        if ($steering_xml_id && $row->steering_xml_id ne $steering_xml_id) {
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
        

        if (! exists($steering{$row->steering_xml_id})) {
            my $ds = $self->db->resultset('Deliveryservice')->search( { xml_id => $row->steering_xml_id } )->single();
            my $client_steering;
            if ($ds->type->name =~ /CLIENT_STEERING/) {
                $client_steering = '1';
            }
            else { 
                $client_steering = '0';
            }
            $steering{$row->steering_xml_id} = {"deliveryService" => $row->steering_xml_id, "clientSteering" => \$client_steering};
        }

        my $steering_entry = $steering{$row->steering_xml_id};

        if (! exists($steering_entry->{'filters'})) {
            $steering_entry->{'filters'} = []
        }

        my $rs_filters = $self->db->resultset('RegexByDeliveryServiceList')->search({'ds_id' => $target_id, 'type' => "STEERING_REGEXP" }, {order_by =>'pattern'} );
        while (my $r2 = $rs_filters->next) {
            push(@{$steering_entry->{'filters'}}, {deliveryService => $row->target_xml_id, pattern => $r2->pattern});
        }

        if (! exists($steering_entry->{"targets"})) {
            $steering_entry->{"targets"} = [];
        }

        my $targets = $steering_entry->{"targets"};

        push(@{$targets},{
            'deliveryService' => $row->target_xml_id,
            'weight' => $row->weight,
        });

    }


    if ($steering_xml_id) {
        my $steering_response = (values %steering)[0];
        if (!$steering_response) {
            return;
        }

        return $steering_response;
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

    my $ds = $self->get_ds_id($steering_xml_id);

    if (!$ds) {
        return $self->render(json => {}, status => 409);
    }

    my $rows = [];

    foreach my $xml_id (@{$target_xml_ids}) {
        my $target_ds = $self->get_ds_id($xml_id);

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

sub get_ds_id {
    my $self = shift;
    my $xml_id = shift;

    return $self->db->resultset('Deliveryservice')->search( { xml_id => $xml_id } )->get_column('id')->single();
}

sub update() {
    my $self = shift;

    if (!(&is_admin($self)) && !(&is_steering($self))) {
        return $self->render(json => {"message" => "unauthorized"}, status => 401);
    }

    my $steering_xml_id  = $self->param('xml_id');

    my $rs_data = $self->db->resultset('SteeringView')->search(
        {'steering_xml_id' => $steering_xml_id},
        {order_by => ['steering_xml_id', 'target_xml_id']});

    my $row = $rs_data->next;

    if (!$row) {
        return $self->render(json => {}, status => 409)
    }

    my $steering_id = $row->steering_id;

    my $name = $self->current_user()->{username};
    my $user_id = $self->db->resultset('TmUser')->search( { username => $name}, {columns => 'id'} )->single->id;
    my $dsu_row = $self->db->resultset('DeliveryserviceTmuser')->search(
        {tm_user_id => $user_id, deliveryservice => $row->steering_id})->single;

    if (!$dsu_row) {
        return $self->render(json => {"message" => "unauthorized"}, status => 401);
    }

    if (!$self->req->json->{'targets'} ) {
        return $self->render(json => {"message" => "please provide a valid json including targets"}, status => 400);
    }

    my $valid_targets = {};

    do {
        $valid_targets->{$row->target_xml_id} = $row->target_id;
    } while ($row = $rs_data->next);

    my $req_targets = $self->req->json->{'targets'};

    foreach my $req_target (@{$req_targets}) {
        if (!$req_target->{'deliveryService'} || !$req_target->{'weight'}) {
           return $self->render(json => {"message" => "please provide a valid json for targets"}, status => 400);
        }
        if (!exists($valid_targets->{$req_target->{'deliveryService'}})) {
            return $self->render(json => {} , status => 409);
        }
    }

    my $req_filters = $self->req->json->{'filters'};

    foreach my $req_filter (@{$req_filters}) {
        if (!$req_filter->{'deliveryService'} || !$req_filter->{'pattern'}) {
            return $self->render(json => {"message" => "please provide a valid json for filters"}, status => 400);
        }
        if (!exists($valid_targets->{$req_filter->{'deliveryService'}})) {
            return $self->render(json => {}, status => 409);
        }
    }

    my $steering_regex_type = $self->db->resultset('Type')->find({name => "STEERING_REGEXP"})->id;

    # Start Transaction
    my $transaction_guard = $self->db->txn_scope_guard;

    foreach my $req_target (@{$req_targets}) {
        my $target_id = $valid_targets->{$req_target->{'deliveryService'}};

        if ($req_target->{'weight'}) {
            my $steering_target_row = $self->db->resultset('SteeringTarget')->find({ deliveryservice => $steering_id, target => $target_id});
            $steering_target_row->weight($req_target->{weight});
            $steering_target_row->update;
        }

        if ($self->req->json->{'filters'}) {
            # delete existing filters
            my $dsr_rs =  $self->db->resultset('DeliveryserviceRegex')->search({deliveryservice => $target_id});

            while (my $dsr_row = $dsr_rs->next) {
                $self->db->resultset('Regex')->search({id => $dsr_row->regex->id, type => $steering_regex_type})->delete;
            }
            # add filters for target
            foreach my $filter (@{$req_filters}) {
                my $filter_ds = $self->get_ds_id($filter->{deliveryService});
                if ($filter_ds eq $target_id) {
                    my $regex_row = $self->db->resultset('Regex')->create({pattern => $filter->{pattern}, type => $steering_regex_type});
                    $self->db->resultset('DeliveryserviceRegex')->create({deliveryservice => $target_id, regex => $regex_row->id})
                }
            }
        }
    }

    # Commit and end transaction
    $transaction_guard->commit;

    return $self->success($self->find_steering($steering_xml_id));
}

1;
