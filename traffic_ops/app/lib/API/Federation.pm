package API::Federation;
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

use Data::Validate::IP qw(is_ipv4 is_ipv6);

sub index {
  my $self = shift;
  my $data;

  my $rs_data = $self->db->resultset('FederationDeliveryservice')
    ->search( {}, { prefetch => [ 'federation', 'deliveryservice' ] } );

  while ( my $row = $rs_data->next ) {
    my $xml_id = $row->deliveryservice->xml_id;
    my $m;

    $m->{'cname'} = $row->federation->cname;
    $m->{'ttl'}   = $row->federation->ttl;

    my $id        = $row->federation->id;
    my @resolvers = $self->db->resultset('FederationResolver')->search(
      { 'federation_federation_resolvers.federation' => $id },
      { prefetch => 'federation_federation_resolvers' }
    )->all();

    for my $r (@resolvers) {
      my $type = lc $r->type->name;
      if ( defined $m->{$type} ) {
        push( $m->{$type}, $r->ip_address );
      }
      else {
        @{ $m->{$type} } = ();
        push( $m->{$type}, $r->ip_address );
      }
    }

    if ( defined $data ) {
      my $x = $self->find_delivery_service( $xml_id, $data );
      if ( !defined $x ) {
        $data = $self->add_delivery_service( $xml_id, $m, $data );
      }
      else {
        my @map = $x->{'mappings'};
        push( @map, $m );
        $x->{'mappings'} = \@map;
      }

    }
    else {
      $data = $self->add_delivery_service( $xml_id, $m, $data );
    }

  }
  $self->success( $data );
}

sub find_delivery_service {
  my $self   = shift;
  my $xml_id = shift;
  my $data   = shift;
  my $hash;

  foreach my $key (@{$data}) {
    if ( $key->{'deliveryService'} eq $xml_id ) {
      $hash = $key;
    }
  }
  return ($hash);
}

sub add_delivery_service {
  my $self   = shift;
  my $xml_id = shift;
  my $m      = shift;
  my $data   = shift;

  push( my @map, $m );

  push(
    @${data},
    {   "deliveryService" => $xml_id,
      "mappings"        => \@map
    }
  );
  return $data;
}


1;
