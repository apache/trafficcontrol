use utf8;
package Schema::Result::Federation;

# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
# 
#   http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.


# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::Federation

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<federation>

=cut

__PACKAGE__->table("federation");

=head1 ACCESSORS

=head2 id

  data_type: 'integer'
  is_auto_increment: 1
  is_nullable: 0

=head2 cname

  data_type: 'varchar'
  is_nullable: 0
  size: 1024

=head2 description

  data_type: 'varchar'
  is_nullable: 1
  size: 1024

=head2 ttl

  data_type: 'integer'
  is_nullable: 0

=head2 last_updated

  data_type: 'timestamp'
  datetime_undef_if_invalid: 1
  default_value: current_timestamp
  is_nullable: 1

=cut

__PACKAGE__->add_columns(
  "id",
  { data_type => "integer", is_auto_increment => 1, is_nullable => 0 },
  "cname",
  { data_type => "varchar", is_nullable => 0, size => 1024 },
  "description",
  { data_type => "varchar", is_nullable => 1, size => 1024 },
  "ttl",
  { data_type => "integer", is_nullable => 0 },
  "last_updated",
  {
    data_type => "timestamp",
    datetime_undef_if_invalid => 1,
    default_value => \"current_timestamp",
    is_nullable => 1,
  },
);

=head1 PRIMARY KEY

=over 4

=item * L</id>

=back

=cut

__PACKAGE__->set_primary_key("id");

=head1 RELATIONS

=head2 federation_deliveryservices

Type: has_many

Related object: L<Schema::Result::FederationDeliveryservice>

=cut

__PACKAGE__->has_many(
  "federation_deliveryservices",
  "Schema::Result::FederationDeliveryservice",
  { "foreign.federation" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 federation_federation_resolvers

Type: has_many

Related object: L<Schema::Result::FederationFederationResolver>

=cut

__PACKAGE__->has_many(
  "federation_federation_resolvers",
  "Schema::Result::FederationFederationResolver",
  { "foreign.federation" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 federation_tmusers

Type: has_many

Related object: L<Schema::Result::FederationTmuser>

=cut

__PACKAGE__->has_many(
  "federation_tmusers",
  "Schema::Result::FederationTmuser",
  { "foreign.federation" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);


# Created by DBIx::Class::Schema::Loader v0.07042 @ 2015-10-01 14:21:23
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:z0P4j54QodUN1uxCvy3iRQ


# You can replace this text with custom code or comments, and it will be preserved on regeneration
1;
