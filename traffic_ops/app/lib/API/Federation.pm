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
  my @data;
  my $foo;
  my $bar;
  my $orderby = $self->param('orderby') || "name";
  my $rs_data = $self->db->resultset("FederationMapping")->search(
    undef,
    {   prefetch => ['federation_resolver'],
      order_by => $orderby
    }
  );

  my $resolverKey;

  while ( my $row = $rs_data->next ) {
    if ( is_ipv4( $row->federation_resolver->ip_address ) ) {
      $resolverKey = "resolvers4";
    }
    else {
      $resolverKey = "resolvers6";
    }

    # push(
    #   @data,
    #   {   "cname"      => $row->cname,
    #     "ttl"        => $row->ttl,
    #     $resolverKey => $row->federation_resolver->ip_address,
    #   }
    # );

    $bar->{"cname"}      = $row->cname;
    $bar->{"ttl"}        = $row->ttl;
    $bar->{$resolverKey} = $row->federation_resolver->ip_address;

    $foo->{"id"}      = $row->id;
    $foo->{"mapping"} = $bar;
  }

  $self->success( $foo );
}
1;
