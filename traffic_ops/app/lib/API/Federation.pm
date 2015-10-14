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
use Net::CIDR;
use JSON;
use Validate::Tiny ':all';
use Data::Validate::IP qw(is_ipv4 is_ipv6);

sub find_tmuser {
  my $self             = shift;
  my $current_username = shift;

  my $tm_user
    = $self->db->resultset('TmUser')
    ->search(
    { username => $current_username, 'role.name' => 'federation' },
    { prefetch => 'role' } )->single();

  return $tm_user;
}

sub find_federations {
  my $self           = shift;
  my $federation_ids = shift;
  my $rs_data;

  if ($federation_ids) {
    $rs_data = $self->db->resultset('FederationDeliveryservice')->search(
      { federation => { -in => $federation_ids } },
      {   prefetch => [ 'federation', 'deliveryservice' ],
        order_by => "deliveryservice.xml_id"
      }
    );
  }
  else {
    $rs_data = $self->db->resultset('FederationDeliveryservice')->search(
      {},
      {   prefetch => [ 'federation', 'deliveryservice' ],
        order_by => "deliveryservice.xml_id"
      }
    );
  }
  return $rs_data;
}

sub find_federation_resolvers {
  my $self          = shift;
  my $federation_id = shift;

  my @resolvers
    = $self->db->resultset('FederationResolver')
    ->search(
    { 'federation_federation_resolvers.federation' => $federation_id },
    { prefetch => 'federation_federation_resolvers' } )->all();

  return @resolvers;
}

sub find_delivery_service {
  my $self   = shift;
  my $xml_id = shift;
  my $data   = shift;
  my $ds;

  foreach my $service ( @{$data} ) {
    if ( $service->{'deliveryService'} eq $xml_id ) {
      $ds = $service;
    }
  }
  return $ds;
}

sub add_delivery_service {
  my $self   = shift;
  my $xml_id = shift;
  my $m      = shift;
  my $data   = shift;

  my $map;
  push( @{$map}, $m );
  push(
    @{$data},
    {   "deliveryService" => $xml_id,
      "mappings"        => $map
    }
  );
  return $data;
}

sub update_delivery_service {
  my $self = shift;
  my $ds   = shift;
  my $m    = shift;

  my $map = $ds->{'mappings'};
  push( @{$map}, $m );
  $ds->{'mappings'} = $map;
}

sub index {
  my $self = shift;
  my $data;

  if ( !&is_admin($self) ) {
    return $self->alert(
      "You must be an ADMIN to perform this operation!");
  }

  my $rs_data = $self->find_federations();
  while ( my $row = $rs_data->next ) {
    my $mapping;
    $mapping->{'cname'} = $row->federation->cname;
    $mapping->{'ttl'}   = $row->federation->ttl;

    my $federation_id = $row->federation->id;
    my @resolvers     = $self->find_federation_resolvers($federation_id);
    for my $resolver (@resolvers) {
      my $type = lc $resolver->type->name;
      if ( defined $mapping->{$type} ) {
        push( $mapping->{$type}, $resolver->ip_address );
      }
      else {
        @{ $mapping->{$type} } = ();
        push( $mapping->{$type}, $resolver->ip_address );
      }
    }

    my $xml_id = $row->deliveryservice->xml_id;
    if ( defined $data ) {
      my $ds = $self->find_delivery_service( $xml_id, $data );
      if ( !defined $ds ) {
        $data
          = $self->add_delivery_service( $xml_id, $mapping, $data );
      }
      else {
        $self->update_delivery_service( $ds, $mapping );
      }
    }
    else {
      $data = $self->add_delivery_service( $xml_id, $mapping, $data );
    }
  }
  $self->success($data);
}

sub find_federation_tmuser {
  my $self             = shift;
  my $current_username = shift;
  my @federation_ids;

  my $tm_user = $self->find_tmuser($current_username);
  if ( defined $tm_user ) {
    @federation_ids = $self->db->resultset('FederationTmuser')->search(
      {   tm_user => $tm_user->id,
        role    => $tm_user->role->id
      },
    )->get_column('federation')->all();
  }
  else {
    return $self->alert(
      "You must be a Federation user to perform this operation!");
  }

  return @federation_ids;
}

sub external_index {
  my $self             = shift;
  my $current_username = $self->current_user()->{username};
  my $data;

  my $rs_data;
  if ( &is_admin($self) ) {
    $rs_data = $self->find_federations();
  }
  else {
    my @federation_ids = $self->find_federation_tmuser($current_username);
    if ( scalar @federation_ids ) {
      $rs_data = $self->find_federations( \@federation_ids );
    }
    else {
      return $self->alert("No federations assigned to user.");
    }
  }

  while ( my $row = $rs_data->next ) {
    my $mapping;
    $mapping->{'cname'} = $row->federation->cname;
    $mapping->{'ttl'}   = $row->federation->ttl;

    my $federation_id = $row->federation->id;
    my @resolvers     = $self->find_federation_resolvers($federation_id);
    for my $resolver (@resolvers) {
      my $type = lc $resolver->type->name;
      if ( defined $mapping->{$type} ) {
        push( $mapping->{$type}, $resolver->ip_address );
      }
      else {
        @{ $mapping->{$type} } = ();
        push( $mapping->{$type}, $resolver->ip_address );
      }
    }

    my $xml_id = $row->deliveryservice->xml_id;
    if ( defined $data ) {
      my $ds = $self->find_delivery_service( $xml_id, $data );
      if ( !defined $ds ) {
        $data
          = $self->add_delivery_service( $xml_id, $mapping, $data );
      }
      else {
        $self->update_delivery_service( $ds, $mapping );
      }
    }
    else {
      $data = $self->add_delivery_service( $xml_id, $mapping, $data );
    }
  }
  $self->success( \@{$data} );
}

