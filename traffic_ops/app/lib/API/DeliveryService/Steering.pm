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

    my %criteria = length $steering_xml_id ? (steering_xml_id => $steering_xml_id) : ();
    my $rs_data = $self->db->resultset('SteeringView')->search(\%criteria, {order_by => ['steering_xml_id', 'target_xml_id']});

    while ( my $row = $rs_data->next ) {

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

        if ( $row->type eq "STEERING_ORDER" ) {
            push(@{$targets},{
            'deliveryService' => $row->target_xml_id,
            'order' => $row->value,
            'weight'  => 0
            });
        }
        elsif ( $row->type eq "STEERING_WEIGHT" ) {
            push(@{$targets},{
            'deliveryService' => $row->target_xml_id,
            'order' => 0,
            'weight'  => $row->value
            });
        }
        elsif ( $row->type eq "STEERING_GEO_ORDER" ) {
            my $coords = get_primary_origin_coordinates($self, $row->target_id);
            push(@{$targets},{
            'deliveryService' => $row->target_xml_id,
            'order' => 0,
            'geoOrder' => $row->value,
            'latitude' => $coords->{lat},
            'longitude' => $coords->{lon},
            'weight'  => 0
            });
        }
        elsif ( $row->type eq "STEERING_GEO_WEIGHT" ) {
            my $coords = get_primary_origin_coordinates($self, $row->target_id);
            push(@{$targets},{
            'deliveryService' => $row->target_xml_id,
            'order' => 0,
            'geoOrder' => 0,
            'latitude' => $coords->{lat},
            'longitude' => $coords->{lon},
            'weight'  => $row->value
            });
        }

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

sub get_primary_origin_coordinates {
    my $self = shift;
    my $ds_id = shift;

    my %coordinates = (lat => 0.0, lon => 0.0);

    my $origin_rs = $self->db->resultset('Origin')->find(
        { deliveryservice => $ds_id, is_primary => 1 },
        { prefetch => 'coordinate' });

    if ( !defined($origin_rs) || !defined($origin_rs->coordinate) ) {
        return \%coordinates;
    }

    $coordinates{lat} = $origin_rs->coordinate->latitude + 0.0;
    $coordinates{lon} = $origin_rs->coordinate->longitude + 0.0;

    return \%coordinates;
}


sub get_ds_id {
    my $self = shift;
    my $xml_id = shift;

    return $self->db->resultset('Deliveryservice')->search( { xml_id => $xml_id } )->get_column('id')->single();
}

sub get_type {
    my $self = shift;
    my $id = shift;
    my $type;
    
    if ( $id =~ /^\d+$/ ) {
        $type = $self->db->resultset('Type')->search( { id => $id } )->get_column('name')->single();
    }
    else {
        $type = $self->db->resultset('Type')->search( { name => $id } )->get_column('id')->single();
    }
    return $type;
}

# NOTE: STEERING_GEO* types are deliberately ignored in the following endpoint b/c it's soon to be deprecated (use
# the non-internal PUT endpoint instead)
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

    if (!$dsu_row && !&is_admin($self) ) {
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
        if (!$req_target->{'deliveryService'} && ( !$req_target->{'weight'} || !$req_target->{'order'} ) || ( $req_target->{'weight'} && $req_target->{'order'} ) ) {
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
            $steering_target_row->value($req_target->{weight});
            $steering_target_row->type($self->get_type("STEERING_WEIGHT"));
            $steering_target_row->update;
        }
        elsif ($req_target->{'order'}) {
            my $steering_target_row = $self->db->resultset('SteeringTarget')->find({ deliveryservice => $steering_id, target => $target_id});
            $steering_target_row->value($req_target->{order});
            $steering_target_row->type($self->get_type("STEERING_ORDER"));
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