sub find_federation_deliveryservice {
  my $self   = shift;
  my $user   = shift;
  my $xml_id = shift;
  my @federation_ids;

  my @ids = $self->db->resultset('FederationTmuser')
    ->search( { tm_user => $user->id } )->get_column('federation')->all();
  my $ds = $self->db->resultset('Deliveryservice')
    ->search( { xml_id => $xml_id } )->get_column('id')->single();

  if ( scalar @ids ) {
    @federation_ids
      = $self->db->resultset('FederationDeliveryservice')->search(
      {   deliveryservice => $ds,
        federation      => { -in => \@ids }
      },
      { prefetch => 'federation' }
      )->get_column('federation.id')->all();
  }

  return @federation_ids;
}

sub add_resolvers {
  my $self          = shift;
  my $resolve4      = shift;
  my $resolve6      = shift;
  my $federation_id = shift;

  if ( defined $resolve4 ) {
    my $invalid_ip
      = $self->add_federation_resolver( $resolve4, $federation_id,
      "resolve4" );
    if ( defined($invalid_ip) ) {
      return $self->alert("$invalid_ip is not a valid ipv4 address.");
    }
  }

  if ( defined $resolve6 ) {
    my $invalid_ip
      = $self->add_federation_resolver( $resolve6, $federation_id,
      "resolve6" );
    if ( defined($invalid_ip) ) {
      return $self->alert("$invalid_ip is not a valid ipv4 address.");
    }
  }
}

sub add_federation_resolver {
  my $self          = shift;
  my $resolvers     = shift;
  my $federation_id = shift;
  my $type_name     = shift;
  my $resolver;

  foreach my $r ( @{$resolvers} ) {
    for my $ip ($r) {
      my $tmp_ip = $ip;
      $ip = Net::CIDR::cidrvalidate($ip);
      if ( !defined $ip ) {
        return $tmp_ip;
      }

      $resolver
        = $self->db->resultset('FederationResolver')->find_or_create(
        {   ip_address => $ip,
          type       => $self->db->resultset('Type')
            ->search( { name => $type_name } )->get_column('id')
            ->single()
        }
        );

      if ( defined $resolver ) {

        $self->db->resultset('FederationFederationResolver')
          ->find_or_create(
          {   federation          => $federation_id,
            federation_resolver => $resolver->id
          }
          );
      }
    }
  }
  return undef;
}

sub add {
  my $self = shift;

  my $current_username = $self->current_user()->{username};
  my $user             = $self->find_tmuser($current_username);
  if ( !defined $user ) {
    return $self->alert(
      "You must be an Federation user to perform this operation!");
  }

  my $federations = $self->req->json->{'federations'};
  foreach my $ds ( @{$federations} ) {
    my $xml_id   = $ds->{'deliveryService'};
    my $mappings = $ds->{'mappings'};

    my @federation_ids
      = $self->find_federation_deliveryservice( $user, $xml_id );

    foreach my $federation_id (@federation_ids) {
      foreach my $map ( @{$mappings} ) {
        my $resolve4 = $map->{'resolve4'};
        my $resolve6 = $map->{'resolve6'};
        $self->add_resolvers( $resolve4, $resolve6, $federation_id );
      }
    }
  }

  $self->success("Successfully created federations");
}

sub delete_federation_resolver {
  my $self = shift;
  my $user = shift;
  my $deleted_federation_resolver;

  my @federation_ids = $self->db->resultset('FederationTmuser')
    ->search( { tm_user => $user->id } )->get_column('federation')->all();

  if ( scalar @federation_ids ) {
    $deleted_federation_resolver
      = $self->db->resultset('FederationResolver')->search(
      {   'federation_federation_resolvers.federation' =>
          { -in => \@federation_ids }
      },
      { prefetch => 'federation_federation_resolvers' }
      )->delete();
  }
  return $deleted_federation_resolver;
}

sub delete {
  my $self             = shift;
  my $current_username = $self->current_user()->{username};

  my $user = $self->find_tmuser($current_username);
  if ( !defined $user ) {
    return $self->alert(
      "You must be an Federation user to perform this operation!");
  }

  my $deleted_federation_resolver
    = $self->delete_federation_resolver($user);
  if ( !$deleted_federation_resolver ) {
    return $self->alert(
      "No federation resolvers to delete for user $current_username.");
  }

  $self->success(
    "Successfully deleted all federation resolvers for user $current_username."
  );
}

sub update {
  my $self             = shift;
  my $current_username = $self->current_user()->{username};

  $self->delete();
  $self->app->log->info(
    "Successfully deleted all federations for user $current_username.");
  $self->add();
}

1;
